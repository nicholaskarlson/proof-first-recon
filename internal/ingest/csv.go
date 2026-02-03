package ingest

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type Record struct {
	ID          string
	Date        string
	AmountCents int64
	Description string
}

func ReadCSV(path string) ([]Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	// strip UTF-8 BOM if present
	if len(header) > 0 {
		header[0] = strings.TrimPrefix(header[0], "\ufeff")
	}

	want := []string{"id", "date", "amount", "description"}
	if len(header) != len(want) {
		return nil, fmt.Errorf("header must have exactly %d columns: %v (got %v)", len(want), want, header)
	}
	for i := range want {
		if header[i] != want[i] {
			return nil, fmt.Errorf("header mismatch at column %d: want %q got %q", i+1, want[i], header[i])
		}
	}

	seen := make(map[string]struct{})
	out := make([]Record, 0, 64)

	rowNum := 1 // header is row 1
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		rowNum++
		if err != nil {
			return nil, fmt.Errorf("read row %d: %w", rowNum, err)
		}
		if len(row) != 4 {
			return nil, fmt.Errorf("row %d: expected 4 columns, got %d", rowNum, len(row))
		}

		id := strings.TrimSpace(row[0])
		date := strings.TrimSpace(row[1])
		amt := strings.TrimSpace(row[2])
		desc := strings.TrimSpace(row[3])

		if id == "" {
			return nil, fmt.Errorf("row %d: id is required", rowNum)
		}
		if _, ok := seen[id]; ok {
			return nil, fmt.Errorf("row %d: duplicate id %q", rowNum, id)
		}
		seen[id] = struct{}{}

		cents, err := ParseAmountCents(amt)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid amount %q: %w", rowNum, amt, err)
		}

		out = append(out, Record{
			ID:          id,
			Date:        date,
			AmountCents: cents,
			Description: desc,
		})
	}

	return out, nil
}
