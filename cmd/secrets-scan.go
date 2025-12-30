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
	lua := plugins.NewLuaVM()
	defer lua.Close()

	importedPatterns, err := plugins.LoadPatterns(lua, pluginPath, log)
	if err != nil {
		log.Error("failed to load patterns", "error", err)
		os.Exit(1)
	}

	log.Info("loaded patterns", "count", len(importedPatterns))

	compiledPatterns, err := plugins.CompilePatterns(importedPatterns, log)
	if err != nil {
		log.Error("failed to compile patterns", "error", err)
		os.Exit(1)
	}

	log.Info("compiled patterns", "count", len(compiledPatterns))

	for _, pattern := range compiledPatterns {
		log.Debug("loaded pattern", "name", pattern.Name, "severity", pattern.Severity)
	}
}
