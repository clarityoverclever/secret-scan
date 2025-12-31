package scanner

import (
	"context"
	"io"
)

// Extractor converts file content into scannable unicode
type Extractor interface {
	Extract(ctx context.Context, r io.Reader) ([]string, error)

	// Supports return true if a file type is supported
	Supports(filename string) bool
}

type ExtractorRegistry struct {
	extractors []Extractor
}

func NewExtractorRegistry() *ExtractorRegistry {
	return &ExtractorRegistry{
		extractors: []Extractor{
			&TextExtractor{},
			&ExcelExtractor{},
		},
	}
}

func (r *ExtractorRegistry) GetExtractor(filename string) Extractor {
	for _, extractor := range r.extractors {
		if extractor.Supports(filename) {
			return extractor
		}
	}
	return nil // no extractor for this filetype
}
