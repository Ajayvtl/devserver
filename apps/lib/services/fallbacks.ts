import type {
  BootstrapDecision,
  DashboardDataV2,
  DevCenterSectionData,
  KnowledgeArticle,
  LoginData,
  ProjectDetailData,
  ProjectFormData,
  ProjectListItem,
  SettingsData,
  SetupWizardData,
} from '../types'

export const setupFallback: SetupWizardData = {
  steps: [
    { key: 'welcome', title: 'Welcome', subtitle: 'Meet DevServer and confirm the baseline' },
    { key: 'license', title: 'License', subtitle: 'Review licensing and usage terms' },
    { key: 'scan', title: 'System Scan', subtitle: 'Validate platform readiness' },
    { key: 'admin', title: 'Administrator', subtitle: 'Create the first admin account' },
    { key: 'ai', title: 'AI Providers', subtitle: 'Configure optional providers' },
    { key: 'install-type', title: 'Installation Type', subtitle: 'Choose a quick or advanced flow' },
    { key: 'summary', title: 'Summary', subtitle: 'Review before changes are applied' },
    { key: 'installing', title: 'Installing', subtitle: 'Provisioning is in progress' },
    { key: 'finished', title: 'Finished', subtitle: 'The platform is ready' },
  ],
  checks: [
    { label: 'Ubuntu', detail: 'Detected 24.04 LTS baseline', status: 'passed' },
    { label: 'CPU', detail: '8 cores available', status: 'passed' },
    { label: 'RAM', detail: '16 GB available', status: 'passed' },
    { label: 'Disk', detail: '120 GB free on root volume', status: 'warning' },
    { label: 'Internet', detail: 'Connectivity is available', status: 'passed' },
    { label: 'DNS', detail: 'Resolver responds to lookups', status: 'passed' },
    { label: 'Swap', detail: 'Swap is enabled and healthy', status: 'passed' },
    { label: 'SSH', detail: 'SSH service is reachable', status: 'passed' },
    { label: 'Firewall', detail: 'Default policy not yet applied', status: 'warning' },
    { label: 'Git', detail: 'git binary found', status: 'passed' },
    { label: 'Node', detail: 'Node can be installed later', status: 'pending' },
    { label: 'Redis', detail: 'Not installed yet', status: 'pending' },
    { label: 'Nginx', detail: 'Not installed yet', status: 'pending' },
    { label: 'Postgres', detail: 'Not installed yet', status: 'pending' },
    { label: 'PM2', detail: 'Optional process manager unavailable', status: 'pending' },
  ],
  providers: ['OpenAI', 'Anthropic', 'Local Model', 'None'],
  installationTypes: [
    { key: 'quick', label: 'Quick install', detail: 'Recommended defaults with sensible production-ready settings.' },
    { key: 'advanced', label: 'Advanced install', detail: 'Choose every module and inspect each generated step.' },
  ],
  knowledgeArticles: [
    { title: 'System scan basics', question: 'What are we checking before install?', summary: 'We validate OS, memory, storage, networking, and service prerequisites.', link: '#learn-scan' },
    { title: 'What is bootstrap?', question: 'Why does DevServer need a setup wizard?', summary: 'Bootstrap is the controlled path from a clean machine to a managed platform.', link: '#learn-bootstrap' },
  ],
  taskId: 'task-boot-setup',
}

export const loginFallback: LoginData = {
  branding: 'DevServer',
  subtitle: 'Control server bootstrap, deployments, and infrastructure workflows from one place.',
  supportEmail: 'support@devserver.local',
  knowledgeArticles: [
    { title: 'JWT sessions', question: 'How does authentication stay stateless?', summary: 'The UI will post credentials, receive a token, and refresh transparently.', link: '#learn-jwt' },
    { title: 'RBAC', question: 'How are permissions grouped?', summary: 'Roles like Admin, Developer, and Viewer are derived from the backend policy.', link: '#learn-rbac' },
  ],
}

