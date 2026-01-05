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
