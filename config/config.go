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

package config

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Silent            bool
	Verbose           bool
	NoDefaultPatterns bool
	OutputFilename    string
	ScanPath          string
	PatternsPath      string
	Threads           int
}

func ParseFlags() Config {
	cfg := Config{}

	flag.BoolVar(&cfg.Silent, "silent", false, "suppress output")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "enable verbose output")
	flag.BoolVar(&cfg.NoDefaultPatterns, "no-default-patterns", false, "disable loading of default patterns")
	flag.StringVar(&cfg.OutputFilename, "out", "", "output file")
	flag.StringVar(&cfg.PatternsPath, "patterns", "", "path to custome patterns file")
	flag.IntVar(&cfg.Threads, "threads", runtime.NumCPU()-1, "number of threads")

	flag.Parse()

	if cfg.Threads < 1 {
		cfg.Threads = 1
	}

	if flag.NArg() > 0 {
		cfg.ScanPath = flag.Arg(0)
	} else {
		cfg.ScanPath = "."
	}

	// expand home path if supplied as part of the PatternsPath
	if cfg.PatternsPath != "" {
		if cfg.PatternsPath[:2] == "~/" {
			home, _ := os.UserHomeDir()
			cfg.PatternsPath = filepath.Join(home, cfg.PatternsPath[2:])
		}
	}

	return cfg
}
