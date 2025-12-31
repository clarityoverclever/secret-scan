package plugins

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"

	lua "github.com/yuin/gopher-lua"
)

type PatternLoader struct {
	logger *slog.Logger
	vm     *lua.LState
}

type PatternDefinition struct {
	Name     string
	Regex    string
	Severity string
}

type CompiledPattern struct {
	Name     string
	Severity string
	Regex    *regexp.Regexp
}

func NewPatternLoader(log *slog.Logger) *PatternLoader {
	return &PatternLoader{
		logger: log,
		vm:     newLuaVM(),
	}
}

func newLuaVM() *lua.LState {
	return lua.NewState()
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

func (pl *PatternLoader) LoadPatterns(path string) ([]PatternDefinition, error) {
	if err := pl.readPatternFile(path); err != nil {
		return nil, err
	}

	importTable := pl.vm.GetGlobal("patterns")
	if importTable == lua.LNil {
		return nil, fmt.Errorf("patterns table not found")
	}

	luaTable, ok := importTable.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("patterns table is not a table")
	}

	patterns := make([]PatternDefinition, 0, luaTable.Len())
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

		next := PatternDefinition{
			Name:     name.String(),
			Regex:    regex.String(),
			Severity: severity.String(),
		}

		patterns = append(patterns, next)
	})

	if skippedCount > 0 {
		pl.logger.Warn("pattern loading complete with warnings", "skipped: ", skippedCount)
	}

	return patterns, nil
}

func (pl *PatternLoader) CompilePatterns(patterns []PatternDefinition) ([]CompiledPattern, error) {
	compiled := make([]CompiledPattern, 0, len(patterns))

	for _, pattern := range patterns {
		regex, err := regexp.Compile(pattern.Regex)
		if err != nil {
			pl.logger.Warn("failed to compile regex", "pattern", pattern.Name, "error", err)
			continue
		}

		compiled = append(compiled, CompiledPattern{
			Name:     pattern.Name,
			Severity: pattern.Severity,
			Regex:    regex,
		})
	}

	pl.logger.Debug("compiled patterns", "count", len(compiled))

	return compiled, nil
}
