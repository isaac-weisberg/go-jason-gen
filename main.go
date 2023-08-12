package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path"
	"strings"

	packages "golang.org/x/tools/go/packages"
)

func main() {
	var args = os.Args
	if len(args) == 1 {
		panic("Please, provide package paths")
	}

	var remainingArgs = args[1:]

	for _, arg := range remainingArgs {
		var packageLocation = arg

		var thePackage, err = loadPackage(packageLocation)

		if err != nil {
			panic(err)
		}

		err = traverseThePackage(packageLocation, thePackage)
		if err != nil {
			panic(err)
		}
	}
}

type StructFieldType int64

const (
	FieldTypeFirstClass StructFieldType = iota
	FieldTypeEmbeddedStruct
)

type FirstClassFieldParsingStrategy int64

const (
	FirstClassFieldParsingStrategyInt64 FirstClassFieldParsingStrategy = iota
	FirstClassFieldParsingStrategyString
	FirstClassFieldParsingStrategyArbitraryStruct
)

func detectFirstClassFieldParsingStrategy(t string) FirstClassFieldParsingStrategy {
	switch t {
	case "int64":
		return FirstClassFieldParsingStrategyInt64
	case "string":
		return FirstClassFieldParsingStrategyString
	default:
		return FirstClassFieldParsingStrategyArbitraryStruct
	}
}

type FirstClassFieldTypeKind int64

const (
	FirstClassFieldTypeKindScalar FirstClassFieldTypeKind = iota
	FirstClassFieldTypeKindArray
)

type FirstClassField struct {
	fieldName        string
	typeName         string
	typeKind         FirstClassFieldTypeKind
	typeParsingStrat FirstClassFieldParsingStrategy
}

type EmbeddedStructField struct {
	embeddedTypeName string
}

type StructField struct {
	fieldType           StructFieldType
	firstClassField     FirstClassField
	embeddedStructField EmbeddedStructField
}

func newFirstClassStructField(
	fieldName string,
	typeName string,
	typeKind FirstClassFieldTypeKind,
	typeParsingStrat FirstClassFieldParsingStrategy,
) StructField {
	return StructField{
		fieldType: FieldTypeFirstClass,
		firstClassField: FirstClassField{
			fieldName:        fieldName,
			typeName:         typeName,
			typeKind:         typeKind,
			typeParsingStrat: typeParsingStrat,
		},
	}
}

func newEmbeddedStructField(embeddedTypeName string) StructField {
	return StructField{
		fieldType: FieldTypeEmbeddedStruct,
		embeddedStructField: EmbeddedStructField{
			embeddedTypeName: embeddedTypeName,
		},
	}
}

type StructDeclaration struct {
	name   string
	fields []StructField
}

func traverseThePackage(packageLocation string, thePackage *packages.Package) error {
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

				var fields = make([]StructField, 0)
				var metGojasonDecodable = false

				for _, field := range structType.Fields.List {
					fieldTypeAsIdent, ok := field.Type.(*ast.Ident)
					if ok {
						hasName := len(field.Names) > 0

						if hasName {
							name := field.Names[0].Name

							var typeKind = FirstClassFieldTypeKindScalar
							var parsingStrategy = detectFirstClassFieldParsingStrategy(fieldTypeAsIdent.Name)

							var firstClassField = newFirstClassStructField(name, fieldTypeAsIdent.Name, typeKind, parsingStrategy)
							fields = append(fields, firstClassField)
						} else {
							var embeddedStructField = newEmbeddedStructField(fieldTypeAsIdent.Name)
							fields = append(fields, embeddedStructField)
						}

						continue
					}

					fieldTypeAsArray, ok := field.Type.(*ast.ArrayType)
					if ok {
						elementTypeIdent, ok := fieldTypeAsArray.Elt.(*ast.Ident)
						if ok {
							if len(field.Names) == 0 {
								panic("when the field is array, the type MUST be present. a go source doesn't compile like this.")
							}

							var fieldName = field.Names[0].Name
							var typeName = elementTypeIdent.Name
							var typeKind = FirstClassFieldTypeKindArray
							var parsingStrategy = detectFirstClassFieldParsingStrategy(elementTypeIdent.Name)

							var firstClassField = newFirstClassStructField(fieldName, typeName, typeKind, parsingStrategy)
							fields = append(fields, firstClassField)
						}
					}

					fieldTypeAsSelector, ok := field.Type.(*ast.SelectorExpr)
					if ok {
						xAsIdent, ok := fieldTypeAsSelector.X.(*ast.Ident)
						if ok {
							if xAsIdent.Name == "gojason" && fieldTypeAsSelector.Sel.Name == "Decodable" {
								metGojasonDecodable = true
							}
						}
					}

				}

				if metGojasonDecodable {
					var structDeclaration = StructDeclaration{
						name:   typeSpec.Name.Name,
						fields: fields,
					}
					structDeclarations = append(structDeclarations, structDeclaration)
				}
			}

			return false
		})
	}

	var packageName = thePackage.Name
	generateStructDeclarations(packageName, packageLocation, structDeclarations)

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

