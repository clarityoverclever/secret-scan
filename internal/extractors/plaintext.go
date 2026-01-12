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
	"strings"
)

type PlainTextAuditExtractor struct{}

func (e *PlainTextAuditExtractor) Supports(filename string) bool {
	// treat all files as plaintext as failover
	// test for binary files in the extractor
	return true
}

func (e *PlainTextAuditExtractor) Extract(ctx context.Context, r io.Reader) ([]string, error) {
	br := bufio.NewReader(r)

	sniff, err := br.Peek(512)

	if err != nil && err != io.EOF {
		return nil, err
	}

	// check for binary indicator (null byte)
	for _, b := range sniff {
		if b == 0 {
			return nil, nil
		}

	}

	// treat file as plaintext
	var lines []string
	scanner := bufio.NewScanner(br)

	const maxLineLength = 10 * 1024 // 10kb
	buf := make([]byte, maxLineLength)
	scanner.Buffer(buf, maxLineLength)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// check for scanner errors
	if err := scanner.Err(); err != nil {
		// ignore "token too long" errors
		if strings.Contains(err.Error(), "token too long") {
			return lines, nil
		}
		// return other errors
		return lines, scanner.Err()
	}

	return lines, nil
}
