# Contributing to CAS

We welcome contributions to the CAS (Computer Algebra System) project! This document provides guidelines for contributing to this Go port of Khan Academy's KAS library.

## How to Contribute

### Reporting Issues

Before creating an issue, please:

1. Search existing issues to avoid duplicates
2. Include a clear, descriptive title
3. Provide detailed steps to reproduce the problem
4. Include relevant system information (Go version, OS, etc.)
5. Add code examples when applicable

### Submitting Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Follow the existing code style** and conventions
3. **Add tests** for any new functionality
4. **Ensure all tests pass** by running `go test ./...`
5. **Update documentation** if you change APIs
6. **Write clear commit messages** describing your changes

### Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/cas.git
cd cas

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build the project
go build ./cmd/cas
```

### Code Style Guidelines

- Follow standard Go conventions and use `go fmt`
- Use meaningful variable and function names
- Add comments for complex algorithms
- Keep functions focused and reasonably sized
- Follow the existing patterns in the codebase

### Testing Guidelines

- Write unit tests for new functionality
- Maintain or improve test coverage
- Include both positive and negative test cases
- Test edge cases and error conditions
- Use table-driven tests when appropriate

### Mathematical Accuracy

This project prioritizes mathematical correctness:

- Use exact arithmetic when possible
- Preserve symbolic representations
- Test against known mathematical results
- Reference original KAS library behavior for compatibility

### Compatibility with Original KAS

When adding features:

- Check if the feature exists in the original KAS library
- Maintain API compatibility where possible
- Document any intentional deviations
- Test against original KAS test cases when available

### Documentation

- Update README.md for new features
- Add code examples for new APIs
- Update TUTORIAL.md for user-facing changes
- Include LaTeX formatting examples when relevant

## Project Structure

```
cas/
├── cmd/cas/           # CLI application
├── pkg/
│   ├── ast/           # Abstract syntax tree
│   ├── parser/        # Expression parsing
│   ├── calculus/      # Differentiation
│   ├── expand/        # Polynomial expansion
│   ├── latex/         # LaTeX formatting
│   ├── simplify/      # Expression simplification
│   └── solve/         # Equation solving
├── examples/          # Usage examples
└── testdata/         # Test data files
```

## Areas for Contribution

### High Priority
- Bug fixes and stability improvements
- Performance optimizations
- Test coverage improvements
- Documentation enhancements

### Medium Priority
- New mathematical functions
- Enhanced LaTeX formatting
- Additional solving algorithms
- CLI improvements

### Long Term
- Symbolic integration
- Matrix operations
- Complex number support
- Web API interface

## Recognition

Contributors will be acknowledged in:
- Project README
- Release notes for significant contributions
- Git commit history

## Questions?

Feel free to:
- Open an issue for questions
- Start a discussion for design decisions
- Reach out to maintainers for guidance

## License

By contributing, you agree that your contributions will be licensed under the same MIT License that covers the project. See [LICENSE](LICENSE) for details.

## Acknowledgments

This project builds upon the excellent work of Khan Academy's KAS library. We encourage contributors to respect this heritage and maintain the high standards of mathematical accuracy established by the original project.