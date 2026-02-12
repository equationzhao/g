# g - Enhanced ls Alternative

**Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

g is a feature-rich, customizable, and cross-platform `ls` alternative written in Go. It provides enhanced file listing with icons, Git integration, multiple layout options, and advanced sorting capabilities.

## Working Effectively

### Prerequisites and Setup
- Requires Go version >= 1.24.0 (project uses Go 1.24.0)
- Works on Linux, Windows, and macOS
- Repository uses `go mod` for dependency management

### Bootstrap and Build Commands
**ALWAYS run these commands in sequence for a fresh setup:**

```bash
# 1. Verify Go version (must be >= 1.24)
go version

# 2. Download dependencies and build (NEVER CANCEL: first build takes ~20 seconds)
time go build -v .
```

**Build timing expectations:**
- **NEVER CANCEL**: Initial build from scratch: ~12 seconds. Set timeout to 60+ seconds.
- **NEVER CANCEL**: Subsequent builds: ~1.5 seconds. Set timeout to 30+ seconds.
- **NEVER CANCEL**: Tests: ~11 seconds. Set timeout to 30+ seconds.
- **NEVER CANCEL**: Linting: ~60 seconds. Set timeout to 180+ seconds.

### Build Variants
The project supports multiple build configurations using Go build tags:

```bash
# Lite build (minimal dependencies, 7.4MB binary)
CGO_ENABLED=0 go build -ldflags="-s -w" -o g-lite .

# Full build (all features, 8.1MB binary)  
CGO_ENABLED=0 go build -ldflags="-s -w" -tags="fuzzy mounts" -o g-full .

# Fuzzy search only
CGO_ENABLED=0 go build -ldflags="-s -w" -tags="fuzzy" -o g-fuzzy .

# Mounts support only
CGO_ENABLED=0 go build -ldflags="-s -w" -tags="mounts" -o g-mounts .
```

**Build tags:**
- `fuzzy`: Enables fuzzy search and path indexing (~500KB size impact)
- `mounts`: Enables mount point detection (~200KB size impact)

### Testing
```bash
# Run all tests (NEVER CANCEL: takes ~11 seconds)
time go test -v ./...
```

### Code Quality and CI Validation
**Always run these before committing changes:**

```bash
# Install formatting tool
go install mvdan.cc/gofumpt@latest

# Install linter (NEVER CANCEL: takes ~30 seconds to install)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Format code (required by CI)
export PATH=$PATH:~/go/bin
gofumpt -l -extra .  # Check formatting
gofumpt -w -extra .  # Fix formatting

# Lint code (NEVER CANCEL: takes ~60 seconds)
golangci-lint run ./... --timeout=3m

# Run single test file/package
go test -v ./internal/cli/                    # Test specific package
go test -v ./internal/cli/ -run TestDive      # Test specific function
go test -bench=. ./internal/cli/              # Run benchmarks
go test -bench=BenchmarkDive ./internal/cli/  # Run specific benchmark
```

**CI Pipeline:** The project has three GitHub workflows that must pass:
- `go.yml`: Multi-platform builds and tests (Linux, macOS, Windows)
- `gofumpt.yml`: Code formatting verification with gofumpt --extra
- `lint.yml`: Static code analysis with golangci-lint

**Performance Testing:** Use hyperfine for rigorous performance comparisons:
```bash
# Compare optimized vs baseline performance
hyperfine --warmup 10 --min-runs 200 \
  './g-baseline /large-test-dir' \
  './g-optimized /large-test-dir'
```

## Application Usage and Validation

### Basic Functionality Testing
**Always test these scenarios after making changes:**

```bash
# 1. Test basic listing
./g .

# 2. Test with icons and formatting
./g --icon --long .

# 3. Test tree view
./g --tree --icon .

# 4. Test Git integration (if in git repo)
./g --git --icon .

# 5. Test table format
./g --table --size --time .

# 6. Test JSON output
./g --json . | head -10
```

### Comprehensive Validation Scenario
**Always run this complete end-to-end test after making significant changes:**

