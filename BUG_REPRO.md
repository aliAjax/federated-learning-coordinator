# Bug reproduction

The gRPC missing-round path returns an internal error instead of not-found. Run `go test ./api/grpc -run '^TestGRPCMissingRoundKeepsNotFound$' -count=1` from the project root.

Root cause: several round and gRPC wrappers use `%v`, severing the domain sentinel error chain before status mapping.
