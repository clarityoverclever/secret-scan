package plugins

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"secret-scan/internal/models"
	"secret-scan/patterns"

	lua "github.com/yuin/gopher-lua"
)

type PatternLoader struct {
	logger *slog.Logger
	vm     *lua.LState
}

func NewPatternLoader(log *slog.Logger) *PatternLoader {
	return &PatternLoader{
		logger: log,
		vm:     newLuaVM(),
	}
}

func newLuaVM() *lua.LState {
	vm := lua.NewState()

	// Security: Disable dangerous Lua functions
	vm.SetGlobal("dofile", lua.LNil)
	vm.SetGlobal("loadfile", lua.LNil)
	vm.SetGlobal("require", lua.LNil)
	vm.SetGlobal("load", lua.LNil)
	vm.SetGlobal("loadstring", lua.LNil)
	vm.SetGlobal("io", lua.LNil)
	vm.SetGlobal("os", lua.LNil)
	vm.SetGlobal("package", lua.LNil)
	vm.SetGlobal("debug", lua.LNil)

	return vm
}

func (pl *PatternLoader) Close() {
	pl.vm.Close()
}

func (pl *PatternLoader) readPatternFile(path string) error {
	if err := pl.vm.DoFile(filepath.Join(path, "patterns.lua")); err != nil {
		return fmt.Errorf("failed to load pattern file: %w", err)
	}
	return nil
}

func (pl *PatternLoader) LoadPatterns(customPath string, noDefaultPatterns bool) ([]models.PatternDefinition, error) {
	var luaFiles []string
	var luaContents []string

	if noDefaultPatterns != true {
		pl.logger.Debug("loading embedded patterns")
		files, contents, err := pl.loadFromEmbedded()
		if err != nil {
			return nil, fmt.Errorf("failed to load embedded patterns: %w", err)
		}
		luaFiles = append(luaFiles, files...)
		luaContents = append(luaContents, contents...)
	}

	if customPath != "" {
		pl.logger.Debug("loading custom patterns", "path", customPath)
		files, contents, err := pl.loadFromDirectory(customPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load custom patterns: %w", err)
		}
		luaFiles = append(luaFiles, files...)
		luaContents = append(luaContents, contents...)
	}

	return pl.extractPatterns(luaFiles, luaContents)
}

func (pl *PatternLoader) loadFromDirectory(path string) ([]string, []string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read pattern directory: %w", err)
	}

	var files []string
	var contents []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".lua" {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())
		content, err := os.ReadFile(fullPath)
		if err != nil {
			pl.logger.Warn("failed to read pattern file", "path", fullPath, "error", err)
			continue
		}

		files = append(files, entry.Name())
		contents = append(contents, string(content))
	}

	if len(files) == 0 {
		return nil, nil, fmt.Errorf("no .lua pattern files found in directory")
	}

	return files, contents, nil
}

func (pl *PatternLoader) loadFromEmbedded() ([]string, []string, error) {
	var files []string
	var contents []string

	err := fs.WalkDir(patterns.Embedded, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) != ".lua" {
			return nil
		}

		content, err := patterns.Embedded.ReadFile(path)
		if err != nil {
			pl.logger.Warn("failed to read embedded pattern file", "path", path, "error", err)
			return nil
		}

		files = append(files, d.Name())
		contents = append(contents, string(content))

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to read embedded patterns: %w", err)
	}

	if len(files) == 0 {
		return nil, nil, fmt.Errorf("no .lua pattern files found in embedded patterns")
	}

	return files, contents, nil
}

func (pl *PatternLoader) extractPatterns(files []string, contents []string) ([]models.PatternDefinition, error) {
	var allPatterns []models.PatternDefinition

	for index, content := range contents {
		filename := files[index]

		if err := pl.vm.DoString(content); err != nil {
			pl.logger.Warn("failed to load pattern file", "path", filename, "error", err)
			continue
		}

		patterns, err := pl.extractPatternsFromVM(filename)
		if err != nil {
			pl.logger.Warn("failed to extract patterns from VM", "path", filename, "error", err)
			continue
		}

		allPatterns = append(allPatterns, patterns...)
		pl.logger.Debug("loaded patterns from file", "path", filename, "count", len(patterns))
	}

	if len(allPatterns) == 0 {
		return nil, fmt.Errorf("no valid patterns found in any pattern files")
	}

	pl.logger.Info("pattern loading complete", "count", len(allPatterns))
	return allPatterns, nil
}

func (pl *PatternLoader) extractPatternsFromVM(path string) ([]models.PatternDefinition, error) {
	importTable := pl.vm.GetGlobal("patterns")

	if importTable == lua.LNil {
		return nil, fmt.Errorf("patterns table not found")
	}

	luaTable, ok := importTable.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("patterns table is not a table")
	}

	patterns := make([]models.PatternDefinition, 0, luaTable.Len())
	skippedCount := 0

	luaTable.ForEach(func(_, value lua.LValue) {
		entry, ok := value.(*lua.LTable)
		if !ok {
			skippedCount++
			pl.logger.Warn("invalid table entry")
			return // skip invalid table entries
		}

		name := entry.RawGetString("name")
		regex := entry.RawGetString("regex")
		severity := entry.RawGetString("severity")

		if name == lua.LNil || regex == lua.LNil || severity == lua.LNil {
			skippedCount++
			pl.logger.Warn("skipping pattern with missing fields",
				"has_name", name != lua.LNil,
				"has_regex", regex != lua.LNil,
				"has_severity", severity != lua.LNil)
			return // skip invalid table entries
		}

		next := models.PatternDefinition{
			Name:     name.String(),
			Regex:    regex.String(),
			Severity: severity.String(),
		}

		patterns = append(patterns, next)
	})

	if skippedCount > 0 {
		pl.logger.Warn("pattern loading complete with warnings", "skipped: ", skippedCount)
	}

	// clear VM patterns for next file
	pl.vm.SetGlobal("patterns", lua.LNil)

	return patterns, nil
}

func (pl *PatternLoader) CompilePatterns(patterns []models.PatternDefinition) ([]models.CompiledPattern, error) {
	compiled := make([]models.CompiledPattern, 0, len(patterns))

	for _, pattern := range patterns {
		regex, err := regexp.Compile(pattern.Regex)
		if err != nil {
			pl.logger.Warn("failed to compile regex", "pattern", pattern.Name, "error", err)
			continue
		}

		compiled = append(compiled, models.CompiledPattern{
			Name:     pattern.Name,
			Severity: pattern.Severity,
			Regex:    regex,
		})
	}

	if len(compiled) == 0 {
		return nil, fmt.Errorf("no patterns successfully compiled")
	}

	pl.logger.Debug("patterns compiled", "count", len(compiled))

	return compiled, nil
}
