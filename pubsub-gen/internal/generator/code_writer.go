package generator

import (
	"fmt"
	"strings"
)

type CodeWriter struct{}

func (w CodeWriter) Package(name string) string {
	return fmt.Sprintf("package %s\n\n", name)
}

func (w CodeWriter) Imports(names []string) string {
	result := "import (\n"
	for _, n := range names {
		result += fmt.Sprintf("  \"%s\"\n", n)
	}
	return fmt.Sprintf("%s)\n", result)
}

func (w CodeWriter) DefineType(travName string) string {
	return fmt.Sprintf("type %s struct{}\n", travName)
}

func (w CodeWriter) Constructor(travName string) string {
	return fmt.Sprintf(" func New%s()%s{ return %s{} }\n", strings.Title(travName), travName, travName)
}

func (w CodeWriter) Traverse(travName, firstField string) string {
	return fmt.Sprintf(`
func (s %s) Traverse(data interface{}, currentPath []string) pubsub.Paths {
	return s._%s(data, currentPath)
}
`, travName, firstField)
}

func (w CodeWriter) Done(travName string) string {
	return fmt.Sprintf(`
	func (s %s) done(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}
`, travName)
}

func (w CodeWriter) FieldStartStruct(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool) string {
	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s == nil {
    return pubsub.NewPathsWithTraverser([]string{""}, pubsub.TreeTraverserFunc(s.done))
  }
		`, castTypeName)
	}

	return fmt.Sprintf(`
func(s %s) %s(data interface{}, currentPath []string) pubsub.Paths {
	%sreturn pubsub.NewPathsWithTraverser([]string{"%s"}, pubsub.TreeTraverserFunc(s.%s_%s))
}
`, travName, prefix, nilCheck, parentFieldName, prefix, fieldName)
}

func (w CodeWriter) FieldStructFunc(travName, prefix, fieldName, nextFieldName, castTypeName string) string {
	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}, currentPath []string) pubsub.Paths {
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%%v", %s.%s)}, pubsub.TreeTraverserFunc(s.%s_%s))
}
`, travName, prefix, fieldName, castTypeName, fieldName, prefix, nextFieldName)
}

func (w CodeWriter) FieldStructFuncLast(travName, prefix, fieldName, castTypeName string) string {
	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}, currentPath []string) pubsub.Paths {
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%%v", %s.%s)}, pubsub.TreeTraverserFunc(s.done))
}
`, travName, prefix, fieldName, castTypeName, fieldName)
}

func (w CodeWriter) FieldPeersBodyEntry(prefix, name, castTypeName, fieldName string) string {
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

func (w CodeWriter) FieldPeersFunc(travName, prefix, fieldName, body string) string {
	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}, currentPath []string) pubsub.Paths {
  return pubsub.PathAndTraversers(
    []pubsub.PathAndTraverser{%s})
}
`, travName, prefix, fieldName, body)
}

func (w CodeWriter) InterfaceTypeBodyEntry(prefix, castTypeName, fieldName, structPkgPrefix string, implementers []string) string {
	body := fmt.Sprintf("switch %s.%s.(type) {", castTypeName, fieldName)
	for _, i := range implementers {
		body += fmt.Sprintf(`
case %s%s:
	return s.%s_%s_%s(data, currentPath)
`, structPkgPrefix, i, prefix, fieldName, i)
	}
	body += `
default:
	return pubsub.NewPathsWithTraverser([]string{""}, pubsub.TreeTraverserFunc(s.done))
}`

	return body
}

func (w CodeWriter) InterfaceTypeFieldsFunc(travName, prefix, fieldName, body string) string {
	return fmt.Sprintf(`
func (s %s) %s_%s (data interface{}, currentPath []string) pubsub.Paths {
%s
}
`, travName, prefix, fieldName, body)
}
