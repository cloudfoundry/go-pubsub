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
// Data is set by the previously returned next. This is done as a courtesy
// for the DataAssigner implementation. The value is not used by the pubsub
// library in any way and the Subscriber will always receive the original
// data value.

// Traversing a path ends when the return len(paths) == 0. If
// len(paths) > 1, then each path will be traversed.
type DataAssigner interface {
	Assign(data interface{}, currentPath []string) (paths []string, next interface{})
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

	children, next := a.Assign(next, l)

	for _, child := range children {
		s.traversePublish(d, next, a, n.FetchChild(child), append(l, child))
	}
}
