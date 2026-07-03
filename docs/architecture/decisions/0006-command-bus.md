# ADR 0006: Command Bus and Task Separation

## Status
Proposed

## Context
Events describe what happened. Commands describe what should happen. DevServer needs both sides of that boundary to stay clear as it grows toward automation and AI-driven workflows.

## Decision
Introduce a command bus that accepts user intent, resolves a capability provider, and hands the resulting work to the task engine. The Event Bus remains the notification channel for lifecycle updates.

## Consequences
- Writes stay explicit and auditable.
- Frontend actions become commands rather than direct API implementation calls.
- The task engine remains the execution layer, while the event bus remains the observation layer.
