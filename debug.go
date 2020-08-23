package astquery

import (
	"fmt"
	"go/ast"
	"strings"
)

var isDebug bool

func debug(v ...interface{}) {
	if isDebug {
		fmt.Print(v...)
	}
}

func debugln(v ...interface{}) {
	if isDebug {
		fmt.Println(v...)
	}
}

func debugf(s string, v ...interface{}) {
	if isDebug {
		fmt.Printf(s, v...)
	}
}

func nodesToStr(ns []ast.Node) string {
	if !isDebug || ns == nil {
		return ""
	}

	s := make([]string, len(ns))
	for i := range s {
		s[i] = strings.TrimPrefix(fmt.Sprintf("%T", ns[i]), "*ast.")
	}
	return strings.Join(s, ",")
}
