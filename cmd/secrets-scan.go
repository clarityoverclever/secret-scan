package main

import (
	"github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	L.DoString("print('hello, from Lua')")
}
