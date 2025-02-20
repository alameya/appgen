package generator

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(protoPath string) ([]*Model, error) {
	// Компилируем proto файл
	cmd := exec.Command("protoc",
		"--descriptor_set_out=/tmp/proto.pb",
		"--include_imports",
		protoPath,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to compile proto file: %w, output: %s", err, out)
	}

	// Читаем дескриптор
	descBytes, err := os.ReadFile("/tmp/proto.pb")
	if err != nil {
		return nil, fmt.Errorf("failed to read descriptor: %w", err)
	}

	var fdSet descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(descBytes, &fdSet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal descriptor: %w", err)
	}

	fd, err := protodesc.NewFiles(&fdSet)
	if err != nil {
		return nil, fmt.Errorf("failed to create file descriptor: %w", err)
	}

	desc, err := fd.FindFileByPath(protoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find proto file: %w", err)
	}

	var models []*Model

	// Parse messages
	messages := desc.Messages()
	for i := 0; i < messages.Len(); i++ {
		message := messages.Get(i)
		name := string(message.Name())

		fmt.Printf("Parsing message: %s\n", name)

		// Пропускаем сообщения запросов и ответов
		if strings.HasSuffix(name, "Request") || strings.HasSuffix(name, "Response") {
			continue
		}

		model := &Model{
			Name:   name,
			Fields: make([]Field, 0, message.Fields().Len()),
		}

		fields := message.Fields()
		for j := 0; j < fields.Len(); j++ {
			field := fields.Get(j)
			f := Field{
				Name:     string(field.Name()),
				Type:     getGoType(field),
				JsonName: field.JSONName(),
				DbName:   string(field.Name()),
				SqlType:  getSqlType(field.Kind().String()),
				Last:     j == fields.Len()-1,
			}
			model.Fields = append(model.Fields, f)
		}

		models = append(models, model)
	}

	return models, nil
}

func getGoType(field protoreflect.FieldDescriptor) string {
	switch field.Kind() {
	case protoreflect.Int64Kind:
		return "int64"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.DoubleKind:
		return "float64"
	case protoreflect.MessageKind:
		return "*" + string(field.Message().Name())
	default:
		return "string"
	}
}

func getValidations(field protoreflect.FieldDescriptor) []string {
	var validations []string

	// Parse field options/annotations for validations
	// Example: required, min_length, max_length, etc.

	return validations
}

func getSqlType(protoType string) string {
	switch protoType {
	case "INT64":
		return "BIGINT"
	case "STRING":
		return "TEXT"
	case "BOOL":
		return "BOOLEAN"
	case "DOUBLE":
		return "DOUBLE PRECISION"
	default:
		return "TEXT"
	}
}
