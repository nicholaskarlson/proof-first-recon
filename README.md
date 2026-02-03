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

## Quick start

### 1) Run tests (the acceptance gate)

```bash
go test ./...
```

### 2) Run the demo fixture

```bash
go run ./cmd/recon demo --out ./out
ls -la ./out
```

### 3) Reconcile your own CSVs

```bash
go run ./cmd/recon run --left path/to/left.csv --right path/to/right.csv --out ./out
```

### Optional convenience

```bash
make test
make demo
```

## Input / output schemas

- Inputs: `schemas/input_schema.md`
- Outputs: `schemas/output_schema.md`

## Safety: expected failure behavior

This tool is designed to **fail fast** rather than guess:

- Duplicate `id` values in a single file → error
- Header mismatch → error
- Amounts with more than 2 decimal places → error

See: `fixtures/input_fail/` and `tests/failures_test.go`.

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
