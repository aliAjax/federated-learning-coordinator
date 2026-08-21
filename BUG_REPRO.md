# Bug reproduction

A request with the default privacy budget panics with `assignment to entry in nil map`. Run `go test -race ./internal/privacy/application -run '^TestDefaultBudgetDoesNotPanic$' -count=1` from the project root.

Root cause: the zero-value privacy memory store has a nil map, while the default consume path treats an empty provider result as writable budget state.
