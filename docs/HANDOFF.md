# HANDOFF — proof-first-recon (v0)

This repo is intentionally small and **handoff-friendly**: a deterministic CSV reconciliation CLI with fixtures, golden tests, and fail-fast validation.

---

## What you’re getting

A command-line tool that reconciles two CSV exports and writes **provable**, deterministic outputs:

- `matched.csv`
- `unmatched_left.csv`
- `unmatched_right.csv`
- `recon_summary.json`

**v0 rule:** exact match on `id`.

---

## Requirements

- Go (recent stable version)
- Optional: `make` (only for convenience targets)

No database. No services. No Docker required.

---

## Canonical commands

```bash
# Proof gate (one command)
make verify

# Proof gates (portable, no Makefile)
go test -count=1 ./...
go run ./cmd/recon demo --out ./out
```

## How to validate (the acceptance gate)

From repo root:

```bash
go test -count=1 ./...
```

Or:

```bash
make test
```

If tests pass, the tool’s behavior matches the checked-in fixtures.

---

## How to run

### Demo (recomputes fixtures and verifies outputs match goldens)

```bash
go run ./cmd/recon demo --out ./out
ls -la ./out/case01
```

### Optional: Python check (stdlib only)

```bash
python3 examples/python/verify_recon_case.py
```

This is an independent verification lane (no third-party deps).

### Your own files

```bash
go run ./cmd/recon run --left path/to/left.csv --right path/to/right.csv --out ./out
```

---

## Input contract (v0)

Both CSV files must be UTF-8 with **exactly** these headers:

- `id` (string)
- `date` (YYYY-MM-DD)
- `amount` (signed, 2-decimal money)
- `description` (string)

See `schemas/input_schema.md`.

---

## Fail-fast safety rules (by design)

The tool rejects inputs rather than guessing:

- Duplicate `id` values in a single file → **error**
- Header mismatch → **error**
- Amount with more than 2 decimal places → **error**

Examples live in `fixtures/input/<case>/` with expected errors in `fixtures/expected/<case>/error.txt` (enforced by `tests/failures_test.go`).

---

## Output contract (v0)

Outputs are written to the directory you specify via `--out`.

- `matched.csv`  
  Rows matched between left/right by `id`.

- `unmatched_left.csv`  
  Rows only present in left.

- `unmatched_right.csv`  
  Rows only present in right.

- `recon_summary.json`  
  Counts and totals (in cents + formatted strings).

See `schemas/output_schema.md`.

---

## Determinism guarantees

- Output rows are sorted by `id` ascending.
- Money is parsed and stored in **cents** (int64) to avoid floating point drift.
- The same inputs produce byte-stable outputs (verified via golden tests).

---

## Troubleshooting (common)

- **“header mismatch”**  
  Your CSV header row doesn’t exactly match `id,date,amount,description`.

- **“duplicate id”**  
  A file contains the same `id` more than once (v0 requires unique ids per file).

- **“decimal” / “more than 2 decimal places”**  
  Money values must be at most 2 decimals.

---

## If you change behavior (how to update fixtures safely)

Golden fixtures define correctness. If you intentionally change behavior:

1) Run the demo to regenerate outputs:
   ```bash
   go run ./cmd/recon demo --out ./out
   ```
2) Update expected fixtures intentionally:
   ```bash
   cp -f out/case01/*.csv fixtures/expected/case01/
   cp -f out/case01/recon_summary.json fixtures/expected/case01/
   ```
3) Re-run tests:
   ```bash
   go test -count=1 ./...
   ```
4) Commit both the code change and the fixture update in the same PR.

---

## Where to look next

- CLI entrypoint: `cmd/recon/main.go`
- Core logic: `internal/core/`
- CSV parsing + money parsing: `internal/ingest/`
- Writers: `internal/report/`
- Tests: `tests/`
- Fixtures: `fixtures/`

This repo is designed so a competent Go developer can extend it without needing the original author.
