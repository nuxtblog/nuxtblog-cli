// Package scaffold generates new plugin projects from templates.
package scaffold

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Options configures plugin scaffolding.
type Options struct {
	Type        string // yaml, js, go, full
	ID          string // e.g. nuxtblog-plugin-my-thing
	Name        string // display name
	Description string
	Author      string
	OutputDir   string
}

// TemplateData is passed to all templates.
type TemplateData struct {
	ID          string
	Name        string
	Description string
	Author      string
	PkgName     string // Go package name derived from ID
	ModulePath  string // Go module path
}

// Generate creates a new plugin from templates.
func Generate(opts Options) error {
	data := TemplateData{
		ID:          opts.ID,
		Name:        opts.Name,
		Description: opts.Description,
		Author:      opts.Author,
		PkgName:     toPkgName(opts.ID),
		ModulePath:  "github.com/nuxtblog/nuxtblog/plugins/" + opts.ID,
	}

	switch opts.Type {
	case "yaml":
		return renderTemplates(yamlTemplates, "templates/yaml", opts.OutputDir, data)
	case "js":
		return renderTemplates(jsTemplates, "templates/js", opts.OutputDir, data)
	case "go":
		return renderTemplates(goTemplates, "templates/go", opts.OutputDir, data)
	case "full":
		return renderTemplates(fullTemplates, "templates/full", opts.OutputDir, data)
	default:
		return fmt.Errorf("unknown plugin type: %s", opts.Type)
	}
}

func renderTemplates(fs embed.FS, dir, outputDir string, data TemplateData) error {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read templates/%s: %w", dir, err)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			// Handle subdirectories (e.g. src/, web/)
			subEntries, err := fs.ReadDir(filepath.Join(dir, e.Name()))
			if err != nil {
				return err
			}
			subDir := filepath.Join(outputDir, e.Name())
			if err := os.MkdirAll(subDir, 0o755); err != nil {
				return err
			}
			for _, se := range subEntries {
				if err := renderOneFile(fs, filepath.Join(dir, e.Name(), se.Name()), subDir, data); err != nil {
					return err
				}
			}
			continue
		}
		if err := renderOneFile(fs, filepath.Join(dir, e.Name()), outputDir, data); err != nil {
			return err
		}
	}
	return nil
}

func renderOneFile(fs embed.FS, tmplPath, outDir string, data TemplateData) error {
	content, err := fs.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", tmplPath, err)
	}

	// Strip .tmpl extension for output filename
	outName := filepath.Base(tmplPath)
	outName = strings.TrimSuffix(outName, ".tmpl")

	tmpl, err := template.New(outName).Parse(string(content))
	if err != nil {
		return fmt.Errorf("parse %s: %w", tmplPath, err)
	}

	outPath := filepath.Join(outDir, outName)
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

// toPkgName converts a plugin ID like "nuxtblog-plugin-my-thing" to "mything".
func toPkgName(id string) string {
	name := strings.TrimPrefix(id, "nuxtblog-plugin-")
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	return strings.ToLower(name)
}
