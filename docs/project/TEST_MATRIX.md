# Test Matrix

This document tracks the verification and validation criteria for the DevServer platform.

| Feature / Subsystem | Unit Tests | Integration Tests | Manual Tests | Performance Tests |
| :--- | :---: | :---: | :---: | :---: |
| **Runtime Coordinator** | `Yes` | `Yes` | `Yes` | `N/A` |
| **Task Engine** | `Yes` | `Yes` | `Pending` | `Pending` |
| **Command Engine** | `Yes` | `Yes` | `Pending` | `N/A` |
| **Workspace Indexer** | `Yes` | `Yes` | `Pending` | `Pending` |
| **Event Stream / WS** | `Yes` | `Pending` | `Yes` | `Pending` |
| **File Explorer UI** | `N/A` | `Pending` | `Yes` | `Pending` |
| **Knowledge Engine** | `Yes` | `Pending` | `Yes` | `Pending` |
| **Editor Orchestrator** | `Pending`| `Pending` | `Pending` | `Pending` |
| **AI Subsystem** | `Pending`| `Pending` | `Pending` | `Pending` |

## Focus Areas for Next Sprint
1. **Performance Testing**: The Workspace Indexer requires stress testing against a 10,000+ file monolithic repository.
2. **Integration Testing**: The Event Stream requires automated WebSocket connection testing to ensure state sync reliability.
3. **Manual Validation**: End-to-end task flows (e.g., triggering a task from the UI, processing it in the backend, and observing the progress bar update via WebSocket) must be formally signed off.
