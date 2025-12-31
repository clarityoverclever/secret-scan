package main

import (
	"GoScanForSecrets/config"
	"GoScanForSecrets/internal/logger"
	"GoScanForSecrets/internal/plugins"
	"GoScanForSecrets/internal/scan"
	"context"
	"encoding/json"
	"os"
)

const pluginPath = "patterns/"

func main() {
	// parse flags
	cfg := config.ParseFlags()

	// init logger
	log := logger.SetupLogger(cfg.Silent, cfg.Verbose)

	// init json encoder
	encoder := json.NewEncoder(os.Stdout)

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
	scanner := scan.NewScanner(compiledPatterns, encoder, log, cfg.Threads)

	// add background context for the scanner
	ctx := context.Background()

	log.Debug("scanning path", "path", cfg.ScanPath)

	for _, pattern := range compiledPatterns {
		log.Debug("ready pattern", "name", pattern.Name, "severity", pattern.Severity)
	}

	log.Info("starting scan")

	if err := scanner.ScanPath(ctx, cfg.ScanPath); err != nil {
		log.Error("scan failed", "error", err)
		os.Exit(1)
	}

	log.Info("scan finished")
}
