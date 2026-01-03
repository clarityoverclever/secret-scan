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
