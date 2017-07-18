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
	treeTraverser *spyDataAssigner
	subscription  *spySubscription
}

func TestPubSub(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)
	o.BeforeEach(func(t *testing.T) TPS {
		spyT := newSpyDataAssigner()

		return TPS{
			T:             t,
			subscription:  newSpySubscrption(),
			p:             pubsub.New(),
			treeTraverser: spyT,
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
		t.p.Subscribe(sub1, []string{"a", "b", "c"})
		t.p.Subscribe(sub2, []string{"a", "b", "d"})
		t.p.Subscribe(sub3, []string{"a", "b"})

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

		Expect(t, t.treeTraverser.data).To(HaveLen(len(t.treeTraverser.keys)))
		Expect(t, t.treeTraverser.data[1]).To(Equal("ome-data"))
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

type spyDataAssigner struct {
	keys      map[string][]string
	locations []string
	data      []interface{}
}

func newSpyDataAssigner() *spyDataAssigner {
	return &spyDataAssigner{}
}

func (s *spyDataAssigner) Assign(data interface{}, location []string) ([]string, interface{}) {
	s.data = append(s.data, data)

	key := strings.Join(location, "-")
	s.locations = append(s.locations, key)
	result, ok := s.keys[key]
	if !ok {
		log.Panicf("unknown location: %s", key)
	}

	if len(data.(string)) == 0 {
		return result, data
	}

	return result, data.(string)[1:]
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
