# Phase 2C Roadmap

## Operating Rule
Work only in vertical slices.
Do not create new architectural packages or abstractions.
Every feature is complete only when frontend, backend, command bus, task engine, event stream, provider layer, persistence, and UI are integrated and verified together.

## Milestone 1 - Frontend Runtime Integration
Replace every implementation-specific API call with `POST /api/commands`.
Remove polling.
Drive progress only from `/ws/events`.
Make the frontend fully event-driven.

## Milestone 2 - Workspace Runtime
Finish Explorer, File Tree, File Preview, Git, Dependencies, Environment, Health, Architecture, Database, Services, Deployments, Domains, Logs, AI Context, Knowledge, Doctor, and Settings.
All sections must read from `WorkspaceProvider` and return real data.

## Milestone 3 - Provider Runtime
Finish provider metadata, version, health, status, capabilities, install, update, validate, and uninstall flows.
Keep the scope to provider health and introspection until the runtime is fully wired.

## Milestone 4 - Command Runtime
Every user action becomes a command.
Commands are queued, executed, cancelable, logged, and streamed via events.

## Milestone 5 - Workspace Explorer
Transform the workspace into the main working surface with tree, search, preview, git status, symbols, environment, database, terminal, and logs.

## Milestone 6 - Production Polish
Remove placeholders, mock data, and temporary UI.
Replace them with loading states, error handling, empty states, notifications, keyboard shortcuts, responsive layouts, and accessibility improvements.

## Milestone 7 - Production Validation
Pass Go tests, TypeScript build, frontend production build, command execution, event streaming, provider discovery, workspace indexing, setup flow, and login flow.
Then complete an end-to-end manual validation of setup, workspace import, indexing, provider actions, command execution, live progress, logs, health, and restart persistence.

## Freeze Policy
Once a milestone is complete, freeze the public API, JSON context schema, event names, command contracts, and capability contracts.
Changes require an ADR and a migration path.
