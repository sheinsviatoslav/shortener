package staticlint

import (
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
)

// ExitMainCheckAnalyzer analyzes if os.Exit was called inside main func
var ExitMainCheckAnalyzer = &analysis.Analyzer{
	Name: "exitmaincheck",
	Doc:  "check for os.Exit call inside main func",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	ast.Inspect(f, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.CallExpr:
			switch fun := x.Fun.(type) {
			case *ast.SelectorExpr:
				if fun.X.(*ast.Ident).Name == "os" && fun.Sel.Name == "Exit" {
					pass.Reportf(x.Fun.Pos(), "Exit call inside main func")
				}
			}
		}
		return true
	})

	return nil, nil
}

// Checker adds code analyzers
func Checker() {
	var checks []*analysis.Analyzer
	//for _, v := range staticcheck.Analyzers {
	//	checks = append(checks, v.Analyzer)
	//}
	//
	//for _, v := range simple.Analyzers {
	//	checks = append(checks, v.Analyzer)
	//}

	checks = append(checks,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		ExitMainCheckAnalyzer,
	)

	multichecker.Main(
		checks...,
	)
}
