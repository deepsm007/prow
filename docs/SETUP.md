# Setup Guide (Beginner Friendly)

This guide will help you set up a development environment for working with Prow.

## Prerequisites

### Required Software

1. **Go** (version 1.24.0 or later)
   ```bash
   # Check your Go version
   go version
   
   # If not installed, download from https://golang.org/dl/
   ```

2. **Git**
   ```bash
   git --version
   # Install via your package manager if needed
   ```

3. **Make**
   ```bash
   make --version
   # Usually pre-installed on Linux/macOS
   ```

4. **Docker** (optional, for building container images)
   ```bash
   docker --version
   ```

5. **kubectl** (for interacting with Kubernetes)
   ```bash
   kubectl version --client
   ```

### Optional but Recommended

- **kind** or **minikube** - For local Kubernetes cluster
- **ko** - For building container images
- **jq** - For JSON processing
- **yq** - For YAML processing

## Installation Steps

### 1. Clone the Repository

```bash
# Clone the repository
git clone https://github.com/kubernetes-sigs/prow.git
cd prow

# If you plan to contribute, fork first and clone your fork
```

### 2. Set Up Go Environment

```bash
# Set Go environment variables (if not already set)
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Verify Go is working
go env
```

### 3. Install Dependencies

```bash
# Download dependencies
go mod download

# Verify dependencies
go mod verify
```

### 4. Build the Components

```bash
# Build all components
make build

# Or install to $GOPATH/bin
go install ./cmd/...

# Build specific component
go build ./cmd/hook
```

### 5. Verify Installation

```bash
# Check that components are installed
which hook
hook --help

# List all available components
ls $GOPATH/bin/ | grep -E "(hook|deck|tide|plank)"
```

## Environment Requirements

### Development Environment

- **Operating System**: Linux or macOS (Windows with WSL2)
- **Memory**: Minimum 8GB RAM (16GB recommended)
- **Disk Space**: At least 10GB free space
- **Network**: Internet connection for downloading dependencies

### For Testing with Kubernetes

- Access to a Kubernetes cluster (or use kind/minikube for local)
- Valid kubeconfig file
- Appropriate cluster permissions

### For Integration Testing

- Access to a Prow cluster (for Prow team members)
- GitHub token with appropriate permissions
- GCS bucket for artifacts (optional)

## How to Run the Project

### Running Components Locally

Most components are services that run continuously. Here are examples:

#### Hook Component

```bash
# Basic usage
hook \
  --config-path=/path/to/config.yaml \
  --plugin-config=/path/to/plugins.yaml \
  --dry-run

# With GitHub token
export GITHUB_TOKEN=your-token
hook \
  --config-path=config.yaml \
  --plugin-config=plugins.yaml \
  --github-token-path=/path/to/token
```

#### Controller Manager

```bash
# Run controller manager
prow-controller-manager \
  --config-path=/path/to/config.yaml \
  --kubeconfig=~/.kube/config \
  --enable-controller=plank

# Enable multiple controllers
prow-controller-manager \
  --config-path=config.yaml \
  --enable-controller=plank \
  --enable-controller=scheduler
```

#### Deck (Web UI)

```bash
# Run Deck locally
deck \
  --config-path=/path/to/config.yaml \
  --tide-url=http://localhost:8888 \
  --hook-url=http://localhost:8888

# Access at http://localhost:8080
```

#### Tide

```bash
# Run Tide
tide \
  --config-path=/path/to/config.yaml \
  --github-token-path=/path/to/token \
  --dry-run
```

### Running Tests

```bash
# Run all unit tests
make test

# Run specific package tests
go test ./pkg/hook/...

# Run with verbose output
go test -v ./pkg/hook/...

# Run integration tests
go test ./test/integration/...
```

### Running Locally with Docker

Some components can be run in containers:

```bash
# Build container image
docker build -t prow:latest .

# Run component in container
docker run --rm prow:latest hook --help
```

## How to Test It

### Unit Testing

