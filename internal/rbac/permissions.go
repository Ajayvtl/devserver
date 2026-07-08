package rbac

const (
	PermissionWildcard = "*"

	PermissionMembersView   = "members.view"
	PermissionMembersInvite = "members.invite"
	PermissionMembersManage = "members.manage"
	PermissionMembersRemove = "members.remove"

	PermissionOrgOwner            = "org.owner"
	PermissionOrgAdmin            = "org.admin"
	PermissionOrganizationsView   = "organizations.view"
	PermissionOrganizationsCreate = "organizations.create"
	PermissionRolesView           = "roles.view"
	PermissionUsersView           = "users.view"

	PermissionSettingsView      = "settings.view"
	PermissionSettingsWrite     = "settings.write"
	PermissionEnvironmentsView  = "environments.view"
	PermissionEnvironmentsWrite = "environments.write"
	PermissionSecretsView       = "secrets.view"
	PermissionSecretsWrite      = "secrets.write"
	PermissionVariablesView     = "variables.view"
	PermissionVariablesWrite    = "variables.write"
	PermissionProvidersView     = "providers.view"
	PermissionProvidersWrite    = "providers.write"
)
