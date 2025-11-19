# Usage Guide

This guide provides practical examples and step-by-step instructions for using Prow.

## Table of Contents

- [Configuration](#configuration)
- [Running Components](#running-components)
- [Job Management](#job-management)
- [Plugin Usage](#plugin-usage)
- [Tide Configuration](#tide-configuration)
- [Deck UI](#deck-ui)

## Configuration

### Basic Configuration

Prow requires two main configuration files:

1. **Prow Config** (`config.yaml`): Main Prow configuration
2. **Plugin Config** (`plugins.yaml`): Plugin configuration

### Example Prow Config

```yaml
prowjob_namespace: default
pod_namespace: test-pods
periodics:
- interval: 24h
  name: periodic-job
  spec:
    containers:
    - image: alpine:latest
      command: ["echo", "Hello Prow"]
presubmits:
  myorg/myrepo:
  - name: test-job
    always_run: true
    spec:
      containers:
      - image: alpine:latest
        command: ["make", "test"]
postsubmits:
  myorg/myrepo:
  - name: deploy-job
    spec:
      containers:
      - image: alpine:latest
        command: ["make", "deploy"]
```

### Example Plugin Config

```yaml
plugins:
  myorg/myrepo:
  - approve
  - lgtm
  - trigger
  - welcome
```

## Running Components

### Hook Component

```bash
# Basic usage
hook \
  --config-path=config.yaml \
  --plugin-config=plugins.yaml \
  --dry-run

# Production usage
hook \
  --config-path=/etc/config/config.yaml \
  --plugin-config=/etc/plugins/plugins.yaml \
  --github-token-path=/etc/github/token \
  --hmac-secret=/etc/webhook/hmac
```

### Controller Manager

```bash
# Run with plank controller
prow-controller-manager \
  --config-path=config.yaml \
  --kubeconfig=~/.kube/config \
  --enable-controller=plank

# Run with multiple controllers
prow-controller-manager \
  --config-path=config.yaml \
  --enable-controller=plank \
  --enable-controller=scheduler
```

### Tide

```bash
# Run Tide
tide \
  --config-path=config.yaml \
  --github-token-path=/path/to/token \
  --dry-run

# Production
tide \
  --config-path=/etc/config/config.yaml \
  --github-token-path=/etc/github/token \
  --github-endpoint=http://ghproxy
```

### Deck (Web UI)

```bash
# Run Deck
deck \
  --config-path=config.yaml \
  --tide-url=http://tide:8888 \
  --hook-url=http://hook:8888

# Access at http://localhost:8080
```

### Crier

```bash
# Run Crier
crier \
  --config-path=config.yaml \
  --github-token-path=/path/to/token

# With multiple reporters
crier \
  --config-path=config.yaml \
  --github-token-path=/path/to/token \
  --slack-token-path=/path/to/slack-token
```

## Job Management

### Creating Jobs

Jobs are defined in the Prow config:

```yaml
presubmits:
  myorg/myrepo:
  - name: unit-tests
    always_run: true
    spec:
      containers:
      - image: golang:1.21
        command: ["go", "test", "./..."]
```

### Triggering Jobs

Jobs can be triggered via:

1. **Automatic**: Based on `always_run` and `run_if_changed`
2. **Manual**: Using `/test <job-name>` comment
3. **Retest**: Using `/retest` comment

### Viewing Job Status

```bash
# List ProwJobs
kubectl get prowjobs

# Get specific ProwJob
kubectl get prowjob <name> -o yaml

# View job logs
kubectl logs <pod-name>
```

## Plugin Usage

### Common Plugins

#### Approve Plugin

```yaml
plugins:
  myorg/myrepo:
  - approve
  - config-updater
```

Usage:
- `/approve` - Approve a PR
- `/approve cancel` - Cancel approval

#### LGTM Plugin

```yaml
plugins:
  myorg/myrepo:
  - lgtm
```

Usage:
- `/lgtm` - Mark PR as LGTM
- `/lgtm cancel` - Cancel LGTM

#### Trigger Plugin

```yaml
plugins:
  myorg/myrepo:
  - trigger
```

Usage:
- `/test <job-name>` - Run specific job
- `/test all` - Run all jobs
- `/retest` - Retest failed jobs

#### Welcome Plugin

```yaml
plugins:
  myorg/myrepo:
  - welcome
```

Automatically welcomes new contributors.

### Custom Plugins

Create custom plugins in `pkg/plugins/`:

```go
package myplugin

func HandleIssueComment(e github.IssueCommentEvent) error {
    // Plugin logic
    return nil
}
```

## Tide Configuration

### Basic Tide Config

```yaml
tide:
  queries:
  - repos:
    - myorg/myrepo
    labels:
    - lgtm
    - approved
    missingLabels:
    - do-not-merge
```

### Merge Requirements

```yaml
tide:
  merge_method: squash
  required_status_checks:
    strict: true
    contexts:
    - continuous-integration/travis-ci/pr
```

## Deck UI

### Accessing Deck

```bash
# Local access
http://localhost:8080

# Production
https://prow.example.com
```

### Features

- **Job History**: View past job runs
- **PR History**: View PR status
- **Tide Dashboard**: View merge queue
- **Spyglass**: View artifacts
- **Plugin Help**: View available plugins

## Command-Line Tools

### mkpj - Create ProwJob

```bash
# Create ProwJob YAML
mkpj --job=test-job --config-path=config.yaml

# With specific refs
mkpj \
  --job=test-job \
  --config-path=config.yaml \
  --refs=org/repo@branch
```

### mkpod - Create Pod

```bash
# Create Pod YAML
mkpod --job=test-job --config-path=config.yaml
```

### checkconfig - Validate Config

```bash
# Validate configuration
checkconfig --config-path=config.yaml

# Validate with plugin config
checkconfig \
  --config-path=config.yaml \
  --plugin-config=plugins.yaml
```

## Best Practices

1. **Use Dry Run**: Always test with `--dry-run` first
2. **Validate Config**: Use `checkconfig` before deploying
3. **Test Locally**: Test changes locally before production
4. **Monitor Logs**: Check component logs regularly
5. **Use Labels**: Organize jobs with labels
6. **Resource Limits**: Set resource limits for jobs
7. **Artifact Retention**: Configure artifact retention policies

## Troubleshooting

### Jobs Not Running

- Check ProwJob status: `kubectl get prowjobs`
- Check Pod status: `kubectl get pods`
- Review controller logs
- Verify configuration

### Webhooks Not Working

- Check hook logs
- Verify HMAC secret
- Check GitHub webhook configuration
- Test webhook delivery

### Tide Not Merging

- Check Tide logs
- Verify PR eligibility
- Check required labels
- Review branch protection

### Deck Not Loading

- Check Deck logs
- Verify config path
- Check network connectivity
- Review browser console

## Getting Help

- Check component help: `component --help`
- Review logs with `--log-level=debug`
- Check [FAQ](FAQ.md) for common issues
- Search GitHub issues
- Ask in #sig-testing on Slack

