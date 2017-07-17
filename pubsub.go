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
	e  SubscriptionEnroller
	a  DataAssigner
	n  *node.Node
}

// SubscriptionEnroller enrolls each subscription. Enroll is called until
// keepTraversing is false. When this happens, path is ignored and the
// subscription is saved at the current level of the tree. Otherwise path
// is used to assign where the subscription is stored.
//
// The passed in subscription will be the same instance. currentPath is a
// slice of paths built up from the returned path value for each level.
// (e.g., currentPath = ["A", "B"] and path = "C". Then the next currentPath
// will be ["A", "B", "C"])
type SubscriptionEnroller interface {
	// Enroll is invoked until keepTraversing is false.
	Enroll(sub Subscription, currentPath []string) (path string, keepTraversing bool)
}

// DataAssigner assigns published data to the correct subscriptions. Each
// data point can be assigned to several subscriptions. As the data traverses
// the given paths, it will write to any subscribers that are assigned there.
// Data can go down multiple paths (i.e., len(paths) > 1).
//
// Traversing a path ends when the return len(paths) == 0. If
// len(paths) > 1, then each path will be traversed.
type DataAssigner interface {
	Assign(data interface{}, currentPath []string) (paths []string)
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
func New(e SubscriptionEnroller, a DataAssigner) *PubSub {
	return &PubSub{
		e: e,
		a: a,
		n: node.New(),
	}
}

// Subscribe will add a subscription using the SubscriptionEnroller to
// the PubSub. It returns a function that can be used to unsubscribe.
func (s *PubSub) Subscribe(subscription Subscription) Unsubscriber {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.traverseSubscribe(subscription, s.n, nil)
}

// Publish writes data using the DataAssigner to the interested subscriptions.
func (s *PubSub) Publish(d interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.traversePublish(d, s.n, nil)
}

func (s *PubSub) traversePublish(d interface{}, n *node.Node, l []string) {
	n.ForEachSubscription(func(ss node.Subscription) {
		ss.Write(d)
	})

	children := s.a.Assign(d, l)

	for _, child := range children {
		s.traversePublish(d, n.FetchChild(child), append(l, child))
	}
}

func (s *PubSub) traverseSubscribe(ss Subscription, n *node.Node, l []string) Unsubscriber {
	child, ok := s.e.Enroll(ss, l)
	if !ok {
		n.AddSubscription(ss)
		return func() {
			// TODO: Delete empty nodes
			s.mu.Lock()
			defer s.mu.Unlock()
			current := s.n
			for _, ll := range l {
				current = current.FetchChild(ll)
			}
			current.DeleteSubscription(ss)
		}
	}

	return s.traverseSubscribe(ss, n.AddChild(child), append(l, child))
}
