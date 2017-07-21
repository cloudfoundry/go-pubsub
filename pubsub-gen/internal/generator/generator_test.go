package generator_test

import (
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"github.com/apoydence/pubsub/pubsub-gen/internal/generator"
	"github.com/apoydence/pubsub/pubsub-gen/internal/inspector"
)

type TG struct {
	*testing.T
	g generator.Generator
}

func TestWriter(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) TG {
		return TG{
			T: t,
			g: generator.New(),
		}
	})

	o.Group("flat struct", func() {
		o.Spec("it generates functions for each field in a struct", func(t TG) {
			m := map[string]inspector.Struct{
				"X": {Fields: []inspector.Field{
					{Name: "A", Type: "string"},
					{Name: "B", Type: "int"},
				}},
			}
			src, err := t.g.Generate(m, "mypack", "myassigner", "X", false)
			Expect(t, err == nil).To(BeTrue())

			_ = src
			// TODO: We need to solve how to better test all this
		})
	})
}
