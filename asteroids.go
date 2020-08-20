package main

import (
	"github.com/grahamjenson/asteroids/desktop"
)

func main() {
	config := &desktop.Config{
		WasmBin: WASM_BIN,
		Width:   1280,
		Height:  720 + 20, // add 20 for the top bar
		Title:   "Asteroids",
	}
	desktop.CreateDesktopApp(config)
}
