package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitanalyzer",
	Doc:  "Checks for the existence os.Exit() in the main package",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		for _, f := range pass.Files {
			ast.Inspect(f, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.File:
					if x.Name.Name != "main" {
						return false
					}
				case *ast.SelectorExpr:
					if ident, ok := x.X.(*ast.Ident); ok {
						if ident.Name == "os" && x.Sel.Name == "Exit" {
							pass.Reportf(ident.NamePos, "os.Exit is called in main package")
						}
					}
				}

				return true
			})
		}

		return nil, nil
	},
}