```bash
# Run all unit tests
make test

# Run tests for specific package
go test ./pkg/config/...

# Run tests with coverage
go test -cover ./pkg/...

# Generate coverage report
go test -coverprofile=coverage.out ./pkg/...
go tool cover -html=coverage.out
```

### Integration Testing

```bash
# Run integration tests
go test ./test/integration/...

# Run specific integration test
go test -v ./test/integration/... -run TestName
```

### Manual Testing

1. **Test Hook Locally:**
   ```bash
   # Start hook server
   hook --config-path=config.yaml --dry-run --port=8888
   
   # Send test webhook (in another terminal)
   curl -X POST http://localhost:8888/hook \
     -H "Content-Type: application/json" \
     -H "X-GitHub-Event: pull_request" \
     -d @test-webhook.json
   ```

2. **Test Configuration:**
   ```bash
   # Validate configuration
   checkconfig --config-path=config.yaml
   ```

3. **Create Test ProwJob:**
   ```bash
   # Create ProwJob YAML
   mkpj --job=test-job --config-path=config.yaml
   ```

## Common Errors and Solutions

### Error: "go: cannot find module"

**Solution:**
```bash
# Ensure you're in the repository root
cd /path/to/prow

# Download dependencies
go mod download

# Or update dependencies
go mod tidy
```

### Error: "command not found: hook"

**Solution:**
```bash
# Install components
go install ./cmd/...

# Or add $GOPATH/bin to PATH
export PATH=$PATH:$GOPATH/bin

# Or build and use directly
go build ./cmd/hook
./hook --help
```

### Error: "permission denied" when accessing cluster

**Solution:**
```bash
# Check kubeconfig
kubectl config current-context

# Verify cluster access
kubectl get nodes

# Check if you need to login
kubectl auth can-i create pods
```

### Error: "cannot load package" or import errors

**Solution:**
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Verify go.mod is correct
go mod verify
```

### Error: "no such file or directory" for test data

**Solution:**
```bash
# Ensure you're running from repository root
pwd  # Should show /path/to/prow

# Check test data exists
ls test/integration/
```

### Error: Build failures due to missing dependencies

**Solution:**
```bash
# Update dependencies
go mod tidy
go mod download
```

### Error: TypeScript compilation errors (for Deck)

**Solution:**
```bash
# Install Node.js dependencies
cd cmd/deck
npm install

# Build frontend
npm run build
```

### Error: "context deadline exceeded" when testing

**Solution:**
- Increase timeout: `go test -timeout=10m ./...`
- Check cluster connectivity
- Verify resource availability

## Development Workflow

### 1. Make Changes

```bash
# Create a feature branch
git checkout -b feature/my-feature

# Make your changes
# ...

# Format code
make update-gofmt

# Run tests
make test
```

### 2. Verify Changes

```bash
# Run linters
make verify

# Verify code generation
make verify-codegen

# Run integration tests
go test ./test/integration/...
```

### 3. Commit Changes

```bash
# Stage changes
git add .

# Commit with descriptive message
git commit -m "Add feature: description of changes"

# Push to your fork
git push origin feature/my-feature
```

## IDE Setup

### VS Code

Recommended extensions:
- Go extension
- YAML extension
- Kubernetes extension

Settings:
```json
{
  "go.useLanguageServer": true,
  "go.formatTool": "goimports",
  "editor.formatOnSave": true
}
```

### GoLand / IntelliJ

- Install Go plugin
- Configure Go SDK
- Enable goimports on save

## Next Steps

After setting up your environment:

1. Read the [Usage Guide](USAGE.md) to learn how to use Prow
2. Explore the [Codebase Walkthrough](CODEBASE_WALKTHROUGH.md) to understand the structure
3. Check the [Contributing Guide](CONTRIBUTING_GUIDE.md) if you want to contribute
4. Review the [Onboarding Guide](ONBOARDING.md) for deep dive

## Getting Help

- Check existing documentation
- Search existing issues on GitHub
- Ask in #sig-testing on Kubernetes Slack
- Create a new issue with detailed error information

