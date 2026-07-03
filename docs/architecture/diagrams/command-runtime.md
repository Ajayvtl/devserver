# Command Runtime Diagram

```mermaid
flowchart LR
  Frontend --> CommandBus
  AI --> CommandBus
  CommandBus --> CapabilityEngine
  CapabilityEngine --> Provider
  Provider --> TaskEngine
  TaskEngine --> EventBus
  EventBus --> EventStream
  EventStream --> Frontend
```
