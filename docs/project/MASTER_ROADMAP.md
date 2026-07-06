# DevServer Master Roadmap

This document is the **Single Source of Truth** for the DevServer project trajectory. 
*Rule: Never rewrite phases. Only update status.*

## Phase 1: Core Architecture & Knowledge Engine
**Status**: `COMPLETED`
* File System Indexing
* AST Parsing & Cross-language Symbol extraction
* Workspace Provider state mapping

## Phase 2: Runtime Foundation & Command Engine
**Status**: `COMPLETED`
* Unified `runtime.Component` layer
* Coordinator, Lifecycle, and Registry
* Async Task Engine & Command Bus routing
* Event Stream WebSocket normalization

## Phase 3: The UI Control Plane & Environment Runtime
**Status**: `IN_PROGRESS`
* Multi-pane IDE layout
* Environment Model Abstraction (Project -> Env -> Resources)
* Global Context Switcher (Local/Staging/Prod)
* Virtualized File Explorer & context menus
* Dashboard integrations (Health, Projects, Settings)

## Phase 4: Editor Orchestrator
**Status**: `PENDING`
* Subsystem: `Editor Provider` (OpenVSCode Server wrapper)
* Subsystem: `Process Manager`
* Subsystem: `Session Manager`
* Reverse Proxy & Authentication integration
* DevServer IPC VS Code Extension

## Phase 5: Terminal & Debug Runtimes
**Status**: `PENDING`
* Subsystem: `Terminal Provider` (Shells, Docker, K8s)
* Subsystem: `Debug Provider` (DAP Orchestration)
* UI: Terminal and Debug panels

## Phase 6: AI Operating Layer
**Status**: `PENDING`
* Subsystem: `AI Runtime` (LLM dispatch & context gathering)
* UI: Chat interface & Inline edit overlays
* Editor IPC bridge (Continuous context streaming)
* AI intents (Explain, Refactor, Generate) via Command Bus

## Phase 7: Polish & Production
**Status**: `PENDING`
* Comprehensive E2E Testing
* Performance & Memory hardening
* Packaging & Installation workflows
* v1.0 Release Candidate

## Blockers & Open Issues
* **Editor Integration**: Awaiting implementation of `internal/orchestrator/editor` following Runtime foundation.
* **UI Placeholders**: Dashboard tabs for Services/Infrastructure require backend Provider bindings.
