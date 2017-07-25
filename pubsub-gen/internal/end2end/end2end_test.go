package end2end_test

import (
	"flag"
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"github.com/apoydence/pubsub"
	. "github.com/apoydence/pubsub/pubsub-gen/internal/end2end"
	"github.com/apoydence/pubsub/pubsub-gen/setters"
)

func TestEnd2End(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)
	flag.Parse()

	o.Spec("routes data as expected", func(t *testing.T) {
		ps := pubsub.New()
		s := StructTraverser{}
		sub1 := &mockSubscription{}
		sub2 := &mockSubscription{}
		sub3 := &mockSubscription{}
		sub4 := &mockSubscription{}
		sub5 := &mockSubscription{}

		ps.Subscribe(sub1, s.CreatePath(nil))
		ps.Subscribe(sub2, s.CreatePath(&XFilter{
			I: setters.Int(1),
			Y1: &YFilter{
				I: setters.Int(1),
				J: setters.String("a"),
			},
		}))
		ps.Subscribe(sub3, s.CreatePath(&XFilter{
			Y1: &YFilter{
				J: setters.String("b"),
			},
		}))
		ps.Subscribe(sub4, s.CreatePath(&XFilter{
			Y2: &YFilter{},
		}))
		ps.Subscribe(sub5, s.CreatePath(&XFilter{
			M_M2: &M2Filter{
				B: setters.Int(2),
			},
		}))

		ps.Publish(&X{I: 1, J: "a", Y1: Y{I: 1, J: "a"}, Y2: &Y{I: 1, J: "a"}}, s)
		ps.Publish(&X{I: 1, J: "a", Y1: Y{I: 2, J: "b"}, Y2: &Y{I: 1, J: "a"}}, s)
		ps.Publish(&X{I: 1, J: "x", Y1: Y{I: 2, J: "b"}}, s)
		ps.Publish(&X{I: 1, J: "x", Y1: Y{I: 2, J: "b"}, M: M2{1, 2}}, s)

		Expect(t, sub1.callCount).To(Equal(4))
		Expect(t, sub2.callCount).To(Equal(1))
		Expect(t, sub3.callCount).To(Equal(3))
		Expect(t, sub4.callCount).To(Equal(2))
		Expect(t, sub5.callCount).To(Equal(1))
	})
}

type mockSubscription struct {
	callCount int
}

func (m *mockSubscription) Write(data interface{}) {
	m.callCount++
}

//go:generate go install github.com/apoydence/pubsub/pubsub-gen
//go:generate $GOPATH/bin/pubsub-gen --struct-name=github.com/apoydence/pubsub/pubsub-gen/internal/end2end.X --package=end2end_test --traverser=StructTraverser --output=$GOPATH/src/github.com/apoydence/pubsub/pubsub-gen/internal/end2end/generated_traverser_test.go --pointer --interfaces={"message":["M1","M2"]} --include-pkg-name=true --imports=github.com/apoydence/pubsub/pubsub-gen/internal/end2end
