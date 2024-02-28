package templates

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

//go:generate go run supported_templates/generate.go

type EmbeddedTemplateProvider struct {
	language string
	variant  string
}

func NewEmbeddedTemplateProvider(language, variant string) EmbeddedTemplateProvider {
	return EmbeddedTemplateProvider{language: language, variant: variant}
}

func (tp EmbeddedTemplateProvider) Get() ([]Template, error) {
	if !slices.Contains(SupportedTemplates, tp.language) {
		return nil, fmt.Errorf("%v is not supported", tp.language)
	}
	return collectTemplates(tp.language)
}

func (tp EmbeddedTemplateProvider) TargetPath() string {
	return filepath.Join(tp.language, tp.variant)
}

func collectTemplates(language string) ([]Template, error) {
	embeddedPath := filepath.Join(templateFSRoot, language)
	pathGlob := filepath.Join(embeddedPath, "*.tpl")

	tmplGlob := template.Must(
		template.ParseFS(
			templateFS, pathGlob,
		),
	)

	var templates []Template
	for _, tmpl := range tmplGlob.Templates() {
		relPath, _ := strings.CutSuffix(tmpl.Name(), ".tpl")
		templates = append(templates, Template{
			RelativePath: relPath,
			GoTemplate:   tmpl,
		})
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("could not find native templates for `%v`", language)
	}

	return templates, nil
}
