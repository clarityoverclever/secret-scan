package scan

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

	// bufio.Scanner has a maximum line length of 64kb per token
	// this will limit the line length the scanner will read to
	// prevent errors on lines which exceed 64kb.

	const maxLineLength = 10 * 1024 // 10kb
	buf := make([]byte, maxLineLength)
	scanner.Buffer(buf, maxLineLength)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) < maxLineLength {
			lines = append(lines, line)
		}

		// ignore "token too long" errors
		if err := scanner.Err(); err != nil {
			if strings.Contains(err.Error(), "token too long") {
				return lines, nil
			}

			// return other errors
			return lines, scanner.Err()
		}
	}

	return lines, nil
}
