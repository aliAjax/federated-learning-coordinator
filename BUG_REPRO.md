# Bug reproduction

Revealing a missing mask can return success and leave a zero-value mask for later consumers. Run `go test -race ./internal/secureagg/application -run '^TestMissingMaskReturnsExplicitError$' -count=1` from the project root.

Root cause: the reveal service fabricates a missing mask instead of returning the sentinel error, and the memory store reports missing entries as present.
