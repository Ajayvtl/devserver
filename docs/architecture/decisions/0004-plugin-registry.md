# ADR 0004: Capability-Based Plugin Registry

## Status
Proposed

## Context
DevServer is moving toward plugins that can install tooling, extend workspace analysis, and trigger automated actions. A simple package list is not enough once plugins need explicit capabilities, event subscriptions, and execution permissions.

## Decision
Plugins will register capabilities explicitly, and the registry will resolve plugins by capability rather than by name alone.

## Consequences
- The platform can safely expose only the capabilities a plugin declares.
- Plugin installation and execution can be driven from the same task/event infrastructure.
- Capability changes become observable events for the frontend and AI context services.
