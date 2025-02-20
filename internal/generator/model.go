package generator

type Model struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name        string
	Type        string
	JsonName    string
	DbName      string
	SqlType     string
	Last        bool
	Validations []string
}
