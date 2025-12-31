package config

import "flag"

type Config struct {
	Silent         bool
	Verbose        bool
	OutputFilename string
	ScanPath       string
}

func ParseFlags() Config {
	cfg := Config{}

	flag.BoolVar(&cfg.Silent, "silent", false, "suppress output")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "enable verbose output")
	flag.StringVar(&cfg.OutputFilename, "output", "", "output file")

	flag.Parse()

	if flag.NArg() > 0 {
		cfg.ScanPath = flag.Arg(0)
	} else {
		cfg.ScanPath = "."
	}

	return cfg
}
