package main

import (
	"github.com/apoydence/pubsub"
	"hash/crc64"
)

func StructTravTraverse(data interface{}) pubsub.Paths {
	return _a(data)
}

func done(data interface{}) pubsub.Paths {
	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		return 0, nil, false
	})
}

func hashBool(data bool) uint64 {
	if data {
		return 1
	}
	return 0
}

var tableECMA = crc64.MakeTable(crc64.ECMA)

func _a(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_b), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).a), tableECMA), pubsub.TreeTraverser(_b), true
		default:
			return 0, nil, false
		}
	})
}

func _b(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return __w_x
				}), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).b), tableECMA),
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return __w_x
				}), true
		default:
			return 0, nil, false
		}
	})
}

func __w_x(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
	switch idx {

	case 0:

		if data.(*someType).w == nil {
			return 0, pubsub.TreeTraverser(done), true
		}

		return 1, pubsub.TreeTraverser(_w_i), true

	case 1:

		if data.(*someType).x == nil {
			return 0, pubsub.TreeTraverser(done), true
		}

		return 2, pubsub.TreeTraverser(_x_i), true

	default:
		return 0, nil, false
	}
}

func _w(data interface{}) pubsub.Paths {

	if data.(*someType).w == nil {
		return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
			switch idx {
			case 0:
				return 0, pubsub.TreeTraverser(done), true
			default:
				return 0, nil, false
			}
		})
	}

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(_w_i), true
		default:
			return 0, nil, false
		}
	})
}

func _w_i(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_w_j), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).w.i), tableECMA), pubsub.TreeTraverser(_w_j), true
		default:
			return 0, nil, false
		}
	})
}

func _w_j(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(done), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).w.j), tableECMA), pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

func _x(data interface{}) pubsub.Paths {

	if data.(*someType).x == nil {
		return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
			switch idx {
			case 0:
				return 0, pubsub.TreeTraverser(done), true
			default:
				return 0, nil, false
			}
		})
	}

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(_x_i), true
		default:
			return 0, nil, false
		}
	})
}

func _x_i(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_x_j), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).x.i), tableECMA), pubsub.TreeTraverser(_x_j), true
		default:
			return 0, nil, false
		}
	})
}

func _x_j(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(done), true
		case 1:
			return crc64.Checksum([]byte(data.(*someType).x.j), tableECMA), pubsub.TreeTraverser(done), true
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

func StructTravCreatePath(f *someTypeFilter) []uint64 {
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

	path = append(path, createPath_w(f.w)...)

	path = append(path, createPath_x(f.x)...)

	for i := len(path) - 1; i >= 1; i-- {
		if path[i] != 0 {
			break
		}
		path = path[:i]
	}

	return path
}

func createPath_w(f *wFilter) []uint64 {
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

func createPath_x(f *xFilter) []uint64 {
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
