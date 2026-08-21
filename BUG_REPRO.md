# Bug reproduction

Loading a missing model loses the not-found sentinel and is reported as an invalid request. Run `go test ./internal/storage/postgres -run '^TestMissingModelPreservesNotFound$' -count=1` from the project root.

Root cause: repository error wrappers format sentinel errors with `%v`, breaking `errors.Is` through the storage boundary.