type InitializationValue struct {
	valueName          string
	needsDereferencing bool
}

func generateStructDeclarations(packageName string, packageLocation string, structDeclarations []StructDeclaration) {
	var builder customBuilder

	builder.WriteLine("// Code generated by \"go-jason-gen\"; DO NOT EDIT.")
	builder.WriteLine()
	builder.WriteLine("package ", packageName)
	builder.WriteLine()

	if len(structDeclarations) > 0 {
		builder.WriteLine("import (")
		builder.WriteLineFI(1, `"errors"`)
		builder.WriteLineFI(1, `"fmt"`)
		builder.WriteLine()
		builder.WriteLineFI(1, `gojason "github.com/isaac-weisberg/go-jason"`)
		builder.WriteLineFI(1, `parser "github.com/isaac-weisberg/go-jason/parser"`)
		builder.WriteLineFI(1, `values "github.com/isaac-weisberg/go-jason/values"`)
		builder.WriteLine(")")
		builder.WriteLine()
	}

	for _, declaration := range structDeclarations {
		var structName = declaration.name
		var structNameCapitalized = firstCapitalized(structName)

		builder.WriteLineFI(0, "func make%sFromJson(bytes []byte) (*%s, error) {", structNameCapitalized, structName)

		builder.WriteLineIndent(1, "var j = errors.Join")
		builder.WriteLineIndent(1, "var e = errors.New")
		builder.WriteLine()
		builder.WriteLineIndent(1, "rootValueAny, err := parser.Parse(bytes)")
		builder.WriteLineIndent(1, "if err != nil {")
		builder.WriteLineIndent(2, `return nil, j(e("parsing json into an object tree failed"), err)`)
		builder.WriteLineIndent(1, "}")
		builder.WriteLine()
		builder.WriteLineIndent(1, "rootObject, err := rootValueAny.AsObject()")
		builder.WriteLineIndent(1, "if err != nil {")
		builder.WriteLineIndent(2, `return nil, j(e("interpreting root json value as an object failed"), err)`)
		builder.WriteLineIndent(1, "}")
		builder.WriteLine()
		builder.WriteLineFI(1, "parsedObject, err := parse%sFromJsonObject(rootObject)", structNameCapitalized)
		builder.WriteLineIndent(1, "if err != nil {")
		builder.WriteLineIndent(2, `return nil, j(e("parsing json into the resulting value failed"), err)`)
		builder.WriteLineIndent(1, "}")
		builder.WriteLine()
		builder.WriteLineIndent(1, "return parsedObject, nil")
		builder.WriteLine("}")
		builder.WriteLine()

		builder.WriteLineFI(0, "func parse%sFromJsonObject(rootObject *values.JsonValueObject) (*%s, error) {", structNameCapitalized, structName)
		builder.WriteLineIndent(1, "var j = errors.Join")
		builder.WriteLineIndent(1, "var e = errors.New")
		builder.WriteLine()
		builder.WriteLineIndent(1, "var stringKeyValues = rootObject.StringKeyedKeyValuesOnly()")
		builder.WriteLineIndent(1, "_ = stringKeyValues")
		builder.WriteLine()

		var keysAndValuesForThem = make(map[string]InitializationValue, 0)

		for _, field := range declaration.fields {
			switch field.fieldType {
			case FieldTypeFirstClass:
				generateFirstClassFieldDeclaration(&builder, keysAndValuesForThem, field.firstClassField)
			case FieldTypeEmbeddedStruct:
				generateEmbeddedStructFieldDeclaration(&builder, keysAndValuesForThem, field.embeddedStructField)
			}
		}

		builder.WriteLineFI(1, `var decodable = gojason.Decodable{}`)
		builder.WriteLineFI(1, `var resultingStruct%s = %s{`, structNameCapitalized, declaration.name)
		builder.WriteLineFI(2, `Decodable: decodable,`)
		for k, value := range keysAndValuesForThem {
			var formatStringRepresentingUsageOfValue string
			if value.needsDereferencing {
				formatStringRepresentingUsageOfValue = "*%s"
			} else {
				formatStringRepresentingUsageOfValue = "%s"
			}
			var valueUsage = fmt.Sprintf(formatStringRepresentingUsageOfValue, value.valueName)
			builder.WriteLineFI(2, `%s: %s,`, k, valueUsage)
		}
		builder.WriteLineFI(1, `}`)
		builder.WriteLineFI(1, `return &resultingStruct%s, nil`, structNameCapitalized)
		builder.WriteLine("}")
		builder.WriteLine()
	}

	var result = builder.String()

	fileNameToWrite := "go_jason_generated.go"
	filePathToWrite := path.Join(packageLocation, fileNameToWrite)

	os.WriteFile(filePathToWrite, []byte(result), 0644)
}

