# Bug reproduction

A failed aggregation batch can log success and retain resources. Run `go test ./internal/aggregation/application -run '^TestBatchAggregationKeepsPrimaryError$' -count=1` from the project root.

Root cause: aggregation and privacy error wrapping uses `%v`, so the primary error chain is lost and cleanup/cancellation handling reports the wrong outcome.
