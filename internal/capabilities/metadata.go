package capabilities

// Metadata describes a provider's capability entry.
type Metadata struct {
	Provider    string       `json:"provider"`
	Capability  Capability   `json:"capability"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	Category    string       `json:"category,omitempty"`
	Version     string       `json:"version,omitempty"`
	Priority    int          `json:"priority"`
	Tags        []string     `json:"tags,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

// Binding is the resolved capability-provider pair.
type Binding struct {
	Provider Provider `json:"-"`
	Metadata Metadata `json:"metadata"`
}
