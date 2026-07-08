# RBAC Permission Matrix

This matrix documents the current RBAC permissions, the resources they protect, the roles that should receive them, and the API/UI surface they guard.

| Resource | Permission | Role(s) | HTTP Method | API Path | Frontend UI | Required Permission |
| --- | --- | --- | --- | --- | --- | --- |
| Organization listing | organizations.view | Owner, Admin, Developer, Operator, Viewer | GET | /api/v1/organizations | Organization switcher | organizations.view |
| Organization creation | organizations.create | Owner, Admin | POST | /api/v1/organizations | Create org workflows | organizations.create |
| Role listing | roles.view | Owner, Admin, Developer, Operator, Viewer | GET | /api/v1/roles | Role management panels | roles.view |
| Member listing | members.view | Owner, Admin, Developer, Operator, Viewer | GET | /api/v1/organizations/members | Members page | members.view |
| Member invitation | members.invite | Owner, Admin | POST | /api/v1/organizations/members | Invite member button | members.invite |
| Member role updates | members.manage | Owner, Admin | PUT/PATCH | /api/v1/organizations/members | Role selector, member actions | members.manage |
| Member removals | members.remove | Owner, Admin | DELETE | /api/v1/organizations/members | Remove member actions | members.remove |
| Transfer ownership | org.owner | Owner | POST | /api/v1/organizations/transfer-ownership | Transfer owner control | org.owner |
| App settings read | settings.view | Owner, Admin, Developer, Operator, Viewer | GET | /api/v1/settings | Settings page read access | settings.view |
| App settings write | settings.write | Owner, Admin | POST/PATCH | /api/v1/settings | Save settings | settings.write |
| Environments read | environments.view | Owner, Admin, Developer, Operator, Viewer | GET | /api/v1/environments | Environments list | environments.view |
| Environments write | environments.write | Owner, Admin, Developer, Operator | POST/PATCH | /api/v1/environments | Create environment | environments.write |
| Secrets read | secrets.view | Owner, Admin, Developer, Operator, Viewer | GET | /api/v1/secrets | Secret list | secrets.view |
| Secrets write | secrets.write | Owner, Admin, Developer, Operator | POST/PATCH | /api/v1/secrets | Add secret | secrets.write |
| Variables read | variables.view | Owner, Admin, Developer, Operator, Viewer | GET | /api/v1/variables | Variable list | variables.view |
| Variables write | variables.write | Owner, Admin, Developer, Operator | POST/PATCH | /api/v1/variables | Add variable | variables.write |
| Providers read | providers.view | Owner, Admin, Developer, Operator, Viewer | GET | /api/v1/providers | Provider list | providers.view |
| Providers write | providers.write | Owner, Admin, Developer, Operator | POST | /api/v1/providers, /api/v1/providers/test | Configure provider | providers.write |
| Superadmin org list | org.admin | Super Admin | GET | /api/v1/superadmin/organizations | Super admin dashboard | org.admin |
