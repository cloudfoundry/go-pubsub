package pubsub_test

import (
	"math/rand"
	"sync"
	"testing"

	"code.cloudfoundry.org/go-pubsub"
)

func BenchmarkPublishing(b *testing.B) {
	b.StopTimer()
	p := pubsub.New()
	for i := 0; i < 100; i++ {
		_, f := newSpySubscrption()
		p.Subscribe(f, pubsub.WithPath(randPath()))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		p.Publish("data", pubsub.LinearTreeTraverser(randPath()))
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
		p.Publish(data[i%len(data)], st.traverse)
	}
}

func BenchmarkSubscriptions(b *testing.B) {
	b.StopTimer()
	p := pubsub.New()
	b.StartTimer()

	b.RunParallel(func(b *testing.PB) {
		i := rand.Int()
		for b.Next() {
			_, f := newSpySubscrption()
			unsub := p.Subscribe(f, pubsub.WithPath(randPath()))
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
	b.StartTimer()

	b.RunParallel(func(b *testing.PB) {
		i := rand.Int()
		for b.Next() {
			p.Publish("data", pubsub.LinearTreeTraverser(randPath()))
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
			p.Publish(data[i%len(data)], st.traverse)
			i++
		}
	})
}

func BenchmarkPublishingWhileSubscribing(b *testing.B) {
	b.StopTimer()
	p := pubsub.New()

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
	b.StartTimer()

	b.RunParallel(func(b *testing.PB) {
		i := rand.Int()
		for b.Next() {
			p.Publish("data", pubsub.LinearTreeTraverser(randPath()))
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
			p.Publish(data[i%len(data)], st.traverse)
			i++
		}
	})
}

func randPath() []uint64 {
	var r []uint64
	for i := 0; i < 10; i++ {
		r = append(r, uint64(rand.Int63n(10)))
	}
	return r
}

type someType struct {
	a uint64
	b uint64
	w *w
	x *x
}

type w struct {
	i uint64
	j uint64
}

type x struct {
	i uint64
	j uint64
}

func randNum(i int64) uint64 {
	return uint64(rand.Int63n(i))
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

func (s StructTraverser) traverse(data interface{}) pubsub.Paths {
	// a
	return pubsub.PathsWithTraverser([]uint64{0, data.(*someType).a}, pubsub.TreeTraverser(s.b))
}

func (s StructTraverser) b(data interface{}) pubsub.Paths {
	return pubsub.PathAndTraversers(
		[]pubsub.PathAndTraverser{
			{
				Path:      0,
				Traverser: pubsub.TreeTraverser(s.w),
			},
			{
				Path:      data.(*someType).b,
				Traverser: pubsub.TreeTraverser(s.w),
			},
			{
				Path:      0,
				Traverser: pubsub.TreeTraverser(s.x),
			},
			{
				Path:      data.(*someType).b,
				Traverser: pubsub.TreeTraverser(s.x),
			},
		},
	)
}

var (
	W = uint64(1)
	X = uint64(2)
)

func (s StructTraverser) w(data interface{}) pubsub.Paths {
	if data.(*someType).w == nil {
		return pubsub.PathsWithTraverser([]uint64{0}, pubsub.TreeTraverser(s.done))
	}

	return pubsub.PathsWithTraverser([]uint64{W}, pubsub.TreeTraverser(s.wi))
}

func (s StructTraverser) wi(data interface{}) pubsub.Paths {
	return pubsub.PathsWithTraverser([]uint64{0, data.(*someType).w.i}, pubsub.TreeTraverser(s.wj))
}

func (s StructTraverser) wj(data interface{}) pubsub.Paths {
	return pubsub.PathsWithTraverser([]uint64{0, data.(*someType).w.j}, pubsub.TreeTraverser(s.done))
}

func (s StructTraverser) x(data interface{}) pubsub.Paths {
	if data.(*someType).x == nil {
		return pubsub.PathsWithTraverser([]uint64{0}, pubsub.TreeTraverser(s.done))
	}

	return pubsub.PathsWithTraverser([]uint64{X}, pubsub.TreeTraverser(s.xi))
}

func (s StructTraverser) xi(data interface{}) pubsub.Paths {
	return pubsub.PathsWithTraverser([]uint64{0, data.(*someType).x.i}, pubsub.TreeTraverser(s.xj))
}

func (s StructTraverser) xj(data interface{}) pubsub.Paths {
	return pubsub.PathsWithTraverser([]uint64{0, data.(*someType).x.j}, pubsub.TreeTraverser(s.done))
}

func (s StructTraverser) done(data interface{}) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}
