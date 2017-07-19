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

	o.Spec("it invokes the TreeTraverser for each level", func(t TPS) {
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
		sub4 := newSpySubscrption()
		t.p.Subscribe(sub1, []string{"a", "b", "c"})
		t.p.Subscribe(sub2, []string{"a", "b", "d"})
		t.p.Subscribe(sub3, []string{"a", "b"})
		t.p.Subscribe(sub4, []string{"j", "k"})

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
		t.p.Publish("some-other-data", pubsub.LinearTreeTraverser([]string{"j", "k"}))

		Expect(t, sub1.data).To(HaveLen(1))
		Expect(t, sub2.data).To(HaveLen(0))
		Expect(t, sub3.data).To(HaveLen(1))
		Expect(t, sub4.data).To(HaveLen(1))

		Expect(t, sub1.data[0]).To(Equal("some-data"))
		Expect(t, sub3.data[0]).To(Equal("some-data"))
		Expect(t, sub4.data[0]).To(Equal("some-other-data"))

		Expect(t, t.treeTraverser.data).To(HaveLen(len(t.treeTraverser.keys)))
		Expect(t, t.treeTraverser.data[1]).To(Equal("some-data"))
	})

	o.Spec("it uses the new TreeTraverser when given one", func(t TPS) {
		sub := newSpySubscrption()
		t.p.Subscribe(sub, []string{"a", "b", "c"})

		f := &fakePaths{}

		t.treeTraverser.ret = func(s *spyTreeTraverser, results []string) pubsub.Paths {
			f.a = s
			f.paths = results
			return f
		}

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

		Expect(t, f.ids).To(Contain(2))
	})

	o.Spec("it does not write to a subscription after it unsubscribes", func(t TPS) {
		sub := newSpySubscrption()
		t.treeTraverser.keys = map[string][]string{
			"": nil,
		}

		unsubscribe := t.p.Subscribe(sub, nil)
		unsubscribe()
		t.p.Publish("some-data", t.treeTraverser)
		Expect(t, sub.data).To(HaveLen(0))
	})
}

func ExamplePubSub() {
	ps := pubsub.New()
	subscription := func(name string) pubsub.SubscriptionFunc {
		return func(data interface{}) {
			fmt.Printf("%s -> %v\n", name, data)
		}
	}

	ps.Subscribe(subscription("sub-1"), []string{"a", "b", "c"})
	ps.Subscribe(subscription("sub-2"), []string{"a", "b", "d"})
	ps.Subscribe(subscription("sub-3"), []string{"a", "b", "e"})

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
