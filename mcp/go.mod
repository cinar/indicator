module github.com/cinar/indicator/mcp

go 1.23

toolchain go1.23.9

require (
	github.com/cinar/indicator/v2 v2.1.12
	github.com/mark3labs/mcp-go v0.31.0
)

// Build against the local tree so mcp/ always sees the indicator library's
// current state instead of trailing the last tagged release.
replace github.com/cinar/indicator/v2 => ../

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
)
