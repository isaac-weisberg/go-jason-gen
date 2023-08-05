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

	var typesOfInterest = []string{"addMoneyRequest"}

	err = traverseThePackage(thePackage, typesOfInterest)
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

func traverseThePackage(thePackage *packages.Package, typesToGenerate []string) error {
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

				var typeName = typeSpec.Name.Name
				if !contains[string](typesToGenerate, typeName) {
					// not interested in this type

					continue
				}

				var structType, ok = typeSpec.Type.(*ast.StructType)
				if !ok {
					panic("cant generate for structs")
				}

				for _, field := range structType.Fields.List {
					fmt.Printf("%+v\n", field)
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
