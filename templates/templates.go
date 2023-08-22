package templates

import (
	"fmt"
	"path/filepath"
	"text/template"

	"slices"
)

//go:generate go run generate.go

// Supported
func GetNativelySupportedTemplates(language, variant string) ([]*template.Template, error) {
	if variants, ok := SupportedTemplates[language]; ok {

		if !slices.Contains(variants, variant) {
			fmt.Printf("%v variant is not supported for %v\n", variant, language)
			return nil, fmt.Errorf("%v variant is not supported for %v", variant, language)
		}
		return findLanguageVariant(language, variant)
	} else {
		fmt.Printf("%v is not supported natively\n", language)
		return nil, fmt.Errorf("%v is not supported natively", language)
	}
}

func findLanguageVariant(language, variant string) ([]*template.Template, error) {
	embeddedPath := filepath.Join(templateFSRoot, language, variant)
	pathGlob := filepath.Join(embeddedPath, "*.tpl")

	tmplGlob := template.Must(
		template.ParseFS(
			templateFS, pathGlob,
		),
	)
	tmpls := tmplGlob.Templates()
	if len(tmpls) == 0 {
		return nil, fmt.Errorf("Could not find native templates for `%v (%v)`", language, variant)
	}
	return tmpls, nil

}
