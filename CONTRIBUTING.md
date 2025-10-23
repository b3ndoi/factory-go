# Contributing to Factory-Go

Thank you for considering contributing to Factory-Go! This document provides guidelines for contributions.

## How to Contribute

### Reporting Issues

- Check existing issues first to avoid duplicates
- Provide a minimal reproduction case
- Include Go version and operating system
- Describe expected vs actual behavior

### Suggesting Features

- Open an issue with the `enhancement` label
- Explain the use case and why it's valuable
- Provide example code showing how it would work
- Consider if it fits the Laravel-inspired philosophy

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Run tests: `go test ./factory/...`
6. Run linter: `golangci-lint run`
7. Commit with clear messages
8. Push and create a Pull Request

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/factory-go.git
cd factory-go

# Run tests
go test -v -cover ./factory/...

# Run linter
golangci-lint run

# Test examples
cd examples/basic && go run main.go
```

## Code Guidelines

### Style

- Follow standard Go conventions (gofmt, goimports)
- Use meaningful variable names
- Add godoc comments for exported functions
- Keep functions focused and testable

### Testing

- **Required:** Tests for all new features
- **Target:** Maintain >85% coverage
- Use table-driven tests where appropriate
- Test both success and error cases

### Documentation

- Update README.md for user-facing changes
- Update CHANGELOG.md following format
- Add examples for complex features
- Keep documentation concise and accurate

## Pull Request Process

1. **CI must pass** - All tests, linting, examples
2. **Tests required** - New features need tests
3. **Documentation** - Update README/CHANGELOG
4. **One feature per PR** - Keep PRs focused
5. **Respond to feedback** - Address review comments

## Versioning

We follow [Semantic Versioning](https://semver.org/):

- **Patch (v1.0.x)** - Bug fixes, documentation
- **Minor (v1.x.0)** - New features, backward compatible
- **Major (vx.0.0)** - Breaking changes

## Code of Conduct

- Be respectful and constructive
- Welcome newcomers
- Focus on code, not people
- Assume good intentions

## Questions?

- Open an issue for discussion
- Check existing issues and examples
- Read the full documentation in README.md

---

Thank you for helping make Factory-Go better! 🎉

