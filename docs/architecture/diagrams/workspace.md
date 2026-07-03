# Workspace Diagram

```mermaid
flowchart TD
  Open[Open Workspace] --> AddIndex[Index Workspace]
  Open --> AddWatch[Watch Workspace]
  AddIndex --> Refresh[Initial Refresh]
  Refresh --> Provider[Provider Cache]
  AddWatch --> Events[Workspace Events]
```
