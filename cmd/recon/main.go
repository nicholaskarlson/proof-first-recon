package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		mode = flag.String("mode", "help", "help|demo (placeholder)")
		out  = flag.String("out", "./out", "output directory")
	)
	flag.Parse()

	switch *mode {
	case "help":
		fmt.Println("proof-first-recon (skeleton)")
		fmt.Println("")
		fmt.Println("Next: implement reconciliation engine + fixtures + golden tests.")
		fmt.Println("Try: go run ./cmd/recon -mode demo -out ./out")
	case "demo":
		_ = os.MkdirAll(*out, 0o755)
		fmt.Printf("Demo placeholder: would write outputs to %s\n", *out)
	default:
		fmt.Println("Unknown mode:", *mode)
		os.Exit(2)
	}
}
