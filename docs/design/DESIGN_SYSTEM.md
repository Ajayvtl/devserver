# Design System 2.0 — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Typography

### Font Stack

| Token | Family | Usage |
|---|---|---|
| `--font-heading` | `'Inter', sans-serif` | Headings, brand, metric values |
| `--font-body` | `'Inter', 'Segoe UI', sans-serif` | Body text, labels, descriptions |
| `--font-mono` | `'JetBrains Mono', 'Fira Code', monospace` | Code blocks, terminal, file paths |

### Type Scale

| Token | Size | Weight | Line Height | Usage |
|---|---|---|---|---|
| `--text-display` | `clamp(2.5rem, 7vw, 5rem)` | 800 | 1.0 | Bootstrap hero, page hero |
| `--text-h1` | `clamp(2rem, 4vw, 3.35rem)` | 700 | 1.05 | Section titles |
| `--text-h2` | `1.3rem` | 700 | 1.2 | Card titles |
| `--text-h3` | `1.15rem` | 700 | 1.3 | Sidebar brand name |
| `--text-h4` | `1rem` | 600 | 1.4 | Timeline item titles |
| `--text-body` | `1rem` | 400 | 1.65 | Default body text |
| `--text-sm` | `0.92rem` | 400 | 1.45 | Subtitles, helper text |
| `--text-xs` | `0.88rem` | 400 | 1.4 | Timestamps, meta labels |
| `--text-eyebrow` | `0.75rem` | 700 | 1.0 | Uppercase group labels |
| `--text-badge` | `0.8rem` | 700 | 1.0 | Badge labels |
| `--text-status` | `0.78rem` | 700 | 1.0 | Status pill labels |

### Letter Spacing

| Token | Value | Usage |
|---|---|---|
| `--tracking-tight` | `-0.03em` | Headings |
| `--tracking-normal` | `0` | Body text |
| `--tracking-wide` | `0.08em` | Status labels |
| `--tracking-caps` | `0.12em` | Eyebrow/group titles |
| `--tracking-ultra` | `0.14em` | Topbar eyebrow |

---

## 2. Spacing

### Scale (8px base)

| Token | Value | Usage |
|---|---|---|
| `--space-0` | `0` | Reset |
| `--space-1` | `4px` | Inline micro gaps |
| `--space-2` | `6px` | Badge/status padding |
| `--space-3` | `8px` | Icon gaps |
| `--space-4` | `10px` | Grid gaps, nav item spacing |
| `--space-5` | `12px` | Card internal gaps, list gaps |
| `--space-6` | `14px` | Stack gaps, timeline gaps |
| `--space-7` | `16px` | Card padding, row padding |
| `--space-8` | `18px` | Section margins, card padding |
| `--space-9` | `22px` | Page gap, section header padding |
| `--space-10` | `24px` | Bootstrap bars margin |
| `--space-11` | `26px` | Sidebar padding |
| `--space-12` | `28px` | Content padding, topbar padding |
| `--space-13` | `34px` | Bootstrap card padding |

---

## 3. Elevation (Shadows)

| Token | Value | Usage |
|---|---|---|
| `--elevation-none` | `none` | Flat elements |
| `--elevation-sm` | `0 4px 12px rgba(0, 0, 0, 0.15)` | Dropdowns, tooltips |
| `--elevation-md` | `0 12px 30px rgba(87, 212, 255, 0.3)` | Brand mark glow |
| `--elevation-lg` | `0 18px 38px rgba(87, 212, 255, 0.18)` | Brand lockup |
| `--elevation-xl` | `0 24px 80px rgba(0, 0, 0, 0.35)` | Cards, modals, bootstrap card |

---

## 4. Motion

### Duration

| Token | Value | Usage |
|---|---|---|
| `--duration-instant` | `0ms` | Immediate feedback |
| `--duration-fast` | `120ms` | Button hover, micro-interactions |
| `--duration-normal` | `140ms` | Nav item hover, border transitions |
| `--duration-slow` | `300ms` | Page transitions, drawer open |
| `--duration-emphasis` | `500ms` | Toast enter/exit |

### Easing

| Token | Value | Usage |
|---|---|---|
| `--ease-default` | `ease` | General transitions |
| `--ease-in-out` | `ease-in-out` | Loading animations |
| `--ease-spring` | `cubic-bezier(0.34, 1.56, 0.64, 1)` | Bounce effects |

### Keyframes

| Name | Usage |
|---|---|
| `barPulse` | Bootstrap loading bars |
| `fadeIn` | Toast appearance, modal entry |
| `slideUp` | Drawer, panel reveal |
| `shimmer` | Skeleton loading animation |

---

## 5. Color Tokens

### Core Palette

| Token | Value | Usage |
|---|---|---|
| `--bg` | `#07111d` | Page background |
| `--bg-elevated` | `rgba(10, 17, 30, 0.82)` | Elevated surfaces |
| `--surface` | `rgba(12, 21, 38, 0.92)` | Card backgrounds |
| `--surface-strong` | `#101b31` | Stronger surface distinction |
| `--text` | `#eef4ff` | Primary text |
| `--muted` | `#97a5bf` | Secondary text |
| `--muted-strong` | `#c7d2e6` | Emphasized secondary text |

### Semantic Colors

