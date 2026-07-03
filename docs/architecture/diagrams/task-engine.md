# Task Engine Diagram

```mermaid
flowchart LR
  Submit[Submit Task] --> Queue[Priority Queue]
  Queue --> Worker[Worker Pool]
  Worker --> Runner[core.Runner]
  Runner --> Rollback[Rollback Handler]
  Worker --> Events[Typed Events]
  Events --> Stream[EventStream]
```
