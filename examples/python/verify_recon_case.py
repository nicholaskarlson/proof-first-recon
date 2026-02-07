#!/usr/bin/env python3
# SPDX-License-Identifier: MIT

from __future__ import annotations

import csv
import json
from decimal import Decimal
from pathlib import Path


def read_text_lf(p: Path) -> str:
    s = p.read_text(encoding="utf-8")
    if "\r\n" in s:
        raise ValueError(f"CRLF detected: {p}")
    return s


def read_csv_dicts(p: Path) -> list[dict[str, str]]:
    text = read_text_lf(p)
    return list(csv.DictReader(text.splitlines()))


def money_to_cents(s: str) -> int:
    # amounts are 2-decimal strings like "100.00" or "-25.50"
    d = Decimal(s)
    return int((d * 100).to_integral_value())


def sum_cents(rows: list[dict[str, str]], field: str) -> int:
    return sum(money_to_cents(r[field]) for r in rows)


def main() -> None:
    repo = Path(".")
    case = "case01"

    out_dir = repo / "out"
    exp_dir = repo / "fixtures" / "expected" / case

    matched = read_csv_dicts(out_dir / "matched.csv")
    ul = read_csv_dicts(out_dir / "unmatched_left.csv")
    ur = read_csv_dicts(out_dir / "unmatched_right.csv")
    summary = json.loads(read_text_lf(out_dir / "recon_summary.json"))

    # Summary shape + counts
    assert summary["schema_version"] == "v0"
    assert summary["rule"] == "exact_id"
    assert summary["matched_count"] == len(matched)
    assert summary["unmatched_left_count"] == len(ul)
    assert summary["unmatched_right_count"] == len(ur)

    # Counts tie out (matched contributes one left + one right)
    left_count = len(matched) + len(ul)
    right_count = len(matched) + len(ur)
    assert summary["left_count"] == left_count
    assert summary["right_count"] == right_count

    # Totals tie out (computed from output artifacts)
    left_total = sum_cents(matched, "left_amount") + sum_cents(ul, "amount")
    right_total = sum_cents(matched, "right_amount") + sum_cents(ur, "amount")
    assert summary["left_total_cents"] == left_total
    assert summary["right_total_cents"] == right_total

    # Matched rows must agree on amount
    for r in matched:
        assert r["id"].strip() != ""
        assert r["left_amount"] == r["right_amount"]

    # Optional: byte-compare against goldens
    for name in [
        "matched.csv",
        "unmatched_left.csv",
        "unmatched_right.csv",
        "recon_summary.json",
    ]:
        got = read_text_lf(out_dir / name)
        exp = read_text_lf(exp_dir / name)
        assert got == exp, f"golden mismatch: {name}"

    print("OK: recon outputs pass invariants + match goldens.")


if __name__ == "__main__":
    main()
