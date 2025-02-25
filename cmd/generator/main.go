package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"generator/internal/generator"
)

func main() {
	protoPath := flag.String("proto", "", "Path to proto files (supports comma-separated list or glob pattern)")
	outputDir := flag.String("output", "out", "Output directory")
	flag.Parse()

	if *protoPath == "" {
		flag.Usage()
		return
	}

	// Обрабатываем список файлов, разделенных запятыми
	var protoFiles []string
	if strings.Contains(*protoPath, ",") {
		protoFiles = strings.Split(*protoPath, ",")
		// Фильтруем пустые строки
		var filteredFiles []string
		for _, file := range protoFiles {
			if file != "" {
				filteredFiles = append(filteredFiles, file)
			}
		}
		protoFiles = filteredFiles
	} else {
		// Если это не список через запятые, обрабатываем как glob
		matches, err := filepath.Glob(*protoPath)
		if err != nil {
			log.Fatalf("Failed to parse glob pattern: %v", err)
		}
		protoFiles = matches
	}

	fmt.Printf("Processing proto files: %v\n", protoFiles)

	g := generator.New()
	if err := g.GenerateFromProtoFiles(protoFiles, *outputDir); err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}
}
