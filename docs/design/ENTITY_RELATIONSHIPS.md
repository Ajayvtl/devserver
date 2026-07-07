# Entity Relationship Model — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Unified Entity Relationship Diagram (ERD)

The following Mermaid diagram maps the database entity constraints, primary-to-foreign keys, and relational cardinality across the platform stack.

```mermaid
erDiagram
    USER ||--o{ MEMBERSHIP : has
    USER ||--o{ SESSION : has_active
    USER ||--o{ AUDIT_LOG : initiates
    
    ORGANIZATION ||--o{ MEMBERSHIP : contains
    ORGANIZATION ||--o{ ENVIRONMENT : defines
    ORGANIZATION ||--o{ PROVIDER : configures
    ORGANIZATION ||--o{ PROJECT : owns
    ORGANIZATION ||--o{ BACKUP_JOB : schedules
    ORGANIZATION ||--o{ AUDIT_LOG : records
    
    ROLE ||--o{ MEMBERSHIP : assigns
    
    PROJECT ||--o{ WORKSPACE : contains
    PROJECT ||--o{ DEPLOYMENT : logs
    
    WORKSPACE ||--o{ TASK : executes
    WORKSPACE ||--o{ SERVICE : hosts
    WORKSPACE ||--o{ AUDIT_LOG : tracks
    
    ENVIRONMENT ||--o{ VARIABLE : exposes
    ENVIRONMENT ||--o{ SECRET : encrypts
    ENVIRONMENT ||--o{ DEPLOYMENT : targets
    
    PROVIDER ||--o{ WORKSPACE : serves_AI_to
    
    TASK }|--|| EXECUTOR : runs_on
    TASK ||--o{ TASK_LOG : streams
    
    DEPLOYMENT ||--o{ DEPLOYMENT_LOG : writes
    
    MEMBERSHIP {
        uuid id PK
        uuid user_id FK
        uuid organization_id FK
        uuid role_id FK
        string status
        timestamp invited_at
    }
    
    USER {
        uuid id PK
        string email UK
        string password_hash
        boolean mfa_enabled
        timestamp created_at
    }
    
    SESSION {
        uuid id PK
        uuid user_id FK
        string token_hash
        string ip_address
        string user_agent
        timestamp expires_at
    }
    
    ORGANIZATION {
        uuid id PK
        string name UK
        string billing_email
        timestamp created_at
    }
    
    ROLE {
        uuid id PK
        string name UK
        string permissions
    }
    
    PROJECT {
        uuid id PK
        uuid organization_id FK
        string name
        string slug UK
        string repo_url
        timestamp created_at
    }
    
    WORKSPACE {
        uuid id PK
        uuid project_id FK
        string name
        string target_executor_id FK
        string status
        string current_branch
        timestamp created_at
    }
    
    ENVIRONMENT {
        uuid id PK
        uuid organization_id FK
        string name
        string type
        timestamp created_at
    }
    
    VARIABLE {
        uuid id PK
        uuid environment_id FK
        string key
        string value
    }
    
    SECRET {
        uuid id PK
        uuid environment_id FK
        string key
        string encrypted_value
    }
    
    PROVIDER {
        uuid id PK
        uuid organization_id FK
        string name
        string type
        string endpoint
        uuid secret_id_ref FK
    }
    
    TASK {
        uuid id PK
        uuid workspace_id FK
        string command
        string status
        integer exit_code
        timestamp started_at
        timestamp finished_at
    }
    
    EXECUTOR {
        uuid id PK
        string name
        string type
        string status
        string config_json
    }
    
    DEPLOYMENT {
        uuid id PK
        uuid project_id FK
        uuid environment_id FK
        string status
        string git_commit_sha
        timestamp created_at
    }
    
    AUDIT_LOG {
        uuid id PK
        uuid user_id FK
        uuid organization_id FK
        uuid workspace_id FK
        string action
        string old_value
        string new_value
        string ip_address
        string correlation_id
        timestamp created_at
    }
```

---

## 2. Cardinality Definitions

* **User to Membership**: One-to-Many (`1:N`). A user can belong to multiple organizations, and each membership relates to a single organization.
* **Organization to Environments**: One-to-Many (`1:N`). An organization owns multiple environment contexts (Dev, Staging, Prod), which are isolated.
* **Project to Workspaces**: One-to-Many (`1:N`). A project acts as a code grouping; developers spawn individual isolated workspaces to edit and test that project.
* **Workspace to Tasks**: One-to-Many (`1:N`). Tasks (builds, tests, index loops) run inside a workspace context and map to an infrastructure **Executor**.
* **Environment to Variables & Secrets**: One-to-Many (`1:N`). Environments isolate plain-text configurations (Variables) and cipher references (Secrets).
* **Audit Log to User / Org / Workspace**: Many-to-One (`N:1`). Audit logs capture actions driven by users within organizations, tracking workspace identifiers where applicable.

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
