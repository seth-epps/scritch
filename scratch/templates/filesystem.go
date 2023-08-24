package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type FilesystemTemplateProvider struct {
	sourcepath string
}

func NewFilesystemTemplateProvider(sourcepath string) FilesystemTemplateProvider {
	return FilesystemTemplateProvider{sourcepath: sourcepath}
}

func (fp FilesystemTemplateProvider) Get() ([]Template, error) {
	var templates []Template

	// Use filepath to ensure we capture nested templates in the structure
	err := filepath.Walk(fp.sourcepath, func(path string, info os.FileInfo, err error) error {
		if path == fp.sourcepath || info.IsDir() && filepath.Ext(path) != "tpl" {
			return nil
		}

		// error can be safely ignored because
		// sourcepath is by construction equal to basepath
		tmplPath, _ := filepath.Rel(fp.sourcepath, path)
		relPath, _ := strings.CutSuffix(tmplPath, ".tpl")
		tmpls := template.Must(
			template.ParseFiles(path),
		)

		for _, tmpl := range tmpls.Templates() {
			templates = append(templates, Template{
				RelativePath: relPath,
				GoTemplate:   tmpl,
			})
		}
		return nil
	})

	if err != nil {
		return templates, fmt.Errorf("Could not find navigate path `%v`: %w", fp.sourcepath, err)
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("Could not find templates to render at path `%v`", fp.sourcepath)
	}

	return templates, nil
}

func (fp FilesystemTemplateProvider) TargetPath() string {
	return filepath.Base(fp.sourcepath)
}
