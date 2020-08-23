package dontpanic

import (
	"github.com/gostaticanalysis/astquery"
	"golang.org/x/tools/go/analysis"
)

const doc = "dontpanic find panic calls"

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "dontpanic",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		astquery.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	e := pass.ResultOf[astquery.Analyzer].(*astquery.Evaluator)
	ns, err := e.Select("//*[@type='CallExpr']/Fun[@type='Ident' and @Name='panic']")
	if err != nil {
		return nil, err
	}

	for _, n := range ns {
		pass.Reportf(n.Pos(), "don't panic")
	}

	return nil, nil
}
