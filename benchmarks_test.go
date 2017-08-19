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
		_, f := newSpySubscrption()
		p.Subscribe(f, pubsub.WithPath(randPath()))
	}
	data := randData()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		p.Publish("data", pubsub.LinearTreeTraverser(data[i%len(data)]))
	}
}

func BenchmarkPublishingStructs(b *testing.B) {
	b.StopTimer()
	p := pubsub.New()
	for i := 0; i < 100; i++ {
		_, f := newSpySubscrption()
		p.Subscribe(f, pubsub.WithPath(randPath()))
	}
	data := randStructs()
	st := StructTraverser{}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		p.Publish(data[i%len(data)], st)
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
			_, f := newSpySubscrption()
			unsub := p.Subscribe(f, pubsub.WithPath(data[i%len(data)]))
			unsub()
			i++
		}
	})
}

func BenchmarkPublishingParallel(b *testing.B) {
	b.StopTimer()
	p := pubsub.New()
	for i := 0; i < 100; i++ {
		_, f := newSpySubscrption()
		p.Subscribe(f, pubsub.WithPath(randPath()))
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

func BenchmarkPublishingParallelStructs(b *testing.B) {
	b.StopTimer()
	p := pubsub.New()
	for i := 0; i < 100; i++ {
		_, f := newSpySubscrption()
		p.Subscribe(f, pubsub.WithPath(randPath()))
	}
	data := randStructs()
	st := StructTraverser{}
	b.StartTimer()

	b.RunParallel(func(b *testing.PB) {
		i := rand.Int()
		for b.Next() {
			p.Publish(data[i%len(data)], st)
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
				_, f := newSpySubscrption()
				unsub := p.Subscribe(f, pubsub.WithPath(data[i%len(data)]))
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

func BenchmarkPublishingWhileSubscribingStructs(b *testing.B) {
	b.StopTimer()
	p := pubsub.New()
	data := randStructs()

	var wg sync.WaitGroup
	for x := 0; x < 5; x++ {
		wg.Add(1)
		go func() {
			wg.Done()
			for i := 0; ; i++ {
				_, f := newSpySubscrption()
				unsub := p.Subscribe(f, pubsub.WithPath(randPath()))
				if i%2 == 0 {
					unsub()
				}
			}
		}()
	}

	wg.Wait()
	st := StructTraverser{}
	b.StartTimer()

	b.RunParallel(func(b *testing.PB) {
		i := rand.Int()
		for b.Next() {
			p.Publish(data[i%len(data)], st)
			i++
		}
	})
}

func randPath() []interface{} {
	var r []interface{}
	for i := 0; i < 10; i++ {
		r = append(r, fmt.Sprintf("%d", rand.Intn(10)))
	}
	return r
}

func randData() [][]interface{} {
	var r [][]interface{}
	for i := 0; i < 100000; i++ {
		r = append(r, nil)
		for j := 0; j < 10; j++ {
			r[i] = append(r[i], fmt.Sprintf("%d", rand.Intn(10)))
		}
	}
	return r
}

type someType struct {
	a string
	b string
	w *w
	x *x
}

type w struct {
	i string
	j string
}

type x struct {
	i string
	j string
}

func randNum(i int) string {
	return fmt.Sprintf("%d", rand.Intn(i))
}

func randStructs() []*someType {
	var r []*someType
	for i := 0; i < 100000; i++ {
		r = append(r, &someType{
			a: randNum(10),
			b: randNum(10),
			x: &x{
				i: randNum(10),
				j: randNum(10),
			},
		})
	}
	return r
}

type StructTraverser struct{}

func (s StructTraverser) Traverse(data interface{}) pubsub.Paths {
	// a
	return pubsub.NewPathsWithTraverser([]interface{}{"", data.(*someType).a}, pubsub.TreeTraverserFunc(s.b))
}

func (s StructTraverser) b(data interface{}) pubsub.Paths {
	return pubsub.PathAndTraversers(
		[]pubsub.PathAndTraverser{
			{
				Path:      "",
				Traverser: pubsub.TreeTraverserFunc(s.w),
			},
			{
				Path:      data.(*someType).b,
				Traverser: pubsub.TreeTraverserFunc(s.w),
			},
			{
				Path:      "",
				Traverser: pubsub.TreeTraverserFunc(s.x),
			},
			{
				Path:      data.(*someType).b,
				Traverser: pubsub.TreeTraverserFunc(s.x),
			},
		},
	)
}

func (s StructTraverser) w(data interface{}) pubsub.Paths {
	if data.(*someType).w == nil {
		return pubsub.NewPathsWithTraverser([]interface{}{""}, pubsub.TreeTraverserFunc(s.done))
	}

	return pubsub.NewPathsWithTraverser([]interface{}{"w"}, pubsub.TreeTraverserFunc(s.wi))
}

func (s StructTraverser) wi(data interface{}) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]interface{}{"", data.(*someType).w.i}, pubsub.TreeTraverserFunc(s.wj))
}

func (s StructTraverser) wj(data interface{}) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]interface{}{"", data.(*someType).w.j}, pubsub.TreeTraverserFunc(s.done))
}

func (s StructTraverser) x(data interface{}) pubsub.Paths {
	if data.(*someType).x == nil {
		return pubsub.NewPathsWithTraverser([]interface{}{""}, pubsub.TreeTraverserFunc(s.done))
	}

	return pubsub.NewPathsWithTraverser([]interface{}{"x"}, pubsub.TreeTraverserFunc(s.xi))
}

func (s StructTraverser) xi(data interface{}) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]interface{}{"", data.(*someType).x.i}, pubsub.TreeTraverserFunc(s.xj))
}

func (s StructTraverser) xj(data interface{}) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]interface{}{"", data.(*someType).x.j}, pubsub.TreeTraverserFunc(s.done))
}

func (s StructTraverser) done(data interface{}) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}
