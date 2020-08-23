package astquery

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"

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

type attr struct {
	parent    ast.Node
	name, val string
}

func (a attr) Pos() token.Pos {
	return a.parent.Pos()
}

func (a attr) End() token.Pos {
	return a.parent.End()
}

// NodeNavigator implements xpath.NodeNavigator.
type NodeNavigator struct {
	in       *Inspector
	fset     *token.FileSet
	root     *pkg
	node     ast.Node
	siblings []ast.Node
	index    int
	attr     int
	attrs    []attr
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
		attr: -1,
	}
}

func (n *NodeNavigator) Node() ast.Node {
	if n.attr != -1 {
		return n.attrs[n.attr]
	}
	return n.node
}

func (n *NodeNavigator) NodeType() xpath.NodeType {
	switch n.node.(type) {
	case *pkg:
		return xpath.RootNode
	}

	if n.attr == -1 {
		return xpath.ElementNode
	} else {
		return xpath.AttributeNode
	}
}

func (n *NodeNavigator) LocalName() string {
	switch node := n.node.(type) {
	case *pkg:
		return ""
	case *ast.File:
		f := n.fset.File(node.Pos())
		return filepath.Base(f.Name())
	}

	if n.attr != -1 {
		return n.attrs[n.attr].name
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

	if n.attr != -1 {
		return n.attrs[n.attr].val
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
		attr:  n.attr,
		attrs: n.attrs,
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
	n.attr = -1
	debugln("#")
}

func (n *NodeNavigator) MoveToParent() bool {
	if n.attr != -1 {
		n.attr = -1
		return true
	}

	switch n.node.(type) {
	case *pkg:
		return false
	}

	parent := n.in.Parent(n.node)
	if parent != nil {
		debugf("^%T(from %T)>", parent, n.node)
		n.node = parent
		switch n.node.(type) {
		case *ast.File:
			n.siblings = n.root.children()
		default:
			n.siblings = n.in.Children(n.in.Parent(n.node))
		}
		n.index = 0
		for i := range n.siblings {
			if n.siblings[i] == n.node {
				n.index = i
				break
			}
		}
		debugf("%T %d %v\n", n.node, n.index, nodesToStr(n.siblings))
		return true
	}

	n.MoveToRoot()
	return true
}

func (n *NodeNavigator) MoveToNextAttribute() bool {
	if n.attr == -1 {
		n.attrs = attributes(n.fset, n.node)
	}

	if n.attr >= len(n.attrs)-1 {
		return false
	}

	n.attr++
	return true
}

func (n *NodeNavigator) MoveToChild() bool {
	if n.attr != -1 {
		return false
	}

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
	debugf("%T[%v]>", n.node, nodesToStr(children))
	if len(children) == 0 {
		debugln("/")
		return false
	}
	n.siblings = children
	n.index = 0
	n.node = n.siblings[0]
	debugf("v%T(%s)[%d]>", n.node, nodesToStr(n.siblings), n.index)

	return true
}

func (n *NodeNavigator) MoveToFirst() bool {
	if n.attr != -1 || len(n.siblings) == 0 {
		return false
	}

	n.index = 0
	n.node = n.siblings[0]
	debugf("<<%T>", n.node)
	return true
}

func (n *NodeNavigator) MoveToNext() bool {
	if n.attr != -1 || len(n.siblings)-1 <= n.index {
		return false
	}
	n.index++
	n.node = n.siblings[n.index]
	debugf("+%T(%s)[%d]>", n.node, nodesToStr(n.siblings), n.index)
	return true
}

func (n *NodeNavigator) MoveToPrevious() bool {
	if n.attr != -1 || n.siblings == nil || n.index <= 0 {
		return false
	}
	n.index--
	n.node = n.siblings[n.index]
	debugf("-%T(%s)[%d]>", n.node, nodesToStr(n.siblings), n.index)
	return true
}

func (n *NodeNavigator) MoveTo(to xpath.NodeNavigator) bool {
	debugln("@")
	_to, _ := to.(*NodeNavigator)
	if _to == nil || n.in != _to.in {
		return false
	}
	n.node = _to.node
	n.siblings = make([]ast.Node, len(_to.siblings))
	copy(n.siblings, _to.siblings)
	n.index = _to.index
	n.attr = _to.attr
	n.attrs = make([]attr, len(_to.attrs))
	copy(n.attrs, _to.attrs)
	return true
}

func attributes(fset *token.FileSet, n ast.Node) []attr {
	switch n.(type) {
	case *pkg:
		return nil
	}

	rv := reflect.ValueOf(n)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	attrs := []attr{
		{
			parent: n,
			name:   "type",
			val:    strings.TrimPrefix(rv.Type().String(), "ast."),
		},
		{
			parent: n,
			name:   "pos",
			val:    fset.Position(n.Pos()).String(),
		},
	}

	if rv.Kind() == reflect.Struct {
		for i := 0; i < rv.NumField(); i++ {
			f := rv.Field(i)
			switch f.Kind() {
			case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String, reflect.UnsafePointer:
				attrs = append(attrs, attr{
					parent: n,
					name: rv.Type().Field(i).Name,
					val:  fmt.Sprintf("%v", rv.Field(i).Interface()),
				})
			}
		}
	}

	return attrs
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
