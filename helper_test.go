package astquery_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/gostaticanalysis/astquery"
	"golang.org/x/tools/txtar"
)

func nodeType(t *testing.T, n ast.Node) string {
	return strings.TrimPrefix(fmt.Sprintf("%T", n), "*ast.")
}

func nodesType(t *testing.T, ns []ast.Node) []string {
	t.Helper()
	if ns == nil {
		return nil
	}
	s := make([]string, len(ns))
	for i := range s {
		s[i] = nodeType(t, ns[i])
	}
	return s
}

func newInspector(t *testing.T, stmtStr string) (ast.Stmt, *astquery.Inspector) {
	t.Helper()
	src := fmt.Sprintf("package a\nfunc f() { %s }", stmtStr)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "a.go", src, 0)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	stmt := f.Decls[0].(*ast.FuncDecl).Body.List[0]
	in := astquery.NewInspector([]*ast.File{f})

	return stmt, in
}

func newSeacher(t *testing.T, path string) *astquery.Seacher {
	t.Helper()
	fset := token.NewFileSet()
	files := parse(t, fset, path)
	return astquery.NewSeacher(fset, files, nil)
}

func parse(t *testing.T, fset *token.FileSet, path string) []*ast.File {
	t.Helper()
	ar, err := txtar.ParseFile(path)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	files := make([]*ast.File, len(ar.Files))
	for i := range ar.Files {
		n, d := ar.Files[i].Name, ar.Files[i].Data
		f, err := parser.ParseFile(fset, n, d, 0)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		files[i] = f
	}

	return files
}
