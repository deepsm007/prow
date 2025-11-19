# Architecture Documentation

## System Architecture

Prow is a Kubernetes-native CI/CD system that uses Kubernetes Custom Resources (ProwJobs) to manage CI/CD workflows. The architecture follows a microservices pattern where each component has a specific responsibility.

### High-Level Architecture Diagram

```mermaid
graph TB
    subgraph "External Systems"
        GitHub[GitHub/Gerrit]
        Slack[Slack]
        PubSub[Pub/Sub]
    end
    
    subgraph "Prow Components"
        Hook[hook<br/>Webhook Server]
        Controller[prow-controller-manager<br/>Main Controller]
        Plank[plank<br/>Pod Manager]
        Scheduler[scheduler<br/>Job Scheduler]
        Tide[tide<br/>Auto Merge]
        Deck[deck<br/>Web UI]
        Crier[crier<br/>Status Reporter]
        Sinker[sinker<br/>Cleanup]
        Horologium[horologium<br/>Periodic Trigger]
    end
    
    subgraph "Kubernetes Cluster"
        ProwJobs[ProwJob CRDs]
        Pods[Job Pods]
        ConfigMaps[ConfigMaps]
        Secrets[Secrets]
    end
    
    subgraph "Storage"
        GCS[Google Cloud Storage]
        Artifacts[Artifacts & Logs]
    end
    
    GitHub --> Hook
    Hook --> Controller
    Controller --> Plank
    Plank --> Scheduler
    Scheduler --> Pods
    Pods --> GCS
    GCS --> Artifacts
    
    Controller --> ProwJobs
    ProwJobs --> Plank
    
    Tide --> GitHub
    Crier --> GitHub
    Crier --> Slack
    Crier --> PubSub
    
    Deck --> GCS
    Deck --> ProwJobs
    
    Sinker --> ProwJobs
    Sinker --> Pods
    
    Horologium --> Controller
    
    style Hook fill:#e1f5ff
    style Controller fill:#fff4e1
    style Plank fill:#e8f5e9
    style Tide fill:#fce4ec
```

## Component Architecture

### Webhook Flow

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant GH as GitHub
    participant Hook as hook
    participant Controller as Controller
    participant Plank as plank
    participant Pod as Kubernetes Pod
    
    Dev->>GH: Opens/Updates PR
    GH->>Hook: Webhook Event
    Hook->>Hook: Validate Signature
    Hook->>Hook: Execute Plugins
    Hook->>Controller: Create ProwJob
    Controller->>Plank: Reconcile ProwJob
    Plank->>Pod: Create Pod
    Pod->>Pod: Execute Job
    Pod->>GCS: Upload Artifacts
    Pod->>Controller: Update Status
    Controller->>Crier: Job Complete
    Crier->>GH: Update Status
```

### Job Execution Flow

```mermaid
graph TD
    A[Webhook/Periodic Trigger] --> B[Create ProwJob]
    B --> C[Controller Reconciles]
    C --> D{Job Type?}
    D -->|Presubmit| E[Run on PR]
    D -->|Postsubmit| F[Run on Merge]
    D -->|Periodic| G[Run on Schedule]
    E --> H[Plank Creates Pod]
    F --> H
    G --> H
    H --> I[Pod Executes]
    I --> J[Upload Artifacts]
    J --> K[Update Status]
    K --> L[Crier Reports]
    
    style H fill:#e1f5ff
    style I fill:#fff4e1
    style L fill:#e8f5e9
```

### Plugin System Architecture

```mermaid
graph LR
    subgraph "hook Component"
        Webhook[Webhook Handler]
        PluginEngine[Plugin Engine]
        Plugins[Plugins]
    end
    
    subgraph "Plugin Types"
        Trigger[Trigger Plugins]
        Review[Review Plugins]
        Label[Label Plugins]
        Milestone[Milestone Plugins]
        Custom[Custom Plugins]
    end
    
    Webhook --> PluginEngine
    PluginEngine --> Plugins
    Plugins --> Trigger
    Plugins --> Review
    Plugins --> Label
    Plugins --> Milestone
    Plugins --> Custom
    
    style PluginEngine fill:#e1f5ff
    style Plugins fill:#fff4e1
```

### Data Flow Diagram

```mermaid
flowchart TD
    subgraph "Input Sources"
        GHWebhooks[GitHub Webhooks]
        GerritEvents[Gerrit Events]
        CronSchedule[Cron Schedule]
    end
    
    subgraph "Processing"
        Hook[hook]
        Controller[Controller Manager]
        Plank[plank]
        Scheduler[scheduler]
    end
    
    subgraph "Storage"
        K8sAPI[Kubernetes API]
        GCS[GCS Artifacts]
        ConfigRepo[Config Repository]
    end
    
    subgraph "Output"
        GitHubStatus[GitHub Status]
        SlackNotif[Slack Notifications]
        DeckUI[Deck UI]
        Spyglass[Spyglass]
    end
    
    GHWebhooks --> Hook
    GerritEvents --> Hook
    CronSchedule --> Horologium
    
    Hook --> Controller
    Horologium --> Controller
    Controller --> Plank
    Plank --> Scheduler
    Scheduler --> K8sAPI
    
    K8sAPI --> GCS
    ConfigRepo --> Hook
    ConfigRepo --> Controller
    
    GCS --> DeckUI
    GCS --> Spyglass
    K8sAPI --> Crier
    Crier --> GitHubStatus
    Crier --> SlackNotif
    
    style Hook fill:#e1f5ff
    style Controller fill:#fff4e1
