package content

import (
	"bytes"
	"html/template"
)

// From https://dev.to/kirklewis/go-text-template-processing-181d

// process applies the data structure 'vars' onto an already
// parsed template 't', and returns the resulting string.
func process(t *template.Template, vars interface{}) (string, error) {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		return "", err
	}
	return tmplBytes.String(), nil
}

func InterpolateString(str string, vars interface{}) (string, error) {
	tmpl, err := template.New("tmpl").Parse(str)

	if err != nil {
		return "", err
	}
	return process(tmpl, vars)
}

