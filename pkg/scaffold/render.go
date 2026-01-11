package scaffold

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type Values struct {
	ENV_REPO_NAME string
	ENV_REPO_URL  string
	ENV_REPO_TAG  string
	Envs          []string
}

func renderArgoAppSet(tmpl string, v Values) (string, error) {
	t, err := template.New("argo").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, v); err != nil {
		return "", err
	}
	return buf.String(), nil
}
