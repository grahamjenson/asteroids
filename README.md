# Asteroids

This is the game asteroids implemented in Golang using Lorca/Webview and WASM.

To run it use:

```
bazel run :asteroids
```

or download the binary from the releases.


Update go dependencies with 

```
bazel run //:gazelle -- update-repos -from_file=go.mod
```
