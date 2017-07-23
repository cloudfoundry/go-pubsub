package generator

import (
	"fmt"
	"strings"

	"github.com/apoydence/pubsub/pubsub-gen/internal/inspector"
)

type Generator struct{}

func New() Generator {
	return Generator{}
}

func (g Generator) Generate(
	m map[string]inspector.Struct,
	packageName string,
	traverserName string,
	structName string,
	isPtr bool,
) (string, error) {
	src := fmt.Sprintf(`package %s

import (
	"github.com/apoydence/pubsub"
	"fmt"
)

type %s struct{}

func New%s()%s{ return %s{} }
`, packageName, traverserName, strings.Title(traverserName), traverserName, traverserName)

	s, ok := m[structName]
	if !ok {
		return "", fmt.Errorf("unknown struct %s", structName)
	}

	if len(s.Fields) == 0 {
		return "", fmt.Errorf("structs with no fields are not yet supported")
	}

	src = fmt.Sprintf(`%s
func (s %s) Traverse(data interface{}, currentPath []string) pubsub.Paths {
	return s._%s(data, currentPath)
}
func (s %s) done(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}
`, src, traverserName, s.Fields[0].Name, traverserName)

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
		var nilCheck string
		if isPtr {
			nilCheck = fmt.Sprintf(`
  if %s == nil {
    return pubsub.NewPathsWithTraverser([]string{""}, pubsub.TreeTraverserFunc(s.done))
  }
		`, castTypeName)
		}

		src = fmt.Sprintf(`%s
func(s %s) %s(data interface{}, currentPaht []string) pubsub.Paths {
	%sreturn pubsub.NewPathsWithTraverser([]string{"%s"}, pubsub.TreeTraverserFunc(s.%s_%s))
}
`, src, traverserName, prefix, nilCheck, parentFieldName, prefix, s.Fields[0].Name)
	}

	for i, f := range s.Fields[:len(s.Fields)-1] {
		src = fmt.Sprintf(`%s
func (s %s) %s_%s(data interface{}, currentPath []string) pubsub.Paths {
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%%v", %s.%s)}, pubsub.TreeTraverserFunc(s.%s_%s))
}
`, src, traverserName, prefix, f.Name, castTypeName, f.Name, prefix, s.Fields[i+1].Name)
	}

	if len(s.PeerTypeFields) == 0 {
		src = fmt.Sprintf(`%s
func (s %s) %s_%s(data interface{}, currentPath []string) pubsub.Paths {
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%%v", %s.%s)}, pubsub.TreeTraverserFunc(s.done))
}
`, src, traverserName, prefix, s.Fields[len(s.Fields)-1].Name, castTypeName, s.Fields[len(s.Fields)-1].Name)

		return src, nil
	}

	var peers string
	for _, pf := range s.PeerTypeFields {
		peers = fmt.Sprintf(`%s
{
  Path:      "",
  Traverser: pubsub.TreeTraverserFunc(s.%s_%s),
},
{
  Path:      fmt.Sprintf("%%v", %s.%s),
  Traverser: pubsub.TreeTraverserFunc(s.%s_%s),
},
`, peers, prefix, pf.Name, castTypeName, s.Fields[len(s.Fields)-1].Name, prefix, pf.Name)
	}

	src = fmt.Sprintf(`%s
func (s %s) %s_%s(data interface{}, currentPath []string) pubsub.Paths {
  return pubsub.PathAndTraversers(
    []pubsub.PathAndTraverser{%s})
}
`, src, traverserName, prefix, s.Fields[len(s.Fields)-1].Name, peers)

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

	return src, nil
}
