package common

// ListOptions defines standard pagination, filtering, and sorting for repositories.
type ListOptions struct {
	Limit  int
	Offset int
	Filter map[string]interface{}
	Sort   []string // e.g., ["-createdAt", "+name"]
}
