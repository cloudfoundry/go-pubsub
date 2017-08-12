package pubsub_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"github.com/apoydence/pubsub"
	"github.com/apoydence/pubsub/pubsub-gen/setters"
)

type TPS struct {
	*testing.T
	p            *pubsub.PubSub
	subscription *spySubscription
	trav         testStructTrav
}

// go:generate pubsub-gen --output=$GOPATH/src/github.com/apoydence/pubsub/gen_struct_test.go --pointer --struct-name=github.com/apoydence/pubsub.testStruct --traverser=testStructTrav --package=pubsub_test
type testStruct struct {
	a int
	b int
}

func TestPubSub(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)
	o.BeforeEach(func(t *testing.T) TPS {
		trav := NewTestStructTrav()

		return TPS{
			T:            t,
			subscription: newSpySubscrption(),
			p:            pubsub.New(),
			trav:         trav,
		}
	})

	o.Spec("it writes to the correct subscription", func(t TPS) {
		sub1 := newSpySubscrption()
		sub2 := newSpySubscrption()
		sub3 := newSpySubscrption()
		sub4 := newSpySubscrption()

		t.p.Subscribe(sub1, pubsub.WithPath(t.trav.CreatePath(&testStructFilter{
			a: setters.Int(1),
			b: setters.Int(2),
		})))
		t.p.Subscribe(sub2, pubsub.WithPath(t.trav.CreatePath(&testStructFilter{
			a: setters.Int(1),
			b: setters.Int(3),
		})))
		t.p.Subscribe(sub3, pubsub.WithPath(t.trav.CreatePath(&testStructFilter{
			a: setters.Int(1),
		})))
		t.p.Subscribe(sub4, pubsub.WithPath(t.trav.CreatePath(&testStructFilter{
			a: setters.Int(6),
			b: setters.Int(8),
		})))

		data := &testStruct{a: 1, b: 2}
		otherData := &testStruct{a: 6, b: 8}
		t.p.Publish(data, t.trav)
		t.p.Publish(otherData, t.trav)

		Expect(t, sub1.data).To(HaveLen(1))
		Expect(t, sub2.data).To(HaveLen(0))
		Expect(t, sub3.data).To(HaveLen(1))
		Expect(t, sub4.data).To(HaveLen(1))

		Expect(t, sub1.data[0]).To(Equal(data))
		Expect(t, sub3.data[0]).To(Equal(data))
		Expect(t, sub4.data[0]).To(Equal(otherData))
	})

	o.Spec("it does not write to a subscription after it unsubscribes", func(t TPS) {
		sub := newSpySubscrption()
		unsubscribe := t.p.Subscribe(sub, pubsub.WithPath(t.trav.CreatePath(&testStructFilter{
			a: setters.Int(1),
			b: setters.Int(2),
		})))

		unsubscribe()

		t.p.Publish(&testStruct{a: 1, b: 2}, t.trav)
		Expect(t, sub.data).To(HaveLen(0))
	})
}

func TestPubSubWithShardID(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)
	o.BeforeEach(func(t *testing.T) TPS {
		trav := NewTestStructTrav()

		return TPS{
			T:            t,
			subscription: newSpySubscrption(),
			p:            pubsub.New(),
			trav:         trav,
		}
	})

	o.Spec("it splits data between same shardIDs", func(t TPS) {
		sub1 := newSpySubscrption()
		sub2 := newSpySubscrption()
		sub3 := newSpySubscrption()
		sub4 := newSpySubscrption()
		sub5 := newSpySubscrption()

		t.p.Subscribe(sub1,
			pubsub.WithShardID("1"),
			pubsub.WithPath(t.trav.CreatePath(&testStructFilter{
				a: setters.Int(1),
			})),
		)

		t.p.Subscribe(sub2,
			pubsub.WithShardID("1"),
			pubsub.WithPath(t.trav.CreatePath(&testStructFilter{
				a: setters.Int(1),
			})),
		)
		t.p.Subscribe(sub3,
			pubsub.WithShardID("2"),
			pubsub.WithPath(t.trav.CreatePath(&testStructFilter{
				a: setters.Int(1),
			})),
		)

		t.p.Subscribe(sub4,
			pubsub.WithPath(t.trav.CreatePath(&testStructFilter{
				a: setters.Int(1),
			})),
		)

		t.p.Subscribe(sub5,
			pubsub.WithPath(t.trav.CreatePath(&testStructFilter{
				a: setters.Int(1),
			})),
		)

		for i := 0; i < 100; i++ {
			t.p.Publish(&testStruct{a: 1, b: 2}, t.trav)
		}

		Expect(t, len(sub1.data)).To(And(BeAbove(0), BeBelow(99)))
		Expect(t, len(sub2.data)).To(And(BeAbove(0), BeBelow(99)))
		Expect(t, len(sub3.data)).To(Equal(100))
		Expect(t, len(sub4.data)).To(Equal(100))
		Expect(t, len(sub5.data)).To(Equal(100))
	})
}

func Example() {
	ps := pubsub.New()
	subscription := func(name string) pubsub.SubscriptionFunc {
		return func(data interface{}) {
			fmt.Printf("%s -> %v\n", name, data)
		}
	}

	ps.Subscribe(subscription("sub-1"), pubsub.WithPath([]string{"a", "b", "c"}))
	ps.Subscribe(subscription("sub-2"), pubsub.WithPath([]string{"a", "b", "d"}))
	ps.Subscribe(subscription("sub-3"), pubsub.WithPath([]string{"a", "b", "e"}))

	ps.Publish("data-1", pubsub.LinearTreeTraverser([]string{"a", "b"}))
	ps.Publish("data-2", pubsub.LinearTreeTraverser([]string{"a", "b", "c"}))
	ps.Publish("data-3", pubsub.LinearTreeTraverser([]string{"a", "b", "d"}))
	ps.Publish("data-3", pubsub.LinearTreeTraverser([]string{"x", "y"}))

	// Output:
	// sub-1 -> data-2
	// sub-2 -> data-3
}

type spySubscription struct {
	mu   sync.Mutex
	data []interface{}
}

func newSpySubscrption() *spySubscription {
	return &spySubscription{}
}

func (s *spySubscription) Write(data interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, data)
}
