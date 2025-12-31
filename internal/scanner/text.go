package scanner

import (
	"context"
	"io"
)

type TextExtractor struct{}

func (e *TextExtractor) Extract(ctx context.Context, r io.Reader) ([]string, error) {
	return nil, nil
}

func (e *TextExtractor) Supports(filename string) bool {
	return true
}
