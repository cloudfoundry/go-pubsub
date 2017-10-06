package main

import (
	"fmt"
	"log"

	"code.cloudfoundry.org/go-pubsub"
	"code.cloudfoundry.org/go-pubsub/pubsub-gen/setters"
)

//go:generate $GOPATH/bin/pubsub-gen --output=$GOPATH/src/code.cloudfoundry.org/go-pubsub/examples/structs/gen_struct.go --pointer --struct-name=code.cloudfoundry.org/go-pubsub/examples/structs.someType --traverser=StructTrav --package=main

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

	ps.Subscribe(Subscription("sub-0"), pubsub.WithPath(StructTravCreatePath(&someTypeFilter{
		a: setters.String("a"),
		b: setters.String("b"),
		w: &wFilter{
			i: setters.String("w.i"),
			j: setters.String("w.j"),
		},
	})))

	ps.Subscribe(Subscription("sub-1"), pubsub.WithPath(StructTravCreatePath(&someTypeFilter{
		a: setters.String("a"),
		b: setters.String("b"),
		x: &xFilter{
			i: setters.String("x.i"),
			j: setters.String("x.j"),
		},
	})))

	ps.Subscribe(Subscription("sub-2"), pubsub.WithPath(StructTravCreatePath(&someTypeFilter{
		b: setters.String("b"),
		x: &xFilter{
			i: setters.String("x.i"),
			j: setters.String("x.j"),
		},
	})))

	ps.Subscribe(Subscription("sub-3"), pubsub.WithPath(StructTravCreatePath(&someTypeFilter{
		x: &xFilter{
			i: setters.String("x.i"),
			j: setters.String("x.j"),
		},
	})))

	ps.Subscribe(Subscription("sub-4"))

	ps.Publish(&someType{a: "a", b: "b", w: &w{i: "w.i", j: "w.j"}, x: &x{i: "x.i", j: "x.j"}}, StructTravTraverse)
	ps.Publish(&someType{a: "a", b: "b", x: &x{i: "x.i", j: "x.j"}}, StructTravTraverse)
	ps.Publish(&someType{a: "a'", b: "b'", x: &x{i: "x.i", j: "x.j"}}, StructTravTraverse)
	ps.Publish(&someType{a: "a", b: "b"}, StructTravTraverse)
}

// Subscription writes any results to stderr
func Subscription(s string) func(interface{}) {
	return func(data interface{}) {
		d := data.(*someType)
		var w string
		if d.w != nil {
			w = fmt.Sprintf("w:{i:%s j:%s}", d.w.i, d.w.j)
		}

		var x string
		if d.x != nil {
			x = fmt.Sprintf("x:{i:%s j:%s}", d.x.i, d.x.j)
		}
		log.Printf("%s <- {a:%s b:%s %s %s}", s, d.a, d.b, w, x)
	}
}
