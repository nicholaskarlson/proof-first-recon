package tests

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholaskarlson/proof-first-recon/internal/core"
)

func TestFail_DuplicateID_Left(t *testing.T) {
	t.Parallel()
	left := filepath.Join("..", "fixtures", "input_fail", "dup_left.csv")
	right := filepath.Join("..", "fixtures", "input_fail", "dup_right.csv")
	_, _, _, err := core.ReconcileFromPaths(left, right)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate id") {
		t.Fatalf("expected duplicate id error, got: %v", err)
	}
}

func TestFail_BadHeader(t *testing.T) {
	t.Parallel()
	left := filepath.Join("..", "fixtures", "input_fail", "bad_header_left.csv")
	right := filepath.Join("..", "fixtures", "input_fail", "bad_header_right.csv")
	_, _, _, err := core.ReconcileFromPaths(left, right)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "header") {
		t.Fatalf("expected header error, got: %v", err)
	}
}

func TestFail_BadAmount(t *testing.T) {
	t.Parallel()
	left := filepath.Join("..", "fixtures", "input_fail", "bad_amount_left.csv")
	right := filepath.Join("..", "fixtures", "input_fail", "bad_amount_right.csv")
	_, _, _, err := core.ReconcileFromPaths(left, right)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "decimal") {
		t.Fatalf("expected decimal error, got: %v", err)
	}
}
