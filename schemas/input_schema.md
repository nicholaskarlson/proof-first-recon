# Input schema (v0)

This tool reconciles two CSV exports: **left** and **right**.

## Minimal required columns
Both input files must contain:

- `id` (string): unique transaction identifier within the file
- `date` (string): YYYY-MM-DD
- `amount` (decimal): signed amount (e.g., -12.34, 99.00)
- `description` (string): free text

## Notes
- Files should be UTF-8 CSV.
- Line endings should be LF (this repo enforces LF via `.gitattributes`).
- Additional columns may be allowed later, but v0 assumes exactly these columns.
