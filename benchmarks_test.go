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
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			unsub := p.Subscribe(newSpySubscrption(), randPath())
			unsub()
		}
	})
}

func BenchmarkPublishing(b *testing.B) {
	p := pubsub.New()
	a := &staticDataAssigner{}
	for i := 0; i < 100; i++ {
		p.Subscribe(newSpySubscrption(), randPath())
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

	b.RunParallel(func(b *testing.PB) {
		buf := make([]byte, 256)
		for b.Next() {
			rand.Read(buf)
			p.Publish(string(buf), a)

			for i := 0; i < 10; i++ {
				unsub := p.Subscribe(newSpySubscrption(), randPath())
				go func() {
					time.Sleep(time.Duration(rand.Intn(int(time.Millisecond))))
					unsub()
				}()
			}
		}
	})
}

func randPath() []string {
	var r []string
	for {
		i := rand.Intn(4)
		if i == 0 {
			return r
		}
		r = append(r, fmt.Sprintf("%d", i))
	}
}

type staticDataAssigner struct{}

func (r *staticDataAssigner) Assign(data interface{}, location []string) (keys []string, next interface{}) {
	return []string{"1", "2", "3", "4"}[:4-len(location)], data
}
