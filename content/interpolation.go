package content

import (
	"bytes"
	"text/template"
)

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
	output, err := process(tmpl, vars)

	if err != nil {
		return "", err
	}
	return output, nil
}

/*
By using text/template instead of html/template,
you avoid the automatic HTML escaping,
preserving the original HTML comments in the output.
*/
