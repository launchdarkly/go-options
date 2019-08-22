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
var imports string

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s <type>:\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "  %s [<option> ... ] <config type> ...\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "  where <option> can be any of:\n\n")
	flag.PrintDefaults()
}

func initFlags() {
	flag.StringVar(&typeName, "type", "", "name of struct to create options for")
	flag.BoolVar(&createNewFunc, "new", true, "with to create a function to return a new config")
	flag.StringVar(&optionInterfaceName, "option", "Option", "name of the interface to use for options")
	flag.StringVar(&imports, "imports", "", "a comma-separated list of packages with optional alias (e.g. time,url=net/url) ")
	flag.StringVar(&outputName, "output", "", "name of output file (default is <type>_options.go)")
	flag.StringVar(&applyFunctionName, "func", "", `name of function created to apply options to <type> (default is "apply<Type>Options")`)
	flag.StringVar(&applyOptionFunctionType, "option_func", "",
		`name of function type created to apply options with pointer receiver to <type> (default is "apply<Option>Func")`)
	flag.StringVar(&optionPrefix, "prefix", "", `name of prefix to use for options (default is the same as "option")`)
	flag.BoolVar(&runGoFmt, "fmt", true, `set to false to skip go format`)
	flag.Usage = Usage
}

type Option struct {
	Name         string
	PublicName   string
	DefaultValue string
	Type         string
	HasPointerReceiver bool
}

type Import struct {
	Alias string
	Path string
}

func main() {
	initFlags()
	flag.Parse()
flag.CommandLine.ErrorHandling()
	types := flag.Args()

	if typeName == "" && len(types) == 0 {
		flag.Usage()
		log.Fatal("missing arguments")
	}

	if typeName != "" {
		types = append(types, typeName)
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
			typeBuf := new(bytes.Buffer)
			if err := printer.Fprint(typeBuf, fset, field.Type); err != nil {
				log.Fatalf("ERROR: unable to print type: %s", err)
			}
			typeStr :=  typeBuf.String()
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
					HasPointerReceiver: strings.HasPrefix(strings.TrimSpace(typeStr), "*"), // TODO: inspect the actual type
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
