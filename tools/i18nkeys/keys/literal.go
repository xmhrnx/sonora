package keys

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
)

// literal returns a string literal's value, or the value of literals
// joined with +.
func literal(e ast.Expr) (string, bool) {
	switch e := e.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return "", false
		}
		s, err := strconv.Unquote(e.Value)
		return s, err == nil
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return "", false
		}
		a, ok1 := literal(e.X)
		b, ok2 := literal(e.Y)
		return a + b, ok1 && ok2
	case *ast.ParenExpr:
		return literal(e.X)
	}
	return "", false
}

// goKeys collects the literal keys of tr calls in the Go files under dir.
func goKeys(dir, root string, add func(text, where string, plural bool)) error {
	fset := token.NewFileSet()
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := sel.X.(*ast.Ident)
			if !ok || pkg.Name != "tr" {
				return true
			}
			arg := 0
			switch sel.Sel.Name {
			case "T", "N", "Msg":
			case "Bind":
				arg = 1
			default:
				return true
			}
			if len(call.Args) <= arg {
				return true
			}
			if s, ok := literal(call.Args[arg]); ok {
				rel, _ := filepath.Rel(root, fset.Position(call.Pos()).Filename)
				add(s, rel, sel.Sel.Name == "N")
			}
			return true
		})
		return nil
	})
}
