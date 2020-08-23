# astquery [![PkgGoDev](https://pkg.go.dev/badge/github.com/gostaticanalysis/astquery)](https://pkg.go.dev/github.com/gostaticanalysis/astquery)

`astquery` selects a node set from AST by XPath.

`astquery` uses [antchfx/xpath](https://github.com/antchfx/xpath).
You can see a document of xpath expressions at [antchfx/xpath's repository](https://github.com/antchfx/xpath#expressions).

`@type` and `@pos` can use as an attribute.
`@type` represents type of a node and `@pos` is `token.Position` of a node in string value.
You can also use field of a node as an attributes such as `@Name` for `*ast.Ident`.

## CLI Tool
### Install

```sh
$ go get -u github.com/gostaticanalysis/astquery/cmd/astquery
```

### How to use

#### Select a node set

```sh
$ astquery '//*[@type="CallExpr"]/Fun[@type="Ident" and @Name="panic"]' fmt
```

#### Select attributes

```sh
# Find calling panic in fmt package
$ astquery '//*[@type="CallExpr"]/Fun[@type="Ident" and @Name="panic"]/@pos' fmt
/usr/local/go/src/fmt/format.go:266:3
/usr/local/go/src/fmt/print.go:553:4
/usr/local/go/src/fmt/scan.go:240:2
/usr/local/go/src/fmt/scan.go:244:2
/usr/local/go/src/fmt/scan.go:253:5
/usr/local/go/src/fmt/scan.go:508:3
/usr/local/go/src/fmt/scan.go:1064:4
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
