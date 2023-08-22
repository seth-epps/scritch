package scratch

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/seth-epps/scritch/templates"
)

const (
	ScritchDirectoryName = ".scritch"
	DestinationFallback  = "/tmp"
	directoryPermissions = os.ModeDir | 0755
)

func GenerateScratch(language, variant string) (scratchLocation string, err error) {

	tmpls, err := templates.GetNativelySupportedTemplates(language, variant)
	if err != nil {
		return scratchLocation, fmt.Errorf("Failed to retrieve template `%v (%v)`: %w", language, variant, err)
	}

	scratchLocation, err = resolveDestinationPath(language, variant)
	if err != nil {
		return scratchLocation, fmt.Errorf("Could not resolve scratch destination: %w", err)
	}

	for _, tmpl := range tmpls {
		err = errors.Join(err, writeScratchFile(tmpl, scratchLocation))
	}

	return scratchLocation, err
}

func resolveDestinationPath(language, variant string) (string, error) {
	//TODO allow configuration of destination
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Could not resolve $HOME; falling back to %v", DestinationFallback)
		home = DestinationFallback
	}

	generatedUUID := uuid.New()
	target := filepath.Join(home, ScritchDirectoryName, language, variant, generatedUUID.String())

	err = os.MkdirAll(target, directoryPermissions)
	if err != nil {
		return "", fmt.Errorf("Could not create directory `%v` to store scratch files: %w", target, err)
	}
	return target, nil
}

func writeScratchFile(tmpl *template.Template, destinationPath string) error {
	errorFormat := "Could not create file %v: %w"
	targetFilename, _ := strings.CutSuffix(tmpl.Name(), ".tpl")
	filePath := filepath.Join(destinationPath, targetFilename)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}
	defer file.Close()

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return nil
}
