package knowledge

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type SymbolScanner struct{}

func NewSymbolScanner() *SymbolScanner {
	return &SymbolScanner{}
}
func exported(name string) bool {
	if name == "" {
		return false
	}

	r := []rune(name)[0]
	return unicode.IsUpper(r)
}

func signature(fset *token.FileSet, node ast.Node) string {

	var buf bytes.Buffer

	if err := format.Node(&buf, fset, node); err != nil {
		return ""
	}

	return strings.TrimSpace(buf.String())
}

func comments(cg *ast.CommentGroup) string {

	if cg == nil {
		return ""
	}

	return strings.TrimSpace(cg.Text())
}
func (s *SymbolScanner) Scan(root string) (*WorkspaceKnowledge, error) {

	result := &WorkspaceKnowledge{}

	fset := token.NewFileSet()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if d.IsDir() {
			switch d.Name() {
			case ".git", ".tmp", "node_modules", "vendor":
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".go" {
			if ext == ".ts" || ext == ".tsx" || ext == ".js" || ext == ".py" {
				scanNonGoFile(path, result)
			}
			return nil
		}

		file, err := parser.ParseFile(
			fset,
			path,
			nil,
			parser.ParseComments,
		)

		if err != nil {
			return nil
		}

		pkg := file.Name.Name
		refs := buildReferences(file)
		result.References.References = append(
			result.References.References,
			refs.References...,
		)
		result.Symbols = append(result.Symbols, Symbol{
			Name:        pkg,
			Kind:        SymbolPackage,
			Package:     pkg,
			Exported:    true,
			Description: "Go package",
			Position: Position{
				File:   path,
				Line:   fset.Position(file.Package).Line,
				Column: fset.Position(file.Package).Column,
			},
		})

		for _, decl := range file.Decls {

			switch node := decl.(type) {

			case *ast.FuncDecl:

				kind := SymbolFunction
				receiver := ""

				if node.Recv != nil {
					kind = SymbolMethod

					if len(node.Recv.List) > 0 {
						switch t := node.Recv.List[0].Type.(type) {
						case *ast.Ident:
							receiver = t.Name
						case *ast.StarExpr:
							if id, ok := t.X.(*ast.Ident); ok {
								receiver = id.Name
							}
						}
					}
				}

				pos := fset.Position(node.Pos())

				result.Symbols = append(result.Symbols, Symbol{
					Name:          node.Name.Name,
					Kind:          kind,
					Package:       pkg,
					Receiver:      receiver,
					Exported:      exported(node.Name.Name),
					Signature:     signature(fset, node),
					Documentation: comments(node.Doc),
					EndLine:       fset.Position(node.End()).Line,
					Position: Position{
						File:   path,
						Line:   pos.Line,
						Column: pos.Column,
					},
				})

			case *ast.GenDecl:

				for _, spec := range node.Specs {

					switch s := spec.(type) {

					case *ast.TypeSpec:

						kind := SymbolType

						switch s.Type.(type) {
						case *ast.StructType:
							kind = SymbolStruct
						case *ast.InterfaceType:
							kind = SymbolInterface
						}

						pos := fset.Position(s.Pos())

						result.Symbols = append(result.Symbols, Symbol{
							Name:          s.Name.Name,
							Kind:          kind,
							Package:       pkg,
							Exported:      exported(s.Name.Name),
							Signature:     signature(fset, s),
							Documentation: comments(node.Doc),
							EndLine:       fset.Position(s.End()).Line,
							Position: Position{
								File:   path,
								Line:   pos.Line,
								Column: pos.Column,
							},
						})

					case *ast.ValueSpec:

						kind := SymbolVar

						if node.Tok == token.CONST {
							kind = SymbolConst
						}

						for _, name := range s.Names {

							pos := fset.Position(name.Pos())

							result.Symbols = append(result.Symbols, Symbol{
								Name:          name.Name,
								Kind:          kind,
								Package:       pkg,
								Exported:      exported(name.Name),
								Signature:     signature(fset, s),
								Documentation: comments(node.Doc),
								EndLine:       fset.Position(s.End()).Line,
								Position: Position{
									File:   path,
									Line:   pos.Line,
									Column: pos.Column,
								},
							})
						}
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func scanNonGoFile(path string, result *WorkspaceKnowledge) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	content := string(data)
	lines := strings.Split(content, "\n")
	ext := filepath.Ext(path)
	pkg := filepath.Base(filepath.Dir(path))

	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if ext == ".ts" || ext == ".tsx" || ext == ".js" || ext == ".jsx" {
			if strings.HasPrefix(trim, "export function ") || strings.HasPrefix(trim, "export class ") || strings.HasPrefix(trim, "export interface ") || strings.HasPrefix(trim, "export type ") || strings.HasPrefix(trim, "export const ") {

				kind := SymbolFunction
				if strings.HasPrefix(trim, "export class") {
					kind = SymbolStruct
				}
				if strings.HasPrefix(trim, "export interface") {
					kind = SymbolInterface
				}
				if strings.HasPrefix(trim, "export type") {
					kind = SymbolType
				}
				if strings.HasPrefix(trim, "export const") {
					kind = SymbolConst
				}

				parts := strings.Fields(trim)
				name := ""
				for j, p := range parts {
					if p == "function" || p == "class" || p == "interface" || p == "type" || p == "const" {
						if len(parts) > j+1 {
							name = parts[j+1]
							idx := strings.IndexAny(name, "(<={:")
							if idx > 0 {
								name = name[:idx]
							}
							break
						}
					}
				}
				if name != "" {
					result.Symbols = append(result.Symbols, Symbol{
						Name:     name,
						Kind:     kind,
						Package:  pkg,
						Exported: true,
						Position: Position{File: path, Line: i + 1, Column: 1},
					})
				}
			}
		} else if ext == ".py" {
			if strings.HasPrefix(trim, "def ") || strings.HasPrefix(trim, "class ") {
				kind := SymbolFunction
				if strings.HasPrefix(trim, "class ") {
					kind = SymbolStruct
				}
				parts := strings.Fields(trim)
				if len(parts) > 1 {
					name := parts[1]
					idx := strings.IndexAny(name, "(:")
					if idx > 0 {
						name = name[:idx]
					}
					result.Symbols = append(result.Symbols, Symbol{
						Name:     name,
						Kind:     kind,
						Package:  pkg,
						Exported: true,
						Position: Position{File: path, Line: i + 1, Column: 1},
					})
				}
			}
		}
	}
}
