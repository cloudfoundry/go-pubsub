package end2end_test

import (
	"code.cloudfoundry.org/go-pubsub"
	"code.cloudfoundry.org/go-pubsub/pubsub-gen/internal/end2end"
	"hash/crc64"
)

func StructTraverserTraverse(data interface{}) pubsub.Paths {
	return _I(data)
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

func _I(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_J), true
		case 1:

			return uint64(data.(*end2end.X).I) + 1, pubsub.TreeTraverser(_J), true
		default:
			return 0, nil, false
		}
	})
}

func _J(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_Repeated), true
		case 1:

			return crc64.Checksum([]byte(data.(*end2end.X).J), tableECMA) + 1, pubsub.TreeTraverser(_Repeated), true
		default:
			return 0, nil, false
		}
	})
}

func _Repeated(data interface{}) pubsub.Paths {

	if data.(*end2end.X).Repeated == nil {
		return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
			switch idx {
			case 0:
				return 0, pubsub.TreeTraverser(_RepeatedY), true
			default:
				return 0, nil, false
			}
		})
	}

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_RepeatedY), true
		case 1:

			var total uint64 = 1
			for _, x := range data.(*end2end.X).Repeated {
				total += crc64.Checksum([]byte(x), tableECMA) + 1
			}
			return total, pubsub.TreeTraverser(_RepeatedY), true
		default:
			return 0, nil, false
		}
	})
}

func _RepeatedY(data interface{}) pubsub.Paths {

	if data.(*end2end.X).RepeatedY == nil {
		return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
			switch idx {
			case 0:
				return 0, pubsub.TreeTraverser(_MapY), true
			default:
				return 0, nil, false
			}
		})
	}

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_MapY), true
		case 1:

			var total uint64 = 1
			for _, x := range data.(*end2end.X).RepeatedY {
				total += uint64(x.I) + 1
			}
			return total, pubsub.TreeTraverser(_MapY), true
		default:
			return 0, nil, false
		}
	})
}

func _MapY(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return ___Y1_Y2_E1_E2_M
				}), true
		case 1:

			var total uint64 = 1
			for x := range data.(*end2end.X).MapY {
				total += crc64.Checksum([]byte(x), tableECMA) + 1
			}
			return total,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return ___Y1_Y2_E1_E2_M
				}), true
		default:
			return 0, nil, false
		}
	})
}

func ___Y1_Y2_E1_E2_M(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
	switch idx {

	case 0:

		return 1, pubsub.TreeTraverser(_Y1_I), true

	case 1:

		if data.(*end2end.X).Y2 == nil {
			return 0, pubsub.TreeTraverser(done), true
		}

		return 2, pubsub.TreeTraverser(_Y2_I), true

	case 2:

		// Empty field name (data.(*end2end.X).E1)
		return 3, pubsub.TreeTraverser(done), true

	case 3:

		if data.(*end2end.X).E2 == nil {
			return 0, pubsub.TreeTraverser(done), true
		}

		// Empty field name (data.(*end2end.X).E2)
		return 4, pubsub.TreeTraverser(done), true

	case 4:
		switch data.(*end2end.X).M.(type) {
		case end2end.M1:
			return 5, _M_M1_A, true

		case *end2end.M2:
			return 6, _M_M2_A, true

		case *end2end.M3:
			// Interface implementation with no fields
			return 7, pubsub.TreeTraverser(done), true

		default:
			return 0, pubsub.TreeTraverser(done), true
		}

	default:
		return 0, nil, false
	}
}

func _Y1(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(_Y1_I), true
		default:
			return 0, nil, false
		}
	})
}

func _Y1_I(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_Y1_J), true
		case 1:

			return uint64(data.(*end2end.X).Y1.I) + 1, pubsub.TreeTraverser(_Y1_J), true
		default:
			return 0, nil, false
		}
	})
}

func _Y1_J(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return ___Y1_E1_E2
				}), true
		case 1:

			return crc64.Checksum([]byte(data.(*end2end.X).Y1.J), tableECMA) + 1,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return ___Y1_E1_E2
				}), true
		default:
			return 0, nil, false
		}
	})
}

