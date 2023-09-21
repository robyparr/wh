package template

import (
	"embed"
	"io"
	"text/template"
)

//go:embed *.txt
var templates embed.FS

func Render(out io.Writer, path string, data any) error {
	tmpl, err := template.ParseFS(templates, path)
	if err != nil {
		return err
	}

	return tmpl.Execute(out, data)
}
