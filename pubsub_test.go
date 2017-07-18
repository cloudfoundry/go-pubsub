package pubsub_test

import (
	"log"
	"strings"
	"sync"
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"github.com/apoydence/pubsub"
)

type TPS struct {
	*testing.T
	p             *pubsub.PubSub
	treeBuilder   *spySubscriptionEnroller
	treeTraverser *spyDataAssigner
	subscription  *spySubscription
}

func TestPubSub(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)
	o.BeforeEach(func(t *testing.T) TPS {
		spyB := newSpySubscriptionEnroller()
		spyT := newSpyDataAssigner()

		return TPS{
			T:             t,
			subscription:  newSpySubscrption(),
			p:             pubsub.New(),
			treeBuilder:   spyB,
			treeTraverser: spyT,
		}
	})

	o.Spec("it invokes the SubscriptionEnroller for each level of a subscription", func(t TPS) {
		t.treeBuilder.keys = map[string]string{
			"":    "a",
			"a":   "b",
			"a-b": "",
		}
		t.p.Subscribe(t.subscription, t.treeBuilder)

		Expect(t, t.treeBuilder.locations).To(HaveLen(3))

		for k := range t.treeBuilder.keys {
			Expect(t, t.treeBuilder.locations).To(Contain(k))
		}
		for _, s := range t.treeBuilder.subs {
			Expect(t, s).To(Equal(t.subscription))
		}
	})

	o.Spec("it invokes the DataAssigner for each level", func(t TPS) {
		t.treeTraverser.keys = map[string][]string{
			"":      []string{"a", "b"},
			"a":     []string{"a", "b"},
			"b":     []string{"a", "b"},
			"a-a":   []string{"a", "b"},
			"a-b":   nil,
			"b-a":   nil,
			"b-b":   nil,
			"a-a-a": nil,
			"a-a-b": nil,
		}

		t.p.Publish("x", t.treeTraverser)

		Expect(t, t.treeTraverser.locations).To(HaveLen(9))
		for k := range t.treeTraverser.keys {
			Expect(t, t.treeTraverser.locations).To(Contain(k))
		}
	})

	o.Spec("it writes to the correct subscription", func(t TPS) {
		sub1 := newSpySubscrption()
		sub2 := newSpySubscrption()
		sub3 := newSpySubscrption()
		t.treeBuilder.keys = map[string]string{
			"":      "a",
			"a":     "b",
			"a-b":   "c",
			"a-b-c": "",
		}
		t.p.Subscribe(sub1, t.treeBuilder)

		t.treeBuilder.keys = map[string]string{
			"":      "a",
			"a":     "b",
			"a-b":   "d",
			"a-b-d": "",
		}
		t.p.Subscribe(sub2, t.treeBuilder)

		t.treeBuilder.keys = map[string]string{
			"":  "a",
			"a": "",
		}
		t.p.Subscribe(sub3, t.treeBuilder)

		t.treeTraverser.keys = map[string][]string{
			"":      []string{"a", "x"},
			"x":     nil,
			"a":     []string{"b", "y"},
			"a-y":   nil,
			"a-b":   []string{"c", "z"},
			"a-b-z": nil,
			"a-b-c": nil,
		}
		t.p.Publish("some-data", t.treeTraverser)

		Expect(t, sub1.data).To(HaveLen(1))
		Expect(t, sub2.data).To(HaveLen(0))
		Expect(t, sub3.data).To(HaveLen(1))

		Expect(t, sub1.data[0]).To(Equal("some-data"))
		Expect(t, sub3.data[0]).To(Equal("some-data"))
	})

	o.Spec("it does not write to a subscription after it unsubscribes", func(t TPS) {
		sub := newSpySubscrption()
		t.treeBuilder.keys = map[string]string{
			"": "",
		}
		t.treeTraverser.keys = map[string][]string{
			"": nil,
		}

		unsubscribe := t.p.Subscribe(sub, t.treeBuilder)
		unsubscribe()
		t.p.Publish("some-data", t.treeTraverser)
		Expect(t, sub.data).To(HaveLen(0))
	})
}

type spyDataAssigner struct {
	keys      map[string][]string
	locations []string
	data      []interface{}
}

func newSpyDataAssigner() *spyDataAssigner {
	return &spyDataAssigner{}
}

func (s *spyDataAssigner) Assign(data interface{}, location []string) []string {
	s.data = append(s.data, data)

	key := strings.Join(location, "-")
	s.locations = append(s.locations, key)
	result, ok := s.keys[key]
	if !ok {
		log.Panicf("unknown location: %s", key)
	}

	return result
}

type spySubscriptionEnroller struct {
	keys      map[string]string
	subs      []pubsub.Subscription
	locations []string
}

func newSpySubscriptionEnroller() *spySubscriptionEnroller {
	return &spySubscriptionEnroller{}
}

func (s *spySubscriptionEnroller) Enroll(sub pubsub.Subscription, location []string) (string, bool) {
	s.subs = append(s.subs, sub)

	key := strings.Join(location, "-")
	s.locations = append(s.locations, key)
	result, ok := s.keys[key]

	if !ok {
		log.Panicf("unknown location: %s", key)
	}

	if result == "" {
		return "", false
	}

	return result, true
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
