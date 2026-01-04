package scan

import (
	"context"
	"encoding/json"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"secret-scan/internal/models"
	"secret-scan/internal/validators"
	"sync"
)

type Scanner struct {
	registry   *ExtractorRegistry
	patterns   []models.CompiledPattern
	validators *validators.Registry
	encoder    *json.Encoder
	logger     *slog.Logger
	numWorkers int
	mutex      sync.Mutex
}

type scanJob struct {
	path      string
	extractor Extractor
}

type Finding struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Pattern  string `json:"pattern"`
	Severity string `json:"severity"`
	Match    string `json:"match"`
}

func NewScanner(patterns []models.CompiledPattern, encoder *json.Encoder, log *slog.Logger, numWorkers int) *Scanner {
	return &Scanner{
		registry:   NewExtractorRegistry(),
		patterns:   patterns,
		validators: validators.NewRegistry(),
		encoder:    encoder,
		logger:     log,
		numWorkers: numWorkers,
	}
}

func (s *Scanner) ScanPath(ctx context.Context, root string) error {
	jobs := make(chan scanJob, 100)
	var wg sync.WaitGroup

	s.logger.Debug("starting worker pool", "workers", s.numWorkers)
	for i := 0; i < s.numWorkers; i++ {
		wg.Add(1)
		go s.worker(ctx, i, jobs, &wg)
	}

	walkError := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
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

		select {
		case jobs <- scanJob{path: path, extractor: extractor}:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	})

	close(jobs)
	wg.Wait()

	return walkError
}

func (s *Scanner) worker(ctx context.Context, id int, jobs <-chan scanJob, wg *sync.WaitGroup) {
	defer wg.Done()

	s.logger.Debug("starting worker", "id", id)
	for job := range jobs {
		select {
		case <-ctx.Done():
			s.logger.Debug("worker stopped", "id", id)
			return
		default:
			if err := s.scanFile(ctx, job.path, job.extractor); err != nil {
				s.logger.Error("failed to scan file", "path", job.path, "error", err)
			}
		}
	}

	s.logger.Debug("worker finished", "id", id)
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

	findings := make([]Finding, 0)
	for lineNum, line := range lines {
		for _, pattern := range s.patterns {
			if pattern.Regex.MatchString(line) {
				matches := pattern.Regex.FindAllString(line, -1)
				// run the match against a validator if specified
				for _, match := range matches {
					if pattern.Validator != "" {
						validator := s.validators.Get(pattern.Validator)
						if validator != nil && !validator.Validate(match, line) {
							s.logger.Debug("match failed validation",
								"pattern", pattern.Name,
								"validator", pattern.Validator,
							)
							continue // skip this match
						}
					}
					findings = append(findings, Finding{
						File:     path,
						Line:     lineNum + 1,
						Pattern:  pattern.Name,
						Severity: pattern.Severity,
						Match:    match,
					})
				}
			}
		}
	}

	if len(findings) > 0 {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		for _, finding := range findings {
			if err := s.encoder.Encode(finding); err != nil {
				s.logger.Error("failed to encode finding", "error", err)
				return err
			}
		}
	}
	return nil
}
