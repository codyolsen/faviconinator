package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	)

	flag.StringVar(&outDir, "out", "", "output directory (default: build/<input basename>)")
	flag.StringVar(&color, "color", "", "hex color (e.g. #ff6600) to tint the icon")
	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.BoolVar(&verbose, "verbose", false, "verbose logging")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")

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
	}

	return generate.Generate(ctx, opts)
}
