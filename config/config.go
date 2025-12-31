package config

import (
	"flag"
	"runtime"
)

type Config struct {
	Silent         bool
	Verbose        bool
	OutputFilename string
	ScanPath       string
	Threads        int
}

func ParseFlags() Config {
	cfg := Config{}

	flag.BoolVar(&cfg.Silent, "silent", false, "suppress output")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "enable verbose output")
	flag.StringVar(&cfg.OutputFilename, "out", "", "output file")
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

	return cfg
}
