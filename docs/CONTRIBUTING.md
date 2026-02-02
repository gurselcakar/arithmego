# Contributing to ArithmeGo

Thank you for your interest in contributing to ArithmeGo! This document provides guidelines and instructions for contributing.

## Getting Started

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/YOUR_USERNAME/arithmego.git
   cd arithmego
   ```

2. **Install Go 1.25 or later**
   - Download from [go.dev](https://go.dev/dl/)
   - Verify installation: `go version`

3. **Build and test**
   ```bash
   make build
   make test
   ```

## Development Workflow

1. Create a branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes

3. Run tests and linter:
   ```bash
   make test
   make lint
   ```

4. Commit your changes with a clear message

5. Push and submit a pull request

## Code Style

- Follow standard Go conventions
- Run `gofmt` and `go vet` before committing
- Keep packages decoupled:
  - `internal/game/` contains game logic with no UI or storage imports
  - `internal/storage/` contains persistence logic with no UI or game imports
- Use `internal/` for all packages (no public API)

## Pull Requests

- Provide a clear description of the changes
- Include tests for new functionality
- Keep PRs focused: one feature or fix per PR
- Ensure all tests pass before requesting review

## Reporting Issues

Use GitHub Issues to report bugs or suggest features. Please include:

- Go version (`go version`)
- Operating system and version
- Steps to reproduce the issue
- Expected vs actual behavior

## Adding New Operations

To add a new arithmetic operation:

1. **Create a new file** in `internal/game/operations/`:
   ```go
   // internal/game/operations/myoperation.go
   package operations
   ```

2. **Implement the Operation interface**:
   ```go
   type MyOperation struct{}

   func (o *MyOperation) Generate(difficulty Difficulty) (*Problem, error) {
       // Generate a problem for this operation
   }

   func (o *MyOperation) Name() string {
       return "My Operation"
   }

   func (o *MyOperation) Symbol() string {
       return "?"
   }
   ```

3. **Register via init()**:
   ```go
   func init() {
       Register("myoperation", &MyOperation{})
   }
   ```

4. **Add tests** in `internal/game/operations/myoperation_test.go`

## Questions

If you have questions about contributing, feel free to open an issue for discussion.
