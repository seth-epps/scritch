// Code generated by go generate; DO NOT EDIT.
// This file was generated at
// 2024-03-01 13:44:20.142238 -0500 EST m=+0.001087292

package templates

import "embed"

//go:embed supported_templates/*
var templateFS embed.FS

var templateFSRoot string = "supported_templates"

var SupportedTemplates = []string{
	"go",
	"python",
	"rust",
}
