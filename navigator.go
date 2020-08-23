package astquery

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"github.com/antchfx/xpath"
	"golang.org/x/tools/go/ast/inspector"
)

type pkg struct {
	files []*ast.File
}

var _ ast.Node = (*pkg)(nil)

func (p *pkg) children() []ast.Node {
	ns := make([]ast.Node, len(p.files))
	for i := range p.files {
		ns[i] = p.files[i]
	}
	return ns
}

func (p *pkg) Pos() token.Pos {
	if len(p.files) == 0 {
		return token.NoPos
	}
	return p.files[0].Pos()
}

func (p *pkg) End() token.Pos {
	if len(p.files) == 0 {
		return token.NoPos
	}
	return p.files[len(p.files)-1].End()
}

type NodeNavigator struct {
	in       *Inspector
	fset     *token.FileSet
	root     *pkg
	node     ast.Node
	siblings []ast.Node
	index    int
}

var _ xpath.NodeNavigator = (*NodeNavigator)(nil)

// NewNodeNavigator creates a NodeNavigator.
// If in is not nil NodeNavigator use the given inspector.
func NewNodeNavigator(fset *token.FileSet, files []*ast.File, in *inspector.Inspector) *NodeNavigator {

	if in == nil {
		in = inspector.New(files)
	}

	root := &pkg{files: files}
	return &NodeNavigator{
		in:   &Inspector{in},
		fset: fset,
		node: root,
		root: root,
	}
}

func (n *NodeNavigator) Node() ast.Node {
	return n.node
}

func (n *NodeNavigator) NodeType() xpath.NodeType {
	switch n.node.(type) {
	case *pkg:
		return xpath.RootNode
	}
	return xpath.ElementNode
}

func (n *NodeNavigator) LocalName() string {
	switch node := n.node.(type) {
	case *pkg:
		return ""
	case *ast.File:
		f := n.fset.File(node.Pos())
		return filepath.Base(f.Name())
	}

	return n.in.Name(n.node)
}

func (n *NodeNavigator) Prefix() string {
	return ""
}

func (n *NodeNavigator) Value() string {
	switch node := n.node.(type) {
	case *pkg:
		return ""
	case *ast.File:
		f := n.fset.File(node.Pos())
		return filepath.Base(f.Name())
	}

	return fmt.Sprintf("%v", n.node)
}

func (n *NodeNavigator) Copy() xpath.NodeNavigator {
	copied := &NodeNavigator{
		in:    n.in,
		fset:  n.fset,
		root:  n.root,
		node:  n.node,
		index: n.index,
	}
	if n.siblings != nil {
		copied.siblings = make([]ast.Node, len(n.siblings))
		copy(copied.siblings, n.siblings)
	}
	return copied
}

func (n *NodeNavigator) MoveToRoot() {
	n.node = n.root
	n.index = 0
	n.siblings = nil
}

func (n *NodeNavigator) MoveToParent() bool {
	switch n.node.(type) {
	case *pkg:
		return false
	}

	parent := n.in.Parent(n.node)
	if parent != nil {
		n.node = parent
		n.index = 0
		n.siblings = nil
		return true
	}

	return false
}

func (n *NodeNavigator) MoveToNextAttribute() bool {
	return false
}

func (n *NodeNavigator) MoveToChild() bool {
	switch node := n.node.(type) {
	case *pkg:
		if len(node.files) == 0 {
			return false
		}
		n.siblings = node.children()
		n.index = 0
		n.node = n.siblings[0]
		return true
	}

	children := n.in.Children(n.node)
	if len(children) == 0 {
		return false
	}
	n.siblings = children
	n.index = 0
	n.node = n.siblings[0]

	return true
}

func (n *NodeNavigator) MoveToFirst() bool {
	if len(n.siblings) == 0 {
		return false
	}
	n.index = 0
	n.node = n.siblings[0]
	return true
}

func (n *NodeNavigator) MoveToNext() bool {
	if len(n.siblings)-1 <= n.index {
		return false
	}
	n.index++
	n.node = n.siblings[n.index]
	return true
}

func (n *NodeNavigator) MoveToPrevious() bool {
	if n.siblings == nil || n.index <= 0 {
		return false
	}
	n.index--
	n.node = n.siblings[n.index]
	return true
}

func (n *NodeNavigator) MoveTo(to xpath.NodeNavigator) bool {
	_to, _ := to.(*NodeNavigator)
	if _to == nil || n.in != _to.in {
		return false
	}
	n.node = _to.node
	n.siblings = _to.siblings
	n.index = _to.index
	return true
}

func nodes(iter *xpath.NodeIterator) []ast.Node {
	var ns []ast.Node
	for iter.MoveNext() {
		current, _ := iter.Current().(*NodeNavigator)
		if current != nil && current.Node() != nil {
			switch n := current.Node().(type) {
			case *pkg:
				for _, f := range n.files {
					ns = append(ns, f)
				}
			default:
				ns = append(ns, n)
			}
		}
	}
	return ns
}
