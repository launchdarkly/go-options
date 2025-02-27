package main

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed render.gotmpl
var codeTemplateText string

// store code generating template in constant
var codeTemplate = template.Must(template.New("code").Funcs(funcMap).Parse(codeTemplateText))

var funcMap = template.FuncMap{
	"ToPrivate": func(s string) string {
		return strings.ToLower(s[:1]) + s[1:]
	},
	"ToPublic": func(s string) string {
		return strings.ToUpper(s[:1]) + s[1:]
	},
}
