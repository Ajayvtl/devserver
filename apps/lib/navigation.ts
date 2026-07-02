import type { NavGroup } from './types'

export const navigation: NavGroup[] = [
  {
    title: 'Overview',
    items: [
      {
        label: 'Overview',
        href: '/dashboard',
        subtitle: 'Summary and live system health',
      },
      {
        label: 'Workspace',
        href: '/workspace/devserver',
        subtitle: 'Full workspace runtime explorer',
      },
      {
        label: 'Workspaces',
        href: '/projects',
        subtitle: 'Active workspaces and environments',
      },
      {
        label: 'Settings',
        href: '/settings',
        subtitle: 'Global platform preferences',
      },
    ],
  },
  {
    title: 'Operations',
    items: [
      {
        label: 'Deployments',
        href: '/dashboard#deployments',
        subtitle: 'Release history and rollout state',
      },
      {
        label: 'Servers',
        href: '/dashboard#servers',
        subtitle: 'Registered hosts and health',
      },
      {
        label: 'Domains',
        href: '/dashboard#domains',
        subtitle: 'Public endpoints and SSL',
      },
      {
        label: 'SSL',
        href: '/dashboard#ssl',
        subtitle: 'Certificates and renewal status',
      },
      {
        label: 'Services',
        href: '/dashboard#services',
        subtitle: 'Nginx, Redis, Postgres, and friends',
      },
    ],
  },
  {
    title: 'Platform',
    items: [
      {
        label: 'Databases',
        href: '/dashboard#databases',
        subtitle: 'Platform database and app stores',
      },
      {
        label: 'Storage',
        href: '/dashboard#storage',
        subtitle: 'Volumes, snapshots, and retention',
      },
      {
        label: 'Terminal',
        href: '/dashboard#terminal',
        subtitle: 'Shell access and command execution',
      },
      {
        label: 'Monitoring',
        href: '/dashboard#monitoring',
        subtitle: 'Live metrics and alerts',
      },
      {
        label: 'Logs',
        href: '/dashboard#logs',
        subtitle: 'Execution history and audit trail',
      },
      {
        label: 'Users',
        href: '/dashboard#users',
        subtitle: 'Roles, access, and sessions',
      },
    ],
  },
  {
    title: 'DevCenter',
    items: [
      {
        label: 'Architecture',
        href: '/devcenter/architecture',
        subtitle: 'System structure and boundaries',
      },
      {
        label: 'Tasks',
        href: '/devcenter/tasks',
        subtitle: 'Task progress and rollback design',
      },
      {
        label: 'Knowledge',
        href: '/devcenter/knowledge',
        subtitle: 'Best practices and learning cards',
      },
      {
        label: 'Dependencies',
        href: '/devcenter/dependencies',
        subtitle: 'How services and modules connect',
      },
      {
        label: 'API',
        href: '/devcenter/api',
        subtitle: 'Future API contracts and facades',
      },
      {
        label: 'Database',
        href: '/devcenter/database',
        subtitle: 'Platform schema and data model',
      },
      {
        label: 'Project Doctor',
        href: '/devcenter/project-doctor',
        subtitle: 'Diagnostics and fixes for projects',
      },
    ],
  },
  {
    title: 'Admin',
    items: [
      {
        label: 'Settings',
        href: '/settings',
        subtitle: 'Platform and account preferences',
      },
    ],
  },
]
