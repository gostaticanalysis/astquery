package astquery

import (
	"go/ast"
	"go/token"

	"github.com/antchfx/xpath"
	"golang.org/x/tools/go/ast/inspector"
)

// Seacher seach ast nodes by xpath.
type Seacher struct {
	n *NodeNavigator
}

// NewSeacher creates a Seacher.
// If in is not nil seacher use the given inspector.
func NewSeacher(fset *token.FileSet, files []*ast.File, in *inspector.Inspector) *Seacher {
	return &Seacher{n: NewNodeNavigator(fset, files, in)}
}

// Query search nodes which match the xpath expr.
func (s *Seacher) Query(expr string) []ast.Node {
	n := s.n.Copy()
	iter := xpath.Select(n, expr)
	var ns []ast.Node
	for iter.MoveNext() {
		current, _ := iter.Current().(*NodeNavigator)
		if current != nil && current.Node() != nil {
			switch n := current.Node().(type) {
			case *pkg:
				for _, f := range n.syntax {
					ns = append(ns, f)
				}
			case *fileHolder:
				ns = append(ns, n.syntax)
			default:
				ns = append(ns, n)
			}
		}
	}
	return ns
}

// QueryOne search nodes which match the xpath expr and return the first node.
func (s *Seacher) QueryOne(expr string) ast.Node {
	ns := s.Query(expr)
	if len(ns) == 0 {
		return nil
	}
	return ns[0]
}
