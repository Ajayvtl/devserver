# Release Plan (v1.0)

This document outlines the scope, acceptance criteria, and deferred features for the DevServer v1.0 release.

## v1.0 Scope (The MVP)
The v1.0 release establishes DevServer as a functional, AI-ready development operating system.
* **Core Runtimes**: Workspace, Tasks, Commands, Providers, and Event Bus.
* **UI Shell**: Multi-pane dashboard with File Explorer, Knowledge Explorer, and Provider Dashboards.
* **Editor Integration**: Full embedding of OpenVSCode Server via the `Editor Orchestrator`.
* **AI Integration**: Basic LLM dispatch for "Explain" and "Refactor" commands triggered via the DevServer UI, utilizing context streamed from the Editor.
* **Infrastructure Management**: Functional Provider dashboards for Node, Redis, and Docker.

## Deferred Features (Post v1.0)
The following features are intentionally delayed to prevent scope creep:
* **Multi-User Collaboration**: Live Share equivalents and simultaneous multi-editor session orchestration.
* **Remote Deployment Pipelines**: CI/CD integration and deployment staging (Deployment Manager).
* **Native Debugger UI**: Complex DAP (Debug Adapter Protocol) UI visualization in the DevServer shell (will rely on VS Code's internal debugger for v1.0).
* **Third-Party Plugin API**: Sandboxed third-party DevServer extensions.

## v1.0 Acceptance Criteria
1. **Architectural Purity**: No subsystem bypasses the `runtime.Coordinator`.
2. **Stability**: Memory footprint remains stable after indexing a 10,000 file repository and running for 24 hours.
3. **Testing**: All components in `TEST_MATRIX.md` achieve at least Unit and Integration coverage.
4. **UX**: The user can navigate the entire application, edit code, and trigger tasks without ever feeling they have "left" a single, unified application context.
