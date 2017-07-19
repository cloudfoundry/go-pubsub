package main

import (
	"log"
	"strings"

	"github.com/apoydence/pubsub"
)

//         a     x
//         |      \
//         b       y
//        / \       \
//       c   d       z
func main() {
	ps := pubsub.New()

	ps.Subscribe(Subscription("sub-1"), []string{"a", "b", "c"})
	ps.Subscribe(Subscription("sub-2"), []string{"a", "b", "d"})
	ps.Subscribe(Subscription("sub-3"), []string{"a", "b", "e"})
	ps.Subscribe(Subscription("sub-4"), []string{"a"})
	ps.Subscribe(Subscription("sub-5"), []string{"a", "b"})
	ps.Subscribe(Subscription("sub-6"), []string{"x", "y", "z"})

	dataMap1 := map[string][]string{
		"":      []string{"a"},
		"a":     []string{"b"},
		"a-b":   []string{"c", "d"},
		"a-b-c": nil,
		"a-b-d": nil,
	}
	ps.Publish("data-1", StaticTraverser(dataMap1))

	dataMap2 := map[string][]string{
		"":      []string{"x"},
		"x":     []string{"y"},
		"x-y":   []string{"z"},
		"x-y-z": nil,
	}
	ps.Publish("data-2", StaticTraverser(dataMap2))

	ps.Publish("a-b-cd", StringSplitter("-"))
	ps.Publish("ax-y-z", StringSplitter("-"))

	ps.Publish("linear-1", pubsub.LinearTreeTraverser([]string{"a", "b", "c"}))
	ps.Publish("linear-2", pubsub.LinearTreeTraverser([]string{"a", "b", "d"}))
}

// Subscription writes any results to stderr
type Subscription string

// Write implements pubsub.Subscription
func (s Subscription) Write(data interface{}) {
	log.Printf("%s <- %s", s, data)
}

// StaticTraverser publishes data based on its underlying map and not the data.
// Therefore, it does not look at the data to decide where the data belongs.
// Only the given path.
type StaticTraverser map[string][]string

func (a StaticTraverser) Traverse(data interface{}, currentPath []string) pubsub.Paths {
	path := strings.Join(currentPath, "-")
	ps, ok := a[path]
	if !ok {
		log.Panicf("Unknown path: '%s'", path)
	}

	return pubsub.FlatPaths(ps)
}

// StringSplitter splits on the given string. It then breaks each word up into
// single char strings.
type StringSplitter string

// Traverse implements pubsub.TreeTraverser. It demonstrates how complex/powerful
// Paths can be. In this case, it builds new TreeTraversers for each
// part of the split.
func (s StringSplitter) Traverse(data interface{}, currentPath []string) pubsub.Paths {
	splits := strings.Split(data.(string), string(s))

	// Remove the sepearator
	var stripped []string
	for _, split := range splits {
		if split == string(s) {
			continue
		}
		stripped = append(stripped, split)
	}

	return buildSplitTraverser(stripped)(data, currentPath)
}

func buildSplitTraverser(splits []string) pubsub.TreeTraverserFunc {
	return func(data interface{}, currentPath []string) pubsub.Paths {
		if len(splits) == 0 {
			return pubsub.FlatPaths(nil)
		}

		paths := strings.Split(splits[0], "")
		f := buildSplitTraverser(splits[1:])
		return pubsub.NewPathsWithTraverser(paths, pubsub.TreeTraverserFunc(f))
	}
}
