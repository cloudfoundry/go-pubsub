package generator

import (
	"fmt"
	"sort"
	"strings"

	"code.cloudfoundry.org/go-pubsub/pubsub-gen/internal/inspector"
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

	src, err = g.genPath(src, "", m, genName, structName, genName+"CreatePath", true, 0)
	if err != nil {
		return "", err
	}

	return src, err
}

func (g PathGenerator) genPath(
	src string,
	prefix string,
	m map[string]inspector.Struct,
	genName string,
	structName string,
	funcName string,
	includeMinimize bool,
	enumValue int,
) (string, error) {
	body, err := g.genPathBody(
		m,
		structName,
	)

	if err != nil {
		return "", err
	}

	structName = strings.Trim(structName, "*")
	s, ok := m[structName]
	if !ok {
		return "", fmt.Errorf("unknown struct %s", structName)
	}

	var next string
	for _, pf := range s.PeerTypeFields {
		next += g.genPathNextFunc(m, prefix, pf.Name)
	}

	for f, implementers := range s.InterfaceTypeFields {
		for _, i := range implementers {
			next += g.genPathNextFunc(m, prefix, fmt.Sprintf("%s_%s", f.Name, i))
		}
	}

	var addLabel string
	if enumValue != 0 {
		addLabel = fmt.Sprintf(`path = append(path, %d)`, enumValue)
	}

	var minimize string
	if includeMinimize {
		minimize = `
for i := len(path) - 1; i >= 1; i-- {
	if path[i] != 0 {
		break
	}
	path = path[:i]
}
`
	}

	src += fmt.Sprintf(`
func %s(f *%sFilter) []uint64 {
if f == nil {
	return nil
}
var path []uint64

%s

%s

%s

%s

return path
}
`, funcName, g.sanitizeName(structName), addLabel, body, next, minimize)

	var idx int
	for _, pf := range s.PeerTypeFields {
		src, err = g.genPath(src, fmt.Sprintf("%s_%s", prefix, pf.Name), m, genName, pf.Type, fmt.Sprintf("createPath_%s_%s", prefix, pf.Name), false, idx+1)
		idx++
	}

	for f, implementers := range s.InterfaceTypeFields {
		ii := make([]string, len(implementers))
		copy(ii, implementers)

		for i, v := range ii {
			ii[i] = strings.Trim(v, "*")
		}

		sort.Strings(ii)
		for j, i := range ii {

			i = strings.Replace(i, "*", "", -1)

			src, err = g.genPath(src, fmt.Sprintf("%s_%s", prefix, i), m, genName, i, fmt.Sprintf("createPath_%s_%s_%s", prefix, f.Name, i), false, j+idx+1)
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
	prefix string,
	structName string,
) string {
	structName = strings.Replace(structName, "*", "", -1)
	return fmt.Sprintf(`
path = append(path, createPath_%s_%s(f.%s)...)
`, prefix, structName, structName)
}

func (g PathGenerator) genPathBody(
	m map[string]inspector.Struct,
	structName string,
) (string, error) {
	var src string

	structName = strings.Trim(structName, "*")
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
`, f.Name, strings.Trim(i, "*"))
		}
	}

	onlyOneCheck += `
if count > 1 {
	panic("Only one field can be set")
}
`

	buildPath := ""
	for _, f := range s.Fields {
		var star string
		if !f.Slice.IsSlice && !f.Map.IsMap {
			star = "*"
		}

		// We convert maps to slices for the path. This allows users to enter the
		// desired keys more easily.
		if f.Map.IsMap {
			f.Map.IsMap = false
			f.Slice.IsSlice = true
		}

		dataValue := fmt.Sprintf("%sf.%s", star, f.Name)
		f.Slice.IsBasicType = true
		hashCalc, hashValue := hashSplitFn(f.Type, dataValue, f.Slice, inspector.Map{})

		buildPath += fmt.Sprintf(`
if f.%s != nil {
	%s
	path = append(path, %s)
}else{
	path = append(path, 0)
}
`, f.Name, hashCalc, hashValue)
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

	structName = strings.Trim(structName, "*")
	s, ok := m[structName]
	if !ok {
		return "", fmt.Errorf("unknown struct %s", structName)
	}

	var fields string
	for _, f := range s.Fields {
		if f.Slice.IsSlice || f.Map.IsMap {
			fields += fmt.Sprintf("%s []%s\n", f.Name, f.Type)
			continue
		}

		fields += fmt.Sprintf("%s *%s\n", f.Name, f.Type)
	}

	for _, f := range s.PeerTypeFields {
		fields += fmt.Sprintf("%s *%sFilter\n", f.Name, g.sanitizeName(f.Type))
	}

	for f, implementers := range s.InterfaceTypeFields {
		for _, i := range implementers {
			i = strings.Trim(i, "*")
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
