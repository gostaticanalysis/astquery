package main

import (
	"fmt"
	"go/ast"
	"os"
	"reflect"

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
		e := astquery.New(pkg.Fset, pkg.Syntax, nil)
		v, err := e.Eval(expr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "eval: %v\n", err)
			os.Exit(1)
		}

		switch v := v.(type) {
		case []ast.Node:
			for _, n := range v {
				fmt.Printf("%[1]T %[1]v\n", n)
			}
		default:
			rv := reflect.ValueOf(v)
			switch rv.Kind() {
			case reflect.Array, reflect.Slice:
				for i := 0; i < rv.Len(); i++ {
					fmt.Println(rv.Index(i).Interface())
				}
			case reflect.Map:
				for _, key := range rv.MapKeys() {
					val := rv.MapIndex(key)
					fmt.Printf("%v:%v\n", key.Interface(), val.Interface())
				}
			default:
				fmt.Println(v)
			}
		}
	}
}
