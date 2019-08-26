package main

import (
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"ToTitle": func(s string) string {
		return strings.ToTitle(s[:1]) + s[1:]
	},
}

// store code generating template in constant
//go:generate bash -c "echo -e '// generated code, DO NOT EDIT\npackage main\n\nconst codeTemplateText = `' > template_text.go"
//go:generate bash -c "cat render.gotmpl >> template_text.go"
//go:generate bash -c "echo '`' >> template_text.go"
var codeTemplate = template.Must(template.New("code").Funcs(funcMap).Parse(codeTemplateText))
