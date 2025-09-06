# Testing Guide

This document provides comprehensive information about testing in Zupfmanager.

## Overview

Zupfmanager uses Go's built-in testing framework for all tests. The test suite covers:

- **Unit Tests**: Individual service and component testing
- **Integration Tests**: Cross-service workflow testing
- **API Tests**: HTTP endpoint testing with `httptest`
- **Concurrent Tests**: Multi-goroutine safety testing
- **Resource Management Tests**: Proper cleanup and lifecycle testing

## Quick Start

```bash
# Run all tests
make test

# Run with coverage report
make test-coverage

# Run with race detection
make test-race

# Run benchmarks
make test-bench
```

## Test Structure

```
├── pkg/core/
│   ├── import_test.go      # Import service unit tests
│   ├── project_test.go     # Project service unit tests
│   ├── song_test.go        # Song service unit tests
│   └── services_test.go    # Service integration tests
└── pkg/api/handlers/
    └── import_test.go      # API handler tests
```

## Test Categories

### Unit Tests

Test individual services in isolation:

```go
func TestImportService_ImportFile(t *testing.T) {
    services, err := core.NewServices()
    require.NoError(t, err)
    defer services.Close()
    
    // Test implementation...
}
```

**Coverage**: Core business logic, error handling, edge cases

### Integration Tests

Test complete workflows across multiple services:

```go
func TestServices_Integration(t *testing.T) {
    t.Run("complete workflow", func(t *testing.T) {
        // Import songs -> Create project -> Add songs to project
    })
}
```

**Coverage**: End-to-end functionality, service interactions

### API Tests

Test HTTP endpoints using `httptest`:

```go
func TestImportHandler_DirectoryImport(t *testing.T) {
    // Setup test server
    // Make HTTP requests
    // Verify responses
}
```

**Coverage**: REST API functionality, request/response handling

### Concurrent Tests

Test thread safety and concurrent access:

```go
func TestServices_ConcurrentAccess(t *testing.T) {
    // Multiple goroutines accessing services simultaneously
}
```

**Coverage**: Race conditions, concurrent safety

## Test Patterns

### Service Setup

All tests use the standard service setup pattern:

```go
func TestSomething(t *testing.T) {
    services, err := core.NewServices()
    require.NoError(t, err)
    defer services.Close()
    
    ctx := context.Background()
    // Test implementation...
}
```

### Test Data

Tests create temporary data as needed:

```go
// Create temporary directory
testDir := t.TempDir()

// Create test ABC file
abcContent := `X:1
T:Test Song
M:4/4
K:C
C D E F | G A B c |]`

err := os.WriteFile(filepath.Join(testDir, "test.abc"), []byte(abcContent), 0644)
require.NoError(t, err)
```

### Assertions

Use `testify` for assertions:

```go
// Error assertions
assert.NoError(t, err)
assert.Error(t, err)

// Value assertions
assert.Equal(t, expected, actual)
assert.NotEmpty(t, result)
assert.Len(t, slice, expectedLength)

// Requirements (fail fast)
require.NoError(t, err)
require.NotNil(t, result)
```

## Running Tests

### Basic Commands

```bash
# All tests
go test ./...

# Specific package
go test ./pkg/core

# Verbose output
go test -v ./...

# With coverage
go test -cover ./...
```

### Makefile Targets

```bash
# Standard test run
make test

# Coverage with HTML report
make test-coverage

# Race detection
make test-race

# Benchmarks
make test-bench
```

### Test Filtering

```bash
# Run specific test
go test -run TestImportService_ImportFile ./pkg/core

# Run tests matching pattern
go test -run TestServices ./pkg/core

# Run subtests
go test -run TestServices_Integration/complete_workflow ./pkg/core
```

## Coverage Analysis

### Generating Reports

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

### Current Coverage

- **pkg/core**: ~26% statement coverage
- **pkg/api/handlers**: ~4% statement coverage

### Coverage Goals

- **Core Services**: Target 80%+ coverage
- **API Handlers**: Target 70%+ coverage
- **Critical Paths**: Target 90%+ coverage

## Performance Testing

### Benchmarks

Create benchmark tests for performance-critical code:

