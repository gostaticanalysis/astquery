package astquery_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEvaluator_Select(t *testing.T) {
	t.Parallel()

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
