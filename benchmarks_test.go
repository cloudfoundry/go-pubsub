package pubsub_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/apoydence/pubsub"
)

func BenchmarkSubscriptions(b *testing.B) {
	p := pubsub.New()
	e := &randSubscriptionEnroller{}
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			unsub := p.Subscribe(newSpySubscrption(), e)
			unsub()
		}
	})
}

func BenchmarkPublishing(b *testing.B) {
	p := pubsub.New()
	a := &staticDataAssigner{}
	e := &randSubscriptionEnroller{}
	for i := 0; i < 100; i++ {
		p.Subscribe(newSpySubscrption(), e)
	}

	b.RunParallel(func(b *testing.PB) {
		buf := make([]byte, 256)
		for b.Next() {
			rand.Read(buf)
			p.Publish(string(buf), a)
		}
	})
}

func BenchmarkPublishingWhileSubscribing(b *testing.B) {
	p := pubsub.New()
	a := &staticDataAssigner{}
	e := &randSubscriptionEnroller{}

	b.RunParallel(func(b *testing.PB) {
		buf := make([]byte, 256)
		for b.Next() {
			rand.Read(buf)
			p.Publish(string(buf), a)

			for i := 0; i < 10; i++ {
				unsub := p.Subscribe(newSpySubscrption(), e)
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
