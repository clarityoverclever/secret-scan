// Copyright 2026 Keith Marshall
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package extractors

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
