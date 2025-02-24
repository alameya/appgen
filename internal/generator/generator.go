package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// GenerateFromProto генерирует код из одного proto файла
func (g *Generator) GenerateFromProto(protoPath, outputDir string) error {
	return g.GenerateFromProtoFiles([]string{protoPath}, outputDir)
}

// GenerateFromProtoFiles генерирует код из нескольких proto файлов
func (g *Generator) GenerateFromProtoFiles(protoFiles []string, outputDir string) error {
	var allModels []*Model

	for _, protoPath := range protoFiles {
		models, err := g.parser.Parse(protoPath)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", protoPath, err)
		}
		allModels = append(allModels, models...)
	}

	// Сортируем модели по зависимостям
	sortedModels := g.sortModelsByDependencies(allModels)

	if err := g.generateCommonFiles(allModels, outputDir); err != nil {
		return err
	}

	// Сначала генерируем все файлы кроме миграций
	for i, model := range sortedModels {
		if err := g.generateFilesForModel(model, outputDir, i); err != nil {
			return fmt.Errorf("failed to generate files for model %s: %w", model.Name, err)
		}
	}

	// Затем генерируем миграции в том же порядке что и модели
	for i, model := range sortedModels {
		if err := g.generateMigration(model, outputDir, i); err != nil {
			return fmt.Errorf("failed to generate migration for model %s: %w", model.Name, err)
		}
		// Добавляем задержку между генерацией миграций
		time.Sleep(time.Second)
	}

	return nil
}

func (g *Generator) generateCommonFiles(models []*Model, outputDir string) error {
	commonFiles := map[string]string{
		"main.go.tmpl":            filepath.Join(outputDir, "cmd", "app", "main.go"),
		"go.mod.tmpl":             filepath.Join(outputDir, "go.mod"),
		"docker-compose.yml.tmpl": filepath.Join(outputDir, "docker-compose.yml"),
		"Dockerfile.tmpl":         filepath.Join(outputDir, "Dockerfile"),
		"env.tmpl":                filepath.Join(outputDir, ".env"),
		"repository.go.tmpl":      filepath.Join(outputDir, "internal", "repository", "repository.go"),
		"interfaces.go.tmpl":      filepath.Join(outputDir, "internal", "interfaces", "interfaces.go"),
		"gitlab-ci.yml.tmpl":      filepath.Join(outputDir, ".gitlab-ci.yml"),
		"grpc_test.go.tmpl":       filepath.Join(outputDir, "internal", "tests", "grpc_test.go"),
	}

	for tmpl, outPath := range commonFiles {
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", outPath, err)
		}

		if err := g.template.generateFromTemplateWithVars(tmpl, outPath, nil, models); err != nil {
			return fmt.Errorf("failed to generate %s: %w", outPath, err)
		}
	}

	return nil
}

func (g *Generator) generateFilesForModel(model *Model, outputDir string, modelIndex int) error {
	files := map[string]string{
		"repository_model.go.tmpl": filepath.Join(outputDir, "internal", "repository", strings.ToLower(model.Name), "repository.go"),
		"handler.go.tmpl":          filepath.Join(outputDir, "internal", "handler", strings.ToLower(model.Name), "handler.go"),
		"models.go.tmpl":           filepath.Join(outputDir, "internal", "models", strings.ToLower(model.Name)+".go"),
		"grpc.go.tmpl":             filepath.Join(outputDir, "internal", "grpc", strings.ToLower(model.Name), "server.go"),
	}

	// Генерируем файл сервиса
	servicePath := filepath.Join(outputDir, "internal", "service")
	if err := os.MkdirAll(servicePath, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	if err := g.template.generateFromTemplateWithVars("service.go.tmpl",
		filepath.Join(servicePath, strings.ToLower(model.Name)+".go"), nil, model); err != nil {
		return fmt.Errorf("failed to generate service file: %w", err)
	}

	for tmpl, outPath := range files {
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", outPath, err)
		}

		if err := g.template.generateFromTemplateWithVars(tmpl, outPath, nil, model); err != nil {
			return fmt.Errorf("failed to generate %s: %w", outPath, err)
		}
	}

	return nil
}

