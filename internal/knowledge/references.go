package knowledge

import "go/ast"

type Reference struct {
	Symbol string   `json:"symbol"`
	Calls  []string `json:"calls,omitempty"`
	UsedBy []string `json:"usedBy,omitempty"`
}

type ReferenceIndex struct {
	References []Reference `json:"references"`
}

func buildReferences(file *ast.File) ReferenceIndex {

	index := ReferenceIndex{}

	ast.Inspect(file, func(n ast.Node) bool {

		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		ref := Reference{
			Symbol: fn.Name.Name,
		}

		if fn.Body != nil {

			ast.Inspect(fn.Body, func(n ast.Node) bool {

				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				switch f := call.Fun.(type) {

				case *ast.Ident:
					ref.Calls = append(ref.Calls, f.Name)

				case *ast.SelectorExpr:
					ref.Calls = append(ref.Calls, f.Sel.Name)
				}

				return true
			})
		}

		index.References = append(index.References, ref)

		return false
	})

	return index
}
