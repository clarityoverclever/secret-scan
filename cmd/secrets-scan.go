package main

import (
	"GoScanForSecrets/internal/plugins"
	"fmt"
)

const pluginPath = "patterns/"

func main() {
	// start plugin VM
	lua := plugins.NewLuaVM()
	defer lua.Close()

	importedPatterns, err := plugins.LoadPatterns(lua, pluginPath)
	if err != nil {
	}

	for pattern := range importedPatterns {
		fmt.Println(pattern)
	}

}
