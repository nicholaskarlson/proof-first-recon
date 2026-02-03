package core

import (
	"fmt"
	"sort"

	"github.com/nicholaskarlson/proof-first-recon/internal/ingest"
	"github.com/nicholaskarlson/proof-first-recon/internal/model"
)

// ReconcileFromPaths reads CSVs and reconciles by v0 rule: exact match on id.
// It returns the parsed left/right records and the reconciliation result.
// File writing is handled by internal/report.
func ReconcileFromPaths(leftPath, rightPath string) ([]ingest.Record, []ingest.Record, model.Result, error) {
	left, err := ingest.ReadCSV(leftPath)
	if err != nil {
		return nil, nil, model.Result{}, fmt.Errorf("read left: %w", err)
	}
	right, err := ingest.ReadCSV(rightPath)
	if err != nil {
		return nil, nil, model.Result{}, fmt.Errorf("read right: %w", err)
	}

	res := ReconcileExactID(left, right)
	return left, right, res, nil
}

// ReconcileExactID performs v0 reconciliation: exact match on id.
func ReconcileExactID(left, right []ingest.Record) model.Result {
	rightMap := make(map[string]ingest.Record, len(right))
	for _, r := range right {
		rightMap[r.ID] = r
	}

	var res model.Result
	res.Matched = make([]model.Matched, 0, min(len(left), len(right)))
	res.UnmatchedLeft = make([]ingest.Record, 0, 16)

	for _, l := range left {
		if r, ok := rightMap[l.ID]; ok {
			res.Matched = append(res.Matched, model.Matched{ID: l.ID, Left: l, Right: r})
			delete(rightMap, l.ID)
		} else {
			res.UnmatchedLeft = append(res.UnmatchedLeft, l)
		}
	}

	res.UnmatchedRight = make([]ingest.Record, 0, len(rightMap))
	for _, r := range rightMap {
		res.UnmatchedRight = append(res.UnmatchedRight, r)
	}

	// Deterministic ordering
	sort.Slice(res.Matched, func(i, j int) bool { return res.Matched[i].ID < res.Matched[j].ID })
	sort.Slice(res.UnmatchedLeft, func(i, j int) bool { return res.UnmatchedLeft[i].ID < res.UnmatchedLeft[j].ID })
	sort.Slice(res.UnmatchedRight, func(i, j int) bool { return res.UnmatchedRight[i].ID < res.UnmatchedRight[j].ID })

	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
