package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ProtoGen struct {
	sourceDir string
	outputDir string
}

const commonProtoTemplate = `syntax = "proto3";

package proto;

option go_package = "app/internal/proto";

// EmptyResponse message used for operations that don't return data
message EmptyResponse {}`

const serviceProtoTemplate = `syntax = "proto3";

package proto;

import "google/api/annotations.proto";
import "common.proto";

option go_package = "app/internal/proto";

// {{.ServiceName}} message
message {{.ServiceName}} {
{{- range $line := .MessageFields}}
  {{$line}}
{{- end}}
}

// Create request
message Create{{.ServiceName}}Request {
  {{.ServiceName}} {{toLower .ServiceName}} = 1;
}

// Get request
message Get{{.ServiceName}}Request {
  int64 id = 1;
}

// List request
message List{{.ServiceName}}Request {}

// List response
message List{{.ServiceName}}Response {
  repeated {{.ServiceName}} items = 1;
}

// Update request
message Update{{.ServiceName}}Request {
  int64 id = 1;
  {{.ServiceName}} {{toLower .ServiceName}} = 2;
}

// Delete request
message Delete{{.ServiceName}}Request {
  int64 id = 1;
}

// {{.ServiceName}} service definition
service {{.ServiceName}}Service {
  // Create a new {{.ServiceNameLower}}
  rpc Create(Create{{.ServiceName}}Request) returns ({{.ServiceName}}) {
    option (google.api.http) = {
      post: "/api/v1/{{.ServiceNamePlural}}"
      body: "{{toLower .ServiceName}}"
    };
  }

  // Get {{.ServiceNameLower}} by ID
  rpc Get(Get{{.ServiceName}}Request) returns ({{.ServiceName}}) {
    option (google.api.http) = {
      get: "/api/v1/{{.ServiceNamePlural}}/{id}"
    };
  }

  // List all {{.ServiceNamePlural}}
  rpc List(List{{.ServiceName}}Request) returns (List{{.ServiceName}}Response) {
    option (google.api.http) = {
      get: "/api/v1/{{.ServiceNamePlural}}"
    };
  }

  // Update {{.ServiceNameLower}}
  rpc Update(Update{{.ServiceName}}Request) returns ({{.ServiceName}}) {
    option (google.api.http) = {
      put: "/api/v1/{{.ServiceNamePlural}}/{id}"
      body: "{{toLower .ServiceName}}"
    };
  }

  // Delete {{.ServiceNameLower}}
  rpc Delete(Delete{{.ServiceName}}Request) returns (EmptyResponse) {
    option (google.api.http) = {
      delete: "/api/v1/{{.ServiceNamePlural}}/{id}"
    };
  }
}`

type ServiceData struct {
	ServiceName       string
	ServiceNameLower  string
	ServiceNamePlural string
	MessageFields     []string
}

func NewProtoGen(sourceDir, outputDir string) *ProtoGen {
	return &ProtoGen{
		sourceDir: sourceDir,
		outputDir: outputDir,
	}
}

func (g *ProtoGen) Generate() error {
	// Create output directory
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate common.proto
	if err := g.generateCommonProto(); err != nil {
		return fmt.Errorf("failed to generate common.proto: %w", err)
	}

	// Process all proto files in source directory
	files, err := filepath.Glob(filepath.Join(g.sourceDir, "*.proto"))
	if err != nil {
		return fmt.Errorf("failed to list proto files: %w", err)
	}

	for _, file := range files {
		if err := g.generateServiceProto(file); err != nil {
			return fmt.Errorf("failed to generate service proto for %s: %w", file, err)
		}
	}

	return nil
}

func (g *ProtoGen) generateCommonProto() error {
	outPath := filepath.Join(g.outputDir, "common.proto")
	return ioutil.WriteFile(outPath, []byte(commonProtoTemplate), 0644)
}

func (g *ProtoGen) generateServiceProto(sourcePath string) error {
	// Read source file
	content, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Extract service name from filename
	baseName := filepath.Base(sourcePath)
	serviceName := strings.TrimSuffix(baseName, ".proto")
	serviceName = strings.Title(serviceName)

	// Extract message definition and clean it up
	messageContent := string(content)
	messageContent = strings.ReplaceAll(messageContent, "\r\n", "\n") // Normalize line endings
	messageContent = strings.TrimSpace(messageContent)

	// Format the message definition
	var messageFields []string
	lines := strings.Split(messageContent, "\n")
	inMessage := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "message") {
			inMessage = true
			continue
		}
		if line == "}" {
			inMessage = false
			continue
		}
		if inMessage && line != "" {
			if !strings.Contains(line, "message") && !strings.Contains(line, "}") {
				messageFields = append(messageFields, line)
			}
		}
	}

	data := ServiceData{
		ServiceName:       serviceName,
		ServiceNameLower:  strings.ToLower(serviceName),
		ServiceNamePlural: strings.ToLower(serviceName) + "s",
		MessageFields:     messageFields,
	}

	// Создаем шаблон с нашими вспомогательными функциями
	tmpl := template.New("service")
	tmpl = tmpl.Funcs(template.FuncMap{
		"splitLines": splitLines,
		"contains":   contains,
		"trimPrefix": trimPrefix,
		"toLower":    strings.ToLower,
	})

	// Парсим шаблон
	tmpl, err = tmpl.Parse(serviceProtoTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	outPath := filepath.Join(g.outputDir, baseName)
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if err := tmpl.Execute(outFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// Добавим вспомогательные функции для шаблона
func splitLines(s string) []string {
	return strings.Split(s, "\n")
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Добавляем новую вспомогательную функцию для шаблона
func trimPrefix(s string, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

func main() {
	sourceDir := flag.String("source", "proto", "Source directory containing proto files")
	outputDir := flag.String("output", "out/internal/proto", "Output directory for generated proto files")
	flag.Parse()

	generator := NewProtoGen(*sourceDir, *outputDir)
	if err := generator.Generate(); err != nil {
		log.Fatalf("Failed to generate proto files: %v", err)
	}
}