export const dashboardFallback: DashboardDataV2 = {
  server: 'production-east-1',
  notifications: 4,
  metrics: [
    { label: 'CPU', value: '42%', detail: '32% average over 15 min', trend: '+4%', tone: 'success' },
    { label: 'RAM', value: '68%', detail: 'A little warm on two hosts', trend: '+2%', tone: 'warning' },
    { label: 'Disk', value: '51%', detail: 'Plenty of headroom remains', trend: 'stable', tone: 'info' },
    { label: 'Network', value: '24 ms', detail: 'Median request path latency', trend: '-7%', tone: 'accent' },
  ],
  sections: [
    {
      title: 'Workspaces',
      detail: 'Active workspaces and their health.',
      items: [
        { label: 'Aurora', value: 'healthy', tone: 'success' },
        { label: 'Docs Portal', value: 'deploying', tone: 'warning' },
        { label: 'Internal API', value: 'stable', tone: 'info' },
      ],
    },
    {
      title: 'Domains',
      detail: 'DNS and SSL coverage.',
      items: [
        { label: 'app.devserver.local', value: 'valid', tone: 'success' },
        { label: 'api.devserver.local', value: 'renew in 22 days', tone: 'warning' },
        { label: 'staging.devserver.local', value: 'healthy', tone: 'success' },
      ],
    },
    {
      title: 'Services',
      detail: 'Runtime services and systemd state.',
      items: [
        { label: 'nginx', value: 'running', tone: 'success' },
        { label: 'postgres', value: 'running', tone: 'success' },
        { label: 'redis', value: 'degraded', tone: 'warning' },
      ],
    },
    {
      title: 'Deployments',
      detail: 'Release history and rollout state.',
      items: [
        { label: 'Production', value: '2 hours ago', tone: 'success' },
        { label: 'Staging', value: '12 minutes ago', tone: 'accent' },
        { label: 'Preview', value: 'queued', tone: 'warning' },
      ],
    },
  ],
  activity: [
    { title: 'Redis task queued', detail: 'Prepare memory policy and reload service', when: '2m ago', tone: 'warning' },
    { title: 'Deploy approved', detail: 'Production deployment passed validation', when: '11m ago', tone: 'success' },
    { title: 'SSL renewal warning', detail: 'One certificate will expire soon', when: '25m ago', tone: 'info' },
  ],
  tasks: [
    { title: 'Install nginx', progress: 100, state: 'Done', detail: 'Completed and verified' },
    { title: 'Configure postgres', progress: 72, state: 'Running', detail: 'Creating database roles' },
    { title: 'Prepare redis', progress: 38, state: 'Queued', detail: 'Waiting for package source' },
  ],
  recommendations: [
    { title: 'Enable SSL renewal', detail: 'One domain is nearing expiry and should be automated.', reason: 'Improves uptime and reduces manual maintenance.' },
    { title: 'Review redis memory limits', detail: 'A warning was detected in the monitoring snapshot.', reason: 'Keeps cache pressure from affecting the stack.' },
    { title: 'Add workspace metadata', detail: 'Workspaces can be richer if owners and environments are registered.', reason: 'Improves filtering and auditability.' },
  ],
  knowledge: loginFallback.knowledgeArticles,
}

