package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

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
	outRoot := fs.String("out", "./out", "output directory")
	_ = fs.Parse(args)

	inRoot := filepath.Join("fixtures", "input")
	entries, err := os.ReadDir(inRoot)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	cases := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			cases = append(cases, e.Name())
		}
	}
	sort.Strings(cases)
	if len(cases) == 0 {
		fmt.Println("Error: no fixture cases found under", inRoot)
		os.Exit(1)
	}

	for _, c := range cases {
		leftPath := filepath.Join(inRoot, c, "left.csv")
		rightPath := filepath.Join(inRoot, c, "right.csv")
		expDir := filepath.Join("fixtures", "expected", c)
		outDir := filepath.Join(*outRoot, c)

		_ = os.RemoveAll(outDir)

		wantErrPath := filepath.Join(expDir, "error.txt")
		if wantErr, errRead := os.ReadFile(wantErrPath); errRead == nil {
			// expected-fail case: error text must match fixtures/expected/<case>/error.txt byte-for-byte.
			_, _, _, gotErr := core.ReconcileFromPaths(leftPath, rightPath)
			if gotErr == nil {
				fmt.Printf("MISMATCH: %s expected failure but got success\n", c)
				os.Exit(1)
			}
			gotErrB := []byte(gotErr.Error() + "\n")
			if err := os.MkdirAll(outDir, 0o755); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if err := os.WriteFile(filepath.Join(outDir, "error.txt"), gotErrB, 0o644); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if !bytes.Equal(gotErrB, wantErr) {
				fmt.Printf("MISMATCH: %s (error.txt)\n", c)
				os.Exit(1)
			}
			continue
		} else if errRead != nil && !os.IsNotExist(errRead) {
			fmt.Println("Error:", errRead)
			os.Exit(1)
		}

		if err := run(leftPath, rightPath, outDir); err != nil {
			fmt.Printf("Error: %s: %v\n", c, err)
			os.Exit(1)
		}
		if err := verifyCaseOutputs(c, outDir); err != nil {
			fmt.Println("MISMATCH:", err)
			os.Exit(1)
		}
	}

	fmt.Printf("OK: demo outputs match fixtures (%d case(s))\n", len(cases))
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

func verifyCaseOutputs(caseName, outDir string) error {
	expDir := filepath.Join("fixtures", "expected", caseName)
	files := []string{"matched.csv", "unmatched_left.csv", "unmatched_right.csv", "recon_summary.json"}
	for _, name := range files {
		expPath := filepath.Join(expDir, name)
		gotPath := filepath.Join(outDir, name)
		exp, err := os.ReadFile(expPath)
		if err != nil {
			return fmt.Errorf("%s: read expected %s: %w", caseName, name, err)
		}
		got, err := os.ReadFile(gotPath)
		if err != nil {
			return fmt.Errorf("%s: read output %s: %w", caseName, name, err)
		}
		if !bytes.Equal(got, exp) {
			return fmt.Errorf("%s: %s mismatch", caseName, name)
		}
	}
	return nil
}
