package main

import (
	"context"
	"encoding/json"
	"os"
	"secret-scan/config"
	"secret-scan/internal/logger"
	"secret-scan/internal/plugins"
	"secret-scan/internal/scan"
)

func main() {
	var err error

	// parse flags
	cfg := config.ParseFlags()

	// init logger
	log := logger.SetupLogger(cfg.Silent, cfg.Verbose)

	// create output file if specified
	var outputFile *os.File
	if cfg.OutputFilename != "" {
		log.Debug("creating output file", "path", cfg.OutputFilename)
		outputFile, err = os.Create(cfg.OutputFilename)
		if err != nil {
			log.Error("failed to create output file", "error", err)
			os.Exit(1)
		}
		defer outputFile.Close()
	} else {
		outputFile = os.Stdout
	}

	// init json encoder
	encoder := json.NewEncoder(outputFile)

	// start plugin VM
	log.Debug("initializing lua plugin VM")
	loader := plugins.NewPatternLoader(log)
	defer loader.Close()

	// import patterns
	log.Debug("importing patterns")
	importedPatterns, err := loader.LoadPatterns(cfg.PatternsPath)
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

	log.Info("starting scan")

	if err := scanner.ScanPath(ctx, cfg.ScanPath); err != nil {
		log.Error("scan failed", "error", err)
		os.Exit(1)
	}

	log.Info("scan finished")
}
