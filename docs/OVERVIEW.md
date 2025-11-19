# Prow: High-Level Overview

## What is Prow?

Prow is a Kubernetes-based Continuous Integration and Continuous Deployment (CI/CD) system. It provides automated testing, code review automation, and project management features for Kubernetes and other open-source projects. Prow was originally developed as part of the Kubernetes project and has since become a widely-used CI/CD platform.

## What Problem Does It Solve?

Modern software development, especially in large open-source projects like Kubernetes, requires sophisticated CI/CD infrastructure to:

- **Automate Testing**: Run tests automatically on every pull request and commit
- **Code Review Automation**: Automate repetitive review tasks and enforce project policies
- **Project Management**: Automate issue management, labeling, and project workflows
- **Scalability**: Handle thousands of repositories and millions of test runs
- **Integration**: Integrate with GitHub, Gerrit, and other code hosting platforms
- **Resource Management**: Efficiently manage compute resources across multiple clusters

Prow solves these problems by providing a Kubernetes-native CI/CD platform that scales horizontally, integrates deeply with code hosting platforms, and offers extensive plugin-based automation.

## Who Is It For?

### Primary Users

1. **Kubernetes SIG Testing**: The team responsible for maintaining Kubernetes CI infrastructure
2. **Open Source Project Maintainers**: Teams managing large open-source projects requiring robust CI/CD
3. **CI/CD Engineers**: Engineers building and maintaining CI/CD pipelines
4. **Developers**: Contributors who interact with Prow through GitHub/Gerrit

### Skill Level Requirements

- **Basic Users**: Need to understand YAML configuration and basic CI/CD concepts
- **Advanced Users**: Should be familiar with Kubernetes, Go programming, and CI/CD best practices
- **Contributors**: Need strong Go skills, Kubernetes API knowledge, and understanding of distributed systems

## Key Features

### 1. Job Execution
Prow executes CI jobs as Kubernetes Pods:
- **Presubmit Jobs**: Run on pull requests before merge
- **Postsubmit Jobs**: Run after code is merged
- **Periodic Jobs**: Run on a schedule (cron-based)
- **Batch Jobs**: Run multiple jobs in parallel

### 2. Webhook Handling
Prow's `hook` component processes GitHub/Gerrit webhooks:
- Validates webhook signatures
- Triggers appropriate jobs based on events
- Executes plugins based on webhook content
- Manages job lifecycle

### 3. Plugin System
Extensible plugin architecture for automation:
- **Label Management**: Automatically label PRs and issues
- **Milestone Management**: Track project milestones
- **Approval Automation**: `/approve`, `/lgtm` commands
- **Test Automation**: `/test`, `/retest` commands
- **Issue Management**: Auto-close, auto-assign, etc.
- **Custom Plugins**: Write your own plugins

### 4. Tide - Automated Merging
Tide automatically merges pull requests when:
- All required tests pass
- Code is approved by required reviewers
- Branch protection rules are satisfied
- No merge conflicts exist

### 5. Deck - Web UI
Deck provides a web interface for:
- Viewing job status and history
- Browsing test results and logs
- Managing Prow configuration
- Viewing plugin help
- Spyglass integration for artifact viewing

### 6. Spyglass - Artifact Viewer
Spyglass provides a unified interface for viewing:
- Build logs
- Test results (JUnit XML)
- Coverage reports
- Custom artifacts

### 7. Status Reporting (Crier)
Crier reports job status back to:
- GitHub (commit status, PR comments)
- Gerrit (code review comments)
- Pub/Sub (for external integrations)
- Slack (notifications)

### 8. Resource Management
- **Sinker**: Cleans up old ProwJobs and Pods
- **Scheduler**: Distributes jobs across clusters
- **Plank**: Manages job execution and Pod creation

### 9. Multi-Cluster Support
Prow can distribute jobs across multiple Kubernetes clusters:
- Cluster selection based on job requirements
- Load balancing across clusters
- Cluster-specific configurations

### 10. Integration Support
- **GitHub**: Full integration with GitHub API
- **Gerrit**: Support for Gerrit code review
- **Slack**: Notifications and command execution
- **Pub/Sub**: Event streaming for external systems

## Technology Stack

- **Language**: Primarily Go (94.7%), with TypeScript (2.5%) for frontend
- **Platform**: Kubernetes (native Kubernetes resources)
- **Storage**: GCS (Google Cloud Storage) for artifacts
- **Container Registry**: Various (GCR, Quay, etc.)
- **Frontend**: React/TypeScript for Deck UI
- **API**: REST APIs and gRPC (for Gangway)

## Repository Statistics

- **Total Commits**: 16,583+ commits
- **Language Distribution**: Go 94.7%, TypeScript 2.5%, Shell 1.5%
- **License**: Apache-2.0
- **Maintainers**: Kubernetes SIG Testing

## Key Components

### Core Components

1. **hook** - Webhook server that processes GitHub/Gerrit events
2. **prow-controller-manager** - Main controller managing ProwJobs
3. **plank** - Creates and manages Pods for jobs
4. **scheduler** - Distributes jobs across clusters
5. **tide** - Automated PR merging
6. **deck** - Web UI
7. **crier** - Status reporting
8. **sinker** - Cleanup of old resources
9. **horologium** - Triggers periodic jobs
10. **status-reconciler** - Reconciles GitHub status

### Supporting Components

- **clonerefs** - Clones git repositories
- **initupload** - Initializes job artifacts
- **entrypoint** - Job entrypoint wrapper
- **sidecar** - Sidecar container for log upload
- **gcsupload** - Uploads artifacts to GCS
- **ghproxy** - GitHub API proxy with caching
- **gangway** - OAuth server for Prow access

## Related Documentation

- [Prow Documentation Site](https://docs.prow.k8s.io/)
- [Kubernetes Test Infrastructure](https://github.com/kubernetes/test-infra)
- [Prow Plugin Documentation](https://prow.k8s.io/plugins)

## Getting Started

For new users, we recommend starting with:
1. [Setup Guide](SETUP.md) - Get your development environment ready
2. [Usage Guide](USAGE.md) - Learn how to use Prow
3. [Codebase Walkthrough](CODEBASE_WALKTHROUGH.md) - Understand the structure

For contributors:
1. [Contributing Guide](CONTRIBUTING_GUIDE.md) - Learn how to contribute
2. [Onboarding Guide](ONBOARDING.md) - Deep dive into the codebase