```

## Deployment Architecture

### Production Deployment

```mermaid
graph TB
    subgraph "Prow Cluster"
        subgraph "Prow Namespace"
            HookPod[hook Pod]
            ControllerPod[Controller Pod]
            TidePod[tide Pod]
            DeckPod[deck Pod]
            CrierPod[crier Pod]
        end
    end
    
    subgraph "Build Clusters"
        Build01[build01]
        Build02[build02]
        BuildN[buildN...]
    end
    
    subgraph "External Services"
        GitHub[GitHub]
        GCS[Google Cloud Storage]
        Slack[Slack]
    end
    
    HookPod --> ControllerPod
    ControllerPod --> Build01
    ControllerPod --> Build02
    ControllerPod --> BuildN
    
    GitHub --> HookPod
    Build01 --> GCS
    Build02 --> GCS
    BuildN --> GCS
    
    CrierPod --> GitHub
    CrierPod --> Slack
    
    DeckPod --> GCS
    
    style ControllerPod fill:#e1f5ff
    style Build01 fill:#fff4e1
```

### Controller Architecture

The prow-controller-manager runs multiple controllers:

```mermaid
graph LR
    subgraph "prow-controller-manager"
        PlankCtrl[plank Controller]
        SchedulerCtrl[scheduler Controller]
    end
    
    subgraph "Kubernetes API"
        ProwJobs[ProwJobs]
        Pods[Pods]
        ConfigMaps[ConfigMaps]
    end
    
    PlankCtrl --> ProwJobs
    PlankCtrl --> Pods
    SchedulerCtrl --> ProwJobs
    SchedulerCtrl --> ConfigMaps
    
    style PlankCtrl fill:#e1f5ff
    style SchedulerCtrl fill:#fff4e1
```

## Component Interaction Diagram

```mermaid
graph TB
    subgraph "Event Layer"
        Hook[hook]
        Horologium[horologium]
    end
    
    subgraph "Control Layer"
        Controller[Controller Manager]
        Plank[plank]
        Scheduler[scheduler]
    end
    
    subgraph "Management Layer"
        Tide[tide]
        Sinker[sinker]
        StatusReconciler[status-reconciler]
    end
    
    subgraph "Reporting Layer"
        Crier[crier]
        Deck[deck]
    end
    
    Hook --> Controller
    Horologium --> Controller
    Controller --> Plank
    Plank --> Scheduler
    Scheduler --> Tide
    Tide --> Sinker
    Sinker --> StatusReconciler
    StatusReconciler --> Crier
    Crier --> Deck
    
    style Controller fill:#e1f5ff
    style Plank fill:#fff4e1
```

## Key Design Patterns

### 1. Kubernetes Native
Prow uses Kubernetes Custom Resources (ProwJobs) to represent CI jobs:
- ProwJobs are Kubernetes resources
- Jobs execute as Pods
- Standard Kubernetes patterns (watches, informers)

### 2. Controller Pattern
Components follow the Kubernetes controller pattern:
- Watch for resource changes
- Reconcile desired state
- Handle errors and retries gracefully

### 3. Plugin Architecture
Extensible plugin system:
- Plugins are Go packages
- Hook loads and executes plugins
- Plugins can interact with GitHub, Kubernetes, etc.

### 4. Multi-Cluster
Support for distributing jobs across clusters:
- Cluster selection logic
- Per-cluster configurations
- Load balancing

### 5. Event-Driven
System responds to events:
- GitHub webhooks trigger jobs
- Periodic jobs triggered by cron
- Status updates trigger reporting

## Security Architecture

```mermaid
graph TB
    subgraph "Authentication"
        GitHubOAuth[GitHub OAuth]
        HMAC[HMAC Signatures]
        ServiceAccounts[K8s Service Accounts]
    end
    
    subgraph "Authorization"
        OWNERS[OWNERS Files]
        RBAC[Kubernetes RBAC]
        PluginConfig[Plugin Config]
    end
    
    subgraph "Secrets Management"
        K8sSecrets[Kubernetes Secrets]
        Vault[Vault]
        GSM[Google Secret Manager]
    end
    
    GitHubOAuth --> OWNERS
    HMAC --> PluginConfig
    ServiceAccounts --> RBAC
    
    OWNERS --> K8sSecrets
    RBAC --> K8sSecrets
    PluginConfig --> Vault
    
    Vault --> GSM
    GSM --> K8sSecrets
    
    style HMAC fill:#e1f5ff
    style RBAC fill:#fff4e1
```

## Scalability Considerations

1. **Horizontal Scaling**: Components can be scaled horizontally
2. **Multi-Cluster**: Jobs distributed across multiple clusters
3. **Caching**: ghproxy caches GitHub API responses
4. **Resource Management**: Sinker cleans up old resources
5. **Efficient API Usage**: Informers and watches for Kubernetes API

## Monitoring and Observability

```mermaid
graph LR
    subgraph "Metrics Collection"
        Prometheus[Prometheus]
        Metrics[Metrics Endpoints]
    end
    
    subgraph "Logging"
        Stackdriver[Stackdriver Logs]
        LocalLogs[Local Log Files]
    end
    
    subgraph "Tracing"
        OpenTelemetry[OpenTelemetry]
    end
    
    subgraph "Dashboards"
        Grafana[Grafana]
        DeckUI[Deck UI]
    end
    
    Prometheus --> Grafana
    Metrics --> Prometheus
    Stackdriver --> Grafana
    OpenTelemetry --> Grafana
    DeckUI --> Metrics
    
    style Prometheus fill:#e1f5ff
    style Grafana fill:#fff4e1
```

