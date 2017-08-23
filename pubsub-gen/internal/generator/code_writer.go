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
	return pubsub.Paths( func(idx int, data interface{}) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
	  return nil, nil, false
	})
}
`, travName)
}

func (w CodeWriter) FieldStartStruct(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool) string {
	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s == nil {
		return pubsub.Paths(func(idx int, data interface{}) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverser(s.done), true
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
  return pubsub.Paths(func(idx int, data interface{}) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return "%s", pubsub.TreeTraverser(s.%s_%s), true
			default:
				return nil, nil, false
			}
		})
}
`, travName, prefix, nilCheck, parentFieldName, prefix, fieldName)
}

func (w CodeWriter) FieldSelector(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool) string {
	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s.%s == nil {
		return nil, pubsub.TreeTraverser(s.done), true
  }
		`, castTypeName, parentFieldName)
	}

	return fmt.Sprintf(`
	%s
	return "%s", pubsub.TreeTraverser(s.%s_%s), true
`, nilCheck, parentFieldName, prefix, fieldName)
}

func (w CodeWriter) SelectorFunc(travName, selectorName string, fields []string) string {
	var body string
	for i, f := range fields {
		body += fmt.Sprintf(`
	case %d:
	%s
		`, i, f)
	}

	return fmt.Sprintf(`
	func (s %s) __%s (idx int, data interface{}) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
		switch idx{
	%s
default:
	return nil, nil, false
}
	}
	`, travName, selectorName, body)
}

func (w CodeWriter) FieldStructFunc(travName, prefix, fieldName, nextFieldName, castTypeName string, isPtr bool) string {
	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s.%s == nil {
    return pubsub.Paths(func(idx int, data interface{}) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverser(s.%s_%s), true
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
  return pubsub.Paths(func(idx int, data interface{}) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverser(s.%s_%s), true
			case 1:
				return %s%s.%s, pubsub.TreeTraverser(s.%s_%s), true
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
    return pubsub.Paths(func(idx int, data interface{}) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverser(s.done), true
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
  return pubsub.Paths(func(idx int, data interface{}) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverser(s.done), true
			case 1:
				return %s%s.%s, pubsub.TreeTraverser(s.done), true
			default:
				return nil, nil, false
			}
		})
}
`, travName, prefix, fieldName, nilCheck, star, castTypeName, fieldName)
}

func (w CodeWriter) FieldPeersFunc(travName, prefix, castTypeName, fieldName string, names []string, isPtr bool) string {
	travFunc := fmt.Sprintf(`
    pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
			return s.__%s
 		})`, strings.Join(names, "_"))

	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s.%s == nil {
    return pubsub.Paths(func(idx int, data interface{}) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return nil, %s, true
			default:
				return nil, nil, false
			}
		})
  }
		`, castTypeName, fieldName, travFunc)
	}

	var star string
	if isPtr {
		star = "*"
	}

	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}) pubsub.Paths {
	%s
  return pubsub.Paths(func(idx int, data interface{}) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool){
		switch idx{
		case 0:
				return nil, %s, true
		case 1:
				return %s%s.%s, %s, true
	  default:
			return nil, nil, false
		}
	})
}
`, travName, prefix, fieldName, nilCheck, travFunc, star, castTypeName, fieldName, travFunc)
}

func (w CodeWriter) InterfaceSelector(prefix, castTypeName, fieldName, structPkgPrefix string, implementers map[string]string) string {
	body := fmt.Sprintf("switch %s.%s.(type) {", castTypeName, fieldName)
	for i, f := range implementers {
		body += fmt.Sprintf(`
case %s%s:
	return "%s", s.%s_%s_%s_%s, true
`, structPkgPrefix, i, i, prefix, fieldName, i, f)
	}
	body += `
default:
	return nil, pubsub.TreeTraverser(s.done), true
}`

	return body
}
