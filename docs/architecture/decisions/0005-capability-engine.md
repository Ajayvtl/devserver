# ADR 0005: Capability Engine First

## Status
Proposed

## Context
DevServer has outgrown a plugin-first model. If the system asks for `deploy`, `install`, or `backup`, it should not need to know whether the implementation is Docker, PM2, Systemd, PostgreSQL, or another provider.

## Decision
Capabilities become the primary runtime contract. Providers advertise capability support, the resolver selects an eligible provider, and commands enter the system by capability rather than by concrete provider name.

## Consequences
- AI and the UI can target intent instead of implementation details.
- Multiple providers can satisfy the same capability.
- Plugin registries remain useful, but they are subordinate to capability resolution.
