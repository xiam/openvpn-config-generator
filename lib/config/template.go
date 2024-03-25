package config

import (
	"bytes"
	"fmt"
	"strings"

	"text/template"
)

const (
	templateValue = `
{{- range . -}}
{{- if isString . -}}
{{.Name }} {{ .String | quotedValues }}
{{else if isEmbed . -}}
<{{ .Name }}>
{{.Embed | toString }}
</{{ .Name }}>
{{else -}}
{{ .Name }}
{{end -}}
{{- end -}}
`
)

var templateFuncs = template.FuncMap{
	"isString":     isString,
	"isEmbed":      isEmbed,
	"toString":     toString,
	"quotedValues": quotedValues,
}

func toString(v []byte) string {
	return string(v)
}

func quotedValues(values []string) string {
	res := []string{}

	for i := range values {
		res = append(res, fmt.Sprintf("%q", values[i]))
	}

	return strings.Join(res, " ")
}

func isString(v *configValue) bool {
	return v.Type == configTypeString
}

func isEmbed(v *configValue) bool {
	return v.Type == configTypeEmbed
}

func newTemplate() *template.Template {
	t := template.Must(template.New("").Funcs(templateFuncs).Parse(templateValue))
	return t
}

func compile(values []configValue) ([]byte, error) {
	t := newTemplate()

	buf := bytes.NewBuffer(nil)

	if err := t.Execute(buf, values); err != nil {
		return nil, err
	}

	return bytes.TrimSpace(buf.Bytes()), nil
}
