package generator

import (
	"fmt"
	"strings"

	"github.com/apoydence/pubsub/pubsub-gen/internal/inspector"
)

type PathGenerator struct {
}

func NewPathGenerator() PathGenerator {
	return PathGenerator{}
}

func (g PathGenerator) Generate(
	existingSrc string,
	m map[string]inspector.Struct,
	genName string,
	structName string,
) (string, error) {
	src, err := g.genStruct(existingSrc, m, structName, make(map[string]bool))
	if err != nil {
		return "", err
	}

	src, err = g.genPath(src, m, genName, structName, "CreatePath", "")
	if err != nil {
		return "", err
	}

	return src, err
}
func (g PathGenerator) genPath(
	src string,
	m map[string]inspector.Struct,
	genName string,
	structName string,
	funcName string,
	label string,
) (string, error) {
	body, err := g.genPathBody(
		m,
		structName,
	)

	if err != nil {
		return "", err
	}

	s, ok := m[structName]
	if !ok {
		return "", fmt.Errorf("unknown struct %s", structName)
	}

	var next string
	for _, pf := range s.PeerTypeFields {
		next += g.genPathNextFunc(m, pf.Name)
	}

	for f, implementers := range s.InterfaceTypeFields {
		for _, i := range implementers {
			next += g.genPathNextFunc(m, fmt.Sprintf("%s_%s", f.Name, i))
		}
	}

	var addLabel string
	if label != "" {
		addLabel = fmt.Sprintf(`path = append(path, "%s")`, label)
	}

	src += fmt.Sprintf(`
func (g %s) %s(f *%sFilter) []interface{} {
if f == nil {
	return nil
}
var path []interface{}

%s

%s

%s

return path
}
`, genName, funcName, g.sanitizeName(structName), addLabel, body, next)

	for _, pf := range s.PeerTypeFields {
		src, err = g.genPath(src, m, genName, pf.Type, fmt.Sprintf("createPath_%s", pf.Name), pf.Name)
	}

	for f, implementers := range s.InterfaceTypeFields {
		for _, i := range implementers {
			src, err = g.genPath(src, m, genName, i, fmt.Sprintf("createPath_%s_%s", f.Name, i), i)
			if err != nil {
				return "", err
			}
		}
	}

	return src, nil
}

func (g PathGenerator) sanitizeName(name string) string {
	return strings.Replace(name, ".", "", -1)
}

func (g PathGenerator) genPathNextFunc(
	m map[string]inspector.Struct,
	structName string,
) string {
	return fmt.Sprintf(`
path = append(path, g.createPath_%s(f.%s)...)
`, structName, structName)
}

func (g PathGenerator) genPathBody(
	m map[string]inspector.Struct,
	structName string,
) (string, error) {
	var src string

	s, ok := m[structName]
	if !ok {
		return "", fmt.Errorf("unknown struct %s", structName)
	}

	onlyOneCheck := "var count int"
	for _, f := range s.PeerTypeFields {
		onlyOneCheck += fmt.Sprintf(`
if f.%s != nil{
	count++
}
`, f.Name)
	}

	for f, implementers := range s.InterfaceTypeFields {
		for _, i := range implementers {
			onlyOneCheck += fmt.Sprintf(`
if f.%s_%s != nil{
	count++
}
`, f.Name, i)
		}
	}

	onlyOneCheck += `
if count > 1 {
	panic("Only one field can be set")
}
`

	buildPath := ""
	for _, f := range s.Fields {
		buildPath += fmt.Sprintf(`
if f.%s != nil {
	path = append(path, *f.%s)
}else{
	path = append(path, nil)
}
`, f.Name, f.Name)
	}

	src += fmt.Sprintf(`
%s
%s
`, onlyOneCheck, buildPath)

	return src, nil
}

func (g PathGenerator) genStruct(
	src string,
	m map[string]inspector.Struct,
	structName string,
	history map[string]bool,
) (string, error) {
	if history[structName] {
		return src, nil
	}
	history[structName] = true

	s, ok := m[structName]
	if !ok {
		return "", fmt.Errorf("unknown struct %s", structName)
	}

	var fields string
	for _, f := range s.Fields {
		fields += fmt.Sprintf("%s *%s\n", f.Name, f.Type)
	}

	for _, f := range s.PeerTypeFields {
		fields += fmt.Sprintf("%s *%sFilter\n", f.Name, g.sanitizeName(f.Type))
	}

	for f, implementers := range s.InterfaceTypeFields {
		for _, i := range implementers {
			fields += fmt.Sprintf("%s_%s *%sFilter\n", f.Name, i, i)
		}
	}

	src += fmt.Sprintf(`
type %sFilter struct{
%s
}
`, g.sanitizeName(structName), fields)

	for _, f := range s.PeerTypeFields {
		var err error
		src, err = g.genStruct(src, m, f.Type, history)
		if err != nil {
			return "", err
		}
	}

	for _, implementers := range s.InterfaceTypeFields {
		for _, i := range implementers {
			var err error
			src, err = g.genStruct(src, m, i, history)
			if err != nil {
				return "", err
			}
		}
	}

	return src, nil
}
