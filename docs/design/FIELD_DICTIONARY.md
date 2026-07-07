# Platform Field Dictionary — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Setup & Authentication Forms

### 1.1 First-Time Platform Setup Form
- **Form Path**: `/setup` (Setup Wizard S-003)
- **Permissions**: Root Admin Only

| Field Name | API Field | DB Column | Required | Nullable | Validation Rule | Max Length | Default | Placeholder | Tooltip |
|---|---|---|---|---|---|---|---|---|---|
| **Admin Email** | `adminEmail` | `users.email` | Yes | No | Format: Email (regex below) | 255 | None | `admin@company.com` | Primary root administrator account email address |
| **Password** | `password` | `users.password_hash` | Yes | No | Min 12 chars, 1 uppercase, 1 number, 1 special | 128 | None | `••••••••••••` | Passphrase for admin login (encrypted locally) |
| **Org Name** | `orgName` | `organizations.name` | Yes | No | Min 3 alphanumeric characters | 100 | None | `Acme Corp` | Default organization name to create |
| **DB Host** | `dbHost` | `config.db_host` | Yes | No | Valid hostname or IP address | 255 | `localhost` | `127.0.0.1` | Target database server address |

* **Email Regex**: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

### 1.2 User Login Form
- **Form Path**: `/login` (Login Form S-004)
- **Permissions**: Anonymous

| Field Name | API Field | DB Column | Required | Nullable | Validation Rule | Max Length | Default | Placeholder | Tooltip |
|---|---|---|---|---|---|---|---|---|---|
| **Email** | `email` | `users.email` | Yes | No | Format: Email | 255 | None | `user@company.com` | Registered account email address |
| **Password** | `password` | None | Yes | No | Min 1 char | 128 | None | `••••••••••••` | Password associated with the email account |

---

## 2. Organization & Governance Forms

### 2.1 Invite Member Form
- **Form Path**: Settings ➔ Members ➔ Invite Modal (M-006)
- **Permissions**: Admin+

| Field Name | API Field | DB Column | Required | Nullable | Validation Rule | Max Length | Default | Placeholder | Tooltip |
|---|---|---|---|---|---|---|---|---|---|
| **Email** | `email` | `memberships.invited_email` | Yes | No | Format: Email | 255 | None | `developer@company.com` | Email address of the user you want to invite |
| **Role** | `role` | `memberships.role_id` | Yes | No | One of: `owner`, `admin`, `developer`, `operator`, `viewer` | 50 | `developer` | None | Platform authorization level to assign to the user |
| **Message** | `message` | `memberships.invite_message` | No | Yes | Sanitized text | 500 | None | `Welcome to the team!` | Custom invitation message included in mail notification |

### 2.2 Create Organization Form
- **Form Path**: Settings ➔ Organizations ➔ Create Modal (M-005)
- **Permissions**: Owner Only

| Field Name | API Field | DB Column | Required | Nullable | Validation Rule | Max Length | Default | Placeholder | Tooltip |
|---|---|---|---|---|---|---|---|---|---|
| **Org Name** | `name` | `organizations.name` | Yes | No | Alphanumeric and spaces only | 100 | None | `Acme Staging` | Unique business group identifier |
| **Billing Email**| `billingEmail` | `organizations.billing_email` | Yes | No | Format: Email | 255 | Parent email | `billing@company.com` | Primary inbox for invoices and usage alerts |

---

## 3. Configuration & Infrastructure Forms

### 3.1 Create Environment Form
- **Form Path**: Configuration ➔ Environments ➔ Create Modal (M-016)
- **Permissions**: Admin+

| Field Name | API Field | DB Column | Required | Nullable | Validation Rule | Max Length | Default | Placeholder | Tooltip |
|---|---|---|---|---|---|---|---|---|---|
| **Env Name** | `name` | `environments.name` | Yes | No | Alphanumeric (e.g., `Staging-US`) | 50 | None | `Staging` | Logical environment tag name |
| **Env Type** | `type` | `environments.type` | Yes | No | One of: `development`, `staging`, `production` | 20 | `development` | None | Determines security policies and rollback gates |

### 3.2 Save Environment Variable Form
- **Form Path**: Configuration ➔ Environments ➔ Variables Table (E-004)
- **Permissions**: Developer+

