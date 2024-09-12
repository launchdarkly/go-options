package main

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed render_full.gotmpl
var fullTemplateText string

//go:embed render_simple.gotmpl
var simpleTemplateText string

var funcMap = template.FuncMap{
	"ToPrivate": func(s string) string {
		return strings.ToLower(s[:1]) + s[1:]
	},
	"ToPublic": func(s string) string {
		return strings.ToUpper(s[:1]) + s[1:]
	},
}

var fullTemplate = template.Must(template.New("code").Funcs(funcMap).Parse(fullTemplateText))

var simpleTemplate = template.Must(template.New("code").Funcs(funcMap).Parse(simpleTemplateText))
