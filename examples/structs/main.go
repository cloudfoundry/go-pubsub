package main

import (
	"log"

	"github.com/apoydence/pubsub"
)

type someType struct {
	a string
	b string
	x *x
}

type x struct {
	i string
	j string
}

// Here we demonstrate how powerful a DataAssigner can be. We define a
// StructTraverser that reads each field. Fields can be left blank upon
// subscription meaning the field is optinal.
func main() {
	ps := pubsub.New()

	ps.Subscribe(Subscription("sub-1"), []string{"a", "b", "x", "x.i", "x.j"})
	ps.Subscribe(Subscription("sub-2"), []string{"", "b", "x", "x.i", "x.j"})
	ps.Subscribe(Subscription("sub-3"), []string{"", "", "x", "x.i", "x.j"})
	ps.Subscribe(Subscription("sub-4"), []string{""})

	ps.Publish(&someType{a: "a", b: "b", x: &x{i: "x.i", j: "x.j"}}, StructTraverser{})
	ps.Publish(&someType{a: "a'", b: "b'", x: &x{i: "x.i", j: "x.j"}}, StructTraverser{})
	ps.Publish(&someType{a: "a", b: "b"}, StructTraverser{})
}

// Subscription writes any results to stderr
type Subscription string

// Write implements pubsub.Subscription
func (s Subscription) Write(data interface{}) {
	d := data.(*someType)
	if d.x == nil {
		log.Printf("%s <- {a:%s b:%s}", s, d.a, d.b)
		return
	}
	log.Printf("%s <- {a:%s b:%s x:{i:%s j:%s}}", s, d.a, d.b, d.x.i, d.x.j)
}

// StructTraverser traverses type SomeType.
type StructTraverser struct{}

// Assign implements pubsub.DataAssigner. It demonstrates how complex/powerful
// a AssignedPaths can be. In this case, it builds new DataAssigners for
// each part of the struct. This demonstrates how flexible a DataAssigner
// can be.
//
// In this case, each field (e.g. a or b) are optional.
func (s StructTraverser) Assign(data interface{}, currentPath []string) pubsub.AssignedPaths {
	// a
	return pubsub.NewPathsWithAssigner([]string{"", data.(*someType).a}, pubsub.DataAssignerFunc(s.b))
}

func (s StructTraverser) b(data interface{}, currentPath []string) pubsub.AssignedPaths {
	return pubsub.NewPathsWithAssigner([]string{"", data.(*someType).b}, pubsub.DataAssignerFunc(s.x))
}

func (s StructTraverser) x(data interface{}, currentPath []string) pubsub.AssignedPaths {
	if data.(*someType).x == nil {
		return pubsub.NewPathsWithAssigner([]string{""}, pubsub.DataAssignerFunc(s.done))
	}

	return pubsub.NewPathsWithAssigner([]string{"x"}, pubsub.DataAssignerFunc(s.xi))
}

func (s StructTraverser) xi(data interface{}, currentPath []string) pubsub.AssignedPaths {
	return pubsub.NewPathsWithAssigner([]string{"", data.(*someType).x.i}, pubsub.DataAssignerFunc(s.xj))
}

func (s StructTraverser) xj(data interface{}, currentPath []string) pubsub.AssignedPaths {
	return pubsub.NewPathsWithAssigner([]string{"", data.(*someType).x.j}, pubsub.DataAssignerFunc(s.done))
}

func (s StructTraverser) done(data interface{}, currentPath []string) pubsub.AssignedPaths {
	return pubsub.Paths(nil)
}