func generateFirstClassFieldDeclaration(builder *customBuilder, keysAndValues map[string]InitializationValue, firstClassField FirstClassField) {
	fmt.Printf("ASDF %+v\n", firstClassField)

	var fieldName = firstClassField.fieldName
	var fieldNameCapitalized = firstCapitalized(fieldName)
	var fieldType = firstClassField.typeName
	var fieldTypeCapitalized = firstCapitalized(fieldType)

	builder.WriteLineFI(1, `valueFor%sKey, exists := stringKeyValues["%s"]`, fieldNameCapitalized, fieldName)
	builder.WriteLineFI(1, `if !exists {`)
	builder.WriteLineFI(2, `return nil, j(e("value not found for key '%s'"))`, fieldName)
	builder.WriteLineFI(1, "}")

	switch firstClassField.typeKind {
	case FirstClassFieldTypeKindScalar:
		switch firstClassField.typeParsingStrat {
		case FirstClassFieldParsingStrategyInt64:
			builder.WriteLineFI(1, `valueFor%sKeyAsNumberValue, err := valueFor%sKey.AsNumber()`, fieldNameCapitalized, fieldNameCapitalized)
			builder.WriteLineFI(1, `if err != nil {`)
			builder.WriteLineFI(2, `return nil, j(e("interpreting JsonAny as Number failed for key '%s'"), err)`, fieldName)
			builder.WriteLineFI(1, `}`)

			var resultingValueName = fmt.Sprintf(`parsedInt64For%sKey`, fieldNameCapitalized)

			builder.WriteLineFI(1, `%s, err := valueFor%sKeyAsNumberValue.ParseInt64()`, resultingValueName, fieldNameCapitalized)
			builder.WriteLineFI(1, `if err != nil {`)
			builder.WriteLineFI(2, `return nil, j(e("parsing int64 from Number failed for key '%s'"), err)`, fieldName)
			builder.WriteLineFI(1, `}`)

			keysAndValues[fieldName] = InitializationValue{
				valueName:          resultingValueName,
				needsDereferencing: true,
			}
		case FirstClassFieldParsingStrategyString:
			builder.WriteLineFI(1, `valueFor%sKeyAsStringValue, err := valueFor%sKey.AsString()`, fieldNameCapitalized, fieldNameCapitalized)
			builder.WriteLineFI(1, `if err != nil {`)
			builder.WriteLineFI(2, `return nil, j(e("interpreting JsonAny as String failed for key '%s'"), err)`, fieldName)
			builder.WriteLineFI(1, `}`)

			var resultingValueName = fmt.Sprintf(`parsedStringFor%sKey`, fieldNameCapitalized)

			builder.WriteLineFI(1, `%s := valueFor%sKeyAsStringValue.String`, resultingValueName, fieldNameCapitalized)

			keysAndValues[fieldName] = InitializationValue{
				valueName:          resultingValueName,
				needsDereferencing: false,
			}
		case FirstClassFieldParsingStrategyArbitraryStruct:
			builder.WriteLineFI(1, `valueFor%sKeyAsObjectValue, err := valueFor%sKey.AsObject()`, fieldNameCapitalized, fieldNameCapitalized)
			builder.WriteLineFI(1, `if err != nil {`)
			builder.WriteLineFI(2, `return nil, j(e("interpreting JsonAny as Object failed for key '%s'"), err)`, fieldName)
			builder.WriteLineFI(1, `}`)
			var resultingValueName = fmt.Sprintf("parsedValueFor%sKey", fieldNameCapitalized)
			builder.WriteLineFI(1, `%s, err := parse%sFromJsonObject(valueFor%sKeyAsObjectValue)`, resultingValueName, fieldTypeCapitalized, fieldNameCapitalized)
			builder.WriteLineFI(1, `if err != nil {`)
			builder.WriteLineFI(2, `return nil, j(e("parsing '%s' from 'Object' failed for key '%s'"))`, fieldType, fieldName)
			builder.WriteLineFI(1, `}`)

			keysAndValues[fieldName] = InitializationValue{
				valueName:          resultingValueName,
				needsDereferencing: true,
			}
		default:
			panic("not supposed to happen")
		}
	case FirstClassFieldTypeKindArray:
		builder.WriteLineFI(1, `valueFor%sKeyAsArrayValue, err := valueFor%sKey.AsArray()`, fieldNameCapitalized, fieldNameCapitalized)
		builder.WriteLineFI(1, `if err != nil {`)
		builder.WriteLineFI(2, `return nil, j(e("interpreting JsonAny as Array failed for key '%s'"))`, fieldName)
		builder.WriteLineFI(1, `}`)

		var resultingValueName = fmt.Sprintf("resultingArrayFor%sKey", fieldNameCapitalized)
		builder.WriteLineFI(1, `%s := make([]%s, 0, len(valueFor%sKeyAsArrayValue.Values))`, resultingValueName, fieldType, fieldNameCapitalized)
		builder.WriteLineFI(1, `for index, element := range valueFor%sKeyAsArrayValue.Values {`, fieldNameCapitalized)

		switch firstClassField.typeParsingStrat {
		case FirstClassFieldParsingStrategyInt64:

		case FirstClassFieldParsingStrategyString:

		case FirstClassFieldParsingStrategyArbitraryStruct:
			builder.WriteLineFI(2, `elementAsObject, err := element.AsObject()`)
			builder.WriteLineFI(2, `if err != nil {`)
			builder.WriteLineFI(3, `return nil, j(e(fmt.Sprintf("attempted to interpret value at index '%%v' of array for key '%s' as object, but failed", index)), err)`, fieldName)
			builder.WriteLineFI(2, `}`)

			builder.WriteLineFI(2, `parsedValue, err := parse%sFromJsonObject(elementAsObject)`, fieldTypeCapitalized)
			builder.WriteLineFI(2, `if err != nil {`)
			builder.WriteLineFI(3, `return nil, j(e(fmt.Sprintf("failed to parse element at index '%%v' of array for key '%s'", index)), err)`, fieldName)
			builder.WriteLineFI(2, `}`)
			builder.WriteLineFI(2, `resultingArrayFor%sKey = append(resultingArrayFor%sKey, *parsedValue)`, fieldNameCapitalized, fieldNameCapitalized)
		default:
			panic("oops")
		}

		builder.WriteLineFI(1, `}`)

		keysAndValues[fieldName] = InitializationValue{
			valueName:          resultingValueName,
			needsDereferencing: false,
		}
	default:
		panic("no")
	}

	builder.WriteLine()
}

func generateEmbeddedStructFieldDeclaration(builder *customBuilder, keysAndValues map[string]InitializationValue, embeddedStructField EmbeddedStructField) {
	var embeddedTypeName = embeddedStructField.embeddedTypeName
	var capitalizedEmbeddedTypeName = firstCapitalized(embeddedTypeName)

	var resultingValueName = fmt.Sprintf(`valueForEmbedded%s`, capitalizedEmbeddedTypeName)

	builder.WriteLineFI(1, "%s, err := parse%sFromJsonObject(rootObject)", resultingValueName, capitalizedEmbeddedTypeName)
	builder.WriteLineFI(1, "if err != nil {")
	builder.WriteLineFI(2, `return nil, j(e("parsing embedded struct of type '%s' failed"), err)`, embeddedTypeName)
	builder.WriteLineFI(1, "}")

	keysAndValues[embeddedTypeName] = InitializationValue{
		valueName:          resultingValueName,
		needsDereferencing: true,
	}

	builder.WriteLine()
}
