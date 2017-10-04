package inspector

import (
	"fmt"
	"go/ast"
)

type Field struct {
	Name  string
	Type  string
	Ptr   bool
	Slice Slice
	Map   Map
}

type Slice struct {
	IsSlice     bool
	IsBasicType bool
	FieldName   string
}

type Map struct {
	IsMap bool
}

type Struct struct {
	Name                string
	Fields              []Field
	PeerTypeFields      []Field
	InterfaceTypeFields map[Field][]string
}

type StructFetcher struct {
	blacklist  map[string][]string
	knownTypes map[string]string
	sliceTypes map[string]string
}

func NewStructFetcher(blacklist map[string][]string, knownTypes, sliceTypes map[string]string) StructFetcher {
	return StructFetcher{
		blacklist:  blacklist,
		knownTypes: knownTypes,
		sliceTypes: sliceTypes,
	}
}

func (f StructFetcher) Parse(n ast.Node) ([]Struct, error) {
	var structs []Struct
	var name string

	ast.Inspect(n, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			name = x.Name
		case *ast.StructType:
			if name == "" {
				return true
			}

			fields := f.extractFields(name, x.Fields)
			structs = append(structs, Struct{Name: name, Fields: fields})
			name = ""
		}
		return true
	})

	return structs, nil
}

func (f StructFetcher) extractFields(parentName string, n ast.Node) []Field {
	var fields []Field
	ast.Inspect(n, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Field:
			name, ptr, slice, isMap := f.extractType(x.Type)

			var basicSliceType bool
			var sliceFieldName string
			if slice {
				var ok bool
				ok, basicSliceType, sliceFieldName = f.isOKSliceType(name, parentName, f.firstName(x.Names))
				if !ok {
					return true
				}
			}

			if isMap {
				if !f.isBasicType(name) {
					return true
				}
			}

			ff := Field{
				Name: f.firstName(x.Names),
				Type: name,
				Ptr:  ptr,
				Slice: Slice{
					IsSlice:     slice,
					IsBasicType: basicSliceType,
					FieldName:   sliceFieldName,
				},
				Map: Map{
					IsMap: isMap,
				},
			}

			if ff.Name != "" && ff.Type != "" && !f.inBlacklist(ff.Name, parentName) {
				fields = append(fields, ff)
			}
		}
		return true
	})
	return fields
}

func (f StructFetcher) isBasicType(name string) bool {
	switch name {
	case
		"int",
		"int8",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint32",
		"uint64",
		"string",
		"float32",
		"float64",
		"bool",
		"byte":
		return true
	default:
		return false
	}
}

func (f StructFetcher) isOKSliceType(name, structName, fieldName string) (ok, basicType bool, getFieldName string) {
	if f.isBasicType(name) {
		return true, true, ""
	}

	fn, ok := f.sliceTypes[fmt.Sprintf("%s.%s", structName, fieldName)]
	return ok, false, fn
}

func (f StructFetcher) inBlacklist(name, parentType string) bool {
	for _, n := range f.blacklist[parentType] {
		if n == name {
			return true
		}
	}

	for _, n := range f.blacklist["*"] {
		if n == name {
			return true
		}
	}
	return false
}

func (f StructFetcher) firstName(names []*ast.Ident) string {
	if len(names) == 0 {
		return ""
	}

	return names[0].Name
}

func (f StructFetcher) extractType(n ast.Node) (name string, ptr, slice, isMap bool) {
	switch x := n.(type) {
	case *ast.Ident:
		return x.Name, false, false, false
	case *ast.StarExpr:
		name, _, _, _ := f.extractType(x.X)
		return name, true, false, false
	case *ast.SelectorExpr:
		pkg, _, _, _ := f.extractType(x.X)
		name, _, _, _ := f.extractType(x.Sel)
		fullName := fmt.Sprintf("%s.%s", pkg, name)

		if _, ok := f.knownTypes[fullName]; !ok {
			return "", false, false, false
		}

		return fullName, false, false, false
	case *ast.ArrayType:
		name, _, _, _ := f.extractType(x.Elt)
		return name, false, true, false
	case *ast.MapType:
		name, _, _, _ := f.extractType(x.Key)
		return name, false, false, true
	}

	return "", false, false, false
}
