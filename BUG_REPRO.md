# Bug reproduction

Concurrent update snapshots race with map writes, causing missing updates and a `-race` failure. Run `go test -race ./internal/storage -run '^TestConcurrentUpdateSnapshot$' -count=1` from the project root.

Root cause: `Memory.ListUpdates` and aggregation `Memory.List` read shared maps/slices without holding the read lock or cloning returned values.
