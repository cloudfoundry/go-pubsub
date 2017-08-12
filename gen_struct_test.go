package pubsub_test

import (
  "github.com/apoydence/pubsub"
  "fmt"
)
type testStructTrav struct{}
 func NewTestStructTrav()testStructTrav{ return testStructTrav{} }

func (s testStructTrav) Traverse(data interface{}) pubsub.Paths {
	return s._a(data)
}

	func (s testStructTrav) done(data interface{}) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}

func (s testStructTrav) _a(data interface{}) pubsub.Paths {
	
  return pubsub.PathsFunc(func(idx int) (path string, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return "", pubsub.TreeTraverserFunc(s._b), true
			case 1:
				return fmt.Sprintf("%v", data.(*testStruct).a), pubsub.TreeTraverserFunc(s._b), true
			default:
				return "", nil, false
			}
		})
}

func (s testStructTrav) _b(data interface{}) pubsub.Paths {
	
  return pubsub.PathsFunc(func(idx int) (path string, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return "", pubsub.TreeTraverserFunc(s.done), true
			case 1:
				return fmt.Sprintf("%v", data.(*testStruct).b), pubsub.TreeTraverserFunc(s.done), true
			default:
				return "", nil, false
			}
		})
}

type testStructFilter struct{
a *int
b *int

}

func (g testStructTrav) CreatePath(f *testStructFilter) []string {
if f == nil {
	return nil
}
var path []string




var count int
if count > 1 {
	panic("Only one field can be set")
}


if f.a != nil {
	path = append(path, fmt.Sprintf("%v", *f.a))
}else{
	path = append(path, "")
}

if f.b != nil {
	path = append(path, fmt.Sprintf("%v", *f.b))
}else{
	path = append(path, "")
}





return path
}
