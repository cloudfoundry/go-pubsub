package pubsub_test

import (
	"sync"
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"code.cloudfoundry.org/go-pubsub"
	"code.cloudfoundry.org/go-pubsub/pubsub-gen/setters"
)

type TPS struct {
	*testing.T
	p            *pubsub.PubSub
	subscription *spySubscription
	spy          *spySubscription
	sub          func(interface{})
}

//go:generate go install code.cloudfoundry.org/go-pubsub/pubsub-gen
//go:generate $GOPATH/bin/pubsub-gen --output=$GOPATH/src/code.cloudfoundry.org/go-pubsub/gen_struct_test.go --pointer --struct-name=github.com/apoydence/pubsub.testStruct --traverser=testStructTrav --package=pubsub_test
type testStruct struct {
	a  int
	b  int
	aa *testStructA
	bb *testStructB
}

type testStructA struct {
	a int
}

type testStructB struct {
	b int
}

func TestPubSub(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)
	o.BeforeEach(func(t *testing.T) TPS {
		s, f := newSpySubscrption()

		return TPS{
			T:            t,
			sub:          f,
			subscription: s,
			p:            pubsub.New(),
		}
	})

	o.Spec("it writes to the correct subscription", func(t TPS) {
		sub1, f1 := newSpySubscrption()
		sub2, f2 := newSpySubscrption()
		sub3, f3 := newSpySubscrption()
		sub4, f4 := newSpySubscrption()

		t.p.Subscribe(f1, pubsub.WithPath(testStructTravCreatePath(&testStructFilter{
			a: setters.Int(1),
			b: setters.Int(2),
		})))
		t.p.Subscribe(f2, pubsub.WithPath(testStructTravCreatePath(&testStructFilter{
			a: setters.Int(1),
			b: setters.Int(3),
		})))
		t.p.Subscribe(f3, pubsub.WithPath(testStructTravCreatePath(&testStructFilter{
			a: setters.Int(1),
		})))
		t.p.Subscribe(f4, pubsub.WithPath(testStructTravCreatePath(&testStructFilter{
			a: setters.Int(2),
			b: setters.Int(3),
			aa: &testStructAFilter{
				a: setters.Int(4),
			},
		})))

		data := &testStruct{a: 1, b: 2}
		otherData := &testStruct{a: 2, b: 3, aa: &testStructA{a: 4}}
		t.p.Publish(data, testStructTravTraverse)
		t.p.Publish(otherData, testStructTravTraverse)

		Expect(t, sub1.data).To(HaveLen(1))
		Expect(t, sub2.data).To(HaveLen(0))
		Expect(t, sub3.data).To(HaveLen(1))
		Expect(t, sub4.data).To(HaveLen(1))

		Expect(t, sub1.data[0]).To(Equal(data))
		Expect(t, sub3.data[0]).To(Equal(data))
		Expect(t, sub4.data[0]).To(Equal(otherData))
	})

	o.Spec("it does not write to a subscription after it unsubscribes", func(t TPS) {
		sub, f := newSpySubscrption()
		unsubscribe := t.p.Subscribe(f, pubsub.WithPath(testStructTravCreatePath(&testStructFilter{
			a: setters.Int(1),
			b: setters.Int(2),
		})))

		unsubscribe()

		t.p.Publish(&testStruct{a: 1, b: 2}, testStructTravTraverse)
		Expect(t, sub.data).To(HaveLen(0))
	})
}

func TestPubSubWithShardID(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)
	o.BeforeEach(func(t *testing.T) TPS {
		s, f := newSpySubscrption()

		return TPS{
			T:            t,
			subscription: s,
			sub:          f,
			p:            pubsub.New(),
		}
	})

	o.Spec("it splits data between same shardIDs", func(t TPS) {
		sub1, f1 := newSpySubscrption()
		sub2, f2 := newSpySubscrption()
		sub3, f3 := newSpySubscrption()
		sub4, f4 := newSpySubscrption()
		sub5, f5 := newSpySubscrption()

		t.p.Subscribe(f1,
			pubsub.WithShardID("1"),
			pubsub.WithPath(testStructTravCreatePath(&testStructFilter{
				a: setters.Int(1),
			})),
		)

		t.p.Subscribe(f2,
			pubsub.WithShardID("1"),
			pubsub.WithPath(testStructTravCreatePath(&testStructFilter{
				a: setters.Int(1),
			})),
		)
		t.p.Subscribe(f3,
			pubsub.WithShardID("2"),
			pubsub.WithPath(testStructTravCreatePath(&testStructFilter{
				a: setters.Int(1),
			})),
		)

		t.p.Subscribe(f4,
			pubsub.WithPath(testStructTravCreatePath(&testStructFilter{
				a: setters.Int(1),
			})),
		)

		t.p.Subscribe(f5,
			pubsub.WithPath(testStructTravCreatePath(&testStructFilter{
				a: setters.Int(1),
			})),
		)

		for i := 0; i < 100; i++ {
			t.p.Publish(&testStruct{a: 1, b: 2}, testStructTravTraverse)
		}

		Expect(t, len(sub1.data)).To(And(BeAbove(0), BeBelow(99)))
		Expect(t, len(sub2.data)).To(And(BeAbove(0), BeBelow(99)))
		Expect(t, len(sub3.data)).To(Equal(100))
		Expect(t, len(sub4.data)).To(Equal(100))
		Expect(t, len(sub5.data)).To(Equal(100))
	})
}

type spySubscription struct {
	mu   sync.Mutex
	data []interface{}
}

func newSpySubscrption() (*spySubscription, func(interface{})) {
	s := &spySubscription{}
	return s, func(data interface{}) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.data = append(s.data, data)
	}
}
