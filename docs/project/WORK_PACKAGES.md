# Work Packages

Every development effort must map to a specific Work Package (WP). Phases are immutable; Work Packages are the atomic units of implementation.

### WP-3.1.3 — Services Binding
**Status**: ACCEPTED WITH KNOWN LIMITATIONS
**Acceptance Criteria**:
- [x] Backend UI Contract Mismatch (400 Bad Request) fixed.
- [x] Provider State detection canonicalized.
- [x] Event Bus real-time UI updates (No refresh required).
- [x] ProviderCard UI overhaul (Dense + Metrics).
- [x] Known Technical Debt tracking added.

| WP | Phase | Description | Owner | Depends On | Status | Build | Tests | Manual | Done |
| :--- | :--- | :--- | :--- | :--- | :--- | :---: | :---: | :---: | :---: |
| **WP-1.1** | Phase 1 | File System Indexing | Core | None | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-1.2** | Phase 1 | AST Parsing & Symbol extraction | Core | WP-1.1 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-1.3** | Phase 1 | Workspace Provider state mapping | Core | WP-1.2 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-2.1** | Phase 2 | Unified `runtime.Component` layer | Core | WP-1.3 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-2.2** | Phase 2 | Coordinator, Lifecycle, and Registry | Core | WP-2.1 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-2.3** | Phase 2 | Async Task Engine & Command Bus | Core | WP-2.2 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-2.4** | Phase 2 | Event Stream WebSocket normalization | Core | WP-2.3 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.1** | Phase 3 | Provider Runtime Audit | Frontend | WP-2.4 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.1a** | Phase 3 | Runtime Execution Layer (Local, SSH, Docker, WSL, Kubernetes) | Core | WP-3.1.1 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.2a** | Phase 3 | Runtime Injection (Replace os/exec with executor.Runtime) | Core | WP-3.1.1a | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.2b** | Phase 3 | Redis Provider (Reference service implementation) | Core | WP-3.1.2a | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.2c** | Phase 3 | Node Provider (Reference runtime implementation) | Core | WP-3.1.2b | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.2d** | Phase 3 | Cross-Platform Service Runtime | Core | WP-3.1.2c | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.3** | Phase 3 | Services Binding | Frontend | WP-3.1.2d | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.4** | Phase 3 | Architecture Freeze (Domain Model & Canonical Entities) | Architecture | WP-3.1.3 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.4A** | Phase 3 | Canonical Domain Package (`internal/domain`) | Core | WP-3.1.4 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.4AA** | Phase 3 | Canonical Repository Interfaces (`internal/repository`) | Core | WP-3.1.4A | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.4B** | Phase 3 | Execution Context (Local, SSH, Docker) | Core | WP-3.1.4AA | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.4BA** | Phase 3 | Application Layer Interfaces | Core | WP-3.1.4B | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.4BB** | Phase 3 | Workflow Engine Contracts | Core | WP-3.1.4BA | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.4BC** | Phase 3 | Application Contracts Freeze | Core | WP-3.1.4BB | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.4C** | Phase 3 | Infrastructure Adapter Implementations | Core | WP-3.1.4BC | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.4D** | Phase 3 | Workflow Engine | Core | WP-3.1.4C | `PENDING` | ❌ | ❌ | ❌ | ❌ |
| **WP-3.1.5** | Phase 3 | Environment UI Context Switcher | Frontend | WP-3.1.4D | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-3.1.6** | Phase 3 | Logs & Domains | Frontend | WP-3.1.5 | `PENDING` | ❌ | ❌ | ❌ | ❌ |
| **WP-3.2** | Phase 3 | Architecture Compliance Review | Core | WP-3.1.6 | `PENDING` | ❌ | ❌ | ❌ | ❌ |

---

### Architectural Governance Checklist

Going forward, every completed work package must satisfy all of the following before it can be marked complete:

- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
- [ ] No new deprecated APIs (unless temporary compatibility layer)
- [ ] No circular dependencies
- [ ] Architecture compliance report
- [ ] `WORK_PACKAGES.md` updated
- [ ] `SYSTEM_MAP.md` / `ARCHITECTURE_MAP.md` progress updated
- [ ] `IMPLEMENTATION_STATUS.md` updated
- [ ] `COMPONENT_REGISTRY.md` updated (if new component)
- [ ] Existing functionality remains working (no regressions)
- [ ] Human verification steps included only if UI/API behavior changed
- [ ] No `TODO`/`FIXME` left in completed work packages
- [ ] Every interface has at least one planned implementation
- [ ] Every implementation has at least one interface
- [ ] No package imports `legacy` except temporary adapters
- [ ] Deprecated code references a removal work package
- [ ] Documentation status matches implementation status

| **WP-4.1** | Phase 4 | OpenVSCode Integration | Dev | WP-3.2 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-4.2** | Phase 4 | Editor Orchestration | Core | WP-4.1 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-4.3** | Phase 4 | IPC Bridge | Dev | WP-4.2 | `PENDING` | ❌ | ❌ | ❌ | ❌ |
| **WP-4.4** | Phase 4 | Extension Manager | Core | WP-4.3 | `COMPLETED` | ✅ | ✅ | ✅ | ✅ |
| **WP-5.1** | Phase 5 | AI Providers | Core | WP-4.4 | `PENDING` | ❌ | ❌ | ❌ | ❌ |
| **WP-5.2** | Phase 5 | Deployment Engine | Core | WP-5.1 | `PENDING` | ❌ | ❌ | ❌ | ❌ |
| **WP-5.3** | Phase 5 | Remote Orchestration | Core | WP-5.2 | `PENDING` | ❌ | ❌ | ❌ | ❌ |
| **WP-5.4** | Phase 5 | Multi-Environment Workflows | Core | WP-5.3 | `PENDING` | ❌ | ❌ | ❌ | ❌ |
| **WP-5.5** | Phase 5 | Production Hardening | QA | WP-5.4 | `PENDING` | ❌ | ❌ | ❌ | ❌ |

## Governance Rules
**Architecture Compliance Checklist** (Must be passed for every PR/Work Package):
- [ ] Imports only `internal/domain` (No duplicate entities)
- [ ] No circular dependencies
- [ ] No forbidden imports
- [ ] No new global state
- [ ] No direct UI → executor communication
- [ ] Uses Command Bus & Event Bus
- [ ] Uses Runtime Coordinator
