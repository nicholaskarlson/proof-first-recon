# Output schema (v0)

Outputs are written to the specified output directory.

## Files
- `matched.csv`
  - Rows that matched between left and right by the chosen rule (v0: exact match by `id`).
- `unmatched_left.csv`
  - Rows present only in left.
- `unmatched_right.csv`
  - Rows present only in right.
- `recon_summary.json`
  - Summary counts and totals.

## Determinism rules
- Output rows are sorted by `id` ascending.
- Numeric formatting is stable.
- Same inputs -> same outputs.
