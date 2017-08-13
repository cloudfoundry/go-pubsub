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
		switch idx{
		
case 0:
				return "", 
     pubsub.TreeTraverserFunc(func(data interface{}) pubsub.Paths {
 				return pubsub.CombinePaths(s._aa(data),s._bb(data))
 			}), true
case 1:
				return fmt.Sprintf("%v", data.(*testStruct).b), 
     pubsub.TreeTraverserFunc(func(data interface{}) pubsub.Paths {
 				return pubsub.CombinePaths(s._aa(data),s._bb(data))
 			}), true

	  default:
			return "", nil, false
		}
	})
}

func(s testStructTrav) _aa(data interface{}) pubsub.Paths {
	
  if data.(*testStruct).aa == nil {
		return pubsub.PathsFunc(func(idx int) (path string, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return "", pubsub.TreeTraverserFunc(s.done), true
			default:
				return "", nil, false
			}
		})
  }
		
  return pubsub.PathsFunc(func(idx int) (path string, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return "aa", pubsub.TreeTraverserFunc(s._aa_a), true
			default:
				return "", nil, false
			}
		})
}

func (s testStructTrav) _aa_a(data interface{}) pubsub.Paths {
	
  return pubsub.PathsFunc(func(idx int) (path string, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return "", pubsub.TreeTraverserFunc(s.done), true
			case 1:
				return fmt.Sprintf("%v", data.(*testStruct).aa.a), pubsub.TreeTraverserFunc(s.done), true
			default:
				return "", nil, false
			}
		})
}

func(s testStructTrav) _bb(data interface{}) pubsub.Paths {
	
  if data.(*testStruct).bb == nil {
		return pubsub.PathsFunc(func(idx int) (path string, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return "", pubsub.TreeTraverserFunc(s.done), true
			default:
				return "", nil, false
			}
		})
  }
		
  return pubsub.PathsFunc(func(idx int) (path string, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return "bb", pubsub.TreeTraverserFunc(s._bb_b), true
			default:
				return "", nil, false
			}
		})
}

func (s testStructTrav) _bb_b(data interface{}) pubsub.Paths {
	
  return pubsub.PathsFunc(func(idx int) (path string, nextTraverser pubsub.TreeTraverser, ok bool){
			switch idx {
			case 0:
				return "", pubsub.TreeTraverserFunc(s.done), true
			case 1:
				return fmt.Sprintf("%v", data.(*testStruct).bb.b), pubsub.TreeTraverserFunc(s.done), true
			default:
				return "", nil, false
			}
		})
}

type testStructFilter struct{
a *int
b *int
aa *testStructAFilter
bb *testStructBFilter

}

type testStructAFilter struct{
a *int

}

type testStructBFilter struct{
b *int

}

func (g testStructTrav) CreatePath(f *testStructFilter) []string {
if f == nil {
	return nil
}
var path []string




var count int
if f.aa != nil{
	count++
}

if f.bb != nil{
	count++
}

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




path = append(path, g.createPath_aa(f.aa)...)

path = append(path, g.createPath_bb(f.bb)...)


return path
}

func (g testStructTrav) createPath_aa(f *testStructAFilter) []string {
if f == nil {
	return nil
}
var path []string

path = append(path, "aa")


var count int
if count > 1 {
	panic("Only one field can be set")
}


if f.a != nil {
	path = append(path, fmt.Sprintf("%v", *f.a))
}else{
	path = append(path, "")
}





return path
}

func (g testStructTrav) createPath_bb(f *testStructBFilter) []string {
if f == nil {
	return nil
}
var path []string

path = append(path, "bb")


var count int
if count > 1 {
	panic("Only one field can be set")
}


if f.b != nil {
	path = append(path, fmt.Sprintf("%v", *f.b))
}else{
	path = append(path, "")
}





return path
}
