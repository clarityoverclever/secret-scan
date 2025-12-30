package plugins

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"

	lua "github.com/yuin/gopher-lua"
)

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

func NewLuaVM() *lua.LState {
	return lua.NewState()
}

func readPatternFile(L *lua.LState, path string) error {
	if err := L.DoFile(filepath.Join(path, "patterns.lua")); err != nil {
		return fmt.Errorf("failed to load pattern file: %w", err)
	}
	return nil
}

func LoadPatterns(L *lua.LState, path string, log *slog.Logger) ([]PatternDefinition, error) {
	if err := readPatternFile(L, path); err != nil {
		return nil, err
	}

	importTable := L.GetGlobal("patterns")
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
			log.Warn("invalid table entry")
			return // skip invalid table entries
		}

		name := entry.RawGetString("name")
		regex := entry.RawGetString("regex")
		severity := entry.RawGetString("severity")

		if name == lua.LNil || regex == lua.LNil || severity == lua.LNil {
			skippedCount++
			log.Warn("skipping pattern with missing fields",
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
		log.Info("pattern loading complete with warnings", "skipped: ", skippedCount)
	}

	return patterns, nil
}

func CompilePatterns(patterns []PatternDefinition, log *slog.Logger) ([]CompiledPattern, error) {
	compiled := make([]CompiledPattern, 0, len(patterns))

	for _, pattern := range patterns {
		regex, err := regexp.Compile(pattern.Regex)
		if err != nil {
			log.Warn("failed to compile regex", "pattern", pattern.Name, "error", err)
			continue
		}

		compiled = append(compiled, CompiledPattern{
			Name:     pattern.Name,
			Severity: pattern.Severity,
			Regex:    regex,
		})
	}

	log.Debug("compiled patterns", "count", len(compiled))

	return compiled, nil
}
