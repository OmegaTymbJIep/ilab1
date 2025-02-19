package views

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

const (
	AuthTemplateName = "auth.html"
)

func ReadTemplates(templatesDir string) (*template.Template, error) {
	var templatePaths []string
	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk templates dir: %w", err)
		}

		if !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			templatePaths = append(templatePaths, absPath)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read templates dir: %w", err)
	}

	return template.ParseFiles(templatePaths...)
}
