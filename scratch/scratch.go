package scratch

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/seth-epps/scritch/scratch/templates"
)

const (
	scritchDirectoryName = ".scritch"
	destinationFallback  = "/tmp"
	directoryPermissions = os.ModeDir | 0755
)

type Scratch struct {
	destination string
}

func NewScratch(destination string) Scratch {
	return Scratch{
		destination: destination,
	}
}

func (s Scratch) GenerateScratch(templateProvider templates.TemplateProvider) (scratchLocation string, err error) {
	tmpls, err := templateProvider.Get()
	if err != nil {
		return scratchLocation, fmt.Errorf("Failed to retrieve template: %w", err)
	}

	scratchLocation, err = s.createDestination(templateProvider)
	if err != nil {
		return scratchLocation, fmt.Errorf("Could not resolve scratch destination: %w", err)
	}

	for _, tmpl := range tmpls {
		err = errors.Join(err, writeScratchFile(tmpl, scratchLocation))
	}

	return scratchLocation, err
}

func (s Scratch) createDestination(templateProvider templates.TemplateProvider) (string, error) {
	generatedUUID := uuid.New()
	target := filepath.Join(s.destination, templateProvider.TargetPath(), generatedUUID.String())

	if err := mkdirIfNotExist(target); err != nil {
		return "", err
	}

	return target, nil
}

func writeScratchFile(tmpl templates.Template, destinationPath string) error {
	errorFormat := "Could not create scratch file %v: %w"

	dirPath := filepath.Join(destinationPath, filepath.Dir(tmpl.RelativePath))
	if err := mkdirIfNotExist(dirPath); err != nil {
		return err
	}

	filePath := filepath.Join(destinationPath, tmpl.RelativePath)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}
	defer file.Close()

	var buf bytes.Buffer
	err = tmpl.GoTemplate.Execute(&buf, nil)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return nil
}

func mkdirIfNotExist(path string) error {
	if err := os.MkdirAll(path, directoryPermissions); err != nil {
		return fmt.Errorf("Could not create destination directory `%v` to store scratch files: %w", path, err)
	}
	return nil
}
