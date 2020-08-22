package astquery

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/antchfx/xpath"
	"golang.org/x/tools/go/ast/inspector"
)

type pkg struct {
	files []*fileHolder
	syntax []*ast.File
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

type fileHolder struct {
	name   string
	syntax *ast.File
}

var _ ast.Node = (*fileHolder)(nil)

func (fh *fileHolder) Pos() token.Pos {
	return fh.syntax.Pos()
}

func (fh *fileHolder) End() token.Pos {
	return fh.syntax.End()
}

type NodeNavigator struct {
	in       *Inspector
	root     *pkg
	node     ast.Node
	siblings []ast.Node
	index    int
}

var _ xpath.NodeNavigator = (*NodeNavigator)(nil)

// NewNodeNavigator creates a NodeNavigator.
// If in is not nil NodeNavigator use the given inspector.
func NewNodeNavigator(fset *token.FileSet, files []*ast.File, in *inspector.Inspector) *NodeNavigator {

	holders := make([]*fileHolder, len(files))
	for i := range files {
		f := fset.File(files[i].Pos())
		holders[i] = &fileHolder{
			name:   filepath.Base(f.Name()),
			syntax: files[i],
		}
	}

	if in == nil {
		in = inspector.New(files)
	}

	root := &pkg{
		files: holders,
		syntax: files,
	}

	return &NodeNavigator{
		in:   &Inspector{in},
		node: root,
		root: root,
	}
}

func (n *NodeNavigator) isRoot() bool {
	return n.node == nil
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
	case *fileHolder:
		return node.name
	}

	return strings.TrimPrefix(fmt.Sprintf("%T", n.node), "*ast.")
}

func (n *NodeNavigator) Prefix() string {
	return ""
}

func (n *NodeNavigator) Value() string {
	switch node := n.node.(type) {
	case *pkg:
		return ""
	case *fileHolder:
		return node.name
	}

	return fmt.Sprintf("%v", n.node)
}

func (n *NodeNavigator) Copy() xpath.NodeNavigator {
	copied := &NodeNavigator{
		in:    n.in,
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
	case *fileHolder:
		n.MoveToRoot()
		return true
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
	case *fileHolder:
		n.siblings = nil
		n.index = 0
		n.node = node.syntax
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
