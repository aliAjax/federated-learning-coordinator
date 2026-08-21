# Bug reproduction

An expired model download continues work and can affect a later request. Run `go test ./internal/httpserver -run '^TestRequestDeadlineStopsModelRead$' -count=1` from the project root.

Root cause: request deadline and cancellation are not propagated consistently through the HTTP middleware, configuration, and round service publish path.
