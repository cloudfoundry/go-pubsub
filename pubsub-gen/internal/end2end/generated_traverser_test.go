package end2end_test

import (
	"github.com/apoydence/pubsub"
	"github.com/apoydence/pubsub/pubsub-gen/internal/end2end"
	"hash/crc64"
)

type StructTraverser struct{}

func NewStructTraverser() StructTraverser { return StructTraverser{} }

func (s StructTraverser) Traverse(data interface{}) pubsub.Paths {
	return s._I(data)
}

func (s StructTraverser) done(data interface{}) pubsub.Paths {
	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		return 0, nil, false
	})
}

func (s StructTraverser) hashBool(data bool) uint64 {
	if data {
		return 1
	}
	return 0
}

var tableECMA = crc64.MakeTable(crc64.ECMA)

func (s StructTraverser) hashString(data string) uint64 {
	return crc64.Checksum([]byte(data), tableECMA)
}

func (s StructTraverser) _I(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s._J), true
		case 1:
			return uint64(data.(*end2end.X).I), pubsub.TreeTraverser(s._J), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) _J(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return s.__Y1_Y2_M
				}), true
		case 1:
			return crc64.Checksum([]byte(data.(*end2end.X).J), tableECMA),
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return s.__Y1_Y2_M
				}), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) __Y1_Y2_M(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
	switch idx {

	case 0:

		return 1, pubsub.TreeTraverser(s._Y1_I), true

	case 1:

		if data.(*end2end.X).Y2 == nil {
			return 0, pubsub.TreeTraverser(s.done), true
		}

		return 2, pubsub.TreeTraverser(s._Y2_I), true

	case 2:
		switch data.(*end2end.X).M.(type) {
		case end2end.M1:
			return 1, s._M_M1_A, true

		case end2end.M2:
			return 2, s._M_M2_A, true

		default:
			return 0, pubsub.TreeTraverser(s.done), true
		}

	default:
		return 0, nil, false
	}
}

func (s StructTraverser) _Y1(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(s._Y1_I), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) _Y1_I(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s._Y1_J), true
		case 1:
			return uint64(data.(*end2end.X).Y1.I), pubsub.TreeTraverser(s._Y1_J), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) _Y1_J(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s.done), true
		case 1:
			return crc64.Checksum([]byte(data.(*end2end.X).Y1.J), tableECMA), pubsub.TreeTraverser(s.done), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) _Y2(data interface{}) pubsub.Paths {

	if data.(*end2end.X).Y2 == nil {
		return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
			switch idx {
			case 0:
				return 0, pubsub.TreeTraverser(s.done), true
			default:
				return 0, nil, false
			}
		})
	}

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(s._Y2_I), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) _Y2_I(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s._Y2_J), true
		case 1:
			return uint64(data.(*end2end.X).Y2.I), pubsub.TreeTraverser(s._Y2_J), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) _Y2_J(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s.done), true
		case 1:
			return crc64.Checksum([]byte(data.(*end2end.X).Y2.J), tableECMA), pubsub.TreeTraverser(s.done), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) _M_M1(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(s._M_M1_A), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) _M_M1_A(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s.done), true
		case 1:
			return uint64(data.(*end2end.X).M.(end2end.M1).A), pubsub.TreeTraverser(s.done), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) _M_M2(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(s._M_M2_A), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) _M_M2_A(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s._M_M2_B), true
		case 1:
			return uint64(data.(*end2end.X).M.(end2end.M2).A), pubsub.TreeTraverser(s._M_M2_B), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTraverser) _M_M2_B(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s.done), true
		case 1:
			return uint64(data.(*end2end.X).M.(end2end.M2).B), pubsub.TreeTraverser(s.done), true
		default:
			return 0, nil, false
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

func (s StructTraverser) CreatePath(f *XFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

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
		path = append(path, uint64(*f.I))
	} else {
		path = append(path, 0)
	}

	if f.J != nil {
		path = append(path, crc64.Checksum([]byte(*f.J), tableECMA))
	} else {
		path = append(path, 0)
	}

	path = append(path, s.createPath_Y1(f.Y1)...)

	path = append(path, s.createPath_Y2(f.Y2)...)

	path = append(path, s.createPath_M_M1(f.M_M1)...)

	path = append(path, s.createPath_M_M2(f.M_M2)...)

	for i := len(path) - 1; i >= 1; i-- {
		if path[i] != 0 {
			break
		}
		path = path[:i]
	}

	return path
}

func (s StructTraverser) createPath_Y1(f *YFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 1)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.I != nil {
		path = append(path, uint64(*f.I))
	} else {
		path = append(path, 0)
	}

	if f.J != nil {
		path = append(path, crc64.Checksum([]byte(*f.J), tableECMA))
	} else {
		path = append(path, 0)
	}

	return path
}

func (s StructTraverser) createPath_Y2(f *YFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 2)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.I != nil {
		path = append(path, uint64(*f.I))
	} else {
		path = append(path, 0)
	}

	if f.J != nil {
		path = append(path, crc64.Checksum([]byte(*f.J), tableECMA))
	} else {
		path = append(path, 0)
	}

	return path
}

func (s StructTraverser) createPath_M_M1(f *M1Filter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 1)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.A != nil {
		path = append(path, uint64(*f.A))
	} else {
		path = append(path, 0)
	}

	return path
}

func (s StructTraverser) createPath_M_M2(f *M2Filter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 2)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.A != nil {
		path = append(path, uint64(*f.A))
	} else {
		path = append(path, 0)
	}

	if f.B != nil {
		path = append(path, uint64(*f.B))
	} else {
		path = append(path, 0)
	}

	return path
}
