package scan

import (
	"context"
	"io"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ExcelExtractor struct{}

func (e *ExcelExtractor) Supports(filename string) bool {
	extension := strings.ToLower(filepath.Ext(filename))
	return extension == ".xlsx"
}

func (e *ExcelExtractor) Extract(ctx context.Context, r io.Reader) ([]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string

	for _, sheet := range f.GetSheetList() {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			line := strings.Join(row, "\t") // tab separated
			if strings.TrimSpace(line) != " " {
				lines = append(lines, line)
			}
		}
	}

	return lines, nil
}
