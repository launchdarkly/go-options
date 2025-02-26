// generated code, DO NOT EDIT
package main

const codeTemplateText = `
// Code generated by github.com/launchdarkly/go-options.  DO NOT EDIT.

{{ if and .implementString .options -}}
import "fmt"
{{ end }}

{{ if .imports -}}
import (
{{- range .imports }}
    {{ if .Alias }}  {{ .Alias }} "{{ .Path }}"{{ else }}  "{{ .Path }}"{{ end -}}
{{ end }}
)
{{ end }}

{{ if and .implementEqual .options -}}
import "github.com/google/go-cmp/cmp"
{{ end }}

{{ $applyOptionFuncType := or $.applyOptionFuncType (printf "Apply%sFunc" (ToPublic $.optionTypeName)) }}

type {{ $applyOptionFuncType }} func(c *{{ $.configTypeName }}) {{ if $.returnError -}} error {{ end }}

func (f {{ $applyOptionFuncType }}) apply(c *{{ $.configTypeName }}) {{ if $.returnError -}} error {{ end }} {
    {{ if $.returnError -}} return {{ end }} f(c)
}

{{ $applyFuncName := or $.applyFuncName (printf "apply%sOptions" (ToPublic $.configTypeName)) }}

{{ if $.createNewFunc}}
func {{ if $.newFuncPublic -}}New{{- else -}}new{{- end -}}{{ $.configTypeName | ToPublic}}(options ...{{ $.optionTypeName }}) {{ if $.returnError -}} ({{ $.configTypeName }} , error) {{else}} {{ $.configTypeName }} {{ end }} {
    var c {{ $.configTypeName }}
    {{ if $.returnError -}}
    err := {{ $applyFuncName }}(&c, options...)
    return c, err
    {{- else -}}
    {{ $applyFuncName }}(&c, options...)
    return c
    {{- end }}
}
{{ end }}

func {{ $applyFuncName }}(c *{{ $.configTypeName }}, options ...{{ $.optionTypeName }}) {{ if $.returnError -}} error {{ end }} {
{{- range .options -}}{{ $optionName := .Name }}{{ if .DefaultValue }}
    c.{{ .Name }} = {{ .DefaultValue }}
{{- end }}{{ if .IsStruct }}{{ range .Fields }}{{ if .DefaultValue }}
    c.{{ $optionName }}.{{ .Name }} = {{ .DefaultValue }}
{{- end }}{{ end }}
{{- end }}{{ end }}
{{ if $.returnError -}}
    for _, o := range options {
        if err := o.apply(c); err != nil {
            return err
            }
    }
    return nil
{{- else -}}
    for _, o := range options {
        o.apply(c)
    }
{{- end }}
}

type {{ $.optionTypeName }} interface {
    apply(*{{ $.configTypeName }}) {{ if $.returnError -}} error {{ end }}
}

{{ range .options }}{{ $option := . }}

{{ $name := .PublicName | ToPublic | printf "%s%s" $.optionPrefix }}
{{ if $.optionSuffix }}{{ $name = $.optionSuffix | printf "%s%s" (.PublicName | ToPublic) }}{{ end }}

{{ $implName := $name | printf "%sImpl" | ToPrivate }}

type {{ $implName }} struct {
{{- range .Fields }}
    {{ .ParamName }} {{ .Type }}
{{- end }}
}

func (o {{ $implName }}) apply(c *{{ $.configTypeName }}) {{ if $.returnError -}} error {{ end }} {
{{- if and $option.IsStruct $option.DefaultIsNil }}
    c.{{ $option.Name }} = new({{ $option.Type }})
{{- end }}
{{- range .Fields }}{{ if $option.IsStruct }}
    c.{{ $option.Name }}.{{ .Name }} = o.{{ .ParamName }}
{{- else }}
    c.{{ $option.Name }} = {{ if $option.DefaultIsNil }}&{{ end }}o.{{ .ParamName }}
{{- end }}{{- end }}
{{ if $.returnError -}}
    return nil
{{- end }}
}

{{ if $.implementEqual -}}
func (o {{ $implName }}) Equal(v {{ $implName }}) bool {
    switch {
{{- range .Fields }}
    case !cmp.Equal(o.{{ .ParamName }}, v.{{ .ParamName }}):
        return false
{{- end }}
    }
    return true
}
{{ end }}

{{ if $.implementString -}}
func (o {{ $implName }}) String() string {
    name := "{{ $name }}"
{{ if $option.IsStruct }}
    type stripped {{ $implName }}
    value := stripped(o)
{{- else -}}
{{- range .Fields }}{{/* there should only be one field since this isn't a struct */}}
    // hack to avoid go vet error about passing a function to Sprintf
    var value interface{} = o.{{ .ParamName }}
{{- end }}
{{- end }}
    return fmt.Sprintf("%s: %+v", name, value)
}
{{ end }}

{{ if .Docs }}
{{- range $i, $doc := .Docs }}// {{ if eq $i 0 }}{{ $name }} {{ end }}{{ $doc }}{{ end -}}
{{ end -}}
func {{ $name }}(
{{- range $i, $f := .Fields }}{{ if ne $i 0 }},{{ end }}{{ $f.ParamName }} {{ $f.ParamType }}{{ end -}}
) {{ $.optionTypeName }} {
    return {{ $implName }}{
{{- range .Fields }}
        {{ .ParamName }}: {{ .ParamName }},
{{- end }}
    }
}
{{ end }}
`
