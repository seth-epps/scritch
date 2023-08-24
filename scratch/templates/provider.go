package templates

import (
	"text/template"
)

type TemplateProvider interface {
	Get() ([]Template, error)
	TargetPath() string
}

type Template struct {
	// Filepath relative to source to render the template to
	RelativePath string

	// Underlying Go yemplate to render
	GoTemplate *template.Template
}
