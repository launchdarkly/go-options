package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/structtag"

	"golang.org/x/tools/go/packages"
)

var typeName string
var optionInterfaceName string
var outputName string
var applyFunctionName string
var applyOptionFunctionType string
var createNewFunc bool
var runGoFmt bool
var optionPrefix string
var optionSuffix string
var imports string
var quoteStrings bool

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s <type>:\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "  %s [<option> ... ] <config type> ...\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "  where <option> can be any of:\n\n")
	flag.PrintDefaults()
}

func initFlags() {
	flag.StringVar(&typeName, "type", "", "name of struct to create options for")
	flag.BoolVar(&createNewFunc, "new", true, "whether to create a function to return a new config")
	flag.StringVar(&optionInterfaceName, "option", "Option", "name of the interface to use for options")
	flag.StringVar(&imports, "imports", "", "a comma-separated list of packages with optional alias (e.g. time,url=net/url) ")
	flag.StringVar(&outputName, "output", "", "name of output file (default is <type>_options.go)")
	flag.StringVar(&applyFunctionName, "func", "", `name of function created to apply options to <type> (default is "apply<Type>Options")`)
	flag.StringVar(&applyOptionFunctionType, "option_func", "",
		`name of function type created to apply options with pointer receiver to <type> (default is "apply<Option>Func")`)
	flag.StringVar(&optionPrefix, "prefix", "", `name of prefix to use for options (default is the same as "option")`)
	flag.StringVar(&optionSuffix, "suffix", "", `name of suffix to use for options (forces use of suffix, cannot with used with prefix)`)
	flag.BoolVar(&quoteStrings, "quote-default-strings", true, `set to false to disable automatic quoting of string field defaults`)
	flag.BoolVar(&runGoFmt, "fmt", true, `set to false to skip go format`)
	flag.Usage = Usage
}

type Field struct {
	Name         string
	ParamName    string
	Type         string
	DefaultValue string
}

type Option struct {
	Name         string
	PublicName   string
	DefaultValue string
	Fields       []Field
	Docs         []string
	DefaultIsNil bool
	IsStruct     bool
	Type         string
}

type Import struct {
	Alias string
	Path  string
}

func main() {
	initFlags()
	flag.Parse()
	flag.CommandLine.ErrorHandling()
	types := flag.Args()

	if optionPrefix != "" && optionSuffix != "" {
		log.Fatal("cannot specify both -prefix and -suffix options")
	}

	if typeName == "" && len(types) == 0 {
		flag.Usage()
		log.Fatal("missing arguments")
	}

	if typeName != "" {
		types = append(types, typeName)
	}

	cfg := &packages.Config{
		Mode:  packages.NeedSyntax | packages.NeedTypes | packages.NeedName,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg)
	if err != nil {
		log.Fatalf("ERROR: Failed to load pacakges %s", err)
	}

	if len(pkgs) != 1 {
		log.Fatalf("ERROR: expected a single package but %d packages were found", len(pkgs))
	}

	success := false
	for _, file := range pkgs[0].Syntax {
		ast.Inspect(file, func(node ast.Node) bool {
			found := writeOptionsFile(types, pkgs[0].Name, node, pkgs[0].Fset)
			if found {
				success = true
			}
			return !found
		})
	}

	if !success {
		log.Fatalf(`unable to find type "%s"`, typeName)
	}
}

