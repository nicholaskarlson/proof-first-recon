# proof-first-recon

Deterministic CSV reconciliation engine (Go-first) with fixtures and golden tests.

![ci](https://github.com/nicholaskarlson/proof-first-recon/actions/workflows/ci.yml/badge.svg)
![license](https://img.shields.io/badge/license-MIT-blue.svg)

## What this is
A small command-line tool that reconciles two CSV exports (e.g., bank vs. ledger, payouts vs. transactions)
and produces **provable** outputs:

- `matched.csv`
- `unmatched_left.csv`
- `unmatched_right.csv`
- `recon_summary.json`
- (optional later) `recon_report.html`

**Design goals**
- Same inputs → same outputs (deterministic)
- Clear acceptance tests (golden fixtures)
- Easy handoff (no maintenance trap)

