package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicholaskarlson/proof-first-recon/internal/core"
	"github.com/nicholaskarlson/proof-first-recon/internal/report"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "run":
		runCmd(os.Args[2:])
	case "demo":
		demoCmd(os.Args[2:])
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Println("Unknown command:", os.Args[1])
		fmt.Println()
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println("proof-first-recon")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  recon run  --left <file.csv> --right <file.csv> --out <dir>")
	fmt.Println("  recon demo --out <dir>")
	fmt.Println()
	fmt.Println("v0 rule: exact match by id")
}

func runCmd(args []string) {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	leftPath := fs.String("left", "", "left CSV path")
	rightPath := fs.String("right", "", "right CSV path")
	outDir := fs.String("out", "./out", "output directory")
	_ = fs.Parse(args)

	if *leftPath == "" || *rightPath == "" {
		fmt.Println("Error: --left and --right are required")
		fmt.Println()
		usage()
		os.Exit(2)
	}

	if err := run(*leftPath, *rightPath, *outDir); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Wrote outputs to", *outDir)
}

func demoCmd(args []string) {
	fs := flag.NewFlagSet("demo", flag.ExitOnError)
	outDir := fs.String("out", "./out", "output directory")
	_ = fs.Parse(args)

	leftPath := filepath.Join("fixtures", "input", "case01", "left.csv")
	rightPath := filepath.Join("fixtures", "input", "case01", "right.csv")

	if err := run(leftPath, rightPath, *outDir); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Demo complete. Wrote outputs to", *outDir)
}

func run(leftPath, rightPath, outDir string) error {
	left, right, res, err := core.ReconcileFromPaths(leftPath, rightPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	paths := report.OutputPaths{
		MatchedCSV:        filepath.Join(outDir, "matched.csv"),
		UnmatchedLeftCSV:  filepath.Join(outDir, "unmatched_left.csv"),
		UnmatchedRightCSV: filepath.Join(outDir, "unmatched_right.csv"),
		SummaryJSON:       filepath.Join(outDir, "recon_summary.json"),
	}

	return report.WriteAll(paths, left, right, res)
}
