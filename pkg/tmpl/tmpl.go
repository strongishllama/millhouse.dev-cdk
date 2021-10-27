package tmpl

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

func NewTemplateFromFile(fileSystem embed.FS, path string, data interface{}) ([]byte, error) {
	tmpl := template.New("template")

	fileData, err := fileSystem.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	tmpl, err = tmpl.Parse(string(fileData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	buffer := &bytes.Buffer{}
	if err := tmpl.Execute(buffer, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buffer.Bytes(), nil
}
