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
)

// Options controls how assets are generated.
type Options struct {
	Input     string // path to the source SVG
	OutputDir string // directory to write generated files
	Color     string // optional hex color (e.g. #ff6600) to tint the icon
	Verbose   bool   // emit progress to stderr
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
func Generate(ctx context.Context, opts Options) error {
	if err := opts.validate(); err != nil {
		return err
	}

	if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	magickPath, err := exec.LookPath("magick")
	if err != nil {
		return errors.New("ImageMagick 'magick' binary not found on PATH; install ImageMagick 7+")
	}

	for _, icon := range pngIcons {
		out := filepath.Join(opts.OutputDir, icon.Name)
		if err := renderPNG(ctx, magickPath, opts.Input, out, icon.Size, opts.Color); err != nil {
			return err
		}
		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "wrote %s\n", out)
		}
	}

	icoOut := filepath.Join(opts.OutputDir, "favicon.ico")
	if err := renderICO(ctx, magickPath, opts.Input, icoOut, opts.Color); err != nil {
		return err
	}
	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "wrote %s\n", icoOut)
	}

	svgOut := filepath.Join(opts.OutputDir, "favicon.svg")
	if err := writeFaviconSVG(opts.Input, svgOut); err != nil {
		return fmt.Errorf("write favicon.svg: %w", err)
	}
	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "wrote %s\n", svgOut)
	}

	return nil
}

func renderPNG(ctx context.Context, magickPath, input, output string, size int, color string) error {
	args := []string{
		magickPath,
		"-background", "none",
		"-density", "512",
		input,
	}

	if color != "" {
		args = append(args, "-fill", color, "-colorize", "100")
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

func renderICO(ctx context.Context, magickPath, input, output string, color string) error {
	args := []string{
		magickPath,
		"-background", "none",
		"-density", "512",
		input,
	}

	if color != "" {
		args = append(args, "-fill", color, "-colorize", "100")
	}

	args = append(args,
		"-define", "icon:auto-resize=16,32,48,64",
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

// writeFaviconSVG copies the SVG and strips common white-background fills
// to make the favicon.svg transparent. This is a blunt text transform.
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

	data := buf.Bytes()

	// Replace typical "white background" fills with none.
	// If your foreground uses white, this will nuke that too.
	repls := [][]byte{
		[]byte(`fill="#ffffff"`),
		[]byte(`fill="#FFFFFF"`),
		[]byte(`fill="#fff"`),
		[]byte(`fill="#FFF"`),
		[]byte(`fill="white"`),
	}

	for _, old := range repls {
		data = bytes.ReplaceAll(data, old, []byte(`fill="none"`))
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := out.Write(data); err != nil {
		return err
	}
	return out.Sync()
}
