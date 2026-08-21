# Bug reproduction

A successfully aggregated round remains in `aggregating` instead of reaching `completed`. Run `go test ./internal/round/domain -run '^TestRoundRetryReachesCompleted$' -count=1` from the project root.

Root cause: the round state transition table omits a legal completion path and the serialized status list does not expose the completed state correctly.
