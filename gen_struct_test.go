package pubsub_test

import (
	"code.cloudfoundry.org/go-pubsub"
	"hash/crc64"
)

func testStructTravTraverse(data interface{}) pubsub.Paths {
	return _a(data)
}

func done(data interface{}) pubsub.Paths {
	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		return 0, nil, false
	})
}

func hashBool(data bool) uint64 {
	if data {
		return 2
	}
	return 1
}

var tableECMA = crc64.MakeTable(crc64.ECMA)

func _a(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_b), true
		case 1:

			return uint64(data.(*testStruct).a) + 1, pubsub.TreeTraverser(_b), true
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
					return ___aa_bb
				}), true
		case 1:

			return uint64(data.(*testStruct).b) + 1,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return ___aa_bb
				}), true
		default:
			return 0, nil, false
		}
	})
}

func ___aa_bb(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
	switch idx {

	case 0:

		if data.(*testStruct).aa == nil {
			return 0, pubsub.TreeTraverser(done), true
		}

		return 1, pubsub.TreeTraverser(_aa_a), true

	case 1:

		if data.(*testStruct).bb == nil {
			return 0, pubsub.TreeTraverser(done), true
		}

		return 2, pubsub.TreeTraverser(_bb_b), true

	default:
		return 0, nil, false
	}
}

func _aa(data interface{}) pubsub.Paths {

	if data.(*testStruct).aa == nil {
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
			return 1, pubsub.TreeTraverser(_aa_a), true
		default:
			return 0, nil, false
		}
	})
}

func _aa_a(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(done), true
		case 1:

			return uint64(data.(*testStruct).aa.a) + 1, pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

func _bb(data interface{}) pubsub.Paths {

	if data.(*testStruct).bb == nil {
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
			return 1, pubsub.TreeTraverser(_bb_b), true
		default:
			return 0, nil, false
		}
	})
}

func _bb_b(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(done), true
		case 1:

			return uint64(data.(*testStruct).bb.b) + 1, pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
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

func testStructTravCreatePath(f *testStructFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

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

		path = append(path, uint64(*f.a)+1)
	} else {
		path = append(path, 0)
	}

	if f.b != nil {

		path = append(path, uint64(*f.b)+1)
	} else {
		path = append(path, 0)
	}

	path = append(path, createPath__aa(f.aa)...)

	path = append(path, createPath__bb(f.bb)...)

	for i := len(path) - 1; i >= 1; i-- {
		if path[i] != 0 {
			break
		}
		path = path[:i]
	}

	return path
}

func createPath__aa(f *testStructAFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 1)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.a != nil {

		path = append(path, uint64(*f.a)+1)
	} else {
		path = append(path, 0)
	}

	return path
}

func createPath__bb(f *testStructBFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 2)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.b != nil {

		path = append(path, uint64(*f.b)+1)
	} else {
		path = append(path, 0)
	}

	return path
}