```go
func BenchmarkImportDirectory(b *testing.B) {
    services, err := core.NewServices()
    require.NoError(b, err)
    defer services.Close()
    
    // Setup test data
    testDir := setupLargeTestDirectory(b)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := services.Import.ImportDirectory(context.Background(), testDir)
        require.NoError(b, err)
    }
}
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkImportDirectory ./pkg/core

# With memory profiling
go test -bench=. -benchmem ./...
```

## Race Detection

Go's race detector finds race conditions:

```bash
# Run tests with race detection
go test -race ./...

# Via Makefile
make test-race
```

**Note**: Race detection significantly slows down tests but is crucial for concurrent code.

## Debugging Tests

### Using Delve

```bash
# Debug specific test
dlv test ./pkg/core -- -test.run TestImportService_ImportFile

# Set breakpoints and debug interactively
(dlv) break TestImportService_ImportFile
(dlv) continue
```

### Using IDE

Most Go IDEs support test debugging:

1. Set breakpoints in test code
2. Right-click test function
3. Select "Debug Test"

### Verbose Output

```bash
# See detailed test output
go test -v ./...

# See test logs
go test -v -args -test.v ./...
```

## Continuous Integration

### GitHub Actions

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.23
      
      - name: Run tests
        run: make test
      
      - name: Run tests with race detection
        run: make test-race
      
      - name: Generate coverage
        run: make test-coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

### Pre-commit Hooks

```bash
#!/bin/sh
# .git/hooks/pre-commit

# Run tests before commit
make test || exit 1

# Run race detection
make test-race || exit 1
```

## Best Practices

### Test Organization

1. **One test file per source file**: `service.go` → `service_test.go`
2. **Group related tests**: Use `t.Run()` for subtests
3. **Clear test names**: Describe what is being tested
4. **Arrange-Act-Assert**: Structure tests clearly

### Test Data

1. **Use `t.TempDir()`**: For temporary directories
2. **Clean up resources**: Use `defer` for cleanup
3. **Isolated tests**: Each test should be independent
4. **Realistic data**: Use representative test data

### Assertions

1. **Use `require` for setup**: Fail fast on setup errors
2. **Use `assert` for checks**: Continue testing after failures
3. **Meaningful messages**: Add context to assertions
4. **Test both success and failure**: Cover error paths

### Performance

1. **Parallel tests**: Use `t.Parallel()` when safe
2. **Avoid external dependencies**: Mock external services
3. **Fast feedback**: Keep tests fast for development
4. **Separate slow tests**: Use build tags for integration tests

## Troubleshooting

### Common Issues

**Database locks**: Ensure proper cleanup with `defer services.Close()`

**Race conditions**: Run with `-race` flag to detect

**Flaky tests**: Usually indicate timing issues or shared state

**Memory leaks**: Use benchmarks with `-benchmem` to detect

### Getting Help

1. Check test output for specific error messages
2. Run tests with `-v` flag for verbose output
3. Use race detection to find concurrency issues
4. Check coverage reports to identify untested code

## Recent Improvements

### Enhanced Song Search (v1.x)

The song search functionality has been enhanced to provide more flexible searching:

**Default Behavior**: 
- `song search <query>` now searches both title and filename by default
- Provides better discoverability of songs by filename patterns

**Advanced Options**:
- `--title`: Search only in song titles
- `--filename`: Search only in filenames  
- `--genre`: Search only in genres

**API Changes**:
- `/api/v1/songs/search?q=<query>` searches title + filename by default
- Use specific parameters (`title=true`, `filename=true`, `genre=true`) to limit scope

**Test Coverage**:
- Added filename search tests in `services_test.go`
- API handler tests verify correct search behavior
- CLI tests ensure proper flag handling

## Contributing

When adding new features:

1. **Write tests first**: TDD approach recommended
2. **Test error cases**: Don't just test happy paths
3. **Update documentation**: Keep this guide current
4. **Run full test suite**: Before submitting PRs

### Test Review Checklist

- [ ] Tests cover new functionality
- [ ] Error cases are tested
- [ ] Tests are independent and isolated
- [ ] Test names are descriptive
- [ ] No external dependencies
- [ ] Proper cleanup with `defer`
- [ ] Coverage doesn't decrease significantly
