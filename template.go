package main

import (
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"ToPrivate": func(s string) string {
		return strings.ToLower(s[:1]) + s[1:]
	},
	"ToPublic": func(s string) string {
		return strings.ToUpper(s[:1]) + s[1:]
	},
}

// store code generating template in constant
//go:generate bash -c "echo -e '// generated code, DO NOT EDIT\npackage main\n\nconst codeTemplateText = `' > template_text.go"
//go:generate bash -c "cat render.gotmpl >> template_text.go"
//go:generate bash -c "echo '`' >> template_text.go"
var codeTemplate = template.Must(template.New("code").Funcs(funcMap).Parse(codeTemplateText))
