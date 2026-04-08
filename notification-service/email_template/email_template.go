package emailTemplate

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed *.gohtml
var templateFile embed.FS

var Template *template.Template

func init() {
	Template = template.Must(template.ParseFS(templateFile, "*.gohtml"))
}

func RenderTemplate(name string, data any) (string, error) {
	buf := new(bytes.Buffer)
	err := Template.ExecuteTemplate(buf, name, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
