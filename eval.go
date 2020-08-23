package astquery

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/antchfx/xpath"
	"golang.org/x/tools/go/ast/inspector"
)

// Evaluator evals and selects AST's nodes by XPath.
type Evaluator struct {
	n *NodeNavigator
}

// New creates an Evaluator.
// If the given inspector is not nil macher use it.
func New(fset *token.FileSet, files []*ast.File, in *inspector.Inspector) *Evaluator {
	return &Evaluator{n: NewNodeNavigator(fset, files, in)}
}

// Eval returns the result of the expression.
// The result type of the expression is one of the follow: bool,float64,string,[]ast.Node.
func (e *Evaluator) Eval(expr string) (interface{}, error) {
	n := e.n.Copy()
	_expr, err := xpath.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("expr cannot compile: %w", err)
	}

	v := _expr.Evaluate(n)
	switch v := v.(type) {
	case *xpath.NodeIterator:
		ns := nodes(v)
		vs := make([]interface{}, 0, len(ns))
		for i := range ns {
			switch n := ns[i].(type) {
			case attr:
				vs = append(vs, n.val)
			}
		}
		if len(vs) == len(ns) {
			return vs, nil
		}
		return ns, nil
	}

	return v, nil
}

// Select selects a node set which match the XPath expr.
func (e *Evaluator) Select(expr string) ([]ast.Node, error) {
	n := e.n.Copy()
	_expr, err := xpath.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("expr cannot compile: %w", err)
	}

	return nodes(_expr.Select(n)), nil
}

// SelectOne selects a node set which match the XPath expr and return the first node.
func (e *Evaluator) SelectOne(expr string) (ast.Node, error) {
	ns, err := e.Select(expr)
	if err != nil {
		return nil, err
	}
	if len(ns) == 0 {
		return nil, nil
	}
	return ns[0], nil
}
