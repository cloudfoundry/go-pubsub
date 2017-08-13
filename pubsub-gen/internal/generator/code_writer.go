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
		if n == "" {
			continue
		}
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
func (s %s) Traverse(data interface{}) pubsub.Paths {
	return s._%s(data)
}
`, travName, firstField)
}

func (w CodeWriter) Done(travName string) string {
	return fmt.Sprintf(`
	func (s %s) done(data interface{}) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}
`, travName)
}

func (w CodeWriter) FieldStartStruct(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool) string {
	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s == nil {
		return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverserFunc(s.done), true
			default:
				return nil, nil, false
			}
		})
  }
		`, castTypeName)
	}

	return fmt.Sprintf(`
func(s %s) %s(data interface{}) pubsub.Paths {
	%s
  return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return "%s", pubsub.TreeTraverserFunc(s.%s_%s), true
			default:
				return nil, nil, false
			}
		})
}
`, travName, prefix, nilCheck, parentFieldName, prefix, fieldName)
}

func (w CodeWriter) FieldStructFunc(travName, prefix, fieldName, nextFieldName, castTypeName string, isPtr bool) string {
	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s.%s == nil {
    return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverserFunc(s.%s_%s), true
			default:
				return nil, nil, false
			}
		})
  }
		`, castTypeName, fieldName, prefix, nextFieldName)
	}

	var star string
	if isPtr {
		star = "*"
	}
	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}) pubsub.Paths {
	%s
  return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverserFunc(s.%s_%s), true
			case 1:
				return fmt.Sprintf("%%v", %s%s.%s), pubsub.TreeTraverserFunc(s.%s_%s), true
			default:
				return nil, nil, false
			}
		})
}
`, travName, prefix, fieldName, nilCheck, prefix, nextFieldName, star, castTypeName, fieldName, prefix, nextFieldName)
}

func (w CodeWriter) FieldStructFuncLast(travName, prefix, fieldName, castTypeName string, isPtr bool) string {
	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s.%s == nil {
    return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverserFunc(s.done), true
			default:
				return nil, nil, false
			}
		})
  }
		`, castTypeName, fieldName)
	}

	var star string
	if isPtr {
		star = "*"
	}

	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}) pubsub.Paths {
	%s
  return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverserFunc(s.done), true
			case 1:
				return fmt.Sprintf("%%v", %s%s.%s), pubsub.TreeTraverserFunc(s.done), true
			default:
				return nil, nil, false
			}
		})
}
`, travName, prefix, fieldName, nilCheck, star, castTypeName, fieldName)
}

func (w CodeWriter) FieldPeersBodyEntry(idx int, names []string, prefix, castTypeName, fieldName string) string {
	idx = idx * 2
	idx2 := idx + 1

	var travs []string
	for _, name := range names {
		travs = append(travs, fmt.Sprintf("s.%s_%s(data)", prefix, name))
	}

	travFunc := fmt.Sprintf(`
     pubsub.TreeTraverserFunc(func(data interface{}) pubsub.Paths {
 				return pubsub.CombinePaths(%s)
 			})`, strings.Join(travs, ","))

	return fmt.Sprintf(`
case %d:
				return nil, %s, true
case %d:
				return fmt.Sprintf("%%v", %s.%s), %s, true
`, idx, travFunc, idx2, castTypeName, fieldName, travFunc)
}

func (w CodeWriter) FieldPeersFunc(travName, prefix, fieldName, body string) string {
	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}) pubsub.Paths {
  return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
		switch idx{
		%s
	  default:
			return nil, nil, false
		}
	})
}
`, travName, prefix, fieldName, body)
}

func (w CodeWriter) InterfaceTypeBodyEntry(prefix, castTypeName, fieldName, structPkgPrefix string, implementers []string) string {
	body := fmt.Sprintf("switch %s.%s.(type) {", castTypeName, fieldName)
	for _, i := range implementers {
		body += fmt.Sprintf(`
case %s%s:
	return s.%s_%s_%s(data)
`, structPkgPrefix, i, prefix, fieldName, i)
	}
	body += `
default:
  return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverserFunc(s.done), true
			default:
				return nil, nil, false
			}
		})
}`

	return body
}

func (w CodeWriter) InterfaceTypeFieldsFunc(travName, prefix, fieldName, body string) string {
	return fmt.Sprintf(`
func (s %s) %s_%s (data interface{}) pubsub.Paths {
%s
}
`, travName, prefix, fieldName, body)
}
