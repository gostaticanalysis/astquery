# astquery [![PkgGoDev](https://pkg.go.dev/badge/github.com/gostaticanalysis/astquery)](https://pkg.go.dev/github.com/gostaticanalysis/astquery)

`astquery` selects a node set from AST by XPath.

`astquery` uses [antchfx/xpath](https://github.com/antchfx/xpath).
You can see a document of xpath expressions at [antchfx/xpath's repository](https://github.com/antchfx/xpath#expressions).

You can also use field of a node as an attributes such as `@Name` for `*ast.Ident`.
In addtion you can use the follows as an attribute:

 * `@type`: type of a node
 * `@pos`: `token.Position` of a node in string value
 * `@src`: source code representation of a node with `"go/format".Node`

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

```sh
# Find with src code snipet
$ astquery '//*[starts-with(@src, "panic") and @type="CallExpr"]/@src' fmt
panic("fmt: unknown base; can't happen")
panic(err)
panic(scanError{err})
panic(scanError{errors.New(err)})
panic(e)
panic(io.EOF)
panic(e)
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
