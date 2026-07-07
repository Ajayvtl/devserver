# Notification & Audit Logs Catalogs — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Notification Catalog

Notifications are presented as standard toast messages or stored persistently in the notification drawer popover (M-002).

| ID | Category | Title | Default Message | Triggering Condition | Action Target |
|---|---|---|---|---|---|
| **N-001** | `SUCCESS` | Provider Connected | AI provider `{name}` is online and responding. | Connection test succeeds (M-011) | Navigate to providers list |
| **N-002** | `SUCCESS` | Workspace Created | Workspace `{name}` provisioned successfully. | Git clone & indexing finished | Redirect to workspace IDE |
| **N-003** | `SUCCESS` | Deployment Complete | Build `{version}` deployed successfully. | Release active & health checks pass | Navigate to deployment board |
| **N-004** | `ERROR` | Task Failed | Task `{command}` failed with exit code `{code}`. | Runner exits with non-zero code | Open Task console log pane |
| **N-005** | `ERROR` | Backup Failure | Scheduled database backup job `{job_id}` failed. | Backup agent database dump fails | Navigate to backups panel |
| **N-006** | `ERROR` | Deployment Failed | Compile / test phase failed for build `{version}`. | Deployment pipeline fails | Open deployment error logs |
| **N-007** | `WARNING` | Disk Space Warning | Target executor `{name}` disk space is above 90%. | Telemetry agent returns low disk | Open telemetry monitoring S-009 |
| **N-008** | `WARNING` | Executor Offline | Lost connection to executor `{name}`. | Runner heartbeat timeout | Check executor configuration |
| **N-009** | `WARNING` | Sync Conflict | Local file edits conflict with workspace server. | Git push / sync returns push error | Open merge conflict diff dialog |
| **N-010** | `INFO` | Member Invited | Invitation dispatched to `{email}`. | Member invite dispatched | None |
| **N-011** | `INFO` | Ownership Handover | You are now `{role}` of organization `{org}`. | Organization ownership transfer | Reload page |
| **N-012** | `BACKGROUND` | Indexing Files | Scanning and index workspace directory tree... | Workspace indexing begins | View tree progress indicator |
| **N-013** | `BACKGROUND` | Restoring Database | Restoring database backup from archive... | Database restore triggered | Block viewport with wizard |
| **N-014** | `SYSTEM` | Session Expiring | Your active login session expires in 5 minutes. | Security token expiration limit | Click "Keep me signed in" |

---

## 2. Audit Event Catalog

Audit logs are recorded inside the global database table `audit_logs` via Go middleware interceptors.

### 2.1 Audit Event Record Schema

Every audit log entry contains:
1. **Action**: The unique event key (from table below).
2. **Actor**: User ID and Email.
3. **Timestamp**: ISO 8601 UTC timestamp.
4. **IP Address**: Client request origin IP.
5. **Organization**: Organization ID.
6. **Workspace**: Optional workspace ID (where applicable).
7. **Old Value**: JSON string containing state before mutation.
8. **New Value**: JSON string containing state after mutation.
9. **Reason**: Optional user-provided reason (e.g., for ownership transfer or deletions).
10. **Correlation ID**: Unique trace identifier tracking request lifecycle from client to DB transaction.

### 2.2 Catalog of Logged Actions

| Action Key | Trigger / Mutation Description | Old Value Sample | New Value Sample | Severity |
|---|---|---|---|---|
| `auth.mfa_enabled` | User enables Multi-Factor Auth | `{"mfa": false}` | `{"mfa": true}` | Medium |
| `member.invited` | Admin invites a user to the Org | None | `{"email": "dev@co.com", "role": "developer"}` | Low |
| `member.role_updated` | Admin updates role of a member | `{"role": "developer"}` | `{"role": "operator"}` | Medium |
| `member.removed` | Admin removes member from Org | `{"user_id": "usr-12"}` | None | High |
| `org.owner_transferred` | Owner transfers ownership | `{"owner_id": "usr-1"}` | `{"owner_id": "usr-4"}` | Critical |
| `org.deleted` | Owner deletes the entire tenant | `{"org_id": "org-56"}` | None | Critical |
| `provider.created` | Admin creates an AI provider | None | `{"name": "Ollama", "type": "ollama"}` | Medium |
| `env.variable_updated` | Developer updates env variable | `{"value": "8080"}` | `{"value": "8081"}` | Medium |
| `env.secret_created` | Developer stores a new secret ref | None | `{"secret_ref": "DB_PASS"}` | High |
| `workspace.deleted` | Admin deletes workspace directory | `{"id": "ws-9"}` | None | Critical |
| `deploy.rollback` | Operator triggers version rollback | `{"active": "v1.4"}` | `{"active": "v1.3"}` | High |
| `backup.restore` | Owner restores DB from backup archive | `{"db_state": "dirty"}` | `{"db_state": "restored"}` | Critical |

---

## 3. Correlation ID Specification

To associate frontend client interactions with backend trace logs:
1. Client generates UUIDv4 `X-Correlation-ID` header on request.
2. Go Router extracts and injects this ID into the request context.
3. If background tasks are spawned (e.g., indexing or deployment runs), the task payload inherits this Correlation ID.
4. Database writes to `audit_logs` write the ID into the record.
5. In case of API failure, UI displays this correlation ID inside the error box (`ER-xxx`) so support can query the backend logs.

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
