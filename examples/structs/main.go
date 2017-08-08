package main

import (
	"fmt"
	"log"

	"github.com/apoydence/pubsub"
)

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

// Here we demonstrate how powerful a TreeTraverser can be. We define a
// StructTraverser that reads each field. Fields can be left blank upon
// subscription meaning the field is optinal.
func main() {
	ps := pubsub.New()

	ps.Subscribe(Subscription("sub-0"), pubsub.WithPath([]string{"a", "b", "w", "w.i", "w.j"}))
	ps.Subscribe(Subscription("sub-1"), pubsub.WithPath([]string{"a", "b", "x", "x.i", "x.j"}))
	ps.Subscribe(Subscription("sub-2"), pubsub.WithPath([]string{"", "b", "x", "x.i", "x.j"}))
	ps.Subscribe(Subscription("sub-3"), pubsub.WithPath([]string{"", "", "x", "x.i", "x.j"}))
	ps.Subscribe(Subscription("sub-4"), pubsub.WithPath([]string{""}))

	ps.Publish(&someType{a: "a", b: "b", w: &w{i: "w.i", j: "w.j"}, x: &x{i: "x.i", j: "x.j"}}, StructTraverser{})
	ps.Publish(&someType{a: "a", b: "b", x: &x{i: "x.i", j: "x.j"}}, StructTraverser{})
	ps.Publish(&someType{a: "a'", b: "b'", x: &x{i: "x.i", j: "x.j"}}, StructTraverser{})
	ps.Publish(&someType{a: "a", b: "b"}, StructTraverser{})
}

// Subscription writes any results to stderr
type Subscription string

// Write implements pubsub.Subscription
func (s Subscription) Write(data interface{}) {
	d := data.(*someType)
	var w string
	if d.w != nil {
		w = fmt.Sprintf("w:{i:%s j:%s}", d.w.i, d.w.j)
	}

	var x string
	if d.x != nil {
		x = fmt.Sprintf("x:{i:%s j:%s}", d.x.i, d.x.j)
	}
	log.Printf("%s <- {a:%s b:%s %s %s", s, d.a, d.b, w, x)
}

// StructTraverser traverses type SomeType.
type StructTraverser struct{}

// Traverse implements pubsub.TreeTraverser. It demonstrates how complex/powerful
// Paths can be. In this case, it builds new TreeTraversers for
// each part of the struct. This demonstrates how flexible a TreeTraverser
// can be.
//
// In this case, each field (e.g. a or b) are optional.
func (s StructTraverser) Traverse(data interface{}, currentPath []string) pubsub.Paths {
	// a
	return pubsub.NewPathsWithTraverser([]string{"", data.(*someType).a}, pubsub.TreeTraverserFunc(s.b))
}

func (s StructTraverser) b(data interface{}, currentPath []string) pubsub.Paths {
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

func (s StructTraverser) w(data interface{}, currentPath []string) pubsub.Paths {
	if data.(*someType).w == nil {
		return pubsub.NewPathsWithTraverser([]string{""}, pubsub.TreeTraverserFunc(s.done))
	}

	return pubsub.NewPathsWithTraverser([]string{"w"}, pubsub.TreeTraverserFunc(s.wi))
}

func (s StructTraverser) wi(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"", data.(*someType).w.i}, pubsub.TreeTraverserFunc(s.wj))
}

func (s StructTraverser) wj(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"", data.(*someType).w.j}, pubsub.TreeTraverserFunc(s.done))
}

func (s StructTraverser) x(data interface{}, currentPath []string) pubsub.Paths {
	if data.(*someType).x == nil {
		return pubsub.NewPathsWithTraverser([]string{""}, pubsub.TreeTraverserFunc(s.done))
	}

	return pubsub.NewPathsWithTraverser([]string{"x"}, pubsub.TreeTraverserFunc(s.xi))
}

func (s StructTraverser) xi(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"", data.(*someType).x.i}, pubsub.TreeTraverserFunc(s.xj))
}

func (s StructTraverser) xj(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"", data.(*someType).x.j}, pubsub.TreeTraverserFunc(s.done))
}

func (s StructTraverser) done(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}
