package inspector

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
)

type StructParser interface {
	Parse(n ast.Node) ([]Struct, error)
}

type PackageParser struct {
	s StructParser
}

func NewPackageParser(s StructParser) PackageParser {
	return PackageParser{
		s: s,
	}
}

func (p PackageParser) Parse(packagePath, gopath string) (map[string]Struct, error) {
	pkgPath := filepath.Join(gopath, "src", packagePath)
	files, err := ioutil.ReadDir(pkgPath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]Struct)
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".go" {
			continue
		}

		filePath := filepath.Join(pkgPath, file.Name())
		fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("%s -> %s", filePath, err)
		}

		fset := token.NewFileSet()
		n, err := parser.ParseFile(fset, filePath, fileData, 0)
		if err != nil {
			return nil, fmt.Errorf("%s -> %s", filePath, err)
		}

		ss, err := p.s.Parse(n)

		if err != nil {
			return nil, fmt.Errorf("%s -> %s", filePath, err)
		}

		for _, s := range ss {
			m[s.Name] = s
		}
	}
	return m, nil
}
