# AI Runtime Diagram

```mermaid
flowchart LR
  Context[Workspace Context] --> AI[AI Context Service]
  AI --> Bus[EventBus]
  Bus --> Stream[EventStream]
  Stream --> Frontend[Frontend]
```
