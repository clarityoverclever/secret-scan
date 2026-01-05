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
)

// Extractor converts file content into scannable unicode
type Extractor interface {
	Extract(ctx context.Context, r io.Reader) ([]string, error)

	// Supports return true if a file type is supported
	Supports(filename string) bool
}

type Registry struct {
	extractors []Extractor
}

func NewRegistry() *Registry {
	return &Registry{
		extractors: []Extractor{
			&TextExtractor{},
			&ExcelExtractor{},
		},
	}
}

func (r *Registry) Get(filename string) Extractor {
	for _, extractor := range r.extractors {
		if extractor.Supports(filename) {
			return extractor
		}
	}
	return nil // no extractor for this filetype
}
