package astquery_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestQuery(t *testing.T) {
	t.Parallel()

	S := func(s ...string) []string { return s }
	TD := func(f string) string { return filepath.Join("testdata", "TestQuery", f) }
	cases := map[string]struct {
		path  string
		xpath string
		want  []string
	}{
		"single":   {TD("single.go"), "/*/File/FuncDecl/BlockStmt/*", S("ReturnStmt")},
		"multi":    {TD("multi.go"), "/*/File/FuncDecl/BlockStmt/*", S("ReturnStmt", "ReturnStmt")},
		"filename": {TD("multi.go"), "/a.go/File/FuncDecl/BlockStmt/*", S("ReturnStmt")},
	}

	for n, tt := range cases {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			s := newSeacher(t, tt.path)
			ns := s.Query(tt.xpath)
			got := nodesType(t, ns)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