| Field Name | API Field | DB Column | Required | Nullable | Validation Rule | Max Length | Default | Placeholder | Tooltip |
|---|---|---|---|---|---|---|---|---|---|
| **Key** | `key` | `variables.key` | Yes | No | Regex (see below) | 250 | None | `API_PORT` | Key name for the env variable (uppercase characters only) |
| **Value** | `value` | `variables.value` | Yes | Yes | Any text | 4000 | Empty string| `8080` | Raw value associated with the key |

* **Key Name Regex**: `^[A-Z_][A-Z0-9_]*$`

### 3.3 Save Secret Form
- **Form Path**: Configuration ➔ Environments ➔ Secrets Table (E-005)
- **Permissions**: Developer+

| Field Name | API Field | DB Column | Required | Nullable | Validation Rule | Max Length | Default | Placeholder | Tooltip |
|---|---|---|---|---|---|---|---|---|---|
| **Key** | `key` | `secrets.key` | Yes | No | Regex (same as Key Name Regex) | 250 | None | `DB_PASSWORD` | Cryptographic secret key name |
| **Value** | `value` | `secrets.encrypted_value` | Yes | No | Min 1 char (encrypted locally) | 8000 | None | `••••••••••••` | Private credential value (never exposed in plain text) |

### 3.4 Add AI Provider Form
- **Form Path**: Configuration ➔ AI Providers ➔ Add Drawer (S-014)
- **Permissions**: Admin+

| Field Name | API Field | DB Column | Required | Nullable | Validation Rule | Max Length | Default | Placeholder | Tooltip |
|---|---|---|---|---|---|---|---|---|---|
| **Name** | `name` | `providers.name` | Yes | No | Alphanumeric, spaces, dashes | 100 | None | `Ollama Local` | Human readable label for the provider |
| **Type** | `type` | `providers.type` | Yes | No | One of: `ollama`, `openai`, `anthropic`, `custom` | 50 | `ollama` | None | API target integration type |
| **Endpoint** | `endpoint` | `providers.endpoint` | Yes | No | Valid URI format | 500 | `http://localhost:11434` | `https://api.openai.com/v1` | Root API host endpoint address |
| **API Key** | `apiKey` | `providers.api_key_ref` | No | Yes | Reference to secret key | 250 | None | `DB_PASSWORD` | Associated secret containing access tokens |

---

## 4. Operational & Git Forms

### 4.1 Clone Workspace Form
- **Form Path**: Projects ➔ Workspace Clone Modal (M-013)
- **Permissions**: Developer+

| Field Name | API Field | DB Column | Required | Nullable | Validation Rule | Max Length | Default | Placeholder | Tooltip |
|---|---|---|---|---|---|---|---|---|---|
| **Workspace Name**| `name` | `workspaces.name` | Yes | No | Alphanumeric, spaces, dashes | 100 | None | `My Dev Workspace` | Unique visual label |
| **Repo URL** | `repoUrl` | `workspaces.repo_url` | Yes | No | Valid HTTP/SSH git URL | 1000 | None | `git@github.com:org/repo.git` | Target Git repository clone address |
| **Default Branch**| `branch` | `workspaces.default_branch` | No | No | Valid git branch name format | 250 | `main` | `develop` | Branch target to check out on startup |

### 4.2 Backup Action Wizard
- **Form Path**: Operations ➔ Backup ➔ Backup Dialog (M-012)
- **Permissions**: Operator+

| Field Name | API Field | DB Column | Required | Nullable | Validation Rule | Max Length | Default | Placeholder | Tooltip |
|---|---|---|---|---|---|---|---|---|---|
| **Backup Targets**| `targets` | `backup_jobs.targets` | Yes | No | Array selection (db, files, envs) | 200 | All | None | Specific resources to pack |
| **Compression** | `compression`| `backup_jobs.compression`| Yes | No | One of: `none`, `gzip`, `zstd` | 20 | `zstd` | None | Compression algorithm format |
| **Retention (Days)**| `retention`| `backup_jobs.retention`| Yes | No | Integer: 1 to 365 | 4 | `30` | `30` | Number of days before backup is auto-deleted |

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
