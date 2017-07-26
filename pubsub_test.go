package pubsub_test

import (
	"fmt"
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
	treeTraverser *spyTreeTraverser
	subscription  *spySubscription
}

func TestPubSub(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)
	o.BeforeEach(func(t *testing.T) TPS {
		spyT := newSpyTreeTraverser()

		return TPS{
			T:             t,
			subscription:  newSpySubscrption(),
			p:             pubsub.New(),
			treeTraverser: spyT,
		}
	})

	o.Spec("it stops traversing upon reaching an empty node", func(t TPS) {
		t.treeTraverser.keys = map[string][]string{
			"":      {"a", "b"},
			"a":     {"a", "b"},
			"b":     {"a", "b"},
			"a-a":   {"a", "b"},
			"a-b":   nil,
			"b-a":   nil,
			"b-b":   nil,
			"a-a-a": nil,
			"a-a-b": nil,
		}

		t.p.Publish("x", t.treeTraverser)
		Expect(t, t.treeTraverser.locations).To(HaveLen(1))

		t.p.Subscribe(newSpySubscrption(), pubsub.WithPath("a", "b"))
		t.treeTraverser.locations = nil
		t.p.Publish("x", t.treeTraverser)
		Expect(t, t.treeTraverser.locations).To(HaveLen(3))
		Expect(t, t.treeTraverser.locations).To(Contain("", "a", "a-b"))
	})

	o.Spec("it writes to the correct subscription", func(t TPS) {
		sub1 := newSpySubscrption()
		sub2 := newSpySubscrption()
		sub3 := newSpySubscrption()
		sub4 := newSpySubscrption()
		t.p.Subscribe(sub1, pubsub.WithPath("a", "b", "c"))
		t.p.Subscribe(sub2, pubsub.WithPath("a", "b", "d"))
		t.p.Subscribe(sub3, pubsub.WithPath("a", "b"))
		t.p.Subscribe(sub4, pubsub.WithPath("j", "k"))

		t.treeTraverser.keys = map[string][]string{
			"":      {"a", "x"},
			"x":     nil,
			"a":     {"b", "y"},
			"a-y":   nil,
			"a-b":   {"c", "z"},
			"a-b-z": nil,
			"a-b-c": nil,
		}
		t.p.Publish("some-data", t.treeTraverser)
		t.p.Publish("some-other-data", pubsub.LinearTreeTraverser([]string{"j", "k"}))

		Expect(t, sub1.data).To(HaveLen(1))
		Expect(t, sub2.data).To(HaveLen(0))
		Expect(t, sub3.data).To(HaveLen(1))
		Expect(t, sub4.data).To(HaveLen(1))

		Expect(t, sub1.data[0]).To(Equal("some-data"))
		Expect(t, sub3.data[0]).To(Equal("some-data"))
		Expect(t, sub4.data[0]).To(Equal("some-other-data"))

		Expect(t, t.treeTraverser.data).To(HaveLen(4))
		Expect(t, t.treeTraverser.data[1]).To(Equal("some-data"))
	})

	o.Spec("does not write to the subscription twice", func(t TPS) {
		sub := newSpySubscrption()
		t.p.Subscribe(sub, pubsub.WithPath("a"))

		t.treeTraverser.keys = map[string][]string{
			"":  {"a", "a"},
			"a": nil,
		}
		t.p.Publish("some-data", t.treeTraverser)

		Expect(t, sub.data).To(HaveLen(1))
		Expect(t, sub.data[0]).To(Equal("some-data"))
	})

	o.Spec("it uses the new TreeTraverser when given one", func(t TPS) {
		sub := newSpySubscrption()
		t.p.Subscribe(sub, pubsub.WithPath("a", "b", "c"))

		f := &fakePaths{}

		t.treeTraverser.ret = func(s *spyTreeTraverser, results []string) pubsub.Paths {
			f.a = s
			f.paths = results
			return f
		}

		t.treeTraverser.keys = map[string][]string{
			"":      {"a", "x"},
			"x":     nil,
			"a":     {"b", "y"},
			"a-y":   nil,
			"a-b":   {"c", "z"},
			"a-b-z": nil,
			"a-b-c": nil,
		}
		t.p.Publish("some-data", t.treeTraverser)

		Expect(t, f.ids).To(Contain(2))
	})

	o.Spec("it does not write to a subscription after it unsubscribes", func(t TPS) {
		sub := newSpySubscrption()
		t.treeTraverser.keys = map[string][]string{
			"":  {"a"},
			"a": nil,
		}

		unsubscribe := t.p.Subscribe(sub, pubsub.WithPath("a"))
		unsubscribe()
		t.p.Publish("some-data", t.treeTraverser)
		Expect(t, sub.data).To(HaveLen(0))
	})
}

