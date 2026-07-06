# DevServer: Editor Orchestration Architecture

## 1. Executive Summary
The Editor subsystem elevates code editing from a standalone tool to a deeply integrated DevServer capability. DevServer does not treat the editor as the "application." Instead, the editor acts as one specific execution surface managed by an **Editor Orchestrator**. DevServer remains agnostic to the underlying technology (e.g., code-server, OpenVSCode Server, Theia) while maintaining absolute authority over lifecycle, extensions, sessions, and AI integration.

## 2. Editor Evaluation & Recommendation

A technical evaluation of web-based editors ensures the default implementation matches our enterprise, remote, and AI requirements.

| Criteria | code-server | OpenVSCode Server | Eclipse Theia | Monaco |
| :--- | :--- | :--- | :--- | :--- |
| **Extension Compatibility** | High (Open VSX) | Full (Upstream VS Code API) | High (VS Code API compatible) | None (Syntax only) |
| **Marketplace** | Open VSX | Open VSX | Open VSX | None |
| **Copilot / Proprietary Extensions** | Blocked (Microsoft EULA) | Blocked (Microsoft EULA) | Blocked (Microsoft EULA) | N/A |
| **Cursor-style AI Compatibility** | Requires IPC extension bridge | Requires IPC extension bridge | Requires IPC extension bridge | N/A |
| **Remote Dev / SSH / Containers** | Native support | Native support | Plugin dependent | N/A |
| **Debugging** | Full DAP support | Full DAP support | Full DAP support | None |
| **WebSocket Stability** | Good (custom proxy logic) | Excellent (upstream native) | Good | N/A |
| **Enterprise / Long-term** | Coder (Active) | Gitpod (Active) | Eclipse Foundation (Active) | Microsoft (Active) |

**Recommendation:** **OpenVSCode Server** is the recommended default implementation.
*Justification:* It is maintained directly from the upstream Microsoft VS Code repository by Gitpod, avoiding merge conflicts and proxy quirks associated with the `code-server` fork. It ensures absolute API compatibility for our DevServer IPC extension, robust WebSocket stability, and flawless Dev Container and debugging support. Note: Copilot and proprietary MS extensions are EULA-blocked across all open-source variants, reinforcing the need for DevServer to own the AI orchestration layer externally.

### Iframe vs Native Embedding
OpenVSCode Server is designed to own the entire viewport. 
* **Conclusion**: We will use an `<iframe>` wrapped by the DevServer reverse proxy. To mitigate the "embedded" feeling, DevServer injects custom CSS and configuration (`settings.json`) to hide the Activity Bar, Sidebar, and Status Bar, seamlessly blending the editor's execution surface with the DevServer React UI.

## 3. Editor Orchestrator Architecture

The Editor Provider does not directly manage the workspace. Instead, an **Editor Orchestrator** owns the subsystem.

```text
Workspace
      │
      ▼
Editor Orchestrator
      │
      ├── Session Manager
      ├── Process Manager
      ├── Extension Manager
      ├── Proxy Manager
      ├── Event Bridge
      └── Editor Provider
                │
         OpenVSCode/code-server
```

The Editor Provider itself only wraps the editor instance. It **never** executes OS commands directly; all lifecycle commands flow through the DevServer Command Bus -> Task Engine -> Runner -> Process Manager.

## 4. Extension Management (Extensions as Providers)

DevServer owns the extension lifecycle. The native VS Code extension manager is hidden. Instead, extensions are treated as first-class DevServer Providers (exactly like Redis or Node).

```go
// Providers
// ├── Docker
// ├── Redis
// ├── Node
// ├── Git
// └── Editor Extensions (e.g. Python Extension)
```
Extensions are managed via the DevServer UI (Install, Remove, Update, Enable, Disable, Version, Health). The backend orchestrates these operations via tasks mapped to the Editor Provider.

## 5. Event Model & IPC Bridge

DevServer provisions a thin, native DevServer VS Code extension upon `StartWorkspace`. This extension acts as a bidirectional IPC bridge. It is completely event-driven (no polling).

**Events (Editor -> Event Bus):**
* `editor.documentOpened` / `editor.documentClosed`
* `editor.documentDirty` / `editor.documentFormatted`
* `editor.activeChanged` / `editor.selectionChanged`
* `editor.cursorMoved`
* `editor.workspaceChanged`
* `editor.extensionInstalled` / `editor.extensionRemoved`
* `editor.commandExecuted`
* `editor.terminalStarted` / `editor.terminalExited`

## 6. AI Integration & Continuous Context

The DevServer VS Code extension is **strictly a thin client**. It must *never* execute reasoning, planning, code generation, or agentic logic. 

**Responsibilities:**
1. Execute DevServer commands (e.g., `Apply Edit`).
2. Continuously expose context to the DevServer Event Bus so the backend AI doesn't have to poll.
   * *Context payload includes*: Current File, Current Function, Current Class, Current Selection, Cursor Position, Visible Range, Open Files, Diagnostics, Git Branch, Terminal State.
3. Forward user intents (`Explain Selection`, `Refactor`, `Generate Tests`, `Find Bug`) to the Command Bus. All processing stays inside DevServer.

## 7. Workspace Session Model

Editor state persists as a desktop-like session.

```go
type WorkspaceSession struct {
    WorkspaceID     string
    OpenTabs        []Document
    ActiveTab       string
    PinnedTabs      []string
    CursorPosition  Position
    Selections      []Range
    SplitLayout     LayoutConfig
    ExplorerWidth   int
    PanelSizes      map[string]int
    ExpandedFolders []string
    RecentFiles     []string
    SearchHistory   []string
    CommandHistory  []string
    Terminals       []TerminalInfo
    AIHistory       []Thread
    EditorSettings  map[string]any
}
```

## 8. Future-Proofing (Decoupled Subsystems)

* **Terminals**: The editor's integrated terminal is solely for display. DevServer owns the **Terminal Manager** (PowerShell, WSL, Docker, K8s). Terminals will become independent Providers.
* **Debugger**: Debugging is treated identically to the editor. The **Debug Manager** is a standalone provider orchestrator, not tightly coupled to VS Code.
* **Multi-Instance**: The Orchestrator supports multiple simultaneous editor instances per workspace for multi-user, collaborative editing, and future architectural expansion.

## 9. Implementation Roadmap
1. **Editor Orchestrator**: Scaffolding the orchestrator and Manager components.
2. **Process Manager**: Hooking the Orchestrator into the Task Engine for process spawning.
3. **Reverse Proxy & Auth**: Setting up the HTTP/WS proxy routing.
4. **Session Manager**: Implementing the persistent `WorkspaceSession` state logic.
5. **OpenVSCode Server Integration**: Embedding the iframe and injecting themes/settings.
6. **DevServer VS Code Extension**: Building the thin IPC client (Event & Command bridge).
7. **Knowledge Integration**: Wiring "Jump to Definition" through the Command Bus.
8. **AI Integration**: Enabling continuous context streaming and AI actions.
9. **Extension Manager**: Abstracting VS Code extensions into DevServer Providers.
10. **Production Hardening**: Deprecating the React MVP editor.
