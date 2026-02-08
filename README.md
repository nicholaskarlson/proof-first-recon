# proof-first-recon

Deterministic CSV reconciliation engine (Go-first) with fixtures and golden outputs.

![ci](https://github.com/nicholaskarlson/proof-first-recon/actions/workflows/ci.yml/badge.svg)
![license](https://img.shields.io/badge/license-MIT-blue.svg)

> **Book:** *The Deterministic Finance Toolkit*  
> This repo is **Project 1 of 4**. The exact code referenced in the manuscript is tagged **[`book-v1`](https://github.com/nicholaskarlson/proof-first-recon/tree/book-v1)**.

## Toolkit navigation

- **[proof-first-recon](https://github.com/nicholaskarlson/proof-first-recon)** — deterministic CSV reconciliation (matched/unmatched + summary JSON)
- **[proof-first-auditpack](https://github.com/nicholaskarlson/proof-first-auditpack)** — deterministic audit packs (manifest.json + sha256 + verify)
- **[proof-first-normalizer](https://github.com/nicholaskarlson/proof-first-normalizer)** — deterministic CSV normalize + validate (schema → normalized.csv/errors.csv/report.json)
- **[proof-first-finance-calc](https://github.com/nicholaskarlson/proof-first-finance-calc)** — proof-first finance calc service (Amortization v1 API + demo)

## What it does

Reconciles **two CSV exports** (e.g., bank vs. ledger, payouts vs. transactions) and produces provable outputs:

- `matched.csv`
- `unmatched_left.csv`
- `unmatched_right.csv`
- `recon_summary.json`

## Design goals

- **Deterministic:** same inputs → same outputs
- **Testable:** fixtures + golden outputs define “done”
- **Handoff-friendly:** simple CLI + documented schemas (no maintenance trap)

## Quick start

Requirements:
- Go **1.22+**
- GNU Make (optional, but recommended)

```bash
# One-command proof gate
make verify

# Portable proof gate (no Makefile)
go test -count=1 ./...
go run ./cmd/recon demo --out ./out
```


## Usage

```bash
# Demo: recomputes fixtures and verifies outputs match goldens
go run ./cmd/recon demo --out ./out
# outputs land in ./out/CASE/ (e.g., ./out/case01/)

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

See: `fixtures/input/case02_dup_id/`, `case03_bad_header/`, `case04_bad_amount/` and `fixtures/expected/CASE/error.txt` (enforced by `tests/failures_test.go`).

## Repo layout (overview)

- `cmd/recon/` — CLI entrypoint (`demo`, `run`)
- `internal/core/` — reconciliation logic (no file I/O)
- `internal/report/` — deterministic CSV/JSON writers
- `internal/model/` — shared types (breaks import cycles cleanly)
- `internal/ingest/` — CSV parsing + cents-safe money parsing
- `fixtures/` — test inputs and golden expected outputs
- `tests/` — golden + expected-failure tests

## Status

**v0 rule:** exact match on `id`.

What’s already proven by the proof gate:

- Golden outputs match `fixtures/expected/*` byte-for-byte
- Invalid inputs fail with clear errors (duplicate id, bad header, bad amount)

## Determinism contract

This project is intentionally “boring” in the best way: the same inputs must produce the same outputs.

See: **[`docs/CONVENTIONS.md`](docs/CONVENTIONS.md)** (rounding, ordering, LF, atomic writes, stable JSON, etc.).


## Handoff / maintenance

See: **[`docs/HANDOFF.md`](docs/HANDOFF.md)** (acceptance gates, troubleshooting, and “what to change (and what not to)”).


## License

MIT (see `LICENSE`).

