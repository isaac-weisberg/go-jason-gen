package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	packages "golang.org/x/tools/go/packages"
)

func main() {
	var packageLocation = "./sample"

	var thePackage, err = loadPackage(packageLocation)

	if err != nil {
		panic(err)
	}

	err = traverseThePackage(thePackage)
	if err != nil {
		panic(err)
	}
}

func w(err error, msg string) error {
	return errors.Join(errors.New(msg), err)
}

func e(msg string) error {
	return errors.New(msg)
}

type FirstClassField struct {
	name      string
	fieldType string
}

type StructDeclaration struct {
	structName       string
	firstClassFields []FirstClassField
}

func traverseThePackage(thePackage *packages.Package) error {
	var structDeclarations = make([]StructDeclaration, 0)

	for _, file := range thePackage.Syntax {
		ast.Inspect(file, func(node ast.Node) bool {
			declaration, ok := node.(*ast.GenDecl)
			if !ok {
				return true
			}

			if declaration.Tok != token.TYPE {
				return true
			}

			for _, spec := range declaration.Specs {
				var typeSpec = spec.(*ast.TypeSpec) // succeeds, always

				var structType, ok = typeSpec.Type.(*ast.StructType)
				if !ok {
					panic("cant generate for anything but structs")
				}

				var firstClassFields = make([]FirstClassField, 0)
				var metGojasonDecodable = false

				for _, field := range structType.Fields.List {
					fmt.Printf("%+v, %T\n", field, field.Type)

					fieldTypeAsIdent, ok := field.Type.(*ast.Ident)
					if ok {
						hasName := len(field.Names) > 0

						if hasName {
							name := field.Names[0].Name

							var firstClassField = FirstClassField{
								name:      name,
								fieldType: fieldTypeAsIdent.Name,
							}

							firstClassFields = append(firstClassFields, firstClassField)
						}

						continue
					}

					fieldTypeAsSelector, ok := field.Type.(*ast.SelectorExpr)
					if ok {
						xAsIdent, ok := fieldTypeAsSelector.X.(*ast.Ident)
						if ok {
							if xAsIdent.Name == "gojason" && fieldTypeAsSelector.Sel.Name == "Decodable" {
								metGojasonDecodable = true
							}
						}
						fmt.Printf("%+v, %T\n", fieldTypeAsSelector.X, fieldTypeAsSelector.X)
					}

				}

				if metGojasonDecodable {
					var structDeclaration = StructDeclaration{
						structName: typeSpec.Name.Name,
					}
					structDeclarations = append(structDeclarations, structDeclaration)
				}
			}

			return false
		})
	}

	return nil
}

func loadPackage(locationPattern string) (*packages.Package, error) {
	var tags = make([]string, 0)
	var loadConfig = packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}

	var packages, err = packages.Load(&loadConfig, locationPattern)
	if err != nil {
		return nil, w(err, "load packages failed")
	}

	if len(packages) == 0 {
		return nil, e("found no packages")
	}

	if len(packages) > 1 {
		return nil, e("found more than one package")
	}

	var firstPackage = packages[0]

	return firstPackage, nil
}

func contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
