package knowledge

type Scanner interface {
	Scan(root string) (*WorkspaceKnowledge, error)
}
