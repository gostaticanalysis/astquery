# astquery [![PkgGoDev](https://pkg.go.dev/badge/github.com/gostaticanalysis/astquery)](https://pkg.go.dev/github.com/gostaticanalysis/astquery)

`astquery` selects a node set from AST by XPath.

## CLI Tool
### Install

```sh
$ go get -u github.com/gostaticanalysis/astquery/cmd/astquery
```

### How to use

```sh
$ astquery '//*[@type="CallExpr"]/Fun[@type="Ident" and @Name="panic"]' fmt
```

## Analyzer

see: [examples of analysis.Analyzer](_example)

```go
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
```
