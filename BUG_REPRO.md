# Bug reproduction

Batch publishing a bad message can hang after an error. Run `go test -race ./internal/protocol/application -run '^TestBatchPublishErrorDoesNotHang$' -count=1` from the project root.

Root cause: `Codec.Batch` starts workers and waits in the wrong order, uses an unsafe error channel, and loses the original error while `Bus.Publish` mishandles cancellation.
