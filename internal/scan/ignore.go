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

package scan

import (
	"os"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

type IgnoreManager struct {
	matcher *gitignore.GitIgnore
}

var defaultIgnorePatterns = []string{
	".git/",
	".svn/",
	".hg/",
	"node_modules/",
	"vendor/",
	"*.exe",
	"*.dll",
	"*.so",
	"*.dylib",
	"*.zip",
	"*.tar",
	"*.tar.gz",
	"*.rar",
	"*.7z",
	"*.bin",
	"*.class",
	"*.jar",
	"*.war",
	"*.ear",
	"*.pyc",
	"*.pyo",
	"*.o",
	"*.a",
	"*.lib",
	"*.obj",
	"*.iso",
	"*.img",
	"*.dmg",
}

func NewIgnoreManager(ignorePath string) (*IgnoreManager, error) {
	patterns := make([]string, len(defaultIgnorePatterns))
	copy(patterns, defaultIgnorePatterns)

	// load custom ignore patterns if specified
	if ignorePath != "" {
		customPatterns, err := loadIgnoreFile(ignorePath)
		if err != nil {
			return nil, err
		}

		patterns = append(patterns, customPatterns...)
	}

	matcher := gitignore.CompileIgnoreLines(patterns...)
	return &IgnoreManager{matcher: matcher}, nil
}

// loadIgnoreFile reads ignore patterns from a file
func loadIgnoreFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	patterns := make([]string, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		patterns = append(patterns, line)
	}

	return patterns, nil
}

// ShouldIgnore returns true if the given path should be ignored
func (im *IgnoreManager) ShouldIgnore(path string) bool {
	if im == nil || im.matcher == nil {
		return false
	}

	return im.matcher.MatchesPath(path)
}
