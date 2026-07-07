import type { NavGroup } from './types'

export const navigation: NavGroup[] = [
  {
    title: 'Overview',
    items: [
      {
        label: 'Dashboard',
        href: '/dashboard',
        subtitle: 'Summary and control deck overview',
      },
      {
        label: 'Projects',
        href: '/projects',
        subtitle: 'Workspace listings and management',
      },
      {
        label: 'Workspaces',
        href: '/workspace/devserver',
        subtitle: 'Integrated IDE developer workspace',
      },
    ],
  },
  {
    title: 'Operations & Observe',
    items: [
      {
        label: 'Telemetry Monitor',
        href: '/monitor',
        subtitle: 'Live fleet performance telemetry',
      },
      {
        label: 'Deployments',
        href: '/deploy',
        subtitle: 'Releases, rollback, and pipelines',
      },
      {
        label: 'Backup & Data',
        href: '/backup',
        subtitle: 'Backup, restore, and import logs',
      },
    ],
  },
  {
    title: 'Configuration',
    items: [
      {
        label: 'Environments',
        href: '/config/environments',
        subtitle: 'Variables and encrypted secrets',
      },
      {
        label: 'AI Providers',
        href: '/config/providers',
        subtitle: 'Ollama/OpenAI API configuration',
      },
    ],
  },
  {
    title: 'DevCenter',
    items: [
      {
        label: 'Workspace Doctor',
        href: '/devcenter/project-doctor',
        subtitle: 'Auto-healing workspace diagnostic tools',
      },
      {
        label: 'System Docs',
        href: '/devcenter/architecture',
        subtitle: 'Architecture, DB schemas, and APIs',
      },
      {
        label: 'AI Cost & Usage',
        href: '/devcenter/ai-usage',
        subtitle: 'LLM token tracker and provider metrics',
      },
    ],
  },
  {
    title: 'Governance',
    items: [
      {
        label: 'Admin Settings',
        href: '/settings',
        subtitle: 'Team members, roles, and security',
      },
    ],
  },
]
