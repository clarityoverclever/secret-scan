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
		panic(err)
	}

	fmt.Println("patterns imported")

	compiledPatterns, err := plugins.CompilePatterns(importedPatterns)
	if err != nil {
		panic(err)
	}

	fmt.Println("patterns compiled")

	for _, pattern := range compiledPatterns {
		fmt.Println(pattern.Name)
	}

}
