package inspector_test

import (
	"go/parser"
	"go/token"
	"testing"

	"code.cloudfoundry.org/go-pubsub/pubsub-gen/internal/inspector"
	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
)

type TSF struct {
	*testing.T
	f inspector.StructFetcher
}

func TestStructFetcher(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) TSF {
		return TSF{
			T: t,
			f: inspector.NewStructFetcher(nil, map[string]string{"other.Type": "xx"}, nil),
		}
	})

	o.Group("normal types", func() {
		o.Spec("it parses and returns a single struct", func(t TSF) {
			src := `
package p
type x struct {
	i string
	j int
}
`
			fset := token.NewFileSet()
			n, err := parser.ParseFile(fset, "src.go", src, 0)
			Expect(t, err == nil).To(BeTrue())

			s, err := t.f.Parse(n)
			Expect(t, err == nil).To(BeTrue())
			Expect(t, s).To(HaveLen(1))
			Expect(t, s[0].Name).To(Equal("x"))
			Expect(t, s[0].Fields).To(HaveLen(2))

			Expect(t, s[0].Fields[0].Name).To(Equal("i"))
			Expect(t, s[0].Fields[0].Type).To(Equal("string"))

			Expect(t, s[0].Fields[1].Name).To(Equal("j"))
			Expect(t, s[0].Fields[1].Type).To(Equal("int"))
		})

		o.Spec("it parses multiple structs", func(t TSF) {
			src := `
package p
type x struct {
	i string
	j int
}

type y struct {
	i string
	j int
}
`
			fset := token.NewFileSet()
			n, err := parser.ParseFile(fset, "src.go", src, 0)
			Expect(t, err == nil).To(BeTrue())

			s, err := t.f.Parse(n)
			Expect(t, err == nil).To(BeTrue())

			Expect(t, s).To(HaveLen(2))
			Expect(t, s[0].Name).To(Equal("x"))
			Expect(t, s[1].Name).To(Equal("y"))
		})
	})

	o.Group("pointer type", func() {
		o.Spec("it parses and returns a single struct", func(t TSF) {
			src := `
package p
type x struct {
	i string
	j *Y
}
`
			fset := token.NewFileSet()
			n, err := parser.ParseFile(fset, "src.go", src, 0)
			Expect(t, err == nil).To(BeTrue())

			s, err := t.f.Parse(n)
			Expect(t, err == nil).To(BeTrue())
			Expect(t, s).To(HaveLen(1))
			Expect(t, s[0].Name).To(Equal("x"))
			Expect(t, s[0].Fields).To(HaveLen(2))

			Expect(t, s[0].Fields[0].Name).To(Equal("i"))
			Expect(t, s[0].Fields[0].Type).To(Equal("string"))
			Expect(t, s[0].Fields[0].Ptr).To(BeFalse())

			Expect(t, s[0].Fields[1].Name).To(Equal("j"))
			Expect(t, s[0].Fields[1].Type).To(Equal("Y"))
			Expect(t, s[0].Fields[1].Ptr).To(BeTrue())
		})
	})

	o.Group("sub types", func() {
		o.Spec("it parses and returns a single struct", func(t TSF) {
			src := `
package p
type x struct {
	i other.Type
	j *other.Type
	k dontInclude.Type
}
`
			fset := token.NewFileSet()
			n, err := parser.ParseFile(fset, "src.go", src, 0)
			Expect(t, err == nil).To(BeTrue())

			s, err := t.f.Parse(n)
			Expect(t, err == nil).To(BeTrue())
			Expect(t, s).To(HaveLen(1))
			Expect(t, s[0].Name).To(Equal("x"))
			Expect(t, s[0].Fields).To(HaveLen(2))

			Expect(t, s[0].Fields[0].Name).To(Equal("i"))
			Expect(t, s[0].Fields[0].Type).To(Equal("other.Type"))
			Expect(t, s[0].Fields[0].Ptr).To(BeFalse())

			Expect(t, s[0].Fields[1].Name).To(Equal("j"))
			Expect(t, s[0].Fields[1].Type).To(Equal("other.Type"))
			Expect(t, s[0].Fields[1].Ptr).To(BeTrue())
		})
	})
}

func TestStructFetcherWithBasicSlices(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) TSF {
		return TSF{
			T: t,
			f: inspector.NewStructFetcher(nil, nil, nil),
		}
	})

	o.Spec("it parses and returns a field that is a slice of a basic type", func(t TSF) {
		src := `
package p
type x struct {
	a []int
	b []int8
	c []int32
	d []int64
	e []uint
	f []uint8
	g []uint32
	h []uint64
	i []string
	j []float32
	k []float64
	l []bool
	m []byte
	x []unknown
}
`

		fset := token.NewFileSet()
		n, err := parser.ParseFile(fset, "src.go", src, 0)
		Expect(t, err == nil).To(BeTrue())

		s, err := t.f.Parse(n)
		Expect(t, err == nil).To(BeTrue())
		Expect(t, s).To(HaveLen(1))
		Expect(t, s[0].Name).To(Equal("x"))
		Expect(t, s[0].Fields).To(HaveLen(13))

		for _, f := range s[0].Fields {
			Expect(t, f.Slice.IsSlice).To(BeTrue())
			Expect(t, f.Slice.IsBasicType).To(BeTrue())
			Expect(t, f.Ptr).To(BeFalse())
		}
	})
}

