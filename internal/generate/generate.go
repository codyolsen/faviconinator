package generate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Options controls how assets are generated.
type Options struct {
	Input     string // path to the source SVG
	OutputDir string // directory to write generated files
	Color     string // (unused for now) optional hex color to tint the icon
	Verbose   bool   // emit progress to stderr
	Workers   int    // number of concurrent workers (0 = NumCPU)
	Renderer  string // png renderer: auto (default), magick, or rsvg
}

func (o Options) validate() error {
	if o.Input == "" {
		return errors.New("input SVG path is required")
	}
	info, err := os.Stat(o.Input)
	if err != nil {
		return fmt.Errorf("input: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("input %s is a directory", o.Input)
	}
	if o.OutputDir == "" {
		return errors.New("output directory is required")
	}
	if o.Workers < 0 {
		return errors.New("workers cannot be negative")
	}
	if o.Renderer != "" && o.Renderer != "auto" && o.Renderer != "magick" && o.Renderer != "rsvg" {
		return fmt.Errorf("renderer must be auto, magick, or rsvg (got %q)", o.Renderer)
	}
	return nil
}

type iconSpec struct {
	Size int
	Name string
}

var pngIcons = []iconSpec{
	// Android / Chrome
	{192, "android-chrome-192x192.png"},
	{512, "android-chrome-512x512.png"},
	{36, "android-icon-36x36.png"},
	{48, "android-icon-48x48.png"},
	{72, "android-icon-72x72.png"},
	{96, "android-icon-96x96.png"},
	{144, "android-icon-144x144.png"},
	{192, "android-icon-192x192.png"},

	// Apple icons
	{57, "apple-icon-57x57.png"},
	{60, "apple-icon-60x60.png"},
	{72, "apple-icon-72x72.png"},
	{76, "apple-icon-76x76.png"},
	{114, "apple-icon-114x114.png"},
	{120, "apple-icon-120x120.png"},
	{144, "apple-icon-144x144.png"},
	{152, "apple-icon-152x152.png"},
	{180, "apple-icon-180x180.png"},
	{180, "apple-icon.png"},
	{180, "apple-icon-precomposed.png"},
	{180, "apple-touch-icon.png"},

	// Favicons
	{32, "favicon-32x32.png"},
	{96, "favicon-96x96.png"},
	{96, "favicon.png"},

	// Microsoft tiles
	{70, "ms-icon-70x70.png"},
	{144, "ms-icon-144x144.png"},
	{150, "ms-icon-150x150.png"},
	{310, "ms-icon-310x310.png"},
}

// Generate creates favicons and related assets from the given SVG.
// It requires ImageMagick 7's `magick` binary to be available on PATH.
func Generate(ctx context.Context, opts Options) (Stats, error) {
	start := time.Now()
	stats := Stats{
		Workers: workerCount(opts.Workers),
	}

	if err := opts.validate(); err != nil {
		return stats, err
	}

	if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
		return stats, fmt.Errorf("create output dir: %w", err)
	}

	renderer, paths, err := pickRenderer(opts.Renderer)
	if err != nil {
		return stats, err
	}

	var mu sync.Mutex

	// Build task list so each output is generated concurrently.
	tasks := make([]func() error, 0, len(pngIcons)+2)

	for _, icon := range pngIcons {
		icon := icon // capture
		tasks = append(tasks, func() error {
			out := filepath.Join(opts.OutputDir, icon.Name)
			if err := renderPNG(ctx, renderer, paths, opts.Input, out, icon.Size, opts.Color); err != nil {
				return err
			}
			mu.Lock()
			stats.Outputs = append(stats.Outputs, out)
			mu.Unlock()
			if opts.Verbose {
				fmt.Fprintf(os.Stderr, "wrote %s\n", out)
			}
			return nil
		})
	}

	tasks = append(tasks, func() error {
		icoOut := filepath.Join(opts.OutputDir, "favicon.ico")
		if err := renderICO(ctx, renderer, paths, opts.Input, icoOut, opts.Color); err != nil {
			return err
		}
		mu.Lock()
		stats.Outputs = append(stats.Outputs, icoOut)
		mu.Unlock()
		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "wrote %s\n", icoOut)
		}
		return nil
	})

	tasks = append(tasks, func() error {
		svgOut := filepath.Join(opts.OutputDir, "favicon.svg")
		if err := writeFaviconSVG(opts.Input, svgOut); err != nil {
			return fmt.Errorf("write favicon.svg: %w", err)
		}
		mu.Lock()
		stats.Outputs = append(stats.Outputs, svgOut)
		mu.Unlock()
		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "wrote %s\n", svgOut)
		}
		return nil
	})

	if err := runConcurrent(ctx, tasks, stats.Workers); err != nil {
		return stats, err
	}

	stats.Files = len(stats.Outputs)
	stats.Duration = time.Since(start)
	return stats, nil
}

