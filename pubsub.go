// Package pubsub provides a library that implements the Publish and Subscribe
// model. Subscriptions can subscribe to complex data patterns and data
// will be published to all subscribers that fit the criteria.
//
// Each Subscription when subscribing will walk the underlying decision tree
// to find its place in the tree. The "SubscriptionEnroller" is used to
// analyze the "Subscription" and find the correct node to store it in.
//
// As data is published, the "TreeTraverser" analyzes the data to determine
// what nodes the data belongs to. Data can belong to multiple nodes on the
// same level. This means that when data is published, the system can
// traverse multiple paths for the data.
package pubsub

import (
	"sync"

	"github.com/apoydence/pubsub/internal/node"
)

// PubSub uses the given SubscriptionEnroller to  create the subscription
// tree. It also uses the TreeTraverser to then write to the subscriber. All
// of PubSub's methods safe to access concurrently. PubSub should be
// constructed with New().
type PubSub struct {
	mu sync.RWMutex
	n  *node.Node
}

// New constructs a new PubSub.
func New() *PubSub {
	return &PubSub{
		n: node.New(),
	}
}

// Subscription is a subscription that will have cooresponding data written
// to it.
type Subscription interface {
	Write(data interface{})
}

// SubscriptionFuncis an adapter to allow ordinary functions to be a
// Subscription.
type SubscriptionFunc func(data interface{})

// Write implements Subscription.
func (f SubscriptionFunc) Write(data interface{}) {
	f(data)
}

// Unsubscriber is returned by Subscribe. It should be invoked to
// remove a subscription from the PubSub.
type Unsubscriber func()

// Subscribe will add a subscription using the given path to
// the PubSub. It returns a function that can be used to unsubscribe.
// Path is used to describe the placement of the subscription.
func (s *PubSub) Subscribe(sub Subscription, path []string) Unsubscriber {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := s.n
	for _, p := range path {
		n = n.AddChild(p)
	}
	id := n.AddSubscription(sub)

	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.cleanupSubscriptionTree(s.n, id, path)
	}
}

func (s *PubSub) cleanupSubscriptionTree(n *node.Node, id int64, p []string) {
	if len(p) == 0 {
		n.DeleteSubscription(id)
		return
	}

	child := n.FetchChild(p[0])
	s.cleanupSubscriptionTree(child, id, p[1:])

	if child.ChildLen() == 0 && child.SubscriptionLen() == 0 {
		n.DeleteChild(p[0])
	}
}

// TreeTraverser publishes data to the correct subscriptions. Each
// data point can be published to several subscriptions. As the data traverses
// the given paths, it will write to any subscribers that are assigned there.
// Data can go down multiple paths (i.e., len(paths) > 1).
//
// Traversing a path ends when the return len(paths) == 0. If
// len(paths) > 1, then each path will be traversed.
type TreeTraverser interface {
	// Traverse is used to traverse the subscription tree.
	Traverse(data interface{}, currentPath []string) Paths
}

// TreeTraverserFunc is an adapter to allow ordinary functions to be a
// TreeTraverser.
type TreeTraverserFunc func(data interface{}, currentPath []string) Paths

// Traverse implements TreeTraverser.
func (f TreeTraverserFunc) Traverse(data interface{}, currentPath []string) Paths {
	return f(data, currentPath)
}

// LinearTreeTraverser implements TreeTraverser on behalf of a slice of paths.
// If the data does not traverse multiple paths, then this works well.
type LinearTreeTraverser []string

// Traverse implements TreeTraverser.
func (a LinearTreeTraverser) Traverse(data interface{}, currentPath []string) Paths {
	return a.buildTreeTraverser(a)(data, currentPath)
}

func (a LinearTreeTraverser) buildTreeTraverser(remainingPath []string) TreeTraverserFunc {
	return func(data interface{}, currentPath []string) Paths {
		if len(remainingPath) == 0 {
			return FlatPaths(nil)
		}

		return NewPathsWithTraverser(FlatPaths([]string{remainingPath[0]}), a.buildTreeTraverser(remainingPath[1:]))
	}
}

// Paths is returned by a TreeTraverser. It describes how the data is
// both assigned and how to continue to analyze it.
type Paths interface {
	// At will be called with idx ranging from [0, n] where n is the number
	// of valid paths. This means that the Paths needs to be prepared
	// for an idx that is greater than it has valid data for.
	//
	// If nextTraverser is nil, then the previous TreeTraverser is used.
	At(idx int) (path string, nextTraverser TreeTraverser, ok bool)
}

// FlatPaths implements Paths for a slice of paths. It
// returns nil for all nextTraverser meaning to use the given TreeTraverser.
type FlatPaths []string

// At implements Paths.
func (p FlatPaths) At(idx int) (string, TreeTraverser, bool) {
	if idx >= len(p) {
		return "", nil, false
	}

	return p[idx], nil, true
}

// PathsWithTraverser implements Paths for both a slice of paths and
// a single TreeTraverser. Each path will return the given TreeTraverser.
// It shoudl be constructed with NewPathsWithTraverser().
type PathsWithTraverser struct {
	a TreeTraverser
	p []string
}

// NewPathsWithTraverser creates a new PathsWithTraverser.
func NewPathsWithTraverser(paths []string, a TreeTraverser) PathsWithTraverser {
	return PathsWithTraverser{a: a, p: paths}
}

// At implements Paths.
func (a PathsWithTraverser) At(idx int) (string, TreeTraverser, bool) {
	if idx >= len(a.p) {
		return "", nil, false
	}

	return a.p[idx], a.a, true
}

// Publish writes data using the TreeTraverser to the interested subscriptions.
func (s *PubSub) Publish(d interface{}, a TreeTraverser) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.traversePublish(d, d, a, s.n, nil)
}

func (s *PubSub) traversePublish(d, next interface{}, a TreeTraverser, n *node.Node, l []string) {
	n.ForEachSubscription(func(ss node.Subscription) {
		ss.Write(d)
	})

	paths := a.Traverse(next, l)

	for i := 0; ; i++ {
		child, nextA, ok := paths.At(i)
		if !ok {
			return
		}

		if nextA == nil {
			nextA = a
		}

		s.traversePublish(d, next, nextA, n.FetchChild(child), append(l, child))
	}
}
