// Package pubsub provides a library that implements the Publish and Subscribe
// model. Subscriptions can subscribe to complex data patterns and data
// will be published to all subscribers that fit the criteria.
//
// Each Subscription when subscribing will walk the underlying decision tree
// to find its place in the tree. The "SubscriptionEnroller" is used to
// analyze the "Subscription" and find the correct node to store it in.
//
// As data is published, the "DataAssigner" analyzes the data to determine
// what nodes the data belongs to. Data can belong to multiple nodes on the
// same level. This means that when data is published, the system can
// traverse multiple paths for the data.
package pubsub

import (
	"sync"

	"github.com/apoydence/pubsub/internal/node"
)

// PubSub uses the given SubscriptionEnroller to  create the subscription
// tree. It also uses the DataAssigner to then write to the subscriber. All
// of PubSub's methods safe to access concurrently. PubSub should be
// constructed with New().
type PubSub struct {
	mu sync.RWMutex
	n  *node.Node
}

// DataAssigner assigns published data to the correct subscriptions. Each
// data point can be assigned to several subscriptions. As the data traverses
// the given paths, it will write to any subscribers that are assigned there.
// Data can go down multiple paths (i.e., len(paths) > 1).
//
// Traversing a path ends when the return len(paths) == 0. If
// len(paths) > 1, then each path will be traversed.
type DataAssigner interface {
	Assign(data interface{}, currentPath []string) AssignedPaths
}

// DataAssignerFunc is an adapter to allow ordinary functions to be a
// DataAssigner.
type DataAssignerFunc func(data interface{}, currentPath []string) AssignedPaths

// Assign implements DataAssigner.
func (f DataAssignerFunc) Assign(data interface{}, currentPath []string) AssignedPaths {
	return f(data, currentPath)
}

// AssignedPaths is returned by a DataAssigner. It describes how the data is
// both assigned and how to continue to analyze it.
type AssignedPaths interface {
	// At will be called with idx ranging from [0, n] where n is the number
	// of valid paths. This means that the AssignedPaths needs to be prepared
	// for an idx that is greater than it has valid data for.
	//
	// If nextAssigner is nil, then the previous DataAssigner is used.
	At(idx int) (path string, nextAssigner DataAssigner, ok bool)
}

// Paths implements AssignedPaths for a slice of paths. It
// returns nil for all nextAssigner meaning to use the given DataAssigner.
type Paths []string

// At implements AssignedPaths.
func (p Paths) At(idx int) (string, DataAssigner, bool) {
	if idx >= len(p) {
		return "", nil, false
	}

	return p[idx], nil, true
}

// PathsWithAssigner implements AssignedPaths for both a slice of paths and
// a single DataAssigner. Each path will return the given DataAssigner.
// It shoudl be constructed with NewPathsWithAssigner().
type PathsWithAssigner struct {
	a DataAssigner
	p []string
}

func NewPathsWithAssigner(paths []string, a DataAssigner) PathsWithAssigner {
	return PathsWithAssigner{a: a, p: paths}
}

// At implements AssignedPaths.
func (a PathsWithAssigner) At(idx int) (string, DataAssigner, bool) {
	if idx >= len(a.p) {
		return "", nil, false
	}

	return a.p[idx], a.a, true
}

// Subscription is a subscription that will have cooresponding data written
// to it.
type Subscription interface {
	Write(data interface{})
}

// Unsubscriber is returned by Subscribe. It should be invoked to
// remove a subscription from the PubSub.
type Unsubscriber func()

// New constructs a new PubSub.
func New() *PubSub {
	return &PubSub{
		n: node.New(),
	}
}

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
	n.AddSubscription(sub)

	return func() {
		// TODO: Delete empty nodes
		s.mu.Lock()
		defer s.mu.Unlock()
		current := s.n
		for _, p := range path {
			current = current.FetchChild(p)
		}
		current.DeleteSubscription(sub)
	}
}

// Publish writes data using the DataAssigner to the interested subscriptions.
func (s *PubSub) Publish(d interface{}, a DataAssigner) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.traversePublish(d, d, a, s.n, nil)
}

func (s *PubSub) traversePublish(d, next interface{}, a DataAssigner, n *node.Node, l []string) {
	n.ForEachSubscription(func(ss node.Subscription) {
		ss.Write(d)
	})

	paths := a.Assign(next, l)

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
