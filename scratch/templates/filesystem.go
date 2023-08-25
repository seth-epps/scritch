package templates

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/seth-epps/scritch/scratch/util"
)

type FilesystemTemplateProvider struct {
	sourcepath  string
	searchPaths []string
}

func NewFilesystemTemplateProvider(sourcepath string, searchPaths []string) FilesystemTemplateProvider {
	if len(searchPaths) == 0 {
		fmt.Println("WARNING: Should provide search paths in case of relative sourcepath")
		return FilesystemTemplateProvider{sourcepath: sourcepath}
	}
	return FilesystemTemplateProvider{sourcepath: sourcepath, searchPaths: searchPaths}
}

func (fp FilesystemTemplateProvider) Get() ([]Template, error) {
	var templates []Template

	templateLocation, err := fp.findTemplateLocation()
	if err != nil {
		return templates, err
	}

	// Use filepath to ensure we capture nested templates in the structure
	err = filepath.Walk(templateLocation, func(path string, info os.FileInfo, err error) error {
		if path == templateLocation || info.IsDir() && filepath.Ext(path) != "tpl" {
			return nil
		}

		// error can be safely ignored because
		// sourcepath is by construction equal to basepath
		tmplPath, _ := filepath.Rel(templateLocation, path)
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
		return templates, fmt.Errorf("Could not navigate path `%v`: %w", templateLocation, err)
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("Could not find templates to render at path `%v`", templateLocation)
	}

	return templates, nil
}

func (fp FilesystemTemplateProvider) TargetPath() string {
	return filepath.Base(fp.sourcepath)
}

func (fp FilesystemTemplateProvider) findTemplateLocation() (string, error) {
	path, err := util.ReplaceHomeShortcut(fp.sourcepath)
	if err != nil {
		return fp.sourcepath, fmt.Errorf("could not resolve source: %w", err)
	}
	if filepath.IsAbs(path) {
		return path, nil
	}

	return searchForSource(path, fp.searchPaths)
}

func searchForSource(source string, locations []string) (string, error) {
	var err error
	for _, searchPath := range locations {
		potentialPath := filepath.Join(searchPath, source)
		info, statErr := os.Stat(potentialPath)
		if statErr != nil {
			errors.Join(err, statErr)
			continue
		}
		if info != nil && info.IsDir() {
			return potentialPath, nil
		}
	}
	return "", errors.Join(errors.New("could not find source"), err)
}
