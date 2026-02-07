# Optional Python checks (stdlib only)

Python is **optional** in this repo. The source-of-truth proof gates are in Go:

```bash
go test -count=1 ./...
```

This folder provides an *independent verification lane* (standard library only).

## Recon demo + Python verify

From repo root:

```bash
# Go demo (writes deterministic outputs)
go run ./cmd/recon demo --out ./out

# Python invariants + (optional) golden compare
python3 examples/python/verify_recon_case.py
```
