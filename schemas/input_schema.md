# Input schema (v0)

This tool reconciles two CSV exports: **left** and **right**.

## Required header (exact)
The first row must be exactly:

`id,date,amount,description`

## Column meanings
- `id` (string): unique transaction identifier within the file (v0 requires uniqueness; duplicates are an error)
- `date` (string): `YYYY-MM-DD` (validated lightly in v0; stricter later)
- `amount` (decimal): signed amount with up to 2 decimal places (parsed as integer cents)
- `description` (string): free text

## Notes
- UTF-8 CSV.
- LF line endings recommended (repo enforces LF via `.gitattributes`).
