package astquery_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEvaluator_Select(t *testing.T) {
	t.Parallel()
	//astquery.DebugON(t)

	S := func(s ...string) []string { return s }
	TD := func(f string) string { return filepath.Join("testdata", "TestEvaluator_Select", f) }
	cases := map[string]struct {
		path  string
		xpath string
		want  []string
	}{
		"single":   {TD("single.go"), "/*/Decls[1]/Body/*", S("ReturnStmt")},
		"multi":    {TD("multi.go"), "/*/Decls[1]/Body/*", S("AssignStmt", "ReturnStmt")},
		"filename": {TD("single.go"), "/a.go/Decls[1]/Body/*", S("ReturnStmt")},
		"attr":     {TD("attr.go"), "//*[@type='CallExpr']/Fun[@type='Ident' and @Name='print']", S("Ident", "Ident", "Ident")},
	}

	for n, tt := range cases {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			e := newEvaluator(t, tt.path)
			ns, err := e.Select(tt.xpath)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}
			got := nodesType(t, ns)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEvaluator_Eval(t *testing.T) {
	t.Parallel()
	//astquery.DebugON(t)

	TD := func(f string) string { return filepath.Join("testdata", "TestEvaluator_Eval", f) }
	cases := map[string]struct {
		path  string
		xpath string
		want  interface{}
	}{
		"attr": {TD("attr.go"), "//*[@type='CallExpr']/Fun[@type='Ident']/@Name", []interface{}{"print", "print", "println", "print"}},
		"src": {TD("attr.go"), "//*[@src='print']/@Name", []interface{}{"print", "print", "print"}},
	}

	for n, tt := range cases {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			e := newEvaluator(t, tt.path)
			got, err := e.Eval(tt.xpath)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
