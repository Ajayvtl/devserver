# DevServer Runtime Architecture

## 1. Executive Summary
DevServer is an AI-powered development operating system. It provides a unified control plane that orchestrates infrastructure, artificial intelligence, knowledge graphing, and source code editing into a single, highly cohesive runtime. 

The system is designed so that the **Editor is merely a subsystem**—one of many execution surfaces managed by DevServer. Everything from Docker containers and Node environments to the Editor itself and integrated Terminals are treated as first-class **Providers** orchestrated through a unified Task Engine and Event Bus.

## 2. Canonical System Architecture

```text
React UI (Unified Control Plane)
      │
      ▼  HTTP (REST) / WebSocket
DevServer API
      │
      ▼
Command Bus (Intent Routing)
      │
      ▼
Task Engine (Async Execution & Concurrency)
      │
      ├──▶ Event Bus (State Sync & Live UI Updates)
      │
      ▼
Workspace Runtime (The Operating System)
      │
      ├── Explorer (Virtual File System & Context)
      ├── Knowledge (AST Parsing, Graph Indexing)
      ├── AI (Reasoning, Code Generation, Context Assembly)
      ├── Git (Version Control Lifecycle)
      │
      ├── Provider Manager (Infrastructure Lifecycle)
      │     ├── Docker
      │     ├── Redis
      │     ├── Node
      │     └── System Services
      │
      ├── Editor Orchestrator (Editing Execution Surface)
      │     ├── Session Manager
      │     ├── Process Manager
      │     ├── Extension Manager (Extensions as Providers)
      │     ├── Event/Command Bridge (IPC)
      │     └── Editor Provider (e.g., OpenVSCode Server)
      │
      ├── Terminal Manager (Execution Shells as Providers)
      │     ├── PowerShell / Zsh
      │     ├── WSL
      │     └── Kubernetes / SSH
      │
      ├── Debug Manager (DAP Orchestration)
      │
      └── Deployment Manager (CI/CD Pipeline Execution)
```

## 3. Core Principles

### 3.1. The "Single Application" Mandate
The user must never feel they are switching contexts between a dashboard and an editor. Everything acts as one integrated product. Selecting a file in the DevServer Explorer instantly focuses the Editor Provider; asking the AI to explain a codebase queries the Knowledge Provider and highlights the relevant lines in the Editor Provider.

### 3.2. Asynchronous Command Execution
All state mutations (from spinning up Redis to renaming a file or installing an Editor Extension) follow the strict Command Pattern:
`UI -> POST /api/commands -> Command Bus -> Task Engine -> Provider Runner -> Event Bus`.
The UI never polls; it subscribes to the Event Bus to receive state updates.

### 3.3. Extensions are Providers
VS Code extensions, Python pip packages, and system dependencies are identical conceptually. They are all **Providers**. DevServer owns the installation, health-checking, updating, and disabling of all extensions across all runtimes.

### 3.4. Thin Clients, Heavy Backend
Any client bridging into DevServer (e.g., the DevServer VS Code Extension) must remain an extremely thin IPC layer. It captures continuous context (cursor position, visible range, open files) and streams it to the Event Bus. It accepts commands to manipulate the UI (e.g., `Apply Edit`). It never performs agentic reasoning, code parsing, or AI planning—that logic strictly resides in the DevServer backend.

### 3.5. Pluggable Interfaces
The architecture anticipates swap-outs. If OpenVSCode Server is deprecated, the `Editor Orchestrator` can hot-swap to Eclipse Theia by injecting a new `Editor Provider` without altering the frontend or the Task Engine. The same applies to the Terminal Manager and Debug Manager.

## 4. The Workspace Session
DevServer treats workspaces like desktop sessions. The runtime persists the exact state of the environment—open tabs, terminal history, layout dimensions, running background tasks, and AI threads—ensuring complete environment recovery upon system reboot.

## 5. Subsystem Documentation Matrix
For deep technical specifications on individual components, refer to their dedicated architecture documents:
* **Task Engine & Command Bus**: Refer to the Command Runtime spec.
* **Workspace Knowledge Engine**: Refer to the Knowledge Explorer spec.
* **Editor Provider**: Refer to `editor_provider.md`.
