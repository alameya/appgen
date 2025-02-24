package main

import (
	"flag"
	"log"

	"generator/internal/generator"
)

func main() {
	protoPath := flag.String("proto", "", "Path to proto files (supports glob patterns)")
	outputDir := flag.String("output", "output", "Output directory")
	flag.Parse()

	if *protoPath == "" {
		flag.Usage()
		return
	}

	g := generator.New()
	if err := g.GenerateFromProto(*protoPath, *outputDir); err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}
}
