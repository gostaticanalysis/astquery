package astquery

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/ast/inspector"
)

// Inspector is wrapper of *golang.org/x/go/ast/inspector.Inspector.
type Inspector struct {
	*inspector.Inspector
}

// NewInspector creates an Inspector.
func NewInspector(files []*ast.File) *Inspector {
	return &Inspector{inspector.New(files)}
}

// Children returns children of the given node.
func (in *Inspector) Children(n ast.Node) []ast.Node {
	if n == nil {
		return nil
	}

	var children []ast.Node
	ast.Inspect(n, func(_n ast.Node) bool {
		if _n != nil && n != _n {
			children = append(children, _n)
			return false
		}
		return true
	})
	if len(children) != 0 {
		return children
	}
	return children
}

// Stack returns a stack which contains between the root node and the given node.
// The stack's first element is the root node and last element is the given node.
func (in *Inspector) Stack(n ast.Node) []ast.Node {
	var stack []ast.Node
	in.WithStack(nil, func(_n ast.Node, push bool, _stack []ast.Node) bool {
		if n == _n {
			stack = _stack
			return false
		}
		return true
	})
	return stack
}

// Path returns a path to the given node from the root node.
// The path's first element is the given node and last element is the root node.
func (in *Inspector) Path(n ast.Node) []ast.Node {
	stack := in.Stack(n)
	if stack == nil {
		return nil
	}

	path := make([]ast.Node, len(stack))
	for i := range path {
		path[i] = stack[len(stack)-1-i]
	}

	return path
}

// Parent returns a parent node of the given node.
func (in *Inspector) Parent(n ast.Node) ast.Node {
	path := in.Path(n)
	if len(path) <= 1 {
		return nil
	}
	return path[1]
}

// Name returns parent's field name.
func (in *Inspector) Name(n ast.Node) string {
	var name string
	parent := in.Parent(n)
	astutil.Apply(parent, func(cur *astutil.Cursor) bool {
		if n == cur.Node() {
			name = cur.Name()
			return false
		}
		return true
	}, nil)

	return name
}

// Index returns parent's field index.
func (in *Inspector) Index(n ast.Node) int {
	var index int
	parent := in.Parent(n)
	astutil.Apply(parent, func(cur *astutil.Cursor) bool {
		if n == cur.Node() {
			index = cur.Index()
			return false
		}
		return true
	}, nil)

	return index
}