func writeOptionsFile(types []string, packageName string, node ast.Node, fset *token.FileSet) (found bool) {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.TYPE {
		return false
	}

	for _, spec := range decl.Specs {
		typeSpec := spec.(*ast.TypeSpec)
		t, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		var typeName string
		for _, n := range types {
			if typeSpec.Name.String() == n {
				typeName = n
				break
			}
		}
		if typeName == "" {
			continue
		}

		var options []Option
		for _, field := range t.Fields.List {
			publicName, defaultValue, skip := parseStructTag(field)
			if skip {
				continue
			}
			var docs []string
			if field.Doc != nil {
				docs = append(docs, field.Doc.Text())
			}
			if field.Comment != nil {
				docs = append(docs, field.Comment.Text())
			}

			typeStr := getType(fset, field.Type)

			fieldType := field.Type
			defaultIsNil := false
			if t, isStar := fieldType.(*ast.StarExpr); isStar {
				switch t.X.(type) {
				case *ast.StructType, *ast.ArrayType:
					fieldType = t.X
					defaultIsNil = true
					typeStr = getType(fset, t.X)
				default:
					if strings.HasPrefix(publicName, "*") {
						publicName = publicName[1:]
						typeStr = getType(fset, t.X)
						defaultIsNil = true
					}
				}
			}

			isStruct := false
			var fields []Field
			switch t := fieldType.(type) {
			case *ast.StructType:
				isStruct = true
				for _, sfield := range t.Fields.List {
					paramName, defaultValue, skip := parseStructTag(sfield)
					if skip {
						continue
					}
					typeStr := getType(fset, sfield.Type)
					if strings.HasSuffix(paramName, "...") {
						paramName = paramName[0 : len(paramName)-3]
						switch t := sfield.Type.(type) {
						case *ast.ArrayType:
							typeStr = "..." + getType(fset, t.Elt)
						default:
							log.Fatalf(`expected slice type for "%+v"`, sfield)
						}
					}
					for _, n := range sfield.Names {
						fields = append(fields, Field{
							Name:         n.Name,
							ParamName:    stringsOr(paramName, n.Name),
							Type:         typeStr,
							DefaultValue: defaultValue,
						})
					}
				}
			case *ast.ArrayType:
				if strings.HasSuffix(publicName, "...") {
					publicName = publicName[0 : len(publicName)-3]
					typeStr = "..." + getType(fset, t.Elt)
				}
				fields = append(fields, Field{Name: "", ParamName: "o", Type: typeStr})
			default:
				fields = append(fields, Field{Name: "", ParamName: "o", Type: typeStr})
			}

			if defaultIsNil && defaultValue != "" {
				log.Fatalf(`cannot use pointer value with default value for fields %+v`, field.Names)
			}

			for _, n := range field.Names {
				options = append(options, Option{
					Name:         n.Name,
					PublicName:   stringsOr(publicName, n.Name),
					DefaultValue: defaultValue,
					Fields:       fields,
					Docs:         docs,
					DefaultIsNil: defaultIsNil,
					IsStruct:     isStruct,
					Type:         typeStr,
				})
			}
		}

		var importList []Import
		if imports != "" {
			for _, s := range strings.Split(imports, ",") {
				parts := strings.Split(s, "=")
				if len(parts) == 1 {
					importList = append(importList, Import{Path: parts[0]})
				} else if len(parts) == 2 {
					importList = append(importList, Import{Alias: parts[0], Path: parts[1]})
				} else {
					log.Fatalf(`ERROR: unexpected import description "%s"`, s)
				}
			}
		}

		outputFileName := fmt.Sprintf("%s_options.go", typeSpec.Name)
		if outputName != "" {
			outputFileName = outputName
		}

		buf := bytes.NewBuffer([]byte(fmt.Sprintf("package %s\n\n", packageName)))

		prefix := optionInterfaceName
		if optionPrefix != "" {
			prefix = optionPrefix
		}

		err := codeTemplate.Execute(buf, map[string]interface{}{
			"imports":             importList,
			"options":             options,
			"optionTypeName":      optionInterfaceName,
			"configTypeName":      typeName,
			"optionPrefix":        prefix,
			"optionSuffix":        optionSuffix,
			"applyFuncName":       applyFunctionName,
			"applyOptionFuncName": applyOptionFunctionType,
			"createNewFunc":       createNewFunc,
		})
		if err != nil {
			log.Fatal(fmt.Errorf("template execute failed: %s", err))
		}
		if err := ioutil.WriteFile(outputFileName, buf.Bytes(), 0644); err != nil {
			log.Fatal(fmt.Errorf("write failed: %s", err))
		}
		cmd := exec.Command("gofmt", "-w", outputFileName)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(fmt.Errorf("gofmt failed: %s", err))
		}
	}

	return true
}

func parseStructTag(field *ast.Field) (publicName string, defaultValue string, skip bool) {
	if field.Tag != nil {
		value := field.Tag.Value
		tags, err := structtag.Parse(value[1 : len(value)-1])
		if err == nil {
			tag, err := tags.Get("options")
			if err != nil && err.Error() == "tag does not exist" {
				goto SkipTag
			} else if err != nil {
				log.Fatalf(`ERROR: unable to parse struct tag "%s": %s`, field.Tag.Value, err)
			}
			if tag.Name == "-" {
				return "", "", true
			}
			publicName = tag.Name
			if len(tag.Options) > 0 {
				defaultValue = tag.Options[0]
			}
			if len(tag.Options) > 1 {
				log.Fatalf(`ERROR: format is options:"<name>,<default value>"`)
			}
		}
	}
	SkipTag:
	return publicName, formatDefault(field.Type, defaultValue), false
}

// getType returns a string of the type for a field by looking it up in the original source
func getType(fset *token.FileSet, fieldType ast.Expr) string {
	typeBuf := new(bytes.Buffer)
	if err := printer.Fprint(typeBuf, fset, fieldType); err != nil {
		log.Fatalf("ERROR: unable to print type: %s", err)
	}
	return typeBuf.String()
}

// formatDefault adds quotes to default values for string types
func formatDefault(fieldType ast.Expr, defaultValue string) string {
	switch t := fieldType.(type) {
	case *ast.Ident:
		if t.Name == "string" && defaultValue != "" && quoteStrings {
			return fmt.Sprintf("`%s`", defaultValue)
		}
	}
	return defaultValue
}

// return first non-empty string
func stringsOr(strs ...string) string {
	for _, s := range strs {
		if s != "" {
			return s
		}
	}
	return ""
}