export const settingsFallback: SettingsData = {
  sections: [
    {
      title: 'Platform defaults',
      description: 'Control global behavior for all workspaces and deployments.',
      items: [
        { label: 'Default environment', value: 'production', detail: 'New workspaces inherit this environment unless overridden.' },
        { label: 'Deployment mode', value: 'Rolling', detail: 'Updates roll through the fleet with health checks.' },
        { label: 'Time zone', value: 'Asia/Kolkata', detail: 'Used for audit logs, tasks, and scheduled jobs.' },
      ],
    },
    {
      title: 'Notifications',
      description: 'Route operational updates to the right people.',
      items: [
        { label: 'Email alerts', value: 'Enabled', detail: 'Warnings and failures go to the admin team.' },
        { label: 'Task summaries', value: 'Every 6 hours', detail: 'Digest mode keeps noise manageable.' },
        { label: 'Slack bridge', value: 'Connected', detail: 'Workspace bridge currently points to #devserver-ops.' },
      ],
    },
  ],
  preferences: [
    { label: 'Theme', value: 'Dark / Graphite', detail: 'The product ships in a calm high-contrast palette.' },
    { label: 'Default shell', value: '/bin/bash', detail: 'Used when launching terminal sessions from workspaces.' },
    { label: 'Editor format', value: 'Prettier + ESLint', detail: 'Formatting presets are surfaced in workspace configs.' },
  ],
  security: [
    { label: 'Session TTL', value: '12 hours', detail: 'Admins can reduce this later without changing the UI.' },
    { label: 'MFA enforcement', value: 'Recommended', detail: 'Enable for admin accounts before production use.' },
    { label: 'Audit retention', value: '90 days', detail: 'Audit logs are retained on the platform database.' },
  ],
  aiProviders: [
    { label: 'OpenAI', value: 'Connected', detail: 'Configured and ready' },
    { label: 'Anthropic', value: 'Available', detail: 'Can be enabled from settings' },
    { label: 'Local Model', value: 'Available', detail: 'Local provider stub for offline testing' },
    { label: 'None', value: 'Disabled', detail: 'Skip AI for now' },
  ],
  mcpIntegrations: [
    { label: 'GitHub', value: 'Connected', detail: 'Repository access and pull request metadata' },
    { label: 'Slack', value: 'Connected', detail: 'Notifications and incident routing' },
    { label: 'Linear', value: 'Available', detail: 'Issue tracking integration stub' },
    { label: 'Grafana', value: 'Available', detail: 'Observability integration stub' },
  ],
}

export function projectFallbacks(): ProjectListItem[] {
  return [
    { slug: 'aurora', name: 'Aurora', description: 'Primary workspace for the flagship application.', owner: 'Ava Patel', repository: 'git@github.com:devserver/aurora.git', environment: 'production', status: 'healthy', updatedAt: '5 minutes ago', progress: 88 },
    { slug: 'docs-portal', name: 'Docs Portal', description: 'Documentation and knowledge workspace.', owner: 'Milo Chen', repository: 'git@github.com:devserver/docs-portal.git', environment: 'staging', status: 'warning', updatedAt: '21 minutes ago', progress: 63 },
    { slug: 'internal-api', name: 'Internal API', description: 'Internal APIs powering auth, audit logs, and billing hooks.', owner: 'Priya Singh', repository: 'git@github.com:devserver/internal-api.git', environment: 'production', status: 'maintenance', updatedAt: '1 hour ago', progress: 74 },
  ]
}

export function projectDetailFallback(slug: string): ProjectDetailData {
  const project = projectFallbacks().find((item) => item.slug === slug) ?? projectFallbacks()[0]
  return {
    project,
    repository: { url: project.repository, branch: project.environment === 'production' ? 'main' : 'develop', buildCommand: 'pnpm build', startCommand: 'pnpm start' },
    environment: [
      { key: 'APP_ENV', value: project.environment, visibility: 'public' },
      { key: 'LOG_LEVEL', value: 'info', visibility: 'public' },
      { key: 'DATABASE_URL', value: 'postgres://platform', visibility: 'secret' },
      { key: 'REDIS_URL', value: 'redis://cache', visibility: 'secret' },
    ],
    deployments: [
      { id: `${project.slug}-dep-1`, version: 'v2.3.1', status: 'succeeded', environment: project.environment, deployedAt: '2 hours ago', note: 'Smoke tests passed' },
    ],
    domains: [
      { host: `${project.slug}.devserver.local`, ssl: 'Valid', target: 'nginx / root' },
      { host: `api.${project.slug}.devserver.local`, ssl: 'Renewing', target: 'internal api / upstream' },
    ],
    logs: [
      { id: `${project.slug}-log-1`, title: 'Deployment validated', detail: 'The latest rollout completed smoke checks.', time: '4 minutes ago', tone: 'success' },
    ],
    database: { name: `${project.slug}-platform`, engine: 'PostgreSQL 16', status: 'ready' },
    overview: [
      {
        title: 'Ownership',
        detail: 'Workspace metadata and lifecycle status.',
        items: [
          { label: 'Owner', value: project.owner, tone: 'accent' },
          { label: 'Environment', value: project.environment, tone: 'info' },
          { label: 'Status', value: project.status, tone: 'success' },
        ],
      },
      {
        title: 'Build',
        detail: 'Repository and release commands.',
        items: [
          { label: 'Repository', value: project.repository, tone: 'neutral' },
          { label: 'Branch', value: project.environment === 'production' ? 'main' : 'develop', tone: 'accent' },
          { label: 'Deploy target', value: 'Primary fleet', tone: 'success' },
        ],
      },
    ],
    activity: [
      { title: 'Workspace synced', detail: 'Repository state was refreshed from the latest snapshot.', when: '2m ago', tone: 'success' },
      { title: 'Task queued', detail: 'A deployment workflow is waiting for approval.', when: '15m ago', tone: 'warning' },
    ],
    tasks: [
      { title: 'Build workspace index', progress: 100, state: 'Done', detail: 'Metadata and environments are indexed.' },
      { title: 'Validate secrets', progress: 68, state: 'Running', detail: 'Checking secret references and mounts.' },
      { title: 'Prepare terminal session', progress: 28, state: 'Queued', detail: 'Waiting for a workspace slot.' },
    ],
    recommendations: [
      { title: 'Promote preview to staging', detail: 'This workspace has a solid deployment history and could reuse the same release path.', reason: 'Keeps the release train consistent.' },
      { title: 'Review SSL expiry windows', detail: 'One domain is renewing and should be validated before traffic shifts.', reason: 'Protects the project edge surface.' },
    ],
    knowledge: loginFallback.knowledgeArticles,
  }
}