func ___Y1_E1_E2(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
	switch idx {

	case 0:

		// Empty field name (data.(*end2end.X).Y1.E1)
		return 1, pubsub.TreeTraverser(done), true

	case 1:

		if data.(*end2end.X).Y1.E2 == nil {
			return 0, pubsub.TreeTraverser(done), true
		}

		// Empty field name (data.(*end2end.X).Y1.E2)
		return 2, pubsub.TreeTraverser(done), true

	default:
		return 0, nil, false
	}
}

func _Y1_E1(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

func _Y1_E2(data interface{}) pubsub.Paths {

	if data.(*end2end.X).Y1.E2 == nil {
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
			return 1, pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

func _Y2(data interface{}) pubsub.Paths {

	if data.(*end2end.X).Y2 == nil {
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
			return 1, pubsub.TreeTraverser(_Y2_I), true
		default:
			return 0, nil, false
		}
	})
}

func _Y2_I(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_Y2_J), true
		case 1:

			return uint64(data.(*end2end.X).Y2.I) + 1, pubsub.TreeTraverser(_Y2_J), true
		default:
			return 0, nil, false
		}
	})
}

func _Y2_J(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return ___Y2_E1_E2
				}), true
		case 1:

			return crc64.Checksum([]byte(data.(*end2end.X).Y2.J), tableECMA) + 1,
				pubsub.TreeTraverser(func(data interface{}) pubsub.Paths {
					return ___Y2_E1_E2
				}), true
		default:
			return 0, nil, false
		}
	})
}

func ___Y2_E1_E2(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
	switch idx {

	case 0:

		// Empty field name (data.(*end2end.X).Y2.E1)
		return 1, pubsub.TreeTraverser(done), true

	case 1:

		if data.(*end2end.X).Y2.E2 == nil {
			return 0, pubsub.TreeTraverser(done), true
		}

		// Empty field name (data.(*end2end.X).Y2.E2)
		return 2, pubsub.TreeTraverser(done), true

	default:
		return 0, nil, false
	}
}

func _Y2_E1(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

func _Y2_E2(data interface{}) pubsub.Paths {

	if data.(*end2end.X).Y2.E2 == nil {
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
			return 1, pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

func _E1(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

func _E2(data interface{}) pubsub.Paths {

	if data.(*end2end.X).E2 == nil {
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
			return 1, pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

func _M_M1(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(_M_M1_A), true
		default:
			return 0, nil, false
		}
	})
}

func _M_M1_A(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(done), true
		case 1:

			return uint64(data.(*end2end.X).M.(end2end.M1).A) + 1, pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

func _M_M2(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(_M_M2_A), true
		default:
			return 0, nil, false
		}
	})
}

func _M_M2_A(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_M_M2_B), true
		case 1:

			return uint64(data.(*end2end.X).M.(*end2end.M2).A) + 1, pubsub.TreeTraverser(_M_M2_B), true
		default:
			return 0, nil, false
		}
	})
}

