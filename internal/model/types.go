package model

import "github.com/nicholaskarlson/proof-first-recon/internal/ingest"

type Matched struct {
	ID    string
	Left  ingest.Record
	Right ingest.Record
}

type Result struct {
	Matched        []Matched
	UnmatchedLeft  []ingest.Record
	UnmatchedRight []ingest.Record
}
