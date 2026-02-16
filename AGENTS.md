# AGENTS.md - Agentic Coding Guidelines

This document provides guidelines for agents working on the Indicator Go project.

## Project Overview

Indicator is a Golang library providing technical analysis indicators, strategies, and a backtesting framework. It uses Go 1.22+ with generics and channels for streaming data processing.

## Build/Lint/Test Commands

### Running the Full CI Pipeline
```bash
# Using task (requires go-task)
task

# Or manually:
task fmt
task lint
task test
task docs
```

### Individual Commands
```bash
# Format code
task fmt
# Or: go fix ./...

# Run linters (go vet, gosec, staticcheck, revive)
task lint
# Specific linters:
go vet ./...
go run github.com/securego/gosec/v2/cmd/gosec@v2.20.0 ./...
go run honnef.co/go/tools/cmd/staticcheck@v0.5.1 ./...
go run github.com/mgechev/revive@v1.3.4 -config=revive.toml ./...

# Run tests with coverage
task test
# Or: go test -cover -coverprofile=coverage.out ./...

# Build CLI tools
task build-tools
# Or manually:
go build -o indicator-backtest cmd/indicator-backtest/main.go
go build -o indicator-sync cmd/indicator-sync/main.go
```

### Running a Single Test
```bash
# Run specific test
go test ./trend -run TestApo -v

# Run test in specific file
go test -v ./trend/apo_test.go

# Run with coverage for specific package
go test -cover ./trend
```

### MCP Server
```bash
# Run MCP server
go run ./mcp
```

## Code Style Guidelines

### Project Structure
- **Packages**: Organized by category (trend, momentum, volatility, volume, strategy, helper, etc.)
- **Test files**: Named `*_test.go` in same package with `_test` suffix (e.g., `apo_test.go`)
- **Test data**: CSV files in `testdata/` subdirectories

### Naming Conventions
- **Types**: PascalCase (e.g., `Apo`, `Ema`, `Rsi`)
- **Constants**: PascalCase with descriptive prefix (e.g., `DefaultApoFastPeriod`)
- **Functions**: PascalCase (e.g., `NewApo`, `Compute`)
- **Variables**: camelCase (e.g., `fastPeriod`, `slowPeriod`)
- **Packages**: Single lowercase word (e.g., `trend`, `volume`)

### Generics
- Use generic type parameter `T` constrained to `helper.Number`:
  ```go
  type Apo[T helper.Number] struct {
      FastPeriod int
      FastSmoothing T
  }
  func NewApo[T helper.Number]() *Apo[T] { ... }
  ```

### Channel-Based API
- Functions take channels as input: `<-chan T`
- Functions return channels as output: `<-chan T`
- Use `helper.Buffered` for lookback periods
- Use `helper.Duplicate` for branching data streams

### Error Handling
- Return errors as last return value
- Use `errors.New()` for simple errors
- Use `fmt.Errorf` with %w for wrapped errors
- Handle errors at call site:
  ```go
  result, err := something()
  if err != nil {
      return nil, err
  }
  ```

### Testing Patterns
- Use `package_test` suffix for test packages
- Create struct for test data with CSV header tags:
  ```go
  type ApoData struct {
      Close float64 `header:"Close"`
      Apo   float64 `header:"APO"`
  }
  ```
- Use `helper.ReadFromCsvFile` for test data
- Use `helper.CheckEquals` for assertion
- Use `helper.RoundDigits` for floating-point comparison

### Code Organization
- Put related indicators in same file
- Define default constants at top of file
- Use functional options for optional configuration:
  ```go
  type Option[T any] func(*Config[T])
  func WithOption[T any](value T) Option[T] {
      return func(c *Config[T]) { c.Value = value }
  }
  ```

### Required Copyright Header
Every source file must include:
```go
// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator
```

### Imports
- Standard library first, then external packages
- Use grouped imports:
  ```go
  import (
      "fmt"
      "io"
      
      "github.com/cinar/indicator/v2/helper"
  )
  ```

### Code Quality Requirements
- **100% code coverage** is required for all indicators
- All public types and functions need documentation comments
- Use revive.toml rules (see `revive.toml`)
- Avoid unused parameters (use `_`)
- No dot imports allowed
- Error messages should be lowercase, no punctuation

### Documentation
- Use Go doc comments (starting with type/function name)
- Include formula when applicable:
  ```go
  // Apo represents...
  //
  //	Fast = Ema(values, fastPeriod)
  //	Slow = Ema(values, slowPeriod)
  //	APO = Fast - Slow
  //
  // Example:
  //
  //	apo := trend.NewApo[float64]()
  ```
- Generate docs with: `task docs`

## Key Dependencies

- `github.com/cinar/indicator/v2/helper` - Utility functions for channel operations
- Uses Go 1.22+ with generics
- Testing uses standard `testing` package
