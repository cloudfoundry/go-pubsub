package inspector_test

import (
	"go/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"code.cloudfoundry.org/go-pubsub/pubsub-gen/internal/inspector"
)

type TPP struct {
	*testing.T
	structParser *spyStructParser
	p            inspector.PackageParser
	gopath       string
}

func TestStructPackageParser(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) TPP {
		gopath := writeTestPackage()
		structParser := &spyStructParser{returnValue: []inspector.Struct{
			{Name: "X"}, {Name: "Y"}, {Name: "Z"},
		}}
		return TPP{
			T:            t,
			structParser: structParser,
			p:            inspector.NewPackageParser(structParser),
			gopath:       gopath,
		}
	})

	o.Spec("it opens each file in the given path", func(t TPP) {
		structs, err := t.p.Parse("some-package", t.gopath)
		Expect(t, err == nil).To(BeTrue())
		Expect(t, t.structParser.nodes).To(HaveLen(2))
		Expect(t, structs).To(HaveLen(3))
	})

	o.Spec("it returns an error for an unknown path", func(t TPP) {
		_, err := t.p.Parse("garbage-package", t.gopath)
		Expect(t, err == nil).To(BeFalse())
	})
}

func writeTestPackage() string {
	dir, err := ioutil.TempDir("", "ast-gen")
	if err != nil {
		panic(err)
	}
	os.Mkdir(filepath.Join(dir, "src"), os.ModePerm)
	os.Mkdir(filepath.Join(dir, "src", "some-package"), os.ModePerm)

	ioutil.WriteFile(filepath.Join(dir, "src", "some-package", "test1.go"),
		[]byte(
			`
package p
		`,
		),
		os.ModePerm)

	ioutil.WriteFile(filepath.Join(dir, "src", "some-package", "test2.go"),
		[]byte(
			`
package p
		`,
		),
		os.ModePerm)

	return dir
}

type spyStructParser struct {
	nodes       []ast.Node
	returnValue []inspector.Struct
	err         error
}

func (s *spyStructParser) Parse(n ast.Node) ([]inspector.Struct, error) {
	s.nodes = append(s.nodes, n)
	return s.returnValue, s.err
}
