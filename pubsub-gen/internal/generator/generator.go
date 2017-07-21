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

	var ptr string
	if isPtr {
		ptr = "*"
	}

	src = fmt.Sprintf(`%s
func (s %s) Traverse(data interface{}, currentPath []string) pubsub.Paths {
	return s.%s(data, currentPath)
}

func (s %s) donetrav(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}
`, src, traverserName, s.Fields[0].Name, traverserName)

	for i, f := range s.Fields[:len(s.Fields)-1] {
		src = fmt.Sprintf(`%s
func (s %s) %s(data interface{}, currentPath []string) pubsub.Paths {
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%%v", data.(%s%s).%s)}, pubsub.TreeTraverserFunc(s.%s))
}
`, src, traverserName, f.Name, ptr, structName, f.Name, s.Fields[i+1].Name)
	}

	src = fmt.Sprintf(`%s
func (s %s) %s(data interface{}, currentPath []string) pubsub.Paths {
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%%v", data.(%s%s).%s)}, pubsub.TreeTraverserFunc(s.donetrav))
}
`, src, traverserName, s.Fields[len(s.Fields)-1].Name, ptr, structName, s.Fields[len(s.Fields)-1].Name)

	return src, nil
}
