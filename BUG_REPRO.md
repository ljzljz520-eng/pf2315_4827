# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	example.com/scienceweekly/cmd/weekly	[no test files]
?   	example.com/scienceweekly/internal/api	[no test files]
?   	example.com/scienceweekly/internal/model	[no test files]
ok  	example.com/scienceweekly/internal/archive	0.015s
--- FAIL: Test2315BusinessRegression (0.00s)
    regression_test.go:24: second=[第一条状态 第一条状态]
FAIL
FAIL	example.com/scienceweekly/internal/flow001	0.023s
ok  	example.com/scienceweekly/internal/importer	0.010s
ok  	example.com/scienceweekly/internal/registry	0.009s
ok  	example.com/scienceweekly/internal/review	0.010s
ok  	example.com/scienceweekly/internal/store	0.012s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/weekly): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/weekly): exit `0`
