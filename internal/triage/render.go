package triage

import (
	"bytes"
	"text/template"
)

func RenderTriage(data TriageData) (string, error) {
	tmpl, err := template.New("triage").Parse(triageTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
