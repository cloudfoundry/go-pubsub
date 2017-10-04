package inspector_test

import (
	"testing"

	"code.cloudfoundry.org/go-pubsub/pubsub-gen/internal/inspector"
	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
)

type TL struct {
	*testing.T
	l inspector.Linker
}

func TestLinker(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) TL {
		return TL{
			T: t,
			l: inspector.NewLinker(),
		}
	})

	o.Spec("moves any known field types to PeerTypeFields", func(t TL) {
		m := map[string]inspector.Struct{
			"X": {Fields: []inspector.Field{
				{Name: "A", Type: "string"},
				{Name: "B", Type: "Y"},
				{Name: "C", Type: "Y"},
			}},
			"Y": {Fields: []inspector.Field{
				{Name: "A", Type: "string"},
				{Name: "B", Type: "int"},
			}},
		}
		t.l.Link(m, nil)

		Expect(t, m["X"].Fields).To(HaveLen(1))
		Expect(t, m["X"].Fields[0].Name).To(Equal("A"))

		Expect(t, m["X"].PeerTypeFields).To(HaveLen(2))
		Expect(t, m["X"].PeerTypeFields[0].Name).To(Equal("B"))
		Expect(t, m["X"].PeerTypeFields[1].Name).To(Equal("C"))
	})

	o.Spec("moves any known interface types to InterfaceTypeFields", func(t TL) {
		m := map[string]inspector.Struct{
			"X": {Fields: []inspector.Field{
				{Name: "A", Type: "string"},
				{Name: "B", Type: "Y"},
				{Name: "C", Type: "Y"},
			}},
			"Y": {Fields: []inspector.Field{
				{Name: "A", Type: "string"},
				{Name: "B", Type: "int"},
				{Name: "C", Type: "MyInterfaceThing"},
			}},
		}

		mi := map[string][]string{
			"MyInterfaceThing": {
				"X", "Y",
			},
		}
		c := m["Y"].Fields[2]

		t.l.Link(m, mi)
		Expect(t, m["Y"].Fields).To(HaveLen(2))
		Expect(t, m["Y"].InterfaceTypeFields).To(HaveLen(1))
		Expect(t, m["Y"].InterfaceTypeFields[c]).To(And(
			HaveLen(2),
			Contain("X", "Y"),
		))
	})

	o.Spec("adds non-basic slice types with given field type", func(t TL) {
		m := map[string]inspector.Struct{
			"X": {Fields: []inspector.Field{
				{Name: "A", Type: "Y", Slice: inspector.Slice{
					IsSlice:     true,
					IsBasicType: false,
					FieldName:   "B",
				}},
				{Name: "B", Type: "string", Slice: inspector.Slice{
					IsSlice:     true,
					IsBasicType: true,
				}},
			}},
			"Y": {Fields: []inspector.Field{
				{Name: "A", Type: "string"},
				{Name: "B", Type: "int"},
				{Name: "C", Type: "MyInterfaceThing"},
			}},
		}

		t.l.Link(m, nil)

		Expect(t, m["X"].Fields).To(HaveLen(2))
		Expect(t, m["X"].Fields[0].Name).To(Equal("A"))
		Expect(t, m["X"].Fields[0].Slice.IsSlice).To(BeTrue())
		Expect(t, m["X"].Fields[0].Slice.IsBasicType).To(BeFalse())
		Expect(t, m["X"].Fields[0].Slice.FieldName).To(Equal("B"))
		Expect(t, m["X"].Fields[0].Type).To(Equal("int")) // Comes from Y.B

		// Is left alone because it is a basic type
		Expect(t, m["X"].Fields[1].Name).To(Equal("B"))
		Expect(t, m["X"].Fields[1].Type).To(Equal("string"))
	})
}
