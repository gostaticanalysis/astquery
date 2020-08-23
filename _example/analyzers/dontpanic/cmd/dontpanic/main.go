package main

import (
	"github.com/gostaticanalysis/astquery/_example/analyzers/dontpanic"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(dontpanic.Analyzer) }

