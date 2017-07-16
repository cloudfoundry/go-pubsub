package pubsub_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/apoydence/pubsub"
)

func BenchmarkSubscriptions(b *testing.B) {
	p := pubsub.New(&randTreeBuilder{}, &staticTreeTraverser{})
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			p.Subscribe(newSpySubscrption())
		}
	})
}

func BenchmarkPublishing(b *testing.B) {
	p := pubsub.New(&randTreeBuilder{}, &staticTreeTraverser{})
	for i := 0; i < 100; i++ {
		p.Subscribe(newSpySubscrption())
	}

	b.RunParallel(func(b *testing.PB) {
		buf := make([]byte, 256)
		for b.Next() {
			rand.Read(buf)
			p.Publish(Data(string(buf)))
		}
	})
}

func BenchmarkPublishingWhileSubscribing(b *testing.B) {
	p := pubsub.New(&randTreeBuilder{}, &staticTreeTraverser{})

	b.RunParallel(func(b *testing.PB) {
		buf := make([]byte, 256)
		for b.Next() {
			rand.Read(buf)
			p.Publish(Data(string(buf)))

			for i := 0; i < 10; i++ {
				p.Subscribe(newSpySubscrption())
			}
		}
	})
}

type randTreeBuilder struct {
}

func (r *randTreeBuilder) PlaceSubscription(sub pubsub.Subscription, location []string) (key string, ok bool) {
	i := rand.Intn(4)
	if i == 0 {
		return "", false
	}

	return fmt.Sprintf("%d", i), true
}

type staticTreeTraverser struct{}

func (r *staticTreeTraverser) Traverse(data pubsub.Data, location []string) (keys []string) {
	return []string{"1", "2", "3", "4"}[:4-len(location)]
}
