package keys

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// managerKeys collects the manager's status messages (string literals
// assigned to a Status field) and its error hints (the text field of
// the hints table).
func managerKeys(dir, root string, add func(text, where string, plural bool)) error {
	fset := token.NewFileSet()
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return err
	}
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.AssignStmt:
				for i, lhs := range n.Lhs {
					if sel, ok := lhs.(*ast.SelectorExpr); ok && sel.Sel.Name == "Status" && i < len(n.Rhs) {
						if s, ok := literal(n.Rhs[i]); ok {
							add(s, rel, false)
						}
					}
				}
			case *ast.ValueSpec:
				for _, name := range n.Names {
					if name.Name != "hints" || len(n.Values) == 0 {
						continue
					}
					table, ok := n.Values[0].(*ast.CompositeLit)
					if !ok {
						continue
					}
					for _, el := range table.Elts {
						if row, ok := el.(*ast.CompositeLit); ok && len(row.Elts) > 1 {
							if s, ok := literal(row.Elts[len(row.Elts)-1]); ok {
								add(s, rel, false)
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil
}