// buildDependencyGraph строит граф зависимостей между моделями
func (g *Generator) buildDependencyGraph(models []*Model) map[string][]string {
	dependencies := make(map[string][]string)
	for _, model := range models {
		deps := []string{}
		for _, field := range model.Fields {
			if strings.HasSuffix(field.Name, "_id") {
				referencedModel := strings.TrimSuffix(field.Name, "_id")
				deps = append(deps, referencedModel)
			}
		}
		dependencies[model.Name] = deps
	}
	return dependencies
}

// sortModelsByDependencies сортирует модели так, чтобы зависимые таблицы создавались после зависимостей
func (g *Generator) sortModelsByDependencies(models []*Model) []*Model {
	graph := g.buildDependencyGraph(models)
	visited := make(map[string]bool)
	visiting := make(map[string]bool) // Для обнаружения циклических зависимостей
	sorted := make([]*Model, 0)

	// Выводим дерево зависимостей
	fmt.Println("\nDependency tree:")
	g.printDependencyTree(graph, "", "", make(map[string]bool))
	fmt.Println()

	var visit func(model *Model)
	visit = func(model *Model) {
		if visited[model.Name] {
			return
		}
		if visiting[model.Name] {
			panic(fmt.Sprintf("circular dependency detected: %s", model.Name))
		}

		visiting[model.Name] = true

		for _, dep := range graph[model.Name] {
			for _, m := range models {
				if m.Name == dep {
					visit(m)
				}
			}
		}
		visiting[model.Name] = false
		visited[model.Name] = true
		sorted = append([]*Model{model}, sorted...) // Добавляем в начало списка
	}

	for _, model := range models {
		if !visited[model.Name] {
			visit(model)
		}
	}

	// Выводим порядок генерации
	fmt.Println("Generation order:")
	for i, model := range sorted {
		fmt.Printf("%d. %s\n", i+1, model.Name)
	}
	fmt.Println()

	return sorted
}

// printDependencyTree выводит дерево зависимостей в консоль
func (g *Generator) printDependencyTree(graph map[string][]string, prefix string, name string, visited map[string]bool) {
	if name == "" {
		// Находим корневые узлы (без зависимостей)
		hasIncoming := make(map[string]bool)
		for _, deps := range graph {
			for _, dep := range deps {
				hasIncoming[dep] = true
			}
		}

		for node := range graph {
			if !hasIncoming[node] {
				g.printDependencyTree(graph, "", node, visited)
			}
		}
		return
	}

	if visited[name] {
		fmt.Printf("%s%s (circular)\n", prefix, name)
		return
	}

	visited[name] = true
	fmt.Printf("%s%s\n", prefix, name)

	for i, dep := range graph[name] {
		nextPrefix := prefix + "│   "
		if i == len(graph[name])-1 {
			nextPrefix = prefix + "    "
		}
		g.printDependencyTree(graph, nextPrefix, dep, visited)
	}

	visited[name] = false
}

func (g *Generator) generateMigration(model *Model, outputDir string, index int) error {
	now := time.Now()
	// Используем базовый timestamp и добавляем индекс для сортировки
	baseVersion := now.Format("20060102150405")
	// Добавляем индекс в начало версии для правильной сортировки зависимостей
	version := fmt.Sprintf("%s%02d", baseVersion, index+1)
	migrationPath := filepath.Join(outputDir, "migrations")

	if err := os.MkdirAll(migrationPath, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	filename := fmt.Sprintf("%s_create_%s.sql", version, strings.ToLower(model.Name))
	fullPath := filepath.Join(migrationPath, filename)

	if err := g.template.generateFromTemplateWithVars("migration.sql.tmpl", fullPath, nil, model); err != nil {
		return fmt.Errorf("failed to generate migration file: %w", err)
	}

	return nil
}
