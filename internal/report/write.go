package report

import (
	"bytes"
	"encoding/csv"
	"encoding/json"

	"github.com/nicholaskarlson/proof-first-recon/internal/ingest"
	"github.com/nicholaskarlson/proof-first-recon/internal/model"
)

type OutputPaths struct {
	MatchedCSV        string
	UnmatchedLeftCSV  string
	UnmatchedRightCSV string
	SummaryJSON       string
}

type Summary struct {
	SchemaVersion string `json:"schema_version"`
	Rule          string `json:"rule"`

	LeftCount  int `json:"left_count"`
	RightCount int `json:"right_count"`

	MatchedCount        int `json:"matched_count"`
	UnmatchedLeftCount  int `json:"unmatched_left_count"`
	UnmatchedRightCount int `json:"unmatched_right_count"`

	LeftTotalCents           int64 `json:"left_total_cents"`
	RightTotalCents          int64 `json:"right_total_cents"`
	MatchedLeftTotalCents    int64 `json:"matched_left_total_cents"`
	MatchedRightTotalCents   int64 `json:"matched_right_total_cents"`
	UnmatchedLeftTotalCents  int64 `json:"unmatched_left_total_cents"`
	UnmatchedRightTotalCents int64 `json:"unmatched_right_total_cents"`

	LeftTotal           string `json:"left_total"`
	RightTotal          string `json:"right_total"`
	MatchedLeftTotal    string `json:"matched_left_total"`
	MatchedRightTotal   string `json:"matched_right_total"`
	UnmatchedLeftTotal  string `json:"unmatched_left_total"`
	UnmatchedRightTotal string `json:"unmatched_right_total"`
}

func WriteAll(paths OutputPaths, left, right []ingest.Record, res model.Result) error {
	if err := writeUnmatched(paths.UnmatchedLeftCSV, res.UnmatchedLeft); err != nil {
		return err
	}
	if err := writeUnmatched(paths.UnmatchedRightCSV, res.UnmatchedRight); err != nil {
		return err
	}
	if err := writeMatched(paths.MatchedCSV, res.Matched); err != nil {
		return err
	}
	if err := writeSummary(paths.SummaryJSON, left, right, res); err != nil {
		return err
	}
	return nil
}

func writeUnmatched(path string, rows []ingest.Record) error {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write([]string{"id", "date", "amount", "description"}); err != nil {
		return err
	}
	for _, r := range rows {
		if err := w.Write([]string{r.ID, r.Date, ingest.FormatCents(r.AmountCents), r.Description}); err != nil {
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return writeFileAtomic(path, buf.Bytes(), 0o644)
}

func writeMatched(path string, rows []model.Matched) error {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write([]string{
		"id",
		"left_date", "left_amount", "left_description",
		"right_date", "right_amount", "right_description",
	}); err != nil {
		return err
	}

	for _, m := range rows {
		if err := w.Write([]string{
			m.ID,
			m.Left.Date, ingest.FormatCents(m.Left.AmountCents), m.Left.Description,
			m.Right.Date, ingest.FormatCents(m.Right.AmountCents), m.Right.Description,
		}); err != nil {
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return writeFileAtomic(path, buf.Bytes(), 0o644)
}

func writeSummary(path string, left, right []ingest.Record, res model.Result) error {
	sumLeft := total(left)
	sumRight := total(right)

	matchedLeft := int64(0)
	matchedRight := int64(0)
	for _, m := range res.Matched {
		matchedLeft += m.Left.AmountCents
		matchedRight += m.Right.AmountCents
	}

	unLeft := total(res.UnmatchedLeft)
	unRight := total(res.UnmatchedRight)

	s := Summary{
		SchemaVersion: "v0",
		Rule:          "exact_id",

		LeftCount:  len(left),
		RightCount: len(right),

		MatchedCount:        len(res.Matched),
		UnmatchedLeftCount:  len(res.UnmatchedLeft),
		UnmatchedRightCount: len(res.UnmatchedRight),

		LeftTotalCents:           sumLeft,
		RightTotalCents:          sumRight,
		MatchedLeftTotalCents:    matchedLeft,
		MatchedRightTotalCents:   matchedRight,
		UnmatchedLeftTotalCents:  unLeft,
		UnmatchedRightTotalCents: unRight,

		LeftTotal:           ingest.FormatCents(sumLeft),
		RightTotal:          ingest.FormatCents(sumRight),
		MatchedLeftTotal:    ingest.FormatCents(matchedLeft),
		MatchedRightTotal:   ingest.FormatCents(matchedRight),
		UnmatchedLeftTotal:  ingest.FormatCents(unLeft),
		UnmatchedRightTotal: ingest.FormatCents(unRight),
	}

	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return writeFileAtomic(path, b, 0o644)
}

func total(rows []ingest.Record) int64 {
	var t int64
	for _, r := range rows {
		t += r.AmountCents
	}
	return t
}