func TestPubSubWithShardID(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)
	o.BeforeEach(func(t *testing.T) TPS {
		spyT := newSpyTreeTraverser()

		return TPS{
			T:             t,
			subscription:  newSpySubscrption(),
			p:             pubsub.New(),
			treeTraverser: spyT,
		}
	})

	o.Spec("it splits data between same shardIDs", func(t TPS) {
		sub1 := newSpySubscrption()
		sub2 := newSpySubscrption()
		sub3 := newSpySubscrption()
		sub4 := newSpySubscrption()
		sub5 := newSpySubscrption()
		t.treeTraverser.keys = map[string][]string{
			"":  {"a"},
			"a": nil,
		}

		t.p.Subscribe(
			sub1,
			pubsub.WithShardID("1"),
			pubsub.WithPath("a"),
		)
		t.p.Subscribe(
			sub2,
			pubsub.WithShardID("1"),
			pubsub.WithPath("a"),
		)
		t.p.Subscribe(
			sub3,
			pubsub.WithShardID("2"),
			pubsub.WithPath("a"),
		)
		t.p.Subscribe(sub4, pubsub.WithPath("a"))
		t.p.Subscribe(sub5, pubsub.WithPath("a"))

		for i := 0; i < 100; i++ {
			t.p.Publish("some-data", t.treeTraverser)
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

	ps.Subscribe(subscription("sub-1"), pubsub.WithPath("a", "b", "c"))
	ps.Subscribe(subscription("sub-2"), pubsub.WithPath("a", "b", "d"))
	ps.Subscribe(subscription("sub-3"), pubsub.WithPath("a", "b", "e"))

	ps.Publish("data-1", pubsub.LinearTreeTraverser([]string{"a", "b"}))
	ps.Publish("data-2", pubsub.LinearTreeTraverser([]string{"a", "b", "c"}))
	ps.Publish("data-3", pubsub.LinearTreeTraverser([]string{"a", "b", "d"}))
	ps.Publish("data-3", pubsub.LinearTreeTraverser([]string{"x", "y"}))

	// Output:
	// sub-1 -> data-2
	// sub-2 -> data-3
}

type spyTreeTraverser struct {
	keys      map[string][]string
	locations []string
	data      []interface{}
	id        int
	ret       func(*spyTreeTraverser, []string) pubsub.Paths
}

func newSpyTreeTraverser() *spyTreeTraverser {
	return &spyTreeTraverser{
		ret: func(s *spyTreeTraverser, results []string) pubsub.Paths {
			return pubsub.FlatPaths(results)
		},
	}
}

func (s *spyTreeTraverser) Traverse(data interface{}, location []string) pubsub.Paths {
	s.data = append(s.data, data)

	key := strings.Join(location, "-")
	s.locations = append(s.locations, key)
	result, ok := s.keys[key]
	if !ok {
		log.Panicf("unknown location: %s", key)
	}

	return s.ret(s, result)
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

type fakePaths struct {
	paths []string
	a     *spyTreeTraverser
	ids   []int
}

func (s *fakePaths) At(idx int) (path string, nextTraverser pubsub.TreeTraverser, ok bool) {
	if len(s.paths) <= idx {
		return "", nil, false
	}

	s.ids = append(s.ids, s.a.id)

	return s.paths[idx], &spyTreeTraverser{
		keys:      s.a.keys,
		locations: s.a.locations,
		data:      s.a.data,
		ret:       s.a.ret,
		id:        s.a.id + 1,
	}, true
}
