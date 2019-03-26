package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/fatih/structtag"

	"golang.org/x/tools/go/packages"
)

var typeName string
var optionInterfaceName string
var outputName string
var applyFunctionName string
var createNewFunc bool
var runGoFmt bool

func initFlags() {
	flag.StringVar(&typeName, "type", "", "name of struct to create options for")
	flag.BoolVar(&createNewFunc, "new", true, "with to create a function to return a new config")
	flag.StringVar(&optionInterfaceName, "option", "Option", "name of the interface to use for options")
	flag.StringVar(&outputName, "output", "", "name of output file (default is <type>_options.go)")
	flag.StringVar(&applyFunctionName, "func", "", `name of function created to apply options to <type> (default is "apply<Type>Options")`)
	flag.BoolVar(&runGoFmt, "fmt", true, `set to false to skip go format`)
}

type Option struct {
	Name         string
	PublicName   string
	DefaultValue string
	Type         string
}

func main() {
	initFlags()
	flag.Parse()

	if typeName == "" {
		flag.Usage()
		log.Fatal("missing arguments")
	}

	cfg := &packages.Config{
		Mode:  packages.LoadSyntax,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if len(pkgs) != 1 {
		log.Fatalf("ERROR: expected a single package but %d packages were found", len(pkgs))
	}

	success := false
	for _, file := range pkgs[0].Syntax {
		ast.Inspect(file, func(node ast.Node) bool {
			found := writeOptionsFile(pkgs[0].Name, node)
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

func writeOptionsFile(packageName string, node ast.Node) (found bool) {
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

		if typeSpec.Name.String() != typeName {
			continue
		}

		var options []Option
		for _, field := range t.Fields.List {
			var defaultValue string
			var publicName string
			if field.Tag != nil {
				value := field.Tag.Value
				tags, err := structtag.Parse(value[1 : len(value)-1])
				if err == nil {
					tag, err := tags.Get("options")
					if err != nil {
						log.Fatalf(`ERROR: unable to parse struct tag "%s": %s`, field.Tag.Value, err)
					}
					if tag.Name == "-" {
						continue
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
			typeStr := fmt.Sprintf("%s", field.Type)
			if typeStr == "string" && defaultValue != "" {
				defaultValue = fmt.Sprintf("`%s`", defaultValue)
			}
			for _, n := range field.Names {
				if publicName == "" {
					publicName = n.Name
				}
				options = append(options, Option{
					Name:         n.Name,
					PublicName:   publicName,
					DefaultValue: defaultValue,
					Type:         typeStr,
				})
			}
		}

		outputFileName := fmt.Sprintf("%s_options.go", typeSpec.Name)
		if outputName != "" {
			outputFileName = outputName
		}

		buf := bytes.NewBuffer([]byte(fmt.Sprintf("package %s\n\n", packageName)))

		err := codeTemplate.Execute(buf, map[string]interface{}{
			"options":        options,
			"optionTypeName": optionInterfaceName,
			"configTypeName": typeName,
			"applyFuncName":  applyFunctionName,
			"createNewFunc":  createNewFunc,
		})
		if err != nil {
			log.Fatal(fmt.Errorf("template execute failed: %s", err))
		}
		if err := ioutil.WriteFile(outputFileName, buf.Bytes(), 0644); err != nil {
			log.Fatal(fmt.Errorf("write failed: %s", err))
		}
		if err := exec.Command("gofmt", "-w", outputFileName).Run(); err != nil {
			log.Fatal(fmt.Errorf("gofmt failed: %s", err))
		}
	}

	return true
}
