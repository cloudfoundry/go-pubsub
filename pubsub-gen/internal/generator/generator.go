package generator

import (
	"fmt"

	"github.com/apoydence/pubsub/pubsub-gen/internal/inspector"
)

type Writer interface {
	Package(name string) string
	Imports(names []string) string
	DefineType(travName string) string
	Constructor(travName string) string
	Done(travName string) string
	Traverse(travName, name string) string
	FieldStartStruct(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool) string
	FieldStructFunc(travName, prefix, fieldName, nextFieldName, castTypeName string) string
	FieldStructFuncLast(travName, prefix, fieldName, castTypeName string) string
	FieldPeersBodyEntry(prefix, name, castTypeName, fieldName string) string
	FieldPeersFunc(travName, prefix, fieldName, body string) string
	InterfaceTypeBodyEntry(prefix, castTypeName, fieldName string, implementers []string) string
	InterfaceTypeFieldsFunc(travName, prefix, fieldName, body string) string
}

type Generator struct {
	writer Writer
}

func New(w Writer) Generator {
	return Generator{
		writer: w,
	}
}

func (g Generator) Generate(
	m map[string]inspector.Struct,
	packageName string,
	traverserName string,
	structName string,
	isPtr bool,
) (string, error) {
	src := g.writer.Package(packageName)
	src += g.writer.Imports([]string{"github.com/apoydence/pubsub", "fmt"})
	src += g.writer.DefineType(traverserName)
	src += g.writer.Constructor(traverserName)

	s, ok := m[structName]
	if !ok {
		return "", fmt.Errorf("unknown struct %s", structName)
	}

	if len(s.Fields) == 0 {
		return "", fmt.Errorf("structs with no fields are not yet supported")
	}

	src += g.writer.Traverse(traverserName, s.Fields[0].Name)
	src += g.writer.Done(traverserName)

	var ptr string
	if isPtr {
		ptr = "*"
	}

	return g.generateStructFns(
		src,
		structName,
		traverserName,
		"",
		"",
		fmt.Sprintf("data.(%s%s)", ptr, structName),
		false,
		m,
	)
}

func (g Generator) generateStructFns(
	src string,
	structName string,
	traverserName string,
	prefix string,
	parentFieldName string,
	castTypeName string,
	isPtr bool,
	m map[string]inspector.Struct,
) (string, error) {
	s, ok := m[structName]
	if !ok {
		return "", fmt.Errorf("unknown struct %s", structName)
	}

	if len(s.Fields) == 0 {
		return "", fmt.Errorf("structs with no fields are not yet supported")
	}

	if parentFieldName != "" {
		src += g.writer.FieldStartStruct(
			traverserName,
			prefix,
			s.Fields[0].Name,
			parentFieldName,
			castTypeName,
			isPtr,
		)
	}

	for i, f := range s.Fields[:len(s.Fields)-1] {
		src += g.writer.FieldStructFunc(
			traverserName,
			prefix,
			f.Name,
			s.Fields[i+1].Name,
			castTypeName,
		)
	}

	if len(s.PeerTypeFields) == 0 && len(s.InterfaceTypeFields) == 0 {
		return src + g.writer.FieldStructFuncLast(
			traverserName,
			prefix,
			s.Fields[len(s.Fields)-1].Name,
			castTypeName,
		), nil
	}

	var peers string
	for _, pf := range s.PeerTypeFields {
		peers += g.writer.FieldPeersBodyEntry(
			prefix,
			pf.Name,
			castTypeName,
			s.Fields[len(s.Fields)-1].Name,
		)
	}

	for field, _ := range s.InterfaceTypeFields {
		peers += g.writer.FieldPeersBodyEntry(
			prefix,
			field.Name,
			castTypeName,
			s.Fields[len(s.Fields)-1].Name,
		)
	}

	src += g.writer.FieldPeersFunc(
		traverserName,
		prefix,
		s.Fields[len(s.Fields)-1].Name,
		peers,
	)

	for _, field := range s.PeerTypeFields {
		var err error
		src, err = g.generateStructFns(
			src,
			field.Type,
			traverserName,
			fmt.Sprintf("%s_%s", prefix, field.Name),
			field.Name,
			fmt.Sprintf("%s.%s", castTypeName, field.Name),
			field.Ptr,
			m,
		)
		if err != nil {
			return "", err
		}
	}

	for field, implementers := range s.InterfaceTypeFields {
		body := g.writer.InterfaceTypeBodyEntry(prefix, castTypeName, field.Name, implementers)
		src += g.writer.InterfaceTypeFieldsFunc(traverserName, prefix, field.Name, body)
	}

	for field, implementers := range s.InterfaceTypeFields {
		for _, i := range implementers {
			var err error
			src, err = g.generateStructFns(
				src,
				i,
				traverserName,
				fmt.Sprintf("%s_%s_%s", prefix, field.Name, i),
				field.Name,
				fmt.Sprintf("%s.%s.(%s)", castTypeName, field.Name, i),
				false,
				m,
			)
			if err != nil {
				return "", err
			}
		}
	}

	return src, nil
}
