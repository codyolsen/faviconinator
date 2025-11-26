package generate

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

// Benchmarks full image generation to measure real-world scaling.
func BenchmarkGenerateImages(b *testing.B) {
	if _, err := exec.LookPath("magick"); err != nil {
		b.Skip("ImageMagick 'magick' not found on PATH; install to run benchmarks")
	}

	run := func(name string, workers int) {
		b.Run(name, func(b *testing.B) {
			base := b.TempDir()
			for i := 0; i < b.N; i++ {
				iterDir := filepath.Join(base, strconv.Itoa(i))
				if err := os.MkdirAll(iterDir, 0o755); err != nil {
					b.Fatalf("mkdir: %v", err)
				}

				svgPath := filepath.Join(iterDir, "icon.svg")
				if err := os.WriteFile(svgPath, complexHexagonSVG, 0o644); err != nil {
					b.Fatalf("write svg: %v", err)
				}

				outDir := filepath.Join(iterDir, "out")
				_, err := Generate(context.Background(), Options{
					Input:     svgPath,
					OutputDir: outDir,
					Workers:   workers,
				})
				if err != nil {
					b.Fatalf("generate: %v", err)
				}
			}
		})
	}

	run("workers=1", 1)
	run(fmt.Sprintf("workers=%d", runtime.NumCPU()), 0)
}

//go:embed testdata/complexahexagon.svg
var complexHexagonSVG []byte
