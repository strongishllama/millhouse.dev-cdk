package tmpl

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/gofor-little/xerror"
)

func NewTemplateFromFile(fileSystem embed.FS, path string, data interface{}) ([]byte, error) {
	tmpl := template.New("template")

	fileData, err := fileSystem.ReadFile(path)
	if err != nil {
		return nil, xerror.Wrap("failed to read file data", err)
	}

	tmpl, err = tmpl.Parse(string(fileData))
	if err != nil {
		return nil, xerror.Wrap("failed to parse template", err)
	}

	buffer := &bytes.Buffer{}
	if err := tmpl.Execute(buffer, data); err != nil {
		return nil, xerror.Wrap("failed to execute template", err)
	}

	return buffer.Bytes(), nil
}
