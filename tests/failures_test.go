package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholaskarlson/proof-first-recon/internal/core"
)

func TestExpectedFailCase02DupID(t *testing.T) {
	t.Parallel()
	caseName := "case02_dup_id"
	assertExpectedFailure(t, caseName)
}

func TestExpectedFailCase03BadHeader(t *testing.T) {
	t.Parallel()
	caseName := "case03_bad_header"
	assertExpectedFailure(t, caseName)
}

func TestExpectedFailCase04BadAmount(t *testing.T) {
	t.Parallel()
	caseName := "case04_bad_amount"
	assertExpectedFailure(t, caseName)
}

func assertExpectedFailure(t *testing.T, caseName string) {
	left := filepath.Join("..", "fixtures", "input", caseName, "left.csv")
	right := filepath.Join("..", "fixtures", "input", caseName, "right.csv")
	expPath := filepath.Join("..", "fixtures", "expected", caseName, "error.txt")

	expB, err := os.ReadFile(expPath)
	if err != nil {
		t.Fatalf("read expected error: %v", err)
	}
	exp := strings.TrimSpace(string(expB))

	_, _, _, gotErr := core.ReconcileFromPaths(left, right)
	if gotErr == nil {
		t.Fatalf("expected error, got nil")
	}
	got := strings.TrimSpace(gotErr.Error())
	if got != exp {
		t.Fatalf("error mismatch\n got: %s\n exp: %s", got, exp)
	}
}
