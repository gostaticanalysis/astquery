package astquery_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInspector_Children(t *testing.T) {
	t.Parallel()
	S := func(s ...string) []string {
		return s
	}
	cases := map[string]struct {
		stmt string
		want []string
	}{
		"single":     {"go print()", S("CallExpr")},
		"multi":      {"_ = 1+1", S("Ident", "BinaryExpr")},
		"nochildren": {"return", nil},
	}

	for n, tt := range cases {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			stmt, in := newInspector(t, tt.stmt)
			children := in.Children(stmt)
			got := nodesType(t, children)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestInspector_Stack(t *testing.T) {
	t.Parallel()
	S := func(s ...string) []string {
		return s
	}
	cases := map[string]struct {
		stmt string
		want []string
	}{
		"single": {"return", S("File", "FuncDecl", "BlockStmt", "ReturnStmt")},
	}

	for n, tt := range cases {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			stmt, in := newInspector(t, tt.stmt)
			stack := in.Stack(stmt)
			got := nodesType(t, stack)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestInspector_Path(t *testing.T) {
	t.Parallel()
	S := func(s ...string) []string {
		return s
	}
	cases := map[string]struct {
		stmt string
		want []string
	}{
		"single": {"return", S("ReturnStmt", "BlockStmt", "FuncDecl", "File")},
	}

	for n, tt := range cases {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			stmt, in := newInspector(t, tt.stmt)
			path := in.Path(stmt)
			got := nodesType(t, path)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestInspector_Parent(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		stmt string
		want string
	}{
		"single": {"return", "BlockStmt"},
	}

	for n, tt := range cases {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			stmt, in := newInspector(t, tt.stmt)
			parent := in.Parent(stmt)
			got := nodeType(t, parent)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

/*
func TestInspector_Root(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		stmt string
		want string
	}{
		"single": {"return", "File"},
	}

	for n, tt := range cases {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			stmt, in := newInspector(t, tt.stmt)
			root := in.Root(stmt)
			got := nodeType(t, root)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
*/
