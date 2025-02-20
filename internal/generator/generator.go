package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CloudyKit/jet/v6"
)

type Generator struct {
	parser   *Parser
	template *TemplateGenerator
}

func New() *Generator {
	return &Generator{
		parser:   NewParser(),
		template: NewTemplateGenerator(),
	}
}

func (g *Generator) GenerateFromProtoFiles(protoFiles []string, outputDir string) error {
	var allModels []*Model

	// Собираем модели из всех файлов
	for _, protoPath := range protoFiles {
		models, err := g.parser.Parse(protoPath)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", protoPath, err)
		}
		allModels = append(allModels, models...)
	}

	// Генерируем общие файлы
	if err := g.generateCommonFiles(allModels, outputDir); err != nil {
		return err
	}

	// Генерируем файлы для каждой модели
	for i, model := range allModels {
		if err := g.template.generateFilesForModel(model, outputDir, i); err != nil {
			return fmt.Errorf("failed to generate files for model %s: %w", model.Name, err)
		}
	}

	// Generate gRPC integration tests
	if err := os.MkdirAll(filepath.Join(outputDir, "internal", "integration"), 0755); err != nil {
		return fmt.Errorf("failed to create integration test directory: %w", err)
	}

	vars := make(jet.VarMap)
	outPath := filepath.Join(outputDir, "internal", "integration", "grpc_test.go")
	if err := g.template.generateFromTemplateWithVars("grpc_test.go.tmpl", outPath, vars, allModels); err != nil {
		return fmt.Errorf("failed to generate gRPC integration tests: %w", err)
	}

	return nil
}

func (g *Generator) GenerateFromProto(protoPath, outputDir string) error {
	return g.GenerateFromProtoFiles([]string{protoPath}, outputDir)
}

func (g *Generator) generateCommonFiles(models []*Model, outputDir string) error {
	commonFiles := map[string]string{
		"main.go.tmpl":            filepath.Join(outputDir, "cmd", "app", "main.go"),
		"go.mod.tmpl":             filepath.Join(outputDir, "go.mod"),
		"docker-compose.yml.tmpl": filepath.Join(outputDir, "docker-compose.yml"),
		"Dockerfile.tmpl":         filepath.Join(outputDir, "Dockerfile"),
		"env.tmpl":                filepath.Join(outputDir, ".env"),
		"repository.go.tmpl":      filepath.Join(outputDir, "internal", "repository", "repository.go"),
	}

	// Генерируем каждый файл из шаблона
	for tmpl, outPath := range commonFiles {
		// Создаем необходимые директории
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", outPath, err)
		}

		// Генерируем файл из шаблона
		vars := make(jet.VarMap)
		if err := g.template.generateFromTemplateWithVars(tmpl, outPath, vars, models); err != nil {
			return fmt.Errorf("failed to generate %s: %w", outPath, err)
		}
	}

	return nil
}
