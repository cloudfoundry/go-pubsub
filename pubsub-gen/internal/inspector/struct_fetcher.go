package inspector

import (
	"go/ast"
)

type Field struct {
	Name string
	Type string
	// Ptr  bool
}

type Struct struct {
	Name   string
	Fields []Field
}

type StructFetcher struct{}

func NewStructFetcher() StructFetcher {
	return StructFetcher{}
}

func (f StructFetcher) Parse(n ast.Node) ([]Struct, error) {
	var structs []Struct
	var name string

	ast.Inspect(n, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			name = x.Name
		case *ast.StructType:
			fields := f.extractFields(x.Fields)
			structs = append(structs, Struct{Name: name, Fields: fields})
		}
		return true
	})

	return structs, nil
}

func (f StructFetcher) extractFields(n ast.Node) []Field {
	var fields []Field
	ast.Inspect(n, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Field:
			f := Field{
				Name: x.Names[0].Name,
				Type: f.extractType(x.Type),
			}

			if f.Name != "" || f.Type != "" {
				fields = append(fields, f)
			}
		}
		return true
	})
	return fields
}

func (f StructFetcher) extractType(n ast.Node) string {
	switch x := n.(type) {
	case *ast.Ident:
		return x.Name
	}

	return ""
}
