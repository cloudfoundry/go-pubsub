package generator

import (
	"fmt"
	"sort"
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
	return pubsub.Paths( func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool){
	  return 0, nil, false
	})
}
`, travName)
}

func (w CodeWriter) Hashers(travName string) string {
	return fmt.Sprintf(`
func (s %s) hashBool(data bool) uint64 {
	if data {
		return 1
	}
	return 0
}

var tableECMA = crc64.MakeTable(crc64.ECMA)
func (s %s) hashString(data string) uint64 {
	return crc64.Checksum([]byte(data), tableECMA)
}
`, travName, travName)

}

func (w CodeWriter) FieldStartStruct(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool, enumValue int) string {
	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s == nil {
		return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return 0, pubsub.TreeTraverser(s.done), true
			default:
				return 0, nil, false
			}
		})
  }
		`, castTypeName)
	}

	return fmt.Sprintf(`
func(s %s) %s(data interface{}) pubsub.Paths {
	%s
  return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return %d, pubsub.TreeTraverser(s.%s_%s), true
			default:
				return 0, nil, false
			}
		})
}
`, travName, prefix, nilCheck, enumValue, prefix, fieldName)
}

func (w CodeWriter) FieldSelector(travName, prefix, fieldName, parentFieldName, castTypeName string, isPtr bool, enumValue int) string {
	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s.%s == nil {
		return 0, pubsub.TreeTraverser(s.done), true
  }
		`, castTypeName, parentFieldName)
	}

	return fmt.Sprintf(`
	%s
	return %d, pubsub.TreeTraverser(s.%s_%s), true
`, nilCheck, enumValue, prefix, fieldName)
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
	func (s %s) __%s (idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool){
		switch idx{
	%s
default:
	return 0, nil, false
}
	}
	`, travName, selectorName, body)
}

func (w CodeWriter) FieldStructFunc(travName, prefix, fieldName, nextFieldName, castTypeName, hashType string, isPtr bool) string {
	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s.%s == nil {
    return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return 0, pubsub.TreeTraverser(s.%s_%s), true
			default:
				return 0, nil, false
			}
		})
  }
		`, castTypeName, fieldName, prefix, nextFieldName)
	}

	var star string
	if isPtr {
		star = "*"
	}

	dataValue := fmt.Sprintf("%s%s.%s", star, castTypeName, fieldName)
	wrappedHash := hashFn(hashType, dataValue)

	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}) pubsub.Paths {
	%s
  return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return 0, pubsub.TreeTraverser(s.%s_%s), true
			case 1:
				return %s, pubsub.TreeTraverser(s.%s_%s), true
			default:
				return 0, nil, false
			}
		})
}
`, travName, prefix, fieldName, nilCheck, prefix, nextFieldName, wrappedHash, prefix, nextFieldName)
}

func (w CodeWriter) FieldStructFuncLast(travName, prefix, fieldName, castTypeName, hashType string, isPtr bool) string {
	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s.%s == nil {
    return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return 0, pubsub.TreeTraverser(s.done), true
			default:
				return 0, nil, false
			}
		})
  }
		`, castTypeName, fieldName)
	}

	var star string
	if isPtr {
		star = "*"
	}

	dataValue := fmt.Sprintf("%s%s.%s", star, castTypeName, fieldName)
	wrappedHash := hashFn(hashType, dataValue)

	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}) pubsub.Paths {
	%s
  return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return 0, pubsub.TreeTraverser(s.done), true
			case 1:
				return %s, pubsub.TreeTraverser(s.done), true
			default:
				return 0, nil, false
			}
		})
}
`, travName, prefix, fieldName, nilCheck, wrappedHash)
}

func (w CodeWriter) FieldPeersFunc(travName, prefix, castTypeName, fieldName, hashType string, names []string, isPtr bool) string {
	travFunc := fmt.Sprintf(`
    pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
			return s.__%s
 		})`, strings.Join(names, "_"))

	var nilCheck string
	if isPtr {
		nilCheck = fmt.Sprintf(`
  if %s.%s == nil {
    return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return 0, %s, true
			default:
				return 0, nil, false
			}
		})
  }
		`, castTypeName, fieldName, travFunc)
	}

	var star string
	if isPtr {
		star = "*"
	}

	dataValue := fmt.Sprintf("%s%s.%s", star, castTypeName, fieldName)
	wrappedHash := hashFn(hashType, dataValue)

	return fmt.Sprintf(`
func (s %s) %s_%s(data interface{}) pubsub.Paths {
	%s
  return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool){
		switch idx{
		case 0:
				return 0, %s, true
		case 1:
				return %s, %s, true
	  default:
			return 0, nil, false
		}
	})
}
`, travName, prefix, fieldName, nilCheck, travFunc, wrappedHash, travFunc)
}

func (w CodeWriter) InterfaceSelector(prefix, castTypeName, fieldName, structPkgPrefix string, implementers map[string]string) string {
	idxs := orderImpls(implementers)
	body := fmt.Sprintf("switch %s.%s.(type) {", castTypeName, fieldName)
	for i, f := range implementers {
		body += fmt.Sprintf(`
case %s%s:
	return %d, s.%s_%s_%s_%s, true
`, structPkgPrefix, i, idxs[i], prefix, fieldName, i, f)
	}
	body += `
default:
	return 0, pubsub.TreeTraverser(s.done), true
}`

	return body
}

func orderImpls(impls map[string]string) map[string]int {
	m := make(map[string]int)

	var names []string
	for k := range impls {
		names = append(names, k)
	}

	sort.Strings(names)

	for i, s := range names {
		m[s] = i + 1
	}

	return m
}