| Token | Value | Usage |
|---|---|---|
| `--accent` | `#57d4ff` | Primary brand, links, active states |
| `--accent-strong` | `#2ec4ff` | Hover/focus accent |
| `--success` | `#49d08e` | Healthy, connected, passed |
| `--warning` | `#f6c15f` | Degraded, pending, caution |
| `--danger` | `#f08a8a` | Failed, offline, destructive |
| `--info` | `#82a7ff` | Informational, neutral highlight |

### Border Colors

| Token | Value | Usage |
|---|---|---|
| `--border` | `rgba(148, 163, 184, 0.18)` | Default borders |
| `--border-strong` | `rgba(148, 163, 184, 0.28)` | Hover/focus borders |
| `--border-accent` | `rgba(87, 212, 255, 0.4)` | Active nav item border |
| `--border-code` | `rgba(87, 212, 255, 0.22)` | Code block borders |

### Gradient Presets

| Name | Value | Usage |
|---|---|---|
| `gradient-primary` | `linear-gradient(135deg, rgba(87,212,255,0.98), rgba(73,208,142,0.95))` | Primary buttons, brand mark |
| `gradient-card` | `linear-gradient(180deg, rgba(15,25,41,0.94), rgba(9,15,28,0.9))` | Card backgrounds |
| `gradient-sidebar` | `linear-gradient(180deg, rgba(5,10,19,0.82), rgba(5,10,19,0.55))` | Sidebar background |
| `gradient-active-nav` | `linear-gradient(135deg, rgba(87,212,255,0.16), rgba(73,208,142,0.08))` | Active nav item |

---

## 6. Radius

| Token | Value | Usage |
|---|---|---|
| `--radius-xs` | `4px` | Inline badges, small inputs |
| `--radius-sm` | `12px` | Compact UI elements |
| `--radius-md` | `18px` | Cards, panels, rows, code blocks |
| `--radius-lg` | `24px` | Large cards, metric cards |
| `--radius-xl` | `28px` | Bootstrap card, modals |
| `--radius-pill` | `999px` | Badges, status pills, dots |
| `--radius-button` | `14px` | Buttons |
| `--radius-brand` | `16px` | Brand mark |

---

## 7. Icons

### Icon System

| Property | Value |
|---|---|
| Library | Lucide React (recommended) |
| Default Size | `20px` |
| Stroke Width | `1.5` |
| Color Inheritance | `currentColor` |

### Required Icon Set

| Context | Icons Needed |
|---|---|
| Navigation | Dashboard, Folder, Settings, Code, Terminal, Server, Database |
| Actions | Plus, Edit, Trash, Copy, Download, Upload, Refresh |
| Status | Check, AlertTriangle, XCircle, Clock, Loader |
| Auth | Lock, Unlock, User, Users, Shield |
| AI/ML | Brain, Sparkles, MessageSquare, Wand |
| Infrastructure | Server, HardDrive, Globe, Wifi, Container |

---

## 8. Accessibility

### Color Contrast

| Pair | Ratio | WCAG Level |
|---|---|---|
| `--text` on `--bg` | 15.2:1 | AAA ✅ |
| `--muted` on `--bg` | 5.8:1 | AA ✅ |
| `--accent` on `--bg` | 8.4:1 | AAA ✅ |
| Badge text on badge bg | Variable | Must verify per badge |
| Status text on status bg | Dark on light | Must maintain 4.5:1 min |

### Focus Indicators

| Property | Value |
|---|---|
| Focus ring color | `var(--accent)` |
| Focus ring width | `2px` |
| Focus ring offset | `2px` |
| Focus ring style | `solid` |
| Visible on | All interactive elements |

### ARIA Requirements

| Element | Required Attributes |
|---|---|
| Navigation | `aria-label="Sidebar"`, `aria-current="page"` |
| Buttons | `aria-label` when icon-only |
| Modals | `role="dialog"`, `aria-modal="true"`, `aria-labelledby` |
| Tables | `role="table"`, header cells with `scope` |
| Loading states | `aria-busy="true"`, `aria-live="polite"` |
| Error states | `role="alert"`, `aria-live="assertive"` |
| Form inputs | `aria-required`, `aria-invalid`, `aria-describedby` |

---

## 9. Responsive Grid

### Breakpoints

| Token | Value | Target |
|---|---|---|
| `--bp-mobile` | `≤ 720px` | Mobile phones |
| `--bp-tablet` | `721px – 960px` | Tablets |
| `--bp-desktop-sm` | `961px – 1180px` | Small desktops |
| `--bp-desktop` | `≥ 1181px` | Full desktop |

### Layout Grid

| Breakpoint | Shell Columns | Metric Grid | Page Grid |
|---|---|---|---|
| Desktop | `310px + 1fr` | `4 columns` | `2 columns` |
| Sm Desktop | `310px + 1fr` | `2 columns` | `2 columns` |
| Tablet | `1fr` (sidebar stacks) | `1 column` | `1 column` |
| Mobile | `1fr` | `1 column` | `1 column` |

### Content Width

| Token | Value | Usage |
|---|---|---|
| `--content-max` | `1420px` | Page content maximum |
| `--content-prose` | `70ch` | Description text maximum |
| `--sidebar-width` | `310px` | Fixed sidebar width |

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins.
