package reset

import "go/ast"

// generate:reset
// PackageInfo - информация о пакете
type PackageInfo struct {
	Name    string
	Path    string
	Structs []StructInfo
	i       int
	str     string
	strP    *string
	s       []int
	m       map[string]string
	child   *PackageInfo
}

// StructInfo - информация о структуре
type StructInfo struct {
	Name   string
	Fields []StructField
}

// StructField - информация о поле структуры
type StructField struct {
	Name string
	Type ast.Expr
}
