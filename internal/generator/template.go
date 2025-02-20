package generator

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/iancoleman/strcase"
)

func toCamelCase(s string) string {
	return strcase.ToCamel(s)
}

type TemplateGenerator struct {
	templates *jet.Set
}

func NewTemplateGenerator() *TemplateGenerator {
	// Get absolute path to templates directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	templatesPath := filepath.Join(dir, "internal", "templates")
	log.Printf("Templates path: %s", templatesPath)

	views := jet.NewSet(
		jet.NewOSFileSystemLoader(templatesPath),
		jet.InDevelopmentMode(),
	)

	// Add custom functions
	views.AddGlobal("toLower", strings.ToLower)
	views.AddGlobal("toUpper", strings.ToUpper)
	views.AddGlobal("toCamel", toCamelCase)

	views.AddGlobalFunc("toLowerCamel", func(args jet.Arguments) reflect.Value {
		return reflect.ValueOf(strcase.ToLowerCamel(args.Get(0).String()))
	})
	views.AddGlobalFunc("hasSuffix", func(args jet.Arguments) reflect.Value {
		return reflect.ValueOf(strings.HasSuffix(args.Get(0).String(), args.Get(1).String()))
	})

	views.AddGlobal("hasSuffix", strings.HasSuffix)
	views.AddGlobal("trimSuffix", strings.TrimSuffix)
	views.AddGlobal("existsTable", func(name string) bool {
		return name != "" && name != "." && name != ".."
	})

	return &TemplateGenerator{
		templates: views,
	}
}

func (t *TemplateGenerator) generateFilesForModel(model *Model, outputDir string, modelIndex int) error {
	// Создаем директории для модели
	dirs := []string{
		filepath.Join(outputDir, "internal", "handler", strings.ToLower(model.Name)),
		filepath.Join(outputDir, "internal", "service", strings.ToLower(model.Name)),
		filepath.Join(outputDir, "internal", "repository", strings.ToLower(model.Name)),
		filepath.Join(outputDir, "internal", "models"),
		filepath.Join(outputDir, "internal", "grpc", strings.ToLower(model.Name)),
		filepath.Join(outputDir, "migrations"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Получаем текущее время для версии миграции
	now := time.Now()
	version := now.Format("20060102150405")     // YYYYMMDDHHMMSS
	suffix := fmt.Sprintf("%02d", modelIndex+1) // 01, 02, 03, etc.

	// Генерируем файлы для модели
	files := map[string]string{
		"repository_model.go.tmpl": filepath.Join(outputDir, "internal", "repository", strings.ToLower(model.Name), "repository.go"),
		"handler.go.tmpl":          filepath.Join(outputDir, "internal", "handler", strings.ToLower(model.Name), "handler.go"),
		"service.go.tmpl":          filepath.Join(outputDir, "internal", "service", strings.ToLower(model.Name), "service.go"),
		"models.go.tmpl":           filepath.Join(outputDir, "internal", "models", strings.ToLower(model.Name)+".go"),
		"migration.sql.tmpl":       filepath.Join(outputDir, "migrations", fmt.Sprintf("%s%s_create_%s.sql", version, suffix, strings.ToLower(model.Name))),
		"grpc.go.tmpl":             filepath.Join(outputDir, "internal", "grpc", strings.ToLower(model.Name), "server.go"),
	}

	// Генерация файлов для модели
	for name, path := range files {
		vars := make(jet.VarMap)
		if err := t.generateFromTemplateWithVars(name, path, vars, model); err != nil {
			return fmt.Errorf("failed to generate %s: %w", path, err)
		}
	}

	return nil
}

func (t *TemplateGenerator) generateFromTemplateWithVars(templateName, filePath string, vars jet.VarMap, data interface{}) error {
	tmpl, err := t.templates.GetTemplate(templateName)
	if err != nil {
		return fmt.Errorf("template %s not found: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	log.Printf("Generating file: %s", filePath)
	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}
