package generator

import (
	"fmt"
	"strings"
)

type codeWriter struct{}

func (w codeWriter) Package(name string) string {
	return fmt.Sprintf("package %s\n\n", name)
}

func (w codeWriter) Imports(names []string) string {
	result := "import (\n"
	for _, n := range names {
		result += fmt.Sprintf("  \"%s\"\n", n)
	}
	return fmt.Sprintf("%s)\n", result)
}

func (w codeWriter) DefineType(travName string) string {
	return fmt.Sprintf("type %s struct{}\n", travName)
}

func (w codeWriter) Constructor(travName string) string {
	return fmt.Sprintf(" func New%s()%s{ return %s{} }\n", strings.Title(travName), travName, travName)
}

func (w codeWriter) Traverse(travName, firstField string) string {
	return fmt.Sprintf(`
func (s %s) Traverse(data interface{}, currentPath []string) pubsub.Paths {
	return s._%s(data, currentPath)
}
`, travName, firstField)
}

func (w codeWriter) Done(travName string) string {
	return fmt.Sprintf(`
	func (s %s) done(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}
`, travName)
}

func (w codeWriter) FieldStartStruct(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool) string {
	if parentFieldName == "" {
		return ""
	}

	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s == nil {
    return pubsub.NewPathsWithTraverser([]string{""}, pubsub.TreeTraverserFunc(s.done))
  }
		`, castTypeName)
	}

	return fmt.Sprintf(`
func(s %s) %s(data interface{}, currentPaht []string) pubsub.Paths {
	%sreturn pubsub.NewPathsWithTraverser([]string{"%s"}, pubsub.TreeTraverserFunc(s.%s_%s))
}
`, travName, prefix, nilCheck, parentFieldName, prefix, fieldName)
}

func (w codeWriter) FieldStructFunc(travName, prefix, fieldName, nextFieldName, castTypeName string) string {
	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}, currentPath []string) pubsub.Paths {
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%%v", %s.%s)}, pubsub.TreeTraverserFunc(s.%s_%s))
}
`, travName, prefix, fieldName, castTypeName, fieldName, prefix, nextFieldName)
}

func (w codeWriter) FieldStructFuncLast(travName, prefix, fieldName, castTypeName string) string {
	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}, currentPath []string) pubsub.Paths {
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%%v", %s.%s)}, pubsub.TreeTraverserFunc(s.done))
}
`, travName, prefix, fieldName, castTypeName, fieldName)
}

func (w codeWriter) FieldPeersBodyEntry(prefix, name, castTypeName, fieldName string) string {
	return fmt.Sprintf(`
{
  Path:      "",
  Traverser: pubsub.TreeTraverserFunc(s.%s_%s),
},
{
  Path:      fmt.Sprintf("%%v", %s.%s),
  Traverser: pubsub.TreeTraverserFunc(s.%s_%s),
},
`, prefix, name, castTypeName, fieldName, prefix, name)
}

func (w codeWriter) FieldPeersFunc(travName, prefix, fieldName, body string) string {
	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}, currentPath []string) pubsub.Paths {
  return pubsub.PathAndTraversers(
    []pubsub.PathAndTraverser{%s})
}
`, travName, prefix, fieldName, body)
}
