package reset

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ServiceReset - сервис для генерации метода сброса у структур
type ServiceReset struct {
	root        string // Путь до директории относительно к которой запускается генератор
	PackageInfo []PackageInfo
}

func NewServiceReset(root string) *ServiceReset {
	return &ServiceReset{
		root:        root,
		PackageInfo: []PackageInfo{},
	}
}

func (s *ServiceReset) Run() error {
	err := filepath.Walk(s.root, func(path string, info os.FileInfo, err error) error {

		// Пропускаем определенные директории и с префиксом .
		if info.IsDir() {
			if info.Name() == "vendor" ||
				strings.HasPrefix(info.Name(), ".") && info.Name() != ".." {
				return filepath.SkipDir
			}
		}

		// пропускаем файлы не .go
		if filepath.Ext(path) != ".go" {
			return nil
		}

		fsComment := token.NewFileSet()
		fileComment, err := parser.ParseFile(fsComment, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		// проверяем, что есть комментарий для запуска генерации
		isReset := false
		for _, commentGroup := range fileComment.Comments {
			for _, comment := range commentGroup.List {
				if comment.Text == "// generate:reset" {
					isReset = true
					break
				}
			}
			if isReset {
				break
			}
		}
		if !isReset {
			return nil
		}
		packagePath := filepath.Dir(path)

		// парсим файл для поиска структур
		fs := token.NewFileSet()
		file, err := parser.ParseFile(fs, path, nil, 0)
		if err != nil {
			return err
		}

		packageInfo := PackageInfo{
			Name:    fileComment.Name.Name,
			Path:    packagePath,
			Structs: []StructInfo{},
		}

		for _, decl := range file.Decls {
			//  проверяем, что это объявление типа
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			// анализируем спецификацию типа
			for _, spec := range genDecl.Specs {
				// является ли спецификация объявлением типа
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				// является ли объявление типа определением структуры
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				packageInfo.Structs = append(packageInfo.Structs, StructInfo{
					Name:   typeSpec.Name.Name,
					Fields: getStructFields(structType),
				})
			}
		}

		s.appendPackageInfo(packageInfo)
		return nil
	})
	if err != nil {
		return err
	}

	s.Write()
	return nil
}

// appendPackageInfo - добавляем пакет в список
func (s *ServiceReset) appendPackageInfo(PackageInfo PackageInfo) {
	for i, p := range s.PackageInfo {
		// проверяем если ли уже такой пакет
		if p.Name == PackageInfo.Name && p.Path == PackageInfo.Path {
			s.PackageInfo[i].Structs = append(PackageInfo.Structs, p.Structs...)
			return
		}
	}
	s.PackageInfo = append(s.PackageInfo, PackageInfo)
}

func (s *ServiceReset) Write() {

	for _, p := range s.PackageInfo {
		bufFmt := GenerateResetTemplate(p)
		fmt.Println(string(bufFmt))
		fmt.Println("package", p.Name, p.Path)
		resetFilePath := filepath.Join(p.Path, "reset.gen.go")
		_ = os.Remove(resetFilePath)
		// Записываем сгенерированный код в файл
		err := os.WriteFile(resetFilePath, bufFmt, 0644)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

// getStructFields - получаем поля структуры
func getStructFields(structType *ast.StructType) []StructField {
	var fields []StructField
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			fields = append(fields, StructField{
				Name: name.Name,
				Type: field.Type,
			})
		}
	}
	return fields
}
