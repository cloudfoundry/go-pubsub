package pubsub_test

import (
	"github.com/apoydence/pubsub"
)

type testStructTrav struct{}

func NewTestStructTrav() testStructTrav { return testStructTrav{} }

func (s testStructTrav) Traverse(data interface{}) pubsub.Paths {
	return s._a(data)
}

func (s testStructTrav) done(data interface{}) pubsub.Paths {
	return pubsub.Paths(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		return nil, nil, false
	})
}

func (s testStructTrav) _a(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil, pubsub.TreeTraverser(s._b), true
		case 1:
			return data.(*testStruct).a, pubsub.TreeTraverser(s._b), true
		default:
			return nil, nil, false
		}
	})
}

func (s testStructTrav) _b(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return pubsub.CombinePaths(s._aa(data), s._bb(data))
				}), true
		case 1:
			return data.(*testStruct).b,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return pubsub.CombinePaths(s._aa(data), s._bb(data))
				}), true
		default:
			return nil, nil, false
		}
	})
}

func (s testStructTrav) _aa(data interface{}) pubsub.Paths {

	if data.(*testStruct).aa == nil {
		return pubsub.Paths(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverser(s.done), true
			default:
				return nil, nil, false
			}
		})
	}

	return pubsub.Paths(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return "aa", pubsub.TreeTraverser(s._aa_a), true
		default:
			return nil, nil, false
		}
	})
}

func (s testStructTrav) _aa_a(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil, pubsub.TreeTraverser(s.done), true
		case 1:
			return data.(*testStruct).aa.a, pubsub.TreeTraverser(s.done), true
		default:
			return nil, nil, false
		}
	})
}

func (s testStructTrav) _bb(data interface{}) pubsub.Paths {

	if data.(*testStruct).bb == nil {
		return pubsub.Paths(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
			switch idx {
			case 0:
				return nil, pubsub.TreeTraverser(s.done), true
			default:
				return nil, nil, false
			}
		})
	}

	return pubsub.Paths(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return "bb", pubsub.TreeTraverser(s._bb_b), true
		default:
			return nil, nil, false
		}
	})
}

func (s testStructTrav) _bb_b(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int) (path interface{}, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return nil, pubsub.TreeTraverser(s.done), true
		case 1:
			return data.(*testStruct).bb.b, pubsub.TreeTraverser(s.done), true
		default:
			return nil, nil, false
		}
	})
}

type testStructFilter struct {
	a  *int
	b  *int
	aa *testStructAFilter
	bb *testStructBFilter
}

type testStructAFilter struct {
	a *int
}

type testStructBFilter struct {
	b *int
}

func (g testStructTrav) CreatePath(f *testStructFilter) []interface{} {
	if f == nil {
		return nil
	}
	var path []interface{}

	var count int
	if f.aa != nil {
		count++
	}

	if f.bb != nil {
		count++
	}

	if count > 1 {
		panic("Only one field can be set")
	}

	if f.a != nil {
		path = append(path, *f.a)
	} else {
		path = append(path, nil)
	}

	if f.b != nil {
		path = append(path, *f.b)
	} else {
		path = append(path, nil)
	}

	path = append(path, g.createPath_aa(f.aa)...)

	path = append(path, g.createPath_bb(f.bb)...)

	for i := len(path) - 1; i >= 1; i-- {
		if path[i] != nil {
			break
		}
		path = path[:i]
	}

	return path
}

func (g testStructTrav) createPath_aa(f *testStructAFilter) []interface{} {
	if f == nil {
		return nil
	}
	var path []interface{}

	path = append(path, "aa")

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.a != nil {
		path = append(path, *f.a)
	} else {
		path = append(path, nil)
	}

	return path
}

func (g testStructTrav) createPath_bb(f *testStructBFilter) []interface{} {
	if f == nil {
		return nil
	}
	var path []interface{}

	path = append(path, "bb")

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.b != nil {
		path = append(path, *f.b)
	} else {
		path = append(path, nil)
	}

	return path
}
