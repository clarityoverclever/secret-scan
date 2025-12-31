package scanner

import (
	"bufio"
	"context"
	"io"
	"path/filepath"
	"strings"
)

type TextExtractor struct{}

func (e *TextExtractor) Supports(filename string) bool {
	extension := strings.ToLower(filepath.Ext(filename))
	textExtensions := []string{".txt", ".log", ".yaml", ".yml", ".json", ".md", ".conf", ".cfg", ".csv"}

	for _, supportedExtension := range textExtensions {
		if extension == supportedExtension {
			return true
		}
	}
	return false
}

func (e *TextExtractor) Extract(ctx context.Context, r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
