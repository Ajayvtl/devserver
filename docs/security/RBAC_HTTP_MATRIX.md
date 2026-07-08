# RBAC HTTP Matrix

This document maps protected endpoints to required permissions and expected status codes.

| Method | Endpoint | Required Permission | 200 | 401 | 403 | 404 | 409 | 422 | 500 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| GET | /api/v1/organizations | organizations.view | yes | yes | yes | no | no | no | yes |
| POST | /api/v1/organizations | organizations.create | yes | yes | yes | no | no | yes | yes |
| GET | /api/v1/organizations/members | members.view | yes | yes | yes | no | no | no | yes |
| POST | /api/v1/organizations/members | members.invite | yes | yes | yes | no | yes | yes | yes |
| PUT | /api/v1/organizations/members | members.manage | yes | yes | yes | no | yes | yes | yes |
| PATCH | /api/v1/organizations/members | members.manage | yes | yes | yes | no | yes | yes | yes |
| DELETE | /api/v1/organizations/members | members.remove | yes | yes | yes | no | yes | yes | yes |
| POST | /api/v1/organizations/transfer-ownership | org.owner | yes | yes | yes | no | yes | yes | yes |
| GET | /api/v1/roles | roles.view | yes | yes | yes | no | no | no | yes |
| GET | /api/v1/settings | settings.view | yes | yes | yes | no | no | yes | yes |
| POST | /api/v1/settings | settings.write | yes | yes | yes | no | no | yes | yes |
| GET | /api/v1/environments | environments.view | yes | yes | yes | no | no | yes | yes |
| POST | /api/v1/environments | environments.write | yes | yes | yes | no | no | yes | yes |
| GET | /api/v1/secrets | secrets.view | yes | yes | yes | no | no | yes | yes |
| POST | /api/v1/secrets | secrets.write | yes | yes | yes | no | no | yes | yes |
| GET | /api/v1/variables | variables.view | yes | yes | yes | no | no | yes | yes |
| POST | /api/v1/variables | variables.write | yes | yes | yes | no | no | yes | yes |
| GET | /api/v1/providers | providers.view | yes | yes | yes | no | no | yes | yes |
| POST | /api/v1/providers | providers.write | yes | yes | yes | no | no | yes | yes |
| POST | /api/v1/providers/test | providers.write | yes | yes | yes | no | no | yes | yes |
| GET | /api/v1/superadmin/organizations | org.admin | yes | yes | yes | no | no | no | yes |
