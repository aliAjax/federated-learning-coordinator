#!/usr/bin/env sh
set -eu
lines=$(find . -name '*.go' ! -name '*_test.go' -print0 | xargs -0 wc -l | tail -1 | awk '{print $1}')
printf 'non-test Go source lines: %s\n' "$lines"
test "$lines" -ge 2600
