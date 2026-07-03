# DevServer Release Checklist

Use this file as the single source of truth for release readiness.

## Milestone 1 - Frontend Runtime Integration
Objective:
- The frontend becomes a thin runtime client.

Deliverables:
- [x] Replace all implementation-specific APIs with `POST /api/commands`
- [x] Replace all polling with `/ws/events`
- [x] Remove every `setInterval`, polling loop, and manual refresh
- [x] Replace mock services with command/event wrappers

Acceptance Criteria:
- [x] Zero polling remains
- [x] Zero implementation endpoints are called directly from UI
- [x] Progress is streamed live
- [x] All actions execute through the Command Bus

Status: [x] Complete

## Milestone 2 - Workspace Runtime
Objective:
- The Workspace becomes the primary operating surface.

Required Sections:
- [ ] Explorer
- [ ] Files
- [ ] File Preview
- [ ] Git
- [ ] Dependencies
- [ ] Environment
- [ ] Health
- [ ] Architecture
- [ ] Database
- [ ] Services
- [ ] Deployments
- [ ] Domains
- [ ] Logs
- [ ] Knowledge
- [ ] Doctor
- [ ] Settings

Acceptance Criteria:
- [ ] Every section uses `WorkspaceProvider`
- [ ] Every section uses real data
- [ ] No mock content remains
- [ ] No placeholders remain
- [ ] No fake values remain

Status: [ ] Open

## Milestone 3 - Provider Runtime
Providers:
- [ ] Git
- [ ] Go
- [ ] Node
- [ ] Docker
- [ ] PM2
- [ ] Nginx
- [ ] Redis
- [ ] MySQL
- [ ] PostgreSQL
- [ ] PHP
- [ ] Python
- [ ] Systemd

Every provider must expose:
- [ ] Metadata
- [ ] Version
- [ ] Health
- [ ] Status
- [ ] Capabilities
- [ ] Install
- [ ] Update
- [ ] Validate
- [ ] Uninstall

Acceptance Criteria:
- [ ] Every provider can be discovered
- [ ] Every provider can be inspected
- [ ] Every provider can be managed from the UI

Status: [ ] Open

## Milestone 4 - Command Runtime
Commands:
- [ ] Setup
- [ ] Login
- [ ] Workspace Import
- [ ] Clone Repository
- [ ] Deploy
- [ ] Install Provider
- [ ] Update Provider
- [ ] Backup
- [ ] Restore
- [ ] Doctor
- [ ] Restart Service

Acceptance Criteria:
- [ ] Commands are queued
- [ ] Commands are executed
- [ ] Commands are cancelable
- [ ] Commands are logged
- [ ] Commands are event streamed
- [ ] Commands are persisted

Status: [ ] Open

## Milestone 5 - Production UX
Remove:
- [ ] Mock
- [ ] Placeholder
- [ ] Fake Data
- [ ] TODO
- [ ] Coming Soon

Add:
- [ ] Loading states
- [ ] Empty states
- [ ] Error states
- [ ] Keyboard shortcuts
- [ ] Notifications
- [ ] Responsive layouts
- [ ] Accessibility improvements

Status: [ ] Open

## Milestone 6 - Production Validation
Must pass:
- [ ] `go test ./...`
- [ ] `go build ./...`
- [ ] `npm run build`
- [ ] TypeScript checks
- [ ] Command execution
- [ ] Event streaming
- [ ] Provider discovery
- [ ] Workspace indexing
- [ ] Login
- [ ] Setup
- [ ] Restart persistence

End-to-End Validation:
- [ ] Fresh install
- [ ] Create administrator
- [ ] Login
- [ ] Import Git repository
- [ ] Index workspace
- [ ] Browse files
- [ ] View Git state
- [ ] Inspect provider health
- [ ] Execute a command
- [ ] Watch live progress
- [ ] Verify logs
- [ ] Restart DevServer
- [ ] Confirm persistence

Status: [ ] Open

## Definition Of Done
- [ ] Backend implementation is complete
- [ ] API is complete
- [ ] Command Bus integration is complete
- [ ] Task Engine integration is complete
- [ ] Event streaming is complete
- [ ] Persistence is complete
- [ ] Frontend implementation is complete
- [ ] UI is production-ready
- [ ] Tests pass
- [ ] Manual validation passes
- [ ] Documentation is updated

## V1.0 Exit Criteria
- [ ] Every milestone is marked complete
- [ ] No mock data remains
- [ ] No placeholder UI remains
- [ ] No TODOs remain in production paths
- [ ] All acceptance criteria pass
- [ ] Manual validation succeeds
- [ ] Architecture remains unchanged throughout the release phase

## Change Policy
Every pull request or AI-generated change should update the relevant checklist item only when it is actually complete.
If a milestone changes, update the checklist and the linked architecture docs together.
