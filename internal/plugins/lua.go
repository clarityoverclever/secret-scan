package plugins

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type PatternDefinition struct {
	Name     string
	Regex    string
	Severity string
}

func NewLuaVM() *lua.LState {
	return lua.NewState()
}

func readPatternFile(L *lua.LState, path string) {
	if err := L.DoFile(path + "patterns.lua"); err != nil {
		panic(err)
	}
}

func LoadPatterns(L *lua.LState, path string) ([]PatternDefinition, error) {
	readPatternFile(L, path)

	importTable := L.GetGlobal("patterns")
	if importTable == lua.LNil {
		return nil, fmt.Errorf("patterns table not found")
	}

	luaTable := importTable.(*lua.LTable)
	patterns := make([]PatternDefinition, 0, luaTable.Len())

	luaTable.ForEach(func(_, value lua.LValue) {
		entry := value.(*lua.LTable)

		next := PatternDefinition{
			Name:     entry.RawGetString("name").String(),
			Regex:    entry.RawGetString("regex").String(),
			Severity: entry.RawGetString("severity").String(),
		}

		patterns = append(patterns, next)
	})

	return patterns, nil
}