```bash
# 1. Create test directory and files
mkdir -p /tmp/g-validation-test && cd /tmp/g-validation-test
mkdir -p subdir
echo "test content" > file1.txt  
echo "another test" > subdir/file2.txt

# 2. Test core functionality with actual files
/path/to/g --tree --icon --size .  # Should show tree structure with sizes
/path/to/g --table --time .         # Should show tabular output with timestamps
/path/to/g --json . | jq '.'        # Should produce valid JSON
/path/to/g --recurse .              # Should list all files recursively

# 3. Verify output contains expected elements
# - Icons should be displayed if terminal supports them
# - File sizes should be human-readable (B, KiB, etc.)
# - Tree structure should use box-drawing characters
# - JSON should be valid and parseable

# 4. Clean up
cd /tmp && rm -rf g-validation-test
```

### Advanced Feature Validation
```bash
# Test fuzzy search (requires fuzzy build tag)
./g-full --fuzzy pattern

# Test mount information (requires mounts build tag)  
./g-full --mounts .

# Test recursive listing
./g --recurse directory/

# Test shell integration generation
./g --init bash
./g --init zsh
./g --init fish
```

### Performance Expectations
- Basic listing: ~0.015 seconds
- Tree view: < 1 second for typical directories
- Recursive operations: varies by directory size
- Application startup: Nearly instant

## Key Project Structure

### Repository Layout
```
├── main.go              # Main entry point with panic handling and config loading
├── go.mod               # Go module definition (requires Go 1.24.0+)
├── justfile             # Build automation (just command) 
├── internal/            # Internal Go packages (modular architecture)
│   ├── cli/            # Command line interface and main logic
│   │   ├── g.go        # Core CLI logic with adaptive optimization strategy
│   │   ├── dive.go     # Optimized recursive directory traversal
│   │   └── helpers.go  # Strategy selection helpers
│   ├── util/           # Performance optimization utilities
│   │   ├── adaptive_strategy.go  # Intelligent processing strategy selection
│   │   ├── dirreader.go          # Batch file information retrieval
│   │   └── processors.go         # Directory processor implementations
│   ├── display/        # Multiple output format implementations (table, JSON, tree, etc.)
│   ├── filter/         # File filtering logic
│   ├── sorter/         # Advanced sorting algorithms (version-sort, etc.)
│   ├── git/            # Git status integration
│   ├── theme/          # Customizable color themes
│   ├── content/        # File content analyzers (mime-type, charset, etc.)
│   ├── item/           # File info abstractions
│   └── render/         # Terminal rendering with icon support
├── .github/workflows/  # CI/CD pipelines (go.yml, gofumpt.yml, lint.yml)
├── completions/        # Shell completions (bash, zsh, fish)
├── docs/               # Documentation including BuildOption.md
├── script/             # Development and release automation scripts
└── man/                # Manual pages
```

### Important Files and Architecture
- `main.go`: Entry point with config loading, panic handling, and argument preprocessing
- `internal/cli/g.go`: Core CLI logic implementing adaptive optimization strategies
- `internal/util/adaptive_strategy.go`: Performance optimization strategy selection (50-file threshold)
- `internal/util/dirreader.go`: Batch file operations for reduced system calls  
- `internal/cli/dive.go`: Optimized recursive directory traversal with memory pre-allocation
- `internal/display/`: Modular output formatters (grid, table, tree, JSON, markdown)
- `docs/BuildOption.md`: Build tags and feature configuration
- `CONTRIBUTING.md`: Development workflow and commit standards

### Performance Architecture
The project implements an **adaptive optimization system** that intelligently selects processing strategies:
- **Traditional processing**: Used for small directories (<50 files) and JSON output to avoid overhead
- **Batch processing**: Used for large directories (≥50 files) with optimized system calls and memory allocation
- **Recursive optimization**: Enhanced `dive.go` with pre-allocated memory and efficient string building

## Common Tasks

