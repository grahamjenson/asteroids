loadWebASM = () => {
	const go = new Go();
	getWASM().then( (b64) => {
  	// Decode and convert to ArrayBuffer
    buf = Uint8Array.from(atob(b64), c => c.charCodeAt(0)).buffer
    // WASM init
 		return WebAssembly.instantiate(buf, go.importObject)
	}).then((result) => {
    go.run(result.instance);
  }).catch((err) => {
    console.error("loading wasm failed: " + err);
  });
}

loadWebASM()