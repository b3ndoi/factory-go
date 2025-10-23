# GitHub Actions CI/CD

This directory contains GitHub Actions workflows for automated testing, linting, and releases.

## 🔄 Workflows

### 1. CI Workflow (`ci.yml`)

**Triggers:**
- Push to `main` branch
- Pull requests to `main`

**What it does:**
- ✅ Runs tests on Go 1.21, 1.22, and 1.23
- ✅ Race condition detection (`-race` flag)
- ✅ Coverage reporting (requires >85%)
- ✅ Uploads coverage to Codecov
- ✅ Caches Go modules for faster builds

**Status Badge:**
```markdown
![CI](https://github.com/b3ndoi/factory-go/actions/workflows/ci.yml/badge.svg)
```

---

### 2. Lint Workflow (`ci.yml` - lint job)

**Triggers:**
- Same as CI workflow

**What it does:**
- ✅ Runs golangci-lint with multiple linters
- ✅ Checks code quality and style
- ✅ Security scanning with gosec
- ✅ Detects common issues

**Linters enabled:**
- errcheck, gosimple, govet, ineffassign
- staticcheck, gofmt, goimports
- misspell, revive, bodyclose
- gosec, gocritic

---

### 3. Examples Workflow (`ci.yml` - examples job)

**Triggers:**
- Same as CI workflow

**What it does:**
- ✅ Verifies all examples compile
- ✅ Runs all examples to ensure they work
- ✅ Timeout protection (5s per example)

**Examples tested:**
- basic, api_testing, google_calendar_mock
- database_seeding, complete_app, faker_integration

---

### 4. Release Workflow (`release.yml`)

**Triggers:**
- Push of version tags (v1.0.0, v1.2.3, etc.)

**What it does:**
- ✅ Runs full test suite
- ✅ Calculates coverage
- ✅ Creates GitHub Release automatically
- ✅ Generates release notes
- ✅ Notifies pkg.go.dev for immediate indexing

**Usage:**
```bash
git tag v1.0.0
git push origin v1.0.0
# Release created automatically!
```

---

### 5. CodeQL Workflow (`codeql.yml`)

**Triggers:**
- Push to `main`
- Pull requests
- Weekly schedule (Monday 00:00 UTC)

**What it does:**
- ✅ Advanced security analysis
- ✅ Vulnerability scanning
- ✅ Code quality checks
- ✅ Automatic security alerts

---

### 6. Dependabot (`dependabot.yml`)

**What it does:**
- ✅ Automatically updates GitHub Actions versions
- ✅ Keeps dependencies up to date
- ✅ Creates PRs for updates weekly

---

## 📊 Status Badges

Add these to your README.md:

```markdown
[![CI](https://github.com/b3ndoi/factory-go/actions/workflows/ci.yml/badge.svg)](https://github.com/b3ndoi/factory-go/actions/workflows/ci.yml)
[![CodeQL](https://github.com/b3ndoi/factory-go/actions/workflows/codeql.yml/badge.svg)](https://github.com/b3ndoi/factory-go/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/b3ndoi/factory-go)](https://goreportcard.com/report/github.com/b3ndoi/factory-go)
[![codecov](https://codecov.io/gh/b3ndoi/factory-go/branch/main/graph/badge.svg)](https://codecov.io/gh/b3ndoi/factory-go)
```

---

## 🎯 What Happens on Each Event

### On Push to `main`:
1. CI runs (tests on 3 Go versions)
2. Lint job runs
3. Examples job runs
4. Build job runs
5. CodeQL security scan runs

### On Pull Request:
1. All CI checks must pass
2. Coverage must be >85%
3. No linting errors allowed
4. All examples must compile

### On Tag Push (v1.x.x):
1. Release workflow runs
2. Tests run one final time
3. GitHub Release created automatically
4. pkg.go.dev notified
5. Package available within minutes

---

## 🔧 Local Testing

Test workflows locally before pushing:

```bash
# Run tests like CI does
go test -v -race -coverprofile=coverage.out ./...

# Check coverage
go tool cover -func=coverage.out | grep total

# Run golangci-lint (install first: brew install golangci-lint)
golangci-lint run

# Test all examples
for dir in examples/*/; do
  (cd "$dir" && go build -o /dev/null .)
done
```

---

## 📈 Coverage Reporting

### Codecov Setup (Optional)

1. Sign up at https://codecov.io with GitHub
2. Enable for your repository
3. Coverage uploads automatically on CI runs
4. Get coverage badge from Codecov dashboard

---

## 🎉 Benefits

With this CI/CD setup, you get:

✅ **Automatic testing** on every push/PR
✅ **Multi-version support** (Go 1.21, 1.22, 1.23)
✅ **Code quality** enforcement
✅ **Security scanning** with CodeQL
✅ **Automatic releases** when tagging
✅ **Coverage tracking** with thresholds
✅ **Example validation** ensures they always work
✅ **Dependency updates** via Dependabot

**Professional-grade CI/CD!** 🚀