### Development Workflow
```bash
# 1. Make changes to source files
# 2. Format code (required by CI)
export PATH=$PATH:~/go/bin
gofumpt -w -extra .

# 3. Build and test (NEVER CANCEL: full workflow takes ~90 seconds)
go build .
go test -v ./...

# 4. Validate with lint (NEVER CANCEL: takes ~60 seconds)
golangci-lint run ./... --timeout=3m

# 5. Test functionality with real scenario
mkdir -p /tmp/validation && cd /tmp/validation
echo "test" > sample.txt
/path/to/g --tree --icon --size .
cd - && rm -rf /tmp/validation

# Complete workflow timing: ~90 seconds total
```

### Performance Optimization Development
The project includes adaptive optimization strategies. When working on performance:

```bash
# Create benchmark comparison between implementations
go test -bench=. -benchmem ./internal/cli/ > benchmark_results.txt

# Test with hyperfine for statistical validation
hyperfine --warmup 10 --min-runs 200 \
  './g-baseline -l /test-dir' \
  './g-optimized -l /test-dir'

# Test specific optimization scenarios
./g --tree --recurse large-directory/     # Test recursive optimizations
./g --json small-directory/               # Verify traditional processing
./g --table medium-directory/             # Test adaptive threshold (50+ files)
```

### Build Tools Reference
The project uses `justfile` for build automation, but core Go commands work directly:

```bash
# Instead of: just build-full
CGO_ENABLED=0 go build -ldflags="-s -w" -tags="fuzzy mounts" -o g-full .

# Instead of: just test  
go test -v ./...

# Instead of: just precheck
gofumpt -w -extra . && golangci-lint run ./...

# Check available build targets
just --list
```

### Shell Integration Setup
```bash
# Generate shell aliases
./g --init bash    # For bash
./g --init zsh     # For zsh  
./g --init fish    # For fish
./g --init powershell  # For PowerShell
```

## Troubleshooting

## Troubleshooting

### Common Issues
- **"Go version too low"**: Ensure Go >= 1.24.0 is installed
- **Build fails**: Run `go mod tidy` to sync dependencies
- **Tests fail**: Ensure working directory is repository root
- **Lint fails**: Install golangci-lint compatible with Go 1.24+
- **Performance regression**: Check if adaptive strategy threshold needs adjustment (see `internal/util/adaptive_strategy.go`)
- **Benchmark inconsistency**: Ensure proper warmup and sufficient test runs for statistical validity

### Platform-Specific Notes
- **Linux**: Full functionality available
- **macOS**: All features work, CGO enabled for Darwin builds in CI  
- **Windows**: Core functionality works, some file attribute features limited

### Performance Debugging
When investigating performance issues:

```bash
# Profile CPU usage
go build . && ./g -cpuprofile=cpu.prof /large-directory/
go tool pprof cpu.prof

# Profile memory allocation  
./g -memprofile=mem.prof /large-directory/
go tool pprof mem.prof

# Compare strategies manually
go test -bench=BenchmarkDive -benchmem ./internal/cli/
```

### Binary Size Expectations
- Default build: ~10.7MB
- Lite build (`-ldflags="-s -w"`): 7.4MB  
- Full build with all tags: 8.1MB

## Critical Reminders
- **NEVER CANCEL** any build, test, or lint command - full workflow takes up to 90 seconds
- Always set timeouts of 60+ seconds for builds, 30+ seconds for tests, 180+ seconds for linting
- Code formatting with gofumpt is mandatory - CI will fail without proper formatting
- Git status integration only works when run from within a Git repository
- Always test both lite and full build variants when making significant changes
- Use `/path/to/g` in validation scripts to reference your built binary location
- JSON output validation: always pipe through `jq '.'` to verify structure
- Test with actual files in `/tmp` directories for realistic validation scenarios
- **Performance considerations**: Changes to `internal/util/adaptive_strategy.go` require benchmark validation
- **Adaptive optimization**: The 50-file threshold in strategy selection is empirically determined - don't change without extensive testing
- **Memory allocation**: Pre-allocation strategies in `dive.go` and `dirreader.go` are critical for performance - preserve allocation patterns

## Version Information
- Current version: v0.31.0
- Go version requirement: >= 1.24.0
- License: MIT License
- Repository: https://github.com/Equationzhao/g