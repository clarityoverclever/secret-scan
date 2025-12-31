package scanner

import (
	"context"
	"io"
)

type ExcelExtractor struct{}

func (e *ExcelExtractor) Extract(ctx context.Context, r io.Reader) ([]string, error) {
	return nil, nil
}

func (e *ExcelExtractor) Supports(filename string) bool {
	return true
}
