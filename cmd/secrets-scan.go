package main

import (
	"GoScanForSecrets/internal/logger"
	"GoScanForSecrets/internal/plugins"
	"flag"
	"os"
)

const pluginPath = "patterns/"

func main() {
	// parse flags
	silent := flag.Bool("silent", false, "suppress output")
	verbose := flag.Bool("verbose", false, "enable verbose output")
	flag.Parse()

	// init logger
	log := logger.SetupLogger(*silent, *verbose)

	// start plugin VM
	loader := plugins.NewPatternLoader(log)
	defer loader.Close()

	log.Info("loading patterns")

	importedPatterns, err := loader.LoadPatterns(pluginPath)
	if err != nil {
		log.Error("failed to load patterns", "error", err)
		os.Exit(1)
	}

	log.Info("patterns loaded")

	log.Info("compiling patterns")

	compiledPatterns, err := loader.CompilePatterns(importedPatterns)
	if err != nil {
		log.Error("failed to compile patterns", "error", err)
		os.Exit(1)
	}

	log.Info("pattern compilation complete")

	for _, pattern := range compiledPatterns {
		log.Debug("loaded pattern", "name", pattern.Name, "severity", pattern.Severity)
	}
}
