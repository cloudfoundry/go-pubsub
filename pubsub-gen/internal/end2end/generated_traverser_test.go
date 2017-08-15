package end2end_test

import (
	"github.com/apoydence/pubsub"
	"github.com/apoydence/pubsub/pubsub-gen/internal/end2end"
)

type StructTraverser struct{}

func NewStructTraverser() StructTraverser { return StructTraverser{} }

func (s StructTraverser) Traverse(data interface{}) pubsub.Paths {
	return s._I(data)
}

func (s StructTraverser) done(data interface{}) pubsub.Paths {
	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		return nil, nil, false
	})
}

func (s StructTraverser) _I(data interface{}) pubsub.Paths {

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil, pubsub.TreeTraverserFunc(s._J), true
		case 1:
			return data.(*end2end.X).I, pubsub.TreeTraverserFunc(s._J), true
		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _J(data interface{}) pubsub.Paths {
	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {

		case 0:
			return nil,
				pubsub.TreeTraverserFunc(func(data interface{}) pubsub.Paths {
					return pubsub.CombinePaths(s._Y1(data), s._Y2(data), s._M(data))
				}), true
		case 1:
			return data.(*end2end.X).J,
				pubsub.TreeTraverserFunc(func(data interface{}) pubsub.Paths {
					return pubsub.CombinePaths(s._Y1(data), s._Y2(data), s._M(data))
				}), true

		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _Y1(data interface{}) pubsub.Paths {

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return "Y1", pubsub.TreeTraverserFunc(s._Y1_I), true
		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _Y1_I(data interface{}) pubsub.Paths {

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil, pubsub.TreeTraverserFunc(s._Y1_J), true
		case 1:
			return data.(*end2end.X).Y1.I, pubsub.TreeTraverserFunc(s._Y1_J), true
		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _Y1_J(data interface{}) pubsub.Paths {

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil, pubsub.TreeTraverserFunc(s.done), true
		case 1:
			return data.(*end2end.X).Y1.J, pubsub.TreeTraverserFunc(s.done), true
		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _Y2(data interface{}) pubsub.Paths {

	if data.(*end2end.X).Y2 == nil {
		return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverserFunc(s.done), true
			default:
				return nil, nil, false
			}
		})
	}

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return "Y2", pubsub.TreeTraverserFunc(s._Y2_I), true
		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _Y2_I(data interface{}) pubsub.Paths {

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil, pubsub.TreeTraverserFunc(s._Y2_J), true
		case 1:
			return data.(*end2end.X).Y2.I, pubsub.TreeTraverserFunc(s._Y2_J), true
		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _Y2_J(data interface{}) pubsub.Paths {

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil, pubsub.TreeTraverserFunc(s.done), true
		case 1:
			return data.(*end2end.X).Y2.J, pubsub.TreeTraverserFunc(s.done), true
		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _M(data interface{}) pubsub.Paths {
	switch data.(*end2end.X).M.(type) {
	case end2end.M1:
		return s._M_M1(data)

	case end2end.M2:
		return s._M_M2(data)

	default:
		return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverserFunc(s.done), true
			default:
				return nil, nil, false
			}
		})
	}
}

func (s StructTraverser) _M_M1(data interface{}) pubsub.Paths {

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return "M1", pubsub.TreeTraverserFunc(s._M_M1_A), true
		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _M_M1_A(data interface{}) pubsub.Paths {

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil, pubsub.TreeTraverserFunc(s.done), true
		case 1:
			return data.(*end2end.X).M.(end2end.M1).A, pubsub.TreeTraverserFunc(s.done), true
		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _M_M2(data interface{}) pubsub.Paths {

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return "M2", pubsub.TreeTraverserFunc(s._M_M2_A), true
		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _M_M2_A(data interface{}) pubsub.Paths {

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil, pubsub.TreeTraverserFunc(s._M_M2_B), true
		case 1:
			return data.(*end2end.X).M.(end2end.M2).A, pubsub.TreeTraverserFunc(s._M_M2_B), true
		default:
			return nil, nil, false
		}
	})
}

func (s StructTraverser) _M_M2_B(data interface{}) pubsub.Paths {

	return pubsub.PathsFunc(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil, pubsub.TreeTraverserFunc(s.done), true
		case 1:
			return data.(*end2end.X).M.(end2end.M2).B, pubsub.TreeTraverserFunc(s.done), true
		default:
			return nil, nil, false
		}
	})
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

func (g StructTraverser) CreatePath(f *XFilter) []interface{} {
	if f == nil {
		return nil
	}
	var path []interface{}

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
		path = append(path, *f.I)
	} else {
		path = append(path, nil)
	}

	if f.J != nil {
		path = append(path, *f.J)
	} else {
		path = append(path, nil)
	}

	path = append(path, g.createPath_Y1(f.Y1)...)

	path = append(path, g.createPath_Y2(f.Y2)...)

	path = append(path, g.createPath_M_M1(f.M_M1)...)

	path = append(path, g.createPath_M_M2(f.M_M2)...)

	return path
}

func (g StructTraverser) createPath_Y1(f *YFilter) []interface{} {
	if f == nil {
		return nil
	}
	var path []interface{}

	path = append(path, "Y1")

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.I != nil {
		path = append(path, *f.I)
	} else {
		path = append(path, nil)
	}

	if f.J != nil {
		path = append(path, *f.J)
	} else {
		path = append(path, nil)
	}

	return path
}

func (g StructTraverser) createPath_Y2(f *YFilter) []interface{} {
	if f == nil {
		return nil
	}
	var path []interface{}

	path = append(path, "Y2")

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.I != nil {
		path = append(path, *f.I)
	} else {
		path = append(path, nil)
	}

	if f.J != nil {
		path = append(path, *f.J)
	} else {
		path = append(path, nil)
	}

	return path
}

func (g StructTraverser) createPath_M_M1(f *M1Filter) []interface{} {
	if f == nil {
		return nil
	}
	var path []interface{}

	path = append(path, "M1")

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.A != nil {
		path = append(path, *f.A)
	} else {
		path = append(path, nil)
	}

	return path
}

func (g StructTraverser) createPath_M_M2(f *M2Filter) []interface{} {
	if f == nil {
		return nil
	}
	var path []interface{}

	path = append(path, "M2")

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.A != nil {
		path = append(path, *f.A)
	} else {
		path = append(path, nil)
	}

	if f.B != nil {
		path = append(path, *f.B)
	} else {
		path = append(path, nil)
	}

	return path
}
