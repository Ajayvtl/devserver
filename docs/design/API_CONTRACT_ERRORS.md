# API Contract & Error Code Specifications

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Error Codes Catalog (Gap #13)

The table below maps standard API HTTP status codes to their respective UI elements, error banners, recovery flows, and support actions.

| HTTP Code | Error Classification | Visual UI Element | Retry Rule | Support / Logs Actions |
|---|---|---|---|---|
| **401** | Unauthorized / Session Expired | Redirect to Login S-004 + Toast alert | None (Requires re-auth) | Clear local storage session tokens; log auth expiration event |
| **403** | Forbidden / Action Denied | Block mutation; show overlay banner; disable buttons | None | Render "Contact Administrator" banner; log security permission block |
| **404** | Resource Not Found | Full-page "Workspace not found" page (`ER-005`) | Manual click "Go back to projects" | Provide quick link back to root path; print correlation ID in logs |
| **409** | Conflict / Duplicate Resource | Field-level inline red message (e.g., "Email already in use") | Edit field & click submit | None |
| **422** | Unprocessable / Validation Fail | Highlight form inputs with red border; show helper labels | Edit field inputs & retry | Validate regex patterns locally before submitting next API call |
| **429** | Too Many Requests (Rate limit) | Banner: "Rate limit exceeded. Please wait..." | Auto-retry with jitter after 5s | Lock UI action button for 5s; log rate limit warning |
| **500** | Internal Server Error | Card-level error banner with red alert icon | Manual retry button | Display correlation ID `X-Correlation-ID` to user; write stack trace to logs |
| **502** | Bad Gateway | Page banner: "DevServer engine offline. Retrying..." | Auto-retry 5 times with exponential backoff (1s, 2s, 4s, 8s, 16s) | Offline status indicator in topbar header; show system status link |
| **503** | Service Unavailable | Toast alert: "System busy. Please try again." | Auto-retry 3 times with exponential backoff | Log queue status capacity issues |

---

## 2. API Contract Versioning Rules (Gap #14)

To prevent code breaking between the Go backend router and the Next.js frontend client components, all endpoints must declare their stability tier in comments and URL structures.

```
/api/v1/projects      ➔ Stable production endpoint
/api/experimental/ai ➔ Draft API; subject to change without deprecation warnings
/api/internal/debug  ➔ Excluded from open developer integrations
```

### API Stability Tiers

1. **v1 (Stable)**: Active production contracts. Cannot introduce breaking structural changes without creating a `v2` router workspace. Non-breaking additions are allowed.
2. **Deprecated**: Marked for removal. Frontend client components receive console warnings. Supported for a minimum of 2 major release iterations.
3. **Breaking**: Requires frontend components refactoring. Modifies fields, changes datatypes, or removes endpoints. Permitted only across major version jumps.
4. **Experimental**: Sandbox APIs. Subject to change or deletion without deprecation warnings. Used for drafting features.
5. **Internal**: Reserved for backend engine communication (e.g., executor heartbeats, log indexers). Excluded from open developer integrations.

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
