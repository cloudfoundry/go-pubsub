package generator

import (
	"fmt"
	"strings"

	"github.com/apoydence/pubsub/pubsub-gen/internal/inspector"
)

type TraverserWriter interface {
	Package(name string) string
	Imports(names []string) string
	Done(travName string) string
	Traverse(travName, name string) string
	Hashers(travName string) string

	FieldSelector(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool, enumValue int) string
	InterfaceSelector(prefix, castTypeName, fieldName, structPkgPrefix string, implementers map[string]string, startIdx int) string
	SelectorFunc(travName, selectorName string, fields []string) string

	FieldStartStruct(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool, enumValue int) string
	FieldStructFunc(travName, prefix, fieldName, nextFieldName, castTypeName, hashFn string, isPtr bool) string
	FieldStructFuncLast(travName, prefix, fieldName, castTypeName, hashFn string, isPtr bool) string
	FieldPeersFunc(travName, prefix, castTypeName, fieldName, hashFn string, names []string, isPtr bool) string
}

type TraverserGenerator struct {
	writer TraverserWriter
}

func NewTraverserGenerator(w TraverserWriter) TraverserGenerator {
	return TraverserGenerator{
		writer: w,
	}
}

func (g TraverserGenerator) Generate(
	m map[string]inspector.Struct,
	packageName string,
	traverserName string,
	structName string,
	isPtr bool,
	structPkgPrefix string,
	imports []string,
) (string, error) {
	src := g.writer.Package(packageName)
	src += g.writer.Imports(append([]string{"github.com/apoydence/pubsub", "hash/crc64"}, imports...))

	s, ok := m[structName]
	if !ok {
		return "", fmt.Errorf("unknown struct %s", structName)
	}

	if len(s.Fields) == 0 {
		s.Fields = append(s.Fields, inspector.Field{Name: "empty", Type: "int"})
		m[structName] = s
	}

	src += g.writer.Traverse(traverserName, s.Fields[0].Name)
	src += g.writer.Done(traverserName)
	src += g.writer.Hashers(traverserName)

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
		fmt.Sprintf("data.(%s%s%s)", ptr, structPkgPrefix, structName),
		false,
		structPkgPrefix,
		m,
	)
}

func (g TraverserGenerator) generateStructFns(
	src string,
	structName string,
	traverserName string,
	prefix string,
	parentFieldName string,
	castTypeName string,
	isPtr bool,
	structPkgPrefix string,
	m map[string]inspector.Struct,
) (string, error) {
	s, ok := m[structName]
	if !ok {
		return "", fmt.Errorf("unknown struct %s", structName)
	}

	if len(s.Fields) == 0 {
		s.Fields = append(s.Fields, inspector.Field{Name: "empty", Type: "int"})
		m[structName] = s
	}

	if parentFieldName != "" {
		src += g.writer.FieldStartStruct(
			traverserName,
			prefix,
			s.Fields[0].Name,
			parentFieldName,
			castTypeName,
			isPtr,
			1,
		)
	}

	for i, f := range s.Fields[:len(s.Fields)-1] {
		src += g.writer.FieldStructFunc(
			traverserName,
			prefix,
			f.Name,
			s.Fields[i+1].Name,
			castTypeName,
			f.Type,
			f.Ptr,
		)
	}

	if len(s.PeerTypeFields) == 0 && len(s.InterfaceTypeFields) == 0 {
		return src + g.writer.FieldStructFuncLast(
			traverserName,
			prefix,
			s.Fields[len(s.Fields)-1].Name,
			castTypeName,
			s.Fields[len(s.Fields)-1].Type,
			s.Fields[len(s.Fields)-1].Ptr,
		), nil
	}

	var names []string
	for _, pf := range s.PeerTypeFields {
		names = append(names, pf.Name)
	}

	for field := range s.InterfaceTypeFields {
		names = append(names, field.Name)
	}

	var peerFields []string
	var fieldNames []string

	// Struct Peers
	var i int
	for _, f := range s.PeerTypeFields {
		x, ok := m[f.Type]
		if !ok {
			continue
		}

		fieldNames = append(fieldNames, f.Name)
		peerFields = append(peerFields, g.writer.FieldSelector(
			traverserName,
			fmt.Sprintf("%s_%s", prefix, f.Name),
			x.Fields[0].Name,
			f.Name,
			castTypeName,
			f.Ptr,
			i+1,
		))
		i++
	}

	// Interface Peers
	for field, implementers := range s.InterfaceTypeFields {
		implementersWithFields := make(map[string]string)
		for _, impl := range implementers {
			implementersWithFields[impl] = m[impl].Fields[0].Name
		}

		fieldNames = append(fieldNames, field.Name)
		peerFields = append(peerFields, g.writer.InterfaceSelector(
			prefix,
			castTypeName,
			field.Name,

			structPkgPrefix,
			implementersWithFields,
			i,
		))
	}

	src += g.writer.FieldPeersFunc(
		traverserName,
		prefix,
		castTypeName,
		s.Fields[len(s.Fields)-1].Name,
		s.Fields[len(s.Fields)-1].Type,
		fieldNames,
		s.Fields[len(s.Fields)-1].Ptr,
	)

	if len(peerFields) != 0 {
		src += g.writer.SelectorFunc(traverserName, strings.Join(fieldNames, "_"), peerFields)
	}

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
			structPkgPrefix,
			m,
		)
		if err != nil {
			return "", err
		}
	}

	for field, implementers := range s.InterfaceTypeFields {
		for _, i := range implementers {
			var err error
			src, err = g.generateStructFns(
				src,
				i,
				traverserName,
				fmt.Sprintf("%s_%s_%s", prefix, field.Name, i),
				i,
				fmt.Sprintf("%s.%s.(%s%s)", castTypeName, field.Name, structPkgPrefix, i),
				false,
				structPkgPrefix,
				m,
			)
			if err != nil {
				return "", err
			}
		}
	}

	return src, nil
}

func hashFn(t, dataValue string) string {
	switch t {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "float32", "float64":
		return fmt.Sprintf("uint64(%s)", dataValue)
	case "string":
		return fmt.Sprintf("crc64.Checksum([]byte(%s), tableECMA)", dataValue)
	case "bool":
		return fmt.Sprintf("hashBool(%s)", dataValue)
	default:
		return dataValue
	}
}
