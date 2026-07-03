# Architecture Freeze v1

These components are considered stable unless there is a strong architectural reason to change them.

## Frozen Components
- Event Bus
- Command Bus
- Task Engine
- Workspace Manager
- Workspace Provider
- Capability Engine
- Filesystem Interface
- Context JSON Schema
- API Contract

## Rule
Changes to frozen components require an ADR, a clear rationale, and a migration path where the public contract changes.

## Freeze Policy
When a milestone is complete, freeze the public API, JSON context schema, event names, command contracts, and capability contracts.

## Purpose
This document exists to keep the platform coherent as it grows. The runtime can evolve, but the core contracts should not drift casually.
