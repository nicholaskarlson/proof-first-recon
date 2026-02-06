package tests

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholaskarlson/proof-first-recon/internal/core"
	"github.com/nicholaskarlson/proof-first-recon/internal/report"
)

func TestGoldenCase01(t *testing.T) {
	t.Parallel()

	leftPath := filepath.Join("..", "fixtures", "input", "case01", "left.csv")
	rightPath := filepath.Join("..", "fixtures", "input", "case01", "right.csv")
	expectedDir := filepath.Join("..", "fixtures", "expected", "case01")

	outDir := t.TempDir()

	left, right, res, err := core.ReconcileFromPaths(leftPath, rightPath)
	if err != nil {
		t.Fatalf("ReconcileFromPaths: %v", err)
	}

	paths := report.OutputPaths{
		MatchedCSV:        filepath.Join(outDir, "matched.csv"),
		UnmatchedLeftCSV:  filepath.Join(outDir, "unmatched_left.csv"),
		UnmatchedRightCSV: filepath.Join(outDir, "unmatched_right.csv"),
		SummaryJSON:       filepath.Join(outDir, "recon_summary.json"),
	}

	if err := report.WriteAll(paths, left, right, res); err != nil {
		t.Fatalf("WriteAll: %v", err)
	}

	wantFiles := []string{
		"matched.csv",
		"unmatched_left.csv",
		"unmatched_right.csv",
		"recon_summary.json",
	}

	for _, name := range wantFiles {
		wantPath := filepath.Join(expectedDir, name)
		gotPath := filepath.Join(outDir, name)

		want, err := os.ReadFile(wantPath)
		if err != nil {
			t.Fatalf("read expected %s: %v", name, err)
		}
		got, err := os.ReadFile(gotPath)
		if err != nil {
			t.Fatalf("read got %s: %v", name, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("mismatch in %s\n--- expected ---\n%s\n--- got ---\n%s\n",
				name, string(want), string(got))
		}
	}
}
