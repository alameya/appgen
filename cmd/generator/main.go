package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"generator/internal/generator"
)

func main() {
	protoPath := flag.String("proto", "", "Path to proto files (supports glob patterns)")
	outputDir := flag.String("output", "output", "Output directory")
	flag.Parse()

	if *protoPath == "" {
		log.Fatal("Proto file path is required")
	}

	// Получаем список файлов по маске
	protoFiles, err := filepath.Glob(*protoPath)
	if err != nil {
		log.Fatalf("Failed to find proto files: %v", err)
	}

	if len(protoFiles) == 0 {
		log.Fatalf("No proto files found matching pattern: %s", *protoPath)
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	gen := generator.New()
	if err := gen.GenerateFromProtoFiles(protoFiles, *outputDir); err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}
}
