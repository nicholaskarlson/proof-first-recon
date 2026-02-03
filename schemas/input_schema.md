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

## Validation / errors (v0)
This tool fails fast rather than guessing:

- Duplicate `id` values in a single file → error
- Header mismatch → error
- Amount with more than 2 decimal places → error

## Notes
- UTF-8 CSV.
- LF line endings recommended (repo enforces LF via `.gitattributes`).
