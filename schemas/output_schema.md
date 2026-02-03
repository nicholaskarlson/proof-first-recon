# Output schema (v0)

Outputs are written to the chosen output directory.

## Files
- `matched.csv`
  - Rows that matched between left and right by v0 rule: **exact match on `id`**.
- `unmatched_left.csv`
  - Rows present only in left.
- `unmatched_right.csv`
  - Rows present only in right.
- `recon_summary.json`
  - Summary counts and totals (stored as cents and formatted strings).

## matched.csv columns
`id,left_date,left_amount,left_description,right_date,right_amount,right_description`

## Determinism rules
- Output rows are sorted by `id` ascending.
- Amounts are formatted to 2 decimal places.
- Summary JSON is stable (no timestamps, fixed key order).
- Same inputs -> same outputs.
