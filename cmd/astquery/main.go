package main

import (
	"fmt"
	"os"

	"github.com/gostaticanalysis/astquery"
	"golang.org/x/tools/go/packages"
)

func main() {
	expr := "/"
	pattern := os.Args[1:]
	if len(os.Args) > 1 {
		expr = os.Args[1]
		pattern = os.Args[2:]
	}

	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps}
	pkgs, err := packages.Load(cfg, pattern...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		s := astquery.NewSeacher(pkg.Fset, pkg.Syntax, nil)
		ns := s.Query(expr)
		for _, n := range ns {
			fmt.Printf("%[1]T %[1]v\n", n)
		}
	}
}
