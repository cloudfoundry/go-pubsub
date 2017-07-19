package pubsub_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/apoydence/pubsub"
)

func BenchmarkPublishing(b *testing.B) {
	b.StopTimer()
	p := pubsub.New()
	for i := 0; i < 100; i++ {
		p.Subscribe(newSpySubscrption(), randPath())
	}
	data := randData()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		p.Publish("data", pubsub.LinearTreeTraverser(data[i%len(data)]))
	}
}

func BenchmarkSubscriptions(b *testing.B) {
	b.StopTimer()
	p := pubsub.New()
	data := randData()
	b.StartTimer()

	b.RunParallel(func(b *testing.PB) {
		i := rand.Int()
		for b.Next() {
			unsub := p.Subscribe(newSpySubscrption(), data[i%len(data)])
			unsub()
			i++
		}
	})
}

func BenchmarkPublishingParallel(b *testing.B) {
	b.StopTimer()
	p := pubsub.New()
	for i := 0; i < 100; i++ {
		p.Subscribe(newSpySubscrption(), randPath())
	}
	data := randData()
	b.StartTimer()

	b.RunParallel(func(b *testing.PB) {
		i := rand.Int()
		for b.Next() {
			p.Publish("data", pubsub.LinearTreeTraverser(data[i%len(data)]))
			i++
		}
	})
}

func BenchmarkPublishingWhileSubscribing(b *testing.B) {
	b.StopTimer()
	p := pubsub.New()
	data := randData()

	var wg sync.WaitGroup
	for x := 0; x < 5; x++ {
		wg.Add(1)
		go func() {
			wg.Done()
			for i := 0; ; i++ {
				unsub := p.Subscribe(newSpySubscrption(), data[i%len(data)])
				if i%2 == 0 {
					unsub()
				}
			}
		}()
	}

	wg.Wait()
	b.StartTimer()

	b.RunParallel(func(b *testing.PB) {
		i := rand.Int()
		for b.Next() {
			p.Publish("data", pubsub.LinearTreeTraverser(data[i%len(data)]))
			i++
		}
	})
}

func randPath() []string {
	var r []string
	for i := 0; i < 10; i++ {
		r = append(r, fmt.Sprintf("%d", rand.Intn(10)))
	}
	return r
}

func randData() [][]string {
	var r [][]string
	for i := 0; i < 100000; i++ {
		r = append(r, nil)
		for j := 0; j < 10; j++ {
			r[i] = append(r[i], fmt.Sprintf("%d", rand.Intn(10)))
		}
	}
	return r
}
