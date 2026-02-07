# proof-first-recon

Deterministic CSV reconciliation engine (Go-first) with fixtures and golden tests.

![ci](https://github.com/nicholaskarlson/proof-first-recon/actions/workflows/ci.yml/badge.svg)
![license](https://img.shields.io/badge/license-MIT-blue.svg)

## What this is

A small command-line tool that reconciles **two CSV exports** (e.g., bank vs. ledger, payouts vs. transactions) and produces **provable** outputs:

- `matched.csv`
- `unmatched_left.csv`
- `unmatched_right.csv`
- `recon_summary.json`
- *(optional later)* `recon_report.html`

## Design goals

- **Deterministic:** same inputs → same outputs
- **Testable:** fixtures + golden tests define “done”
- **Handoff-friendly:** simple CLI + documented schemas (no maintenance trap)

## Canonical commands

```bash
# Proof gate (one command)
make verify

# Proof gates (portable, no Makefile)
go test -count=1 ./...
go run ./cmd/recon demo --out ./out
```

## Usage

```bash
# Demo (writes outputs and verifies they match fixtures)
go run ./cmd/recon demo --out ./out

# Reconcile your own files
go run ./cmd/recon run --left path/to/left.csv --right path/to/right.csv --out ./out
```

## Input / output schemas

- Inputs: `schemas/input_schema.md`
- Outputs: `schemas/output_schema.md`

## Safety: expected failure behavior

This tool is designed to **fail fast** rather than guess:

- Duplicate `id` values in a single file → error
- Header mismatch → error
- Amounts with more than 2 decimal places → error

See: `fixtures/input/case02_dup_id/`, `case03_bad_header/`, `case04_bad_amount/` and `fixtures/expected/<case>/error.txt` (enforced by `tests/failures_test.go`).

## Repo layout (overview)

- `cmd/recon/` — CLI entrypoint (`demo`, `run`)
- `internal/core/` — reconciliation logic (no file I/O)
- `internal/report/` — deterministic CSV/JSON writers
- `internal/model/` — shared types (breaks import cycles cleanly)
- `internal/ingest/` — CSV parsing + cents-safe money parsing
- `fixtures/` — test inputs and golden expected outputs
- `tests/` — golden test(s) + expected-failure tests

## Status

**v0 rule:** exact match on `id`.

What’s already proven by tests:

- Golden outputs match `fixtures/expected/*` byte-for-byte
- Invalid inputs fail with clear errors (duplicate id, bad header, bad amount)

## License

MIT (see `LICENSE`).
