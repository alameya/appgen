package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/iancoleman/strcase"
)

func toCamelCase(s string) string {
	return strcase.ToCamel(s)
}

type TemplateGenerator struct {
	templatesPath string
	templates     *template.Template
}

func NewTemplateGenerator() *TemplateGenerator {
	// Определяем путь к шаблонам
	_, filename, _, _ := runtime.Caller(0)
	templatesPath := filepath.Join(filepath.Dir(filename), "..", "templates")
	log.Printf("Templates path: %s", templatesPath)

	// Создаем FuncMap с пользовательскими функциями
	funcMap := template.FuncMap{
		"toLower":      strings.ToLower,
		"toCamel":      strcase.ToCamel,
		"toLowerCamel": strcase.ToLowerCamel,
		"hasSuffix":    strings.HasSuffix,
		"trimSuffix":   strings.TrimSuffix,
		"idx":          func(i int) int { return i + 1 },
	}

	// Загружаем все шаблоны
	tmpl := template.New("").Funcs(funcMap)
	pattern := filepath.Join(templatesPath, "*.tmpl")
	tmpl, err := tmpl.ParseGlob(pattern)
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	return &TemplateGenerator{
		templatesPath: templatesPath,
		templates:     tmpl,
	}
}

func (t *TemplateGenerator) generateFromTemplateWithVars(templateName, outputPath string, vars map[string]interface{}, data interface{}) error {
	// Создаем директорию для файла если её нет
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Открываем файл для записи
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Получаем шаблон и выполняем его
	tmpl := t.templates.Lookup(templateName)
	if tmpl == nil {
		return fmt.Errorf("template %s not found", templateName)
	}

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
