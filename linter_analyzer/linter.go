package linter_analyzer

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "mylinter",
	Doc:  "reports usage of ...",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	pkgName := pass.Pkg.Name()
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Ищем вызовы функций
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Получаем имя функции вызова
			fun := callExpr.Fun

			switch f := fun.(type) {
			case *ast.Ident:
				// Ищем: panic
				if f.Name == "panic" {
					pass.Reportf(callExpr.Pos(), "usage of panic is discouraged")
				}
			case *ast.SelectorExpr:
				// Ищем: log.Fatal или os.Exit вне main
				pkgIdent, ok := f.X.(*ast.Ident)
				if !ok {
					return true
				}
				objName := f.Sel.Name
				pkgNameId := pkgIdent.Name
				fmt.Println("pkgNameId", pkgNameId)
				if (pkgNameId == "log" && objName == "Fatal") || (pkgNameId == "os" && objName == "Exit") {
					// находим окружающую функцию
					enclosingFunc := enclosingFuncName(pass, callExpr)
					if !(pkgName == "main" && enclosingFunc == "main") {
						pass.Reportf(callExpr.Pos(), "call to %s.%s outside main.main function", pkgNameId, objName)
					}
				}
			}

			return true
		})
	}
	return nil, nil
}

// поиск имени функции, в которой находится вызов
func enclosingFuncName(pass *analysis.Pass, node ast.Node) string {
	fmt.Println(pass.Files)
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			// проверяем является ли функцией
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				// находится ли узел внутри функции
				if node.Pos() >= funcDecl.Pos() && node.End() <= funcDecl.End() {
					return funcDecl.Name.Name
				}
			}
		}
	}
	return ""
}
