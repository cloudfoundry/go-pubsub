package inspector

import (
	"go/ast"
)

type Field struct {
	Name string
	Type string
	Ptr  bool
}

type Struct struct {
	Name                string
	Fields              []Field
	PeerTypeFields      []Field
	InterfaceTypeFields map[Field][]string
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
			name, ptr := f.extractType(x.Type)
			f := Field{
				Name: f.firstName(x.Names),
				Type: name,
				Ptr:  ptr,
			}

			if f.Name != "" || f.Type != "" {
				fields = append(fields, f)
			}
		}
		return true
	})
	return fields
}

func (f StructFetcher) firstName(names []*ast.Ident) string {
	if len(names) == 0 {
		return ""
	}

	return names[0].Name
}

func (f StructFetcher) extractType(n ast.Node) (string, bool) {
	switch x := n.(type) {
	case *ast.Ident:
		return x.Name, false
	case *ast.StarExpr:
		if n, ok := x.X.(*ast.Ident); ok {
			return n.Name, true
		}
	}

	return "", false
}
