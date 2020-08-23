package astquery

import (
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer provides *astquery.Evaluator as a result.
//
// Example:
//	func run(pass *analysis.Pass) (interface{}, error) {
//		e := pass.ResultOf[astquery.Analyzer].(*astquery.Evaluator)
//		ns, err := e.Select("//*[@type='CallExpr']/Fun[@type='Ident' and @Name='panic']")
//		if err != nil {
//			return nil, err
//		}
//		
//		for _, n := range ns {
//			pass.Reportf(n.Pos(), "don't panic")
//		}
//		
//		return nil, nil
//	}
var Analyzer = &analysis.Analyzer{
	Name: "astquery",
	Doc:  "search nodes by xpath",
	Run:  new(analyzer).run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	ResultType: reflect.TypeOf(new(Evaluator)),
}

type analyzer struct{}

func (analyzer) run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	return New(pass.Fset, pass.Files, inspect), nil
}
