package main

import (
	"github.com/apoydence/pubsub"
	"hash/crc64"
)

type StructTrav struct{}

func NewStructTrav() StructTrav { return StructTrav{} }

func (s StructTrav) Traverse(data interface{}) pubsub.Paths {
	return s._a(data)
}

func (s StructTrav) done(data interface{}) pubsub.Paths {
	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		return 0, nil, false
	})
}

func (s StructTrav) hashBool(data bool) uint64 {
	if data {
		return 1
	}
	return 0
}

var tableECMA = crc64.MakeTable(crc64.ECMA)

func (s StructTrav) hashString(data string) uint64 {
	return crc64.Checksum([]byte(data), tableECMA)
}

func (s StructTrav) _a(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s._b), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).a), tableECMA), pubsub.TreeTraverser(s._b), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTrav) _b(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return s.__w_x
				}), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).b), tableECMA),
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return s.__w_x
				}), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTrav) __w_x(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
	switch idx {

	case 0:

		if data.(*someType).w == nil {
			return 0, pubsub.TreeTraverser(s.done), true
		}

		return 1, pubsub.TreeTraverser(s._w_i), true

	case 1:

		if data.(*someType).x == nil {
			return 0, pubsub.TreeTraverser(s.done), true
		}

		return 2, pubsub.TreeTraverser(s._x_i), true

	default:
		return 0, nil, false
	}
}

func (s StructTrav) _w(data interface{}) pubsub.Paths {

	if data.(*someType).w == nil {
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
			return 1, pubsub.TreeTraverser(s._w_i), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTrav) _w_i(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s._w_j), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).w.i), tableECMA), pubsub.TreeTraverser(s._w_j), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTrav) _w_j(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s.done), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).w.j), tableECMA), pubsub.TreeTraverser(s.done), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTrav) _x(data interface{}) pubsub.Paths {

	if data.(*someType).x == nil {
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
			return 1, pubsub.TreeTraverser(s._x_i), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTrav) _x_i(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s._x_j), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).x.i), tableECMA), pubsub.TreeTraverser(s._x_j), true
		default:
			return 0, nil, false
		}
	})
}

func (s StructTrav) _x_j(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(s.done), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).x.j), tableECMA), pubsub.TreeTraverser(s.done), true
		default:
			return 0, nil, false
		}
	})
}

type someTypeFilter struct {
	a *string
	b *string
	w *wFilter
	x *xFilter
}

type wFilter struct {
	i *string
	j *string
}

type xFilter struct {
	i *string
	j *string
}

func (s StructTrav) CreatePath(f *someTypeFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	var count int
	if f.w != nil {
		count++
	}

	if f.x != nil {
		count++
	}

	if count > 1 {
		panic("Only one field can be set")
	}

	if f.a != nil {
		path = append(path, crc64.Checksum([]byte(*f.a), tableECMA))
	} else {
		path = append(path, 0)
	}

	if f.b != nil {
		path = append(path, crc64.Checksum([]byte(*f.b), tableECMA))
	} else {
		path = append(path, 0)
	}

	path = append(path, s.createPath_w(f.w)...)

	path = append(path, s.createPath_x(f.x)...)

	for i := len(path) - 1; i >= 1; i-- {
		if path[i] != 0 {
			break
		}
		path = path[:i]
	}

	return path
}

func (s StructTrav) createPath_w(f *wFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 1)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.i != nil {
		path = append(path, crc64.Checksum([]byte(*f.i), tableECMA))
	} else {
		path = append(path, 0)
	}

	if f.j != nil {
		path = append(path, crc64.Checksum([]byte(*f.j), tableECMA))
	} else {
		path = append(path, 0)
	}

	return path
}

func (s StructTrav) createPath_x(f *xFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 2)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.i != nil {
		path = append(path, crc64.Checksum([]byte(*f.i), tableECMA))
	} else {
		path = append(path, 0)
	}

	if f.j != nil {
		path = append(path, crc64.Checksum([]byte(*f.j), tableECMA))
	} else {
		path = append(path, 0)
	}

	return path
}
