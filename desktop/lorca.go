// +build windows linux darwin,!cgo

package desktop

import (
	"fmt"
	"log"
	"net/url"

	"github.com/zserge/lorca"
)

func CreateDesktopApp(config *Config) {
	// Create UI with basic HTML passed via data URI
	indexHTML := STATIC_STRINGS["desktop/index.html"]
	wasmEXEC := STATIC_STRINGS["desktop/wasm_exec.js"]
	initJS := STATIC_STRINGS["desktop/init.js"]

	url := fmt.Sprintf(
		"data:text/html,%s",
		url.PathEscape(indexHTML),
	)

	ui, err := lorca.New(url, "", config.Width, config.Height)

	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// Create a JS function that returns the WASM binary in base64
	ui.Bind("getWASM", func() []byte { return config.WasmBin })

	// Initialize the wasm_exec.js script
	ui.Eval(wasmEXEC)

	// Call the initial JS functions to load the WASM
	ui.Eval(initJS)

	// Wait until UI window is closed
	<-ui.Done()
}
