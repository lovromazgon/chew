package funcmap

const (
	simpleFuncDocTemplate = "simpleFuncDoc"
	nestedFuncDocTemplate = "nestedFuncDoc"
)

type FuncDoc struct {
	Name    string
	Text    string
	Example string

	NestedFuncs []FuncDoc
}

func (fd FuncDoc) Template() string {
	if fd.NestedFuncs == nil || len(fd.NestedFuncs) == 0 {
		return simpleFuncDocTemplate
	} else {
		return nestedFuncDocTemplate
	}
}

var (
	FuncDocTemplates = `
{{ define "simpleFuncDoc.tmpl" -}}
──────────────────────────────────

Function '
{{- if exists . "parent" }}
  {{- .parent.Name }}.
{{- end }}
{{- .Name }}'
--  {{ .Example }}

{{ .Text }}

{{ end }}

{{ define "nestedFuncDoc.tmpl" -}}
──────────────────────────────────────

Function package '{{ .Name }}'

{{ .Text }}

    Nested functions:
{{ range .NestedFuncs }}
  {{- indentTemplate .Template . $ 4 }}
{{ end }}
{{ end }}
	`
)
