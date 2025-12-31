package scanner

import (
	"GoScanForSecrets/internal/models"
	"context"
	"encoding/json"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

type Scanner struct {
	registry *ExtractorRegistry
	patterns []models.CompiledPattern
	encoder  *json.Encoder
	logger   *slog.Logger
}

type Finding struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Pattern  string `json:"pattern"`
	Severity string `json:"severity"`
	Match    string `json:"match"`
}

func NewScanner(patterns []models.CompiledPattern, encoder *json.Encoder, log *slog.Logger) *Scanner {
	return &Scanner{
		registry: NewExtractorRegistry(),
		patterns: patterns,
		encoder:  encoder,
		logger:   log,
	}
}

func (s *Scanner) ScanPath(ctx context.Context, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			s.logger.Error("error accessing path", "path", path, "error", err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		// check if an extractor is available for this file
		extractor := s.registry.GetExtractor(path)
		if extractor == nil {
			return nil
		}

		return s.scanFile(ctx, path, extractor)
	})
}

func (s *Scanner) scanFile(ctx context.Context, path string, extractor Extractor) error {
	s.logger.Debug("scanning file", "path", path)

	f, err := os.Open(path)
	if err != nil {
		s.logger.Warn("failed to open file", "path", path, "error", err)
		return nil
	}
	defer f.Close()

	lines, err := extractor.Extract(ctx, f)
	if err != nil {
		s.logger.Warn("failed to extract lines", "path", path, "error", err)
		return nil
	}

	for lineNum, line := range lines {
		for _, pattern := range s.patterns {
			if pattern.Regex.MatchString(line) {
				finding := Finding{
					File:     path,
					Line:     lineNum + 1,
					Pattern:  pattern.Name,
					Severity: pattern.Severity,
					Match:    line,
				}

				if err := s.encoder.Encode(finding); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