func TestStructFetcherWithNonBasicSlices(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) TSF {
		return TSF{
			T: t,
			f: inspector.NewStructFetcher(nil, nil, map[string]string{"x.a": "myField", "x.c": ""}),
		}
	})

	o.Spec("it parses and returns a field that is a slice of a non-basic type", func(t TSF) {
		src := `
package p
type x struct {
	a []known
	b []unknown
	c []known
}
`

		fset := token.NewFileSet()
		n, err := parser.ParseFile(fset, "src.go", src, 0)
		Expect(t, err == nil).To(BeTrue())

		s, err := t.f.Parse(n)
		Expect(t, err == nil).To(BeTrue())
		Expect(t, s).To(HaveLen(1))
		Expect(t, s[0].Name).To(Equal("x"))
		Expect(t, s[0].Fields).To(HaveLen(2))

		Expect(t, s[0].Fields[0].Slice.IsSlice).To(BeTrue())
		Expect(t, s[0].Fields[0].Slice.IsBasicType).To(BeFalse())
		Expect(t, s[0].Fields[0].Slice.FieldName).To(Equal("myField"))
		Expect(t, s[0].Fields[0].Ptr).To(BeFalse())
		Expect(t, s[0].Fields[0].Type).To(Equal("known"))

		Expect(t, s[0].Fields[1].Slice.IsSlice).To(BeTrue())
		Expect(t, s[0].Fields[1].Slice.IsBasicType).To(BeFalse())
		Expect(t, s[0].Fields[1].Slice.FieldName).To(Equal(""))
		Expect(t, s[0].Fields[1].Ptr).To(BeFalse())
		Expect(t, s[0].Fields[1].Type).To(Equal("known"))
	})
}

func TestStructFetcherWithMapsWithBasicKeys(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) TSF {
		return TSF{
			T: t,
			f: inspector.NewStructFetcher(nil, nil, nil),
		}
	})

	o.Spec("it parses and returns a field that is a map with a basic key type", func(t TSF) {
		src := `
package p
type x struct {
	a map[int]bool
	b map[int8]bool
	c map[int32]bool
	d map[int64]bool
	e map[uint]bool
	f map[uint8]bool
	g map[uint32]bool
	h map[uint64]bool
	i map[string]bool
	j map[float32]bool
	k map[float64]bool
	l map[bool]bool
	m map[byte]bool
	x map[unknown]bool
}
`

		fset := token.NewFileSet()
		n, err := parser.ParseFile(fset, "src.go", src, 0)
		Expect(t, err == nil).To(BeTrue())

		s, err := t.f.Parse(n)
		Expect(t, err == nil).To(BeTrue())
		Expect(t, s).To(HaveLen(1))
		Expect(t, s[0].Name).To(Equal("x"))
		Expect(t, s[0].Fields).To(HaveLen(13))

		for _, f := range s[0].Fields {
			Expect(t, f.Map.IsMap).To(BeTrue())
			Expect(t, f.Ptr).To(BeFalse())
		}
	})
}

func TestStructFetcherWithBlacklist(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) TSF {
		return TSF{
			T: t,
		}
	})

	o.Spec("blacklists the given struct.field combo", func(t TSF) {
		f := inspector.NewStructFetcher(map[string][]string{
			"x": {"a", "b"},
		}, nil, nil)
		src := `
package p
type x struct {
	i string
	j int
	a int
	b int
}
`
		fset := token.NewFileSet()
		n, err := parser.ParseFile(fset, "src.go", src, 0)
		Expect(t, err == nil).To(BeTrue())

		s, err := f.Parse(n)
		Expect(t, err == nil).To(BeTrue())
		Expect(t, s).To(HaveLen(1))
		Expect(t, s[0].Name).To(Equal("x"))
		Expect(t, s[0].Fields).To(HaveLen(2))

		Expect(t, s[0].Fields[0].Name).To(Equal("i"))
		Expect(t, s[0].Fields[0].Type).To(Equal("string"))

		Expect(t, s[0].Fields[1].Name).To(Equal("j"))
		Expect(t, s[0].Fields[1].Type).To(Equal("int"))
	})

	o.Spec("blacklists the given struct.field combo with wildcard structname", func(t TSF) {
		f := inspector.NewStructFetcher(map[string][]string{
			"*": {"a", "b"},
		}, nil, nil)
		src := `
package p
type x struct {
	i string
	j int
	a int
	b int
}
`
		fset := token.NewFileSet()
		n, err := parser.ParseFile(fset, "src.go", src, 0)
		Expect(t, err == nil).To(BeTrue())

		s, err := f.Parse(n)
		Expect(t, err == nil).To(BeTrue())
		Expect(t, s).To(HaveLen(1))
		Expect(t, s[0].Name).To(Equal("x"))
		Expect(t, s[0].Fields).To(HaveLen(2))

		Expect(t, s[0].Fields[0].Name).To(Equal("i"))
		Expect(t, s[0].Fields[0].Type).To(Equal("string"))

		Expect(t, s[0].Fields[1].Name).To(Equal("j"))
		Expect(t, s[0].Fields[1].Type).To(Equal("int"))
	})
}
