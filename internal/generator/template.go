package generator

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/iancoleman/strcase"
)

type TemplateGenerator struct {
	templates map[string]*template.Template
}

func NewTemplateGenerator() *TemplateGenerator {
	return &TemplateGenerator{
		templates: make(map[string]*template.Template),
	}
}

func (t *TemplateGenerator) Generate(models []*Model, outputDir string) error {
	// Удаляем старые директории
	dirsToClean := []string{
		filepath.Join(outputDir, "internal", "handler"),
		filepath.Join(outputDir, "internal", "service"),
		filepath.Join(outputDir, "internal", "repository"),
		filepath.Join(outputDir, "internal", "models"),
		filepath.Join(outputDir, "migrations"),
	}

	for _, dir := range dirsToClean {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to clean directory %s: %w", dir, err)
		}
	}

	// Загрузка шаблонов
	files := []string{
		"repository_model.go.tmpl",
		"service.go.tmpl",
		"handler.go.tmpl",
		"models.go.tmpl",
		"migration.sql.tmpl",
		"init.sql.tmpl",
		"main.go.tmpl",
		"go.mod.tmpl",
		"docker-compose.yml.tmpl",
		"Dockerfile.tmpl",
		"postman_collection.json.tmpl",
		"env.tmpl",
		"migrate.sh.tmpl",
	}

	log.Printf("Loading templates: %v", files)
	for _, f := range files {
		tmpl, err := template.New(f).
			Funcs(template.FuncMap{
				"toLower":      strings.ToLower,
				"toCamel":      strcase.ToCamel,
				"toLowerCamel": strcase.ToLowerCamel,
				"hasSuffix":    strings.HasSuffix,
			}).
			ParseFiles(filepath.Join("internal/templates", f))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", f, err)
		}
		t.templates[f] = tmpl
		log.Printf("Loaded template: %s", f)
	}

	// Загружаем шаблон репозитория отдельно
	repoTmpl, err := template.New("repository.go.tmpl").
		Funcs(template.FuncMap{
			"toLower":      strings.ToLower,
			"toCamel":      strcase.ToCamel,
			"toLowerCamel": strcase.ToLowerCamel,
			"hasSuffix":    strings.HasSuffix,
		}).
		ParseFiles(filepath.Join("internal/templates/repository", "repository.go.tmpl"))
	if err != nil {
		return fmt.Errorf("failed to parse repository template: %w", err)
	}
	t.templates["repository.go.tmpl"] = repoTmpl
	log.Printf("Loaded template: repository.go.tmpl")

	// Создаем общие директории
	commonDirs := []string{
		filepath.Join(outputDir, "cmd"),
		filepath.Join(outputDir, "scripts"),
		filepath.Join(outputDir, "migrations"),
		filepath.Join(outputDir, "internal", "repository"),
	}

	for _, dir := range commonDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Получаем текущее время для версий миграций
	now := time.Now()
	baseVersion := now.Format("20060102150405") // YYYYMMDDHHMMSS

	// Обновляем маппинг для общих файлов
	commonFiles := map[string]string{
		"init.sql.tmpl":                filepath.Join(outputDir, "migrations", fmt.Sprintf("%s00_init.sql", baseVersion)),
		"main.go.tmpl":                 filepath.Join(outputDir, "cmd", "main.go"),
		"go.mod.tmpl":                  filepath.Join(outputDir, "go.mod"),
		"docker-compose.yml.tmpl":      filepath.Join(outputDir, "docker-compose.yml"),
		"Dockerfile.tmpl":              filepath.Join(outputDir, "Dockerfile"),
		"postman_collection.json.tmpl": filepath.Join(outputDir, "postman_collection.json"),
		"env.tmpl":                     filepath.Join(outputDir, ".env"),
		"migrate.sh.tmpl":              filepath.Join(outputDir, "scripts", "migrate.sh"),
		"repository.go.tmpl":           filepath.Join(outputDir, "internal", "repository", "repository.go"),
	}

	// Генерируем общие файлы
	for name, path := range commonFiles {
		if tmpl, ok := t.templates[name]; ok {
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, models); err != nil {
				return fmt.Errorf("failed to execute template %s: %w", name, err)
			}

			log.Printf("Generating common file: %s", path)
			if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", path, err)
			}

			if name == "migrate.sh.tmpl" {
				if err := os.Chmod(path, 0755); err != nil {
					return fmt.Errorf("failed to make script executable: %w", err)
				}
			}
		}
	}

	// Генерация файлов для каждой модели
	for i, model := range models {
		if err := t.generateFilesForModel(model, outputDir, i); err != nil {
			return fmt.Errorf("failed to generate files for model %s: %w", model.Name, err)
		}
	}

	return nil
}

func (t *TemplateGenerator) generateFilesForModel(model *Model, outputDir string, modelIndex int) error {
	// Создаем директории для модели
	dirs := []string{
		filepath.Join(outputDir, "internal", "handler", strings.ToLower(model.Name)),
		filepath.Join(outputDir, "internal", "service", strings.ToLower(model.Name)),
		filepath.Join(outputDir, "internal", "repository", strings.ToLower(model.Name)),
		filepath.Join(outputDir, "internal", "models"),
		filepath.Join(outputDir, "migrations"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Получаем текущее время для версии миграции
	now := time.Now()
	version := now.Format("20060102150405") // YYYYMMDDHHMMSS

	// Маппинг шаблонов к путям файлов для конкретной модели
	fileMapping := map[string]string{
		"repository_model.go.tmpl": filepath.Join(outputDir, "internal", "repository", strings.ToLower(model.Name), "repository.go"),
		"handler.go.tmpl":          filepath.Join(outputDir, "internal", "handler", strings.ToLower(model.Name), "handler.go"),
		"service.go.tmpl":          filepath.Join(outputDir, "internal", "service", strings.ToLower(model.Name), "service.go"),
		"models.go.tmpl":           filepath.Join(outputDir, "internal", "models", strings.ToLower(model.Name)+".go"),
		"migration.sql.tmpl":       filepath.Join(outputDir, "migrations", fmt.Sprintf("%s%02d_create_%s.sql", version, modelIndex+1, strings.ToLower(model.Name))),
	}

	// Генерация файлов для модели
	for name, path := range fileMapping {
		if tmpl, ok := t.templates[name]; ok {
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, model); err != nil {
				return fmt.Errorf("failed to execute template %s: %w", name, err)
			}

			log.Printf("Generating file: %s", path)
			if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", path, err)
			}
		}
	}

	return nil
}