export const workspaceDoctorFallback = {
  checks: [
    { label: 'Repository', detail: 'Connected and fetchable.', status: 'passed' as const },
    { label: 'Environment', detail: 'One secret is missing from the snapshot.', status: 'warning' as const },
    { label: 'Deployment target', detail: 'Healthy and ready to accept rollout.', status: 'passed' as const },
  ],
  recommendations: [
    'Review environment variables.',
    'Validate SSL renewal coverage.',
    'Inspect queued deployment tasks.',
  ],
}

export function knowledgeFallback(scope: 'setup' | 'login' | 'dashboard'): KnowledgeArticle[] {
  if (scope === 'setup') return setupFallback.knowledgeArticles
  if (scope === 'dashboard') return [
    { title: 'Projects', question: 'How do workspaces differ from servers?', summary: 'Workspaces bundle repository, environment, deployment, and monitoring state.', link: '#learn-projects' },
    { title: 'Terminal', question: 'Will commands run locally or remotely?', summary: 'The executor layer will eventually support SSH and dry-run execution.', link: '#learn-terminal' },
  ]
  return loginFallback.knowledgeArticles
}

export function devcenterFallback(slug: string): DevCenterSectionData {
  const common = {
    checklist: [
      { label: 'Shell layout', detail: 'A consistent shell wraps product screens.', status: 'passed' as const },
      { label: 'Service boundaries', detail: 'Service facades keep the UI replaceable.', status: 'passed' as const },
      { label: 'Task engine', detail: 'Progress and rollback hooks are represented in the UI.', status: 'warning' as const },
    ],
    knowledge: [
      { title: 'Why architecture matters', question: 'What does a good product foundation buy us?', summary: 'It keeps the next 100 screens from becoming one-off exceptions.', link: '#architecture-foundation' },
    ],
  }

  return {
    key: slug,
    title: slug === 'tasks' ? 'Tasks' : slug === 'knowledge' ? 'Knowledge' : slug === 'dependencies' ? 'Dependencies' : slug === 'api' ? 'API' : slug === 'database' ? 'Database' : slug === 'project-doctor' ? 'Project Doctor' : 'Architecture',
    subtitle: 'Generated workspace metadata and docs.',
    description: 'The DevCenter is backed by generated project metadata and workspace context.',
    bullets: ['Generated from `.devserver/` context.', 'Keeps docs and code in sync.', 'Provides contextual learning for the workspace.'],
    examples: ['README', 'architecture.md', 'routes.json', 'database.json'],
    checklist: common.checklist,
    knowledge: common.knowledge,
    codeSample: 'Generated context will appear here.',
    metadata: {},
  }
}

export const bootstrapFallback: BootstrapDecision = {
  destination: 'setup',
  setupRequired: true,
  authenticated: false,
  installationProgress: 0,
  message: 'Starting from a clean installation.',
}