func renderPNG(ctx context.Context, renderer string, paths renderPaths, input, output string, size int, color string) error {
	if renderer == "rsvg" {
		return renderPNGRSVG(ctx, paths.rsvg, input, output, size)
	}

	args := []string{
		paths.magick,
		"-background", "none",
		"-density", "512",
		input,
	}

	args = append(args,
		"-resize", fmt.Sprintf("%dx%d", size, size),
		"-gravity", "center",
		"-extent", fmt.Sprintf("%dx%d", size, size),
		"-define", "png:color-type=6",
		output,
	)

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("render %s: %w", output, err)
	}
	return nil
}

func renderPNGRSVG(ctx context.Context, rsvgPath, input, output string, size int) error {
	args := []string{
		rsvgPath,
		"--background-color=transparent",
		"-w", strconv.Itoa(size),
		"-h", strconv.Itoa(size),
		"-o", output,
		input,
	}
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runConcurrent executes all tasks with a bounded worker pool.
func runConcurrent(ctx context.Context, tasks []func() error, workers int) error {
	workerCount := workerCount(workers)
	taskCh := make(chan func() error)

	var (
		wg   sync.WaitGroup
		once sync.Once
		err  error
	)

	worker := func() {
		defer wg.Done()
		for task := range taskCh {
			if ctx.Err() != nil {
				return
			}
			if e := task(); e != nil {
				once.Do(func() { err = e })
			}
		}
	}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	for _, task := range tasks {
		// Stop queueing new work if an error already occurred.
		if err != nil {
			break
		}
		taskCh <- task
	}
	close(taskCh)
	wg.Wait()

	if ctx.Err() != nil && err == nil {
		err = ctx.Err()
	}
	return err
}

func workerCount(requested int) int {
	if requested > 0 {
		return requested
	}
	n := runtime.NumCPU()
	if n < 2 {
		return 2
	}
	return n
}

// Stats reports generation outcomes.
type Stats struct {
	Files    int
	Duration time.Duration
	Workers  int
	Outputs  []string
}

type renderPaths struct {
	magick string
	rsvg   string
}

func pickRenderer(requested string) (string, renderPaths, error) {
	req := requested
	if req == "" {
		req = "auto"
	}

	var paths renderPaths

	lookup := func(name string) (string, bool) {
		p, err := exec.LookPath(name)
		return p, err == nil
	}

	rsvgPath, haveRSVG := lookup("rsvg-convert")
	magickPath, haveMagick := lookup("magick")

	switch req {
	case "rsvg":
		if !haveRSVG {
			return "", paths, errors.New("renderer 'rsvg' requested but rsvg-convert not found on PATH")
		}
		if !haveMagick {
			return "", paths, errors.New("ImageMagick 'magick' required for ICO output; install ImageMagick 7+")
		}
		return "rsvg", renderPaths{magick: magickPath, rsvg: rsvgPath}, nil
	case "magick":
		if !haveMagick {
			return "", paths, errors.New("renderer 'magick' requested but magick not found on PATH; install ImageMagick 7+")
		}
		return "magick", renderPaths{magick: magickPath}, nil
	case "auto":
		if haveRSVG && haveMagick {
			return "rsvg", renderPaths{magick: magickPath, rsvg: rsvgPath}, nil
		}
		if haveMagick {
			return "magick", renderPaths{magick: magickPath}, nil
		}
		return "", paths, errors.New("neither ImageMagick 'magick' nor 'rsvg-convert' found on PATH")
	default:
		return "", paths, fmt.Errorf("unknown renderer %q (use auto, magick, or rsvg)", requested)
	}
}

func renderICO(ctx context.Context, renderer string, paths renderPaths, input, output string, color string) error {
	tmp, err := os.CreateTemp("", "favicon-base-*.png")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer os.Remove(tmpPath)

	if err := renderPNG(ctx, renderer, paths, input, tmpPath, 512, color); err != nil {
		return err
	}

	args := []string{
		paths.magick,
		tmpPath,
		"-alpha", "on",
		"-define", "icon:auto-resize=16,32,48,64",
		output,
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("render %s: %w", output, err)
	}
	return nil
}

// writeFaviconSVG copies the original SVG verbatim to favicon.svg.
func writeFaviconSVG(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, in); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := out.Write(buf.Bytes()); err != nil {
		return err
	}
	return out.Sync()
}
