package main

import (
	"GoScanForSecrets/internal/logger"
	"GoScanForSecrets/internal/plugins"
	//"GoScanForSecrets/internal/scanner"
	"encoding/json"
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

	// init json encoder
	//encoder := json.NewEncoder(os.Stdout)

	// start plugin VM
	log.Debug("initializing lua plugin VM")
	loader := plugins.NewPatternLoader(log)
	defer loader.Close()

	log.Debug("importing patterns")
	importedPatterns, err := loader.LoadPatterns(pluginPath)
	if err != nil {
		log.Error("failed to load patterns", "error", err)
		os.Exit(1)
	}

	compiledPatterns, err := loader.CompilePatterns(importedPatterns)
	if err != nil {
		log.Error("failed to compile patterns", "error", err)
		os.Exit(1)
	}

	// init scanner
	//scanner := scanner.NewScanner(encoder)

	for _, pattern := range compiledPatterns {
		log.Debug("ready pattern", "name", pattern.Name, "severity", pattern.Severity)
	}
}
