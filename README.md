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

## Quick start (recommended)

```bash
go test ./...

go run ./cmd/recon -mode help
go run ./cmd/recon -mode demo -out ./out
```

### Optional convenience

```bash
make test
make demo
```

## Schemas

- Inputs: `schemas/input_schema.md`
- Outputs: `schemas/output_schema.md`

## Repo layout (overview)

- `cmd/recon/` — CLI entrypoint
- `internal/` — engine + validation + report writers *(planned)*
- `fixtures/` — test inputs + expected outputs *(planned)*
- `tests/` — golden tests + invariants *(planned)*

## Status / next milestone

This repo is currently scaffolding. Next milestone:

- Add a small fixture dataset (including duplicates and mismatches)
- Implement `recon run` (v0: exact match by `id`)
- Add golden tests so `go test ./...` verifies byte-stable outputs

## License

MIT (see `LICENSE`).
