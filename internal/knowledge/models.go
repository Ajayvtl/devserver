package knowledge

type SymbolKind string

const (
	SymbolPackage   SymbolKind = "package"
	SymbolStruct    SymbolKind = "struct"
	SymbolInterface SymbolKind = "interface"
	SymbolFunction  SymbolKind = "function"
	SymbolMethod    SymbolKind = "method"
	SymbolType      SymbolKind = "type"
	SymbolConst     SymbolKind = "const"
	SymbolVar       SymbolKind = "var"
)

type Position struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

type Symbol struct {
	Name        string     `json:"name"`
	Kind        SymbolKind `json:"kind"`

	Package     string     `json:"package"`
	Receiver    string     `json:"receiver,omitempty"`

	Signature   string     `json:"signature,omitempty"`
	Description string     `json:"description,omitempty"`

	Exported    bool       `json:"exported"`

	Documentation string   `json:"documentation,omitempty"`

	Position Position `json:"position"`

	EndLine int `json:"endLine,omitempty"`

	References []string `json:"references,omitempty"`
}

type WorkspaceKnowledge struct {
	Symbols    []Symbol      `json:"symbols"`
	References ReferenceIndex `json:"references"`

	// In-memory indexes for constant-time lookups
	SymbolIndex  map[string][]Symbol `json:"-"`
	PackageIndex map[string][]Symbol `json:"-"`
	FileIndex    map[string][]Symbol `json:"-"`
}

func (wk *WorkspaceKnowledge) BuildIndex() {
	wk.SymbolIndex = make(map[string][]Symbol)
	wk.PackageIndex = make(map[string][]Symbol)
	wk.FileIndex = make(map[string][]Symbol)

	// Build direct symbol indexes
	for _, sym := range wk.Symbols {
		wk.SymbolIndex[sym.Name] = append(wk.SymbolIndex[sym.Name], sym)
		if sym.Package != "" {
			wk.PackageIndex[sym.Package] = append(wk.PackageIndex[sym.Package], sym)
		}
		if sym.Position.File != "" {
			wk.FileIndex[sym.Position.File] = append(wk.FileIndex[sym.Position.File], sym)
		}
	}

	// Build relationship graph (UsedBy / Calls)
	usedByMap := make(map[string][]string)
	for _, ref := range wk.References.References {
		for _, call := range ref.Calls {
			usedByMap[call] = append(usedByMap[call], ref.Symbol)
		}
	}

	// Update UsedBy and References back into the graph
	for i, ref := range wk.References.References {
		if users, ok := usedByMap[ref.Symbol]; ok {
			wk.References.References[i].UsedBy = users
		}
	}
}
