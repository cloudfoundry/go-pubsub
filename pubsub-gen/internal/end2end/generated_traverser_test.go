package end2end_test

import (
	"fmt"
	"github.com/apoydence/pubsub"
	"github.com/apoydence/pubsub/pubsub-gen/internal/end2end"
)

type StructTraverser struct{}

func NewStructTraverser() StructTraverser { return StructTraverser{} }

func (s StructTraverser) Traverse(data interface{}, currentPath []string) pubsub.Paths {
	return s._I(data, currentPath)
}

func (s StructTraverser) done(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}

func (s StructTraverser) _I(data interface{}, currentPath []string) pubsub.Paths {

	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*end2end.X).I)}, pubsub.TreeTraverserFunc(s._J))
}

func (s StructTraverser) _J(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.PathAndTraversers(
		[]pubsub.PathAndTraverser{
			{
				Path:      "",
				Traverser: pubsub.TreeTraverserFunc(s._Y1),
			},
			{
				Path:      fmt.Sprintf("%v", data.(*end2end.X).J),
				Traverser: pubsub.TreeTraverserFunc(s._Y1),
			},

			{
				Path:      "",
				Traverser: pubsub.TreeTraverserFunc(s._Y2),
			},
			{
				Path:      fmt.Sprintf("%v", data.(*end2end.X).J),
				Traverser: pubsub.TreeTraverserFunc(s._Y2),
			},

			{
				Path:      "",
				Traverser: pubsub.TreeTraverserFunc(s._M),
			},
			{
				Path:      fmt.Sprintf("%v", data.(*end2end.X).J),
				Traverser: pubsub.TreeTraverserFunc(s._M),
			},
		})
}

func (s StructTraverser) _Y1(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"Y1"}, pubsub.TreeTraverserFunc(s._Y1_I))
}

func (s StructTraverser) _Y1_I(data interface{}, currentPath []string) pubsub.Paths {

	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*end2end.X).Y1.I)}, pubsub.TreeTraverserFunc(s._Y1_J))
}

func (s StructTraverser) _Y1_J(data interface{}, currentPath []string) pubsub.Paths {

	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*end2end.X).Y1.J)}, pubsub.TreeTraverserFunc(s.done))
}

func (s StructTraverser) _Y2(data interface{}, currentPath []string) pubsub.Paths {

	if data.(*end2end.X).Y2 == nil {
		return pubsub.NewPathsWithTraverser([]string{""}, pubsub.TreeTraverserFunc(s.done))
	}
	return pubsub.NewPathsWithTraverser([]string{"Y2"}, pubsub.TreeTraverserFunc(s._Y2_I))
}

func (s StructTraverser) _Y2_I(data interface{}, currentPath []string) pubsub.Paths {

	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*end2end.X).Y2.I)}, pubsub.TreeTraverserFunc(s._Y2_J))
}

func (s StructTraverser) _Y2_J(data interface{}, currentPath []string) pubsub.Paths {

	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*end2end.X).Y2.J)}, pubsub.TreeTraverserFunc(s.done))
}

func (s StructTraverser) _M(data interface{}, currentPath []string) pubsub.Paths {
	switch data.(*end2end.X).M.(type) {
	case end2end.M1:
		return s._M_M1(data, currentPath)

	case end2end.M2:
		return s._M_M2(data, currentPath)

	default:
		return pubsub.NewPathsWithTraverser([]string{""}, pubsub.TreeTraverserFunc(s.done))
	}
}

func (s StructTraverser) _M_M1(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"M1"}, pubsub.TreeTraverserFunc(s._M_M1_A))
}

func (s StructTraverser) _M_M1_A(data interface{}, currentPath []string) pubsub.Paths {

	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*end2end.X).M.(end2end.M1).A)}, pubsub.TreeTraverserFunc(s.done))
}

func (s StructTraverser) _M_M2(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.NewPathsWithTraverser([]string{"M2"}, pubsub.TreeTraverserFunc(s._M_M2_A))
}

func (s StructTraverser) _M_M2_A(data interface{}, currentPath []string) pubsub.Paths {

	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*end2end.X).M.(end2end.M2).A)}, pubsub.TreeTraverserFunc(s._M_M2_B))
}

func (s StructTraverser) _M_M2_B(data interface{}, currentPath []string) pubsub.Paths {

	return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(*end2end.X).M.(end2end.M2).B)}, pubsub.TreeTraverserFunc(s.done))
}

type XFilter struct {
	I    *int
	J    *string
	Y1   *YFilter
	Y2   *YFilter
	M_M1 *M1Filter
	M_M2 *M2Filter
}

type YFilter struct {
	I *int
	J *string
}

type M1Filter struct {
	A *int
}

type M2Filter struct {
	A *int
	B *int
}

func (g StructTraverser) CreatePath(f *XFilter) []string {
	if f == nil {
		return nil
	}
	var path []string

	var count int
	if f.Y1 != nil {
		count++
	}

	if f.Y2 != nil {
		count++
	}

	if f.M_M1 != nil {
		count++
	}

	if f.M_M2 != nil {
		count++
	}

	if count > 1 {
		panic("Only one field can be set")
	}

	if f.I != nil {
		path = append(path, fmt.Sprintf("%v", *f.I))
	} else {
		path = append(path, "")
	}

	if f.J != nil {
		path = append(path, fmt.Sprintf("%v", *f.J))
	} else {
		path = append(path, "")
	}

	path = append(path, g.createPath_Y1(f.Y1)...)

	path = append(path, g.createPath_Y2(f.Y2)...)

	path = append(path, g.createPath_M_M1(f.M_M1)...)

	path = append(path, g.createPath_M_M2(f.M_M2)...)

	return path
}

func (g StructTraverser) createPath_Y1(f *YFilter) []string {
	if f == nil {
		return nil
	}
	var path []string

	path = append(path, "Y1")

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.I != nil {
		path = append(path, fmt.Sprintf("%v", *f.I))
	} else {
		path = append(path, "")
	}

	if f.J != nil {
		path = append(path, fmt.Sprintf("%v", *f.J))
	} else {
		path = append(path, "")
	}

	return path
}

func (g StructTraverser) createPath_Y2(f *YFilter) []string {
	if f == nil {
		return nil
	}
	var path []string

	path = append(path, "Y2")

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.I != nil {
		path = append(path, fmt.Sprintf("%v", *f.I))
	} else {
		path = append(path, "")
	}

	if f.J != nil {
		path = append(path, fmt.Sprintf("%v", *f.J))
	} else {
		path = append(path, "")
	}

	return path
}

func (g StructTraverser) createPath_M_M1(f *M1Filter) []string {
	if f == nil {
		return nil
	}
	var path []string

	path = append(path, "M1")

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.A != nil {
		path = append(path, fmt.Sprintf("%v", *f.A))
	} else {
		path = append(path, "")
	}

	return path
}

func (g StructTraverser) createPath_M_M2(f *M2Filter) []string {
	if f == nil {
		return nil
	}
	var path []string

	path = append(path, "M2")

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.A != nil {
		path = append(path, fmt.Sprintf("%v", *f.A))
	} else {
		path = append(path, "")
	}

	if f.B != nil {
		path = append(path, fmt.Sprintf("%v", *f.B))
	} else {
		path = append(path, "")
	}

	return path
}
