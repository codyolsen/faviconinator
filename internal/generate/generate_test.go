package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteFaviconSVGStripsWhiteFills(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "icon.svg")
	dst := filepath.Join(tmp, "favicon.svg")

	original := `<svg fill="#fff"><rect fill="#ffffff"/><circle fill="white"/></svg>`
	if err := os.WriteFile(src, []byte(original), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}

	if err := writeFaviconSVG(src, dst); err != nil {
		t.Fatalf("writeFaviconSVG: %v", err)
	}

	out, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}

	got := string(out)
	if containsAny(got, []string{`fill="#fff"`, `fill="#ffffff"`, `fill="white"`}) {
		t.Fatalf("expected white fills to be removed, got %q", got)
	}
	if !containsAny(got, []string{`fill="none"`}) {
		t.Fatalf("expected fill=\"none\" placeholder, got %q", got)
	}
}

func TestOptionsValidate(t *testing.T) {
	tmp := t.TempDir()
	in := filepath.Join(tmp, "in.svg")
	if err := os.WriteFile(in, []byte("<svg/>"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	opts := Options{
		Input:     in,
		OutputDir: filepath.Join(tmp, "out"),
	}

	if err := opts.validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}

	opts.Input = tmp // directory
	if err := opts.validate(); err == nil {
		t.Fatalf("expected directory input to fail validation")
	}

	opts = Options{
		Input:     in,
		OutputDir: filepath.Join(tmp, "out"),
		Workers:   -1,
	}
	if err := opts.validate(); err == nil {
		t.Fatalf("expected negative workers to fail validation")
	}
}

func containsAny(s string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}
