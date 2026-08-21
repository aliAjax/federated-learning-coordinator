# Bug reproduction

Clipping layers before aggregation mutates the previously stored update. Run `go test ./internal/tensor -run '^TestClipDoesNotPolluteOriginalLayers$' -count=1` from the project root.

Root cause: `Clip` and aggregation reuse slice backing arrays, while compatibility checks omit `DType`, allowing shared mutable layer data to be overwritten.