func _M_M2_B(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(done), true
		case 1:

			return uint64(data.(*end2end.X).M.(*end2end.M2).B) + 1, pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

func _M_M3(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 1, pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

type XFilter struct {
	I         *int
	J         *string
	Repeated  []string
	RepeatedY []int
	MapY      []string
	Y1        *YFilter
	Y2        *YFilter
	E1        *EmptyFilter
	E2        *EmptyFilter
	M_M1      *M1Filter
	M_M2      *M2Filter
	M_M3      *M3Filter
}

type YFilter struct {
	I  *int
	J  *string
	E1 *EmptyFilter
	E2 *EmptyFilter
}

type EmptyFilter struct {
}

type M1Filter struct {
	A *int
}

type M2Filter struct {
	A *int
	B *int
}

type M3Filter struct {
}

func StructTraverserCreatePath(f *XFilter) []uint64 {
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

	if f.E1 != nil {
		count++
	}

	if f.E2 != nil {
		count++
	}

	if f.M_M1 != nil {
		count++
	}

	if f.M_M2 != nil {
		count++
	}

	if f.M_M3 != nil {
		count++
	}

	if count > 1 {
		panic("Only one field can be set")
	}

	if f.I != nil {

		path = append(path, uint64(*f.I)+1)
	} else {
		path = append(path, 0)
	}

	if f.J != nil {

		path = append(path, crc64.Checksum([]byte(*f.J), tableECMA)+1)
	} else {
		path = append(path, 0)
	}

	if f.Repeated != nil {

		var total uint64 = 1
		for _, x := range f.Repeated {
			total += crc64.Checksum([]byte(x), tableECMA) + 1
		}
		path = append(path, total)
	} else {
		path = append(path, 0)
	}

	if f.RepeatedY != nil {

		var total uint64 = 1
		for _, x := range f.RepeatedY {
			total += uint64(x) + 1
		}
		path = append(path, total)
	} else {
		path = append(path, 0)
	}

	if f.MapY != nil {

		var total uint64 = 1
		for _, x := range f.MapY {
			total += crc64.Checksum([]byte(x), tableECMA) + 1
		}
		path = append(path, total)
	} else {
		path = append(path, 0)
	}

	path = append(path, createPath__Y1(f.Y1)...)

	path = append(path, createPath__Y2(f.Y2)...)

	path = append(path, createPath__E1(f.E1)...)

	path = append(path, createPath__E2(f.E2)...)

	path = append(path, createPath__M_M1(f.M_M1)...)

	path = append(path, createPath__M_M2(f.M_M2)...)

	path = append(path, createPath__M_M3(f.M_M3)...)

	for i := len(path) - 1; i >= 1; i-- {
		if path[i] != 0 {
			break
		}
		path = path[:i]
	}

	return path
}

func createPath__Y1(f *YFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 1)

	var count int
	if f.E1 != nil {
		count++
	}

	if f.E2 != nil {
		count++
	}

	if count > 1 {
		panic("Only one field can be set")
	}

	if f.I != nil {

		path = append(path, uint64(*f.I)+1)
	} else {
		path = append(path, 0)
	}

	if f.J != nil {

		path = append(path, crc64.Checksum([]byte(*f.J), tableECMA)+1)
	} else {
		path = append(path, 0)
	}

	path = append(path, createPath__Y1_E1(f.E1)...)

	path = append(path, createPath__Y1_E2(f.E2)...)

	return path
}

func createPath__Y1_E1(f *EmptyFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 1)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	return path
}

func createPath__Y1_E2(f *EmptyFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 2)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	return path
}

func createPath__Y2(f *YFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 2)

	var count int
	if f.E1 != nil {
		count++
	}

	if f.E2 != nil {
		count++
	}

	if count > 1 {
		panic("Only one field can be set")
	}

	if f.I != nil {

		path = append(path, uint64(*f.I)+1)
	} else {
		path = append(path, 0)
	}

	if f.J != nil {

		path = append(path, crc64.Checksum([]byte(*f.J), tableECMA)+1)
	} else {
		path = append(path, 0)
	}

	path = append(path, createPath__Y2_E1(f.E1)...)

	path = append(path, createPath__Y2_E2(f.E2)...)

	return path
}

func createPath__Y2_E1(f *EmptyFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 1)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	return path
}

func createPath__Y2_E2(f *EmptyFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 2)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	return path
}

func createPath__E1(f *EmptyFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 3)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	return path
}

func createPath__E2(f *EmptyFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 4)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	return path
}

func createPath__M_M1(f *M1Filter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 5)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.A != nil {

		path = append(path, uint64(*f.A)+1)
	} else {
		path = append(path, 0)
	}

	return path
}

func createPath__M_M2(f *M2Filter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 6)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.A != nil {

		path = append(path, uint64(*f.A)+1)
	} else {
		path = append(path, 0)
	}

	if f.B != nil {

		path = append(path, uint64(*f.B)+1)
	} else {
		path = append(path, 0)
	}

	return path
}

func createPath__M_M3(f *M3Filter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	path = append(path, 7)

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	return path
}
