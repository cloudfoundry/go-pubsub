package pubsub_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/apoydence/pubsub"
)

func BenchmarkSubscriptions(b *testing.B) {
	p := pubsub.New(&randSubscriptionEnroller{}, &staticDataAssigner{})
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			unsub := p.Subscribe(newSpySubscrption())
			unsub()
		}
	})
}

func BenchmarkPublishing(b *testing.B) {
	p := pubsub.New(&randSubscriptionEnroller{}, &staticDataAssigner{})
	for i := 0; i < 100; i++ {
		p.Subscribe(newSpySubscrption())
	}

	b.RunParallel(func(b *testing.PB) {
		buf := make([]byte, 256)
		for b.Next() {
			rand.Read(buf)
			p.Publish(string(buf))
		}
	})
}

func BenchmarkPublishingWhileSubscribing(b *testing.B) {
	p := pubsub.New(&randSubscriptionEnroller{}, &staticDataAssigner{})

	b.RunParallel(func(b *testing.PB) {
		buf := make([]byte, 256)
		for b.Next() {
			rand.Read(buf)
			p.Publish(string(buf))

			for i := 0; i < 10; i++ {
				unsub := p.Subscribe(newSpySubscrption())
				go func() {
					time.Sleep(time.Duration(rand.Intn(int(time.Millisecond))))
					unsub()
				}()
			}
		}
	})
}

type randSubscriptionEnroller struct {
}

func (r *randSubscriptionEnroller) Enroll(sub pubsub.Subscription, location []string) (key string, ok bool) {
	i := rand.Intn(4)
	if i == 0 {
		return "", false
	}

	return fmt.Sprintf("%d", i), true
}

type staticDataAssigner struct{}

func (r *staticDataAssigner) Assign(data interface{}, location []string) (keys []string) {
	return []string{"1", "2", "3", "4"}[:4-len(location)]
}
