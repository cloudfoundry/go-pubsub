package end2end_test

import (
	"fmt"
	"github.com/apoydence/pubsub"
)

type StructTraverser struct{}

func NewStructTraverser() StructTraverser { return StructTraverser{} }

func (s StructTraverser) Traverse(data interface{}, currentPath []string) pubsub.Paths {
	return s._i(data, currentPath)
}
func (s StructTraverser) done(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}

func (s StructTraverser) _i(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*X).i)}, pubsub.TreeTraverserFunc(s._j))
}

func (s StructTraverser) _j(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.PathAndTraversers(
		[]pubsub.PathAndTraverser{
			{
				Path:      "",
				Traverser: pubsub.TreeTraverserFunc(s._y1),
			},
			{
				Path:      fmt.Sprintf("%v", data.(*X).j),
				Traverser: pubsub.TreeTraverserFunc(s._y1),
			},

			{
				Path:      "",
				Traverser: pubsub.TreeTraverserFunc(s._y2),
			},
			{
				Path:      fmt.Sprintf("%v", data.(*X).j),
				Traverser: pubsub.TreeTraverserFunc(s._y2),
			},
		})
}

func (s StructTraverser) _y1(data interface{}, currentPaht []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"y1"}, pubsub.TreeTraverserFunc(s._y1_i))
}

func (s StructTraverser) _y1_i(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*X).y1.i)}, pubsub.TreeTraverserFunc(s._y1_j))
}

func (s StructTraverser) _y1_j(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*X).y1.j)}, pubsub.TreeTraverserFunc(s.done))
}

func (s StructTraverser) _y2(data interface{}, currentPaht []string) pubsub.Paths {

	if data.(*X).y2 == nil {
		return pubsub.NewPathsWithTraverser([]string{""}, pubsub.TreeTraverserFunc(s.done))
	}
	return pubsub.NewPathsWithTraverser([]string{"y2"}, pubsub.TreeTraverserFunc(s._y2_i))
}

func (s StructTraverser) _y2_i(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*X).y2.i)}, pubsub.TreeTraverserFunc(s._y2_j))
}

func (s StructTraverser) _y2_j(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*X).y2.j)}, pubsub.TreeTraverserFunc(s.done))
}
