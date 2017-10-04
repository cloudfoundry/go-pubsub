package generator

import (
	"fmt"
	"strings"

	"code.cloudfoundry.org/go-pubsub/pubsub-gen/internal/inspector"
)

type TraverserWriter interface {
	Package(name string) string
	Imports(names []string) string
	Done(travName string) string
	Traverse(travName, name string) string
	Hashers(travName string) string

	FieldSelector(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool, enumValue int) string
	InterfaceSelector(prefix, castTypeName, fieldName, structPkgPrefix string, implementers map[string]string, startIdx int) string
	SelectorFunc(travName, prefix, selectorName string, fields []string) string

	FieldStartStruct(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool, enumValue int) string
	FieldStructFunc(travName, prefix, fieldName, nextFieldName, castTypeName, hashFn string, isPtr bool, slice inspector.Slice, m inspector.Map) string
	FieldStructFuncLast(travName, prefix, fieldName, castTypeName, hashFn string, isPtr bool, slice inspector.Slice, m inspector.Map) string
	FieldPeersFunc(travName, prefix, castTypeName, fieldName, hashFn string, names []string, isPtr bool, slice inspector.Slice, m inspector.Map) string
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
	src += g.writer.Imports(append([]string{"code.cloudfoundry.org/go-pubsub", "hash/crc64"}, imports...))

	s, ok := m[structName]
	if !ok {
		return "", fmt.Errorf("unknown struct %s", structName)
	}

	var name string
	if len(s.Fields) > 0 {
		name = s.Fields[0].Name
	}

	src += g.writer.Traverse(traverserName, name)
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

	if parentFieldName != "" {
		var name string
		if len(s.Fields) > 0 {
			name = s.Fields[0].Name
		}

		src += g.writer.FieldStartStruct(
			traverserName,
			prefix,
			name,
			parentFieldName,
			castTypeName,
			isPtr,
			1,
		)
	}

	if len(s.Fields) > 0 {
		for i, f := range s.Fields[:len(s.Fields)-1] {
			src += g.writer.FieldStructFunc(
				traverserName,
				prefix,
				f.Name,
				s.Fields[i+1].Name,
				castTypeName,
				f.Type,
				f.Ptr,
				f.Slice,
				f.Map,
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
				s.Fields[len(s.Fields)-1].Slice,
				s.Fields[len(s.Fields)-1].Map,
			), nil
		}
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

		var name string
		if len(x.Fields) > 0 {
			name = x.Fields[0].Name
		}

		fieldNames = append(fieldNames, f.Name)
		peerFields = append(peerFields, g.writer.FieldSelector(
			traverserName,
			fmt.Sprintf("%s_%s", prefix, f.Name),
			name,
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
			if len(m[impl].Fields) == 0 {
				implementersWithFields[impl] = ""
				continue
			}

			var name string
			if len(m[impl].Fields) > 0 {
				name = m[impl].Fields[0].Name
			}

			implementersWithFields[impl] = name
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

	if len(s.Fields) > 0 {
		src += g.writer.FieldPeersFunc(
			traverserName,
			prefix,
			castTypeName,
			s.Fields[len(s.Fields)-1].Name,
			s.Fields[len(s.Fields)-1].Type,
			fieldNames,
			s.Fields[len(s.Fields)-1].Ptr,
			s.Fields[len(s.Fields)-1].Slice,
			s.Fields[len(s.Fields)-1].Map,
		)
	}

	if len(peerFields) != 0 {
		src += g.writer.SelectorFunc(traverserName, prefix, strings.Join(fieldNames, "_"), peerFields)
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

func hashSplitFn(t, dataValue string, slice inspector.Slice, m inspector.Map) (calc, value string) {
	if slice.IsSlice {
		x := "x"
		if !slice.IsBasicType {
			x = fmt.Sprintf("x.%s", slice.FieldName)
		}

		_, value := hashSplitFn(t, x, inspector.Slice{}, inspector.Map{})
		return fmt.Sprintf(`
	var total uint64 = 1
	for _, x := range %s{ 
		total += %s
	}`, dataValue, value), "total"
	}

	if m.IsMap {
		_, value := hashSplitFn(t, "x", inspector.Slice{}, inspector.Map{})
		return fmt.Sprintf(`
	var total uint64 = 1
	for x := range %s{ 
		total += %s
	}`, dataValue, value), "total"
	}

	switch t {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "byte", "uint16", "uint32", "float32", "float64":
		return "", fmt.Sprintf("uint64(%s)+1", dataValue)
	case "string":
		return "", fmt.Sprintf("crc64.Checksum([]byte(%s), tableECMA)+1", dataValue)
	case "bool":
		return "", fmt.Sprintf("hashBool(%s)", dataValue)
	default:
		return "", dataValue
	}
}
