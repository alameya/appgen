package generator

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/types/descriptorpb"
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

	return g.template.Generate(allModels, outputDir)
}

func (g *Generator) GenerateFromProto(protoPath, outputDir string) error {
	return g.GenerateFromProtoFiles([]string{protoPath}, outputDir)
}

func (g *Generator) getGoType(field *descriptorpb.FieldDescriptorProto) string {
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		return "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return "float64"
	default:
		return "string"
	}
}

func (g *Generator) parseField(field *descriptorpb.FieldDescriptorProto) Field {
	name := field.GetName()
	fieldType := field.GetType().String()

	var sqlType string
	switch fieldType {
	case "TYPE_INT64":
		sqlType = "BIGINT"
	case "TYPE_STRING":
		sqlType = "TEXT"
	case "TYPE_BOOL":
		sqlType = "BOOLEAN"
	case "TYPE_DOUBLE":
		sqlType = "DOUBLE PRECISION"
	default:
		sqlType = "TEXT"
	}

	return Field{
		Name:     strcase.ToCamel(name),
		Type:     g.getGoType(field),
		JsonName: name,
		DbName:   name,
		SqlType:  sqlType,
	}
}

func (g *Generator) parseMessage(message *descriptorpb.DescriptorProto) (*Model, error) {
	var fields []Field
	for i, field := range message.GetField() {
		f := g.parseField(field)
		f.Last = i == len(message.GetField())-1
		fields = append(fields, f)
	}

	return &Model{
		Name:   message.GetName(),
		Fields: fields,
	}, nil
}

func (g *Generator) generateForModel(model *Model, outputDir string) error {
	if err := g.template.generateModelFiles(model, outputDir); err != nil {
		return fmt.Errorf("failed to generate files for model %s: %w", model.Name, err)
	}

	return nil
}
