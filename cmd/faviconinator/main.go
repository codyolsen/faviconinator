package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/codyolsen/faviconinator/internal/generate"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "faviconinator: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		outDir      string
		color       string
		verbose     bool
		showVersion bool
		jobs        int
		jsonOut     bool
		renderer    string
	)

	flag.StringVar(&outDir, "out", "", "output directory (default: build/<input basename>)")
	flag.StringVar(&color, "color", "", "(reserved) hex color to tint the icon")
	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.BoolVar(&verbose, "verbose", false, "verbose logging")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.IntVar(&jobs, "jobs", 0, "number of concurrent workers (0 = NumCPU)")
	flag.BoolVar(&jsonOut, "json", false, "print result stats as JSON")
	flag.StringVar(&renderer, "renderer", "auto", "png renderer: auto (default), magick, or rsvg (uses rsvg-convert)")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] input.svg\n\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Generate a full favicon set from a single SVG. ImageMagick 7+ must be installed.")
		fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if showVersion {
		fmt.Println(version)
		return nil
	}

	if flag.NArg() != 1 {
		flag.Usage()
		return fmt.Errorf("input SVG is required")
	}

	input := flag.Arg(0)
	if outDir == "" {
		base := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
		outDir = filepath.Join("build", base)
	}
	ctx := context.Background()

	opts := generate.Options{
		Input:     input,
		OutputDir: outDir,
		Color:     color,
		Verbose:   verbose,
		Workers:   jobs,
		Renderer:  renderer,
	}

	start := time.Now()
	stats, err := generate.Generate(ctx, opts)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	if jsonOut {
		type payload struct {
			Files      int           `json:"files"`
			Outputs    []string      `json:"outputs"`
			Workers    int           `json:"workers"`
			DurationMs float64       `json:"duration_ms"`
			Elapsed    time.Duration `json:"elapsed_ns"`
			Version    string        `json:"version"`
		}
		out := payload{
			Files:      stats.Files,
			Outputs:    stats.Outputs,
			Workers:    stats.Workers,
			DurationMs: float64(elapsed.Microseconds()) / 1000.0,
			Elapsed:    elapsed,
			Version:    version,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	fmt.Fprintf(os.Stderr, "Generated %d assets to %s in %s (workers=%d)\n",
		stats.Files, opts.OutputDir, elapsed.Truncate(time.Millisecond), stats.Workers)
	return nil
}
