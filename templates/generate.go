//go:build exclude

// This program generates supported.go.
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

const supportedTemplatesDirectory = "supported_templates"

func main() {
	supportedTemplates, err := getSupportedLanguageTemplates()
	logAndDie(err)

	var buf bytes.Buffer

	err = packageTemplate.Execute(&buf, struct {
		Timestamp time.Time
		Supported map[string][]string
	}{
		Timestamp: time.Now(),
		Supported: supportedTemplates,
	})
	logAndDie(err)

	formatted, err := format.Source(buf.Bytes())
	logAndDie(err)

	file, err := os.Create("supported.go")
	logAndDie(err)

	_, err = file.Write(formatted)
	logAndDie(err)

	defer file.Close()
}

func logAndDie(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Supported
func getSupportedLanguageTemplates() (map[string][]string, error) {

	items, err := os.ReadDir(supportedTemplatesDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to open directory `%v`: %w", supportedTemplatesDirectory, err)
	}

	supported := make(map[string][]string)
	for _, item := range items {
		// First level is languages
		if item.IsDir() {
			// second level contains variants
			language := item.Name()
			variants, err := getSubdirectoryNames(language)
			if err != nil {
				return nil, fmt.Errorf("failed to collect supported variants for `%v`: %w", language, err)
			}
			supported[language] = variants
		}
	}

	return supported, nil
}

func getSubdirectoryNames(directory string) ([]string, error) {
	path := filepath.Join(supportedTemplatesDirectory, directory)
	items, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open subdirectory `%v`: %w", path, err)
	}

	var names []string

	for _, item := range items {
		if item.IsDir() {
			names = append(names, item.Name())
		}
	}
	return names, nil
}

var packageTemplate = template.Must(template.New("").Parse(`
// Code generated by go generate; DO NOT EDIT.
//This file was generated at 
// {{ .Timestamp }}

package templates

import "embed"

//go:embed supported_templates/*
var templateFS embed.FS

var templateFSRoot string = "supported_templates"

var SupportedTemplates = map[string][]string{
	{{ range $language, $variants := .Supported }}"{{ $language }}":{
	{{- range $variant := $variants}}
		"{{ $variant}}",{{ end }}
	},
	{{ end }}
}
`))
