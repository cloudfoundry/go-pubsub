// Package pubsub provides a library that implements the Publish and Subscribe
// model. Subscriptions can subscribe to complex data patterns and data
// will be published to all subscribers that fit the criteria.
//
// Each Subscription when subscribing will walk the underlying subscription
// tree to find its place in the tree. The given path when subscribing is used
// to analyze the Subscription and find the correct node to store it in.
//
// As data is published, the TreeTraverser analyzes the data to determine
// what nodes the data belongs to. Data is written to multiple subscriptions.
// This means that when data is published, the system can
// traverse multiple paths for the data.
package pubsub

import (
	"math/rand"
	"sync"

	"github.com/apoydence/pubsub/internal/node"
)

// PubSub uses the given SubscriptionEnroller to  create the subscription
// tree. It also uses the TreeTraverser to then write to the subscriber. All
// of PubSub's methods safe to access concurrently. PubSub should be
// constructed with New().
type PubSub struct {
	mu rlocker
	n  *node.Node
}

// New constructs a new PubSub.
func New(opts ...PubSubOption) *PubSub {
	p := &PubSub{
		n:  node.New(),
		mu: &sync.RWMutex{},
	}

	for _, o := range opts {
		o.configure(p)
	}

	return p
}

// PubSubOption is used to configure a PubSub.
type PubSubOption interface {
	configure(*PubSub)
}

type pubsubConfigFunc func(*PubSub)

func (f pubsubConfigFunc) configure(p *PubSub) {
	f(p)
}

// WithNoMutex configures a PubSub that does not have any internal mutexes.
// This is useful if more complex or custom locking is required. For example,
// if a subscription needs to subscribe while being published to.
func WithNoMutex() PubSubOption {
	return pubsubConfigFunc(func(p *PubSub) {
		p.mu = nopLock{}
	})
}

// Subscription is a subscription that will have corresponding data written
// to it.
type Subscription func(data interface{})

// Unsubscriber is returned by Subscribe. It should be invoked to
// remove a subscription from the PubSub.
type Unsubscriber func()

// SubscribeOption is used to configure a subscription while subscribing.
type SubscribeOption interface {
	configure(*subscribeConfig)
}

// WithShardID configures a subscription to have a shardID. Subscriptions with
// a shardID are sharded to any subscriptions with the same shardID and path.
// Defaults to an empty shardID (meaning it does not shard).
func WithShardID(shardID string) SubscribeOption {
	return subscribeConfigFunc(func(c *subscribeConfig) {
		c.shardID = shardID
	})
}

// WithPath configures a subscription to reside at a path. The path determines
// what data the subscription is interested in. This value should be
// correspond to what the publishing TreeTraverser yields.
// It defaults to nil (meaning it gets everything).
func WithPath(path []interface{}) SubscribeOption {
	return subscribeConfigFunc(func(c *subscribeConfig) {
		c.path = path
	})
}

type subscribeConfig struct {
	shardID string
	path    []interface{}
}

type subscribeConfigFunc func(*subscribeConfig)

func (f subscribeConfigFunc) configure(c *subscribeConfig) {
	f(c)
}

// Subscribe will add a subscription  to the PubSub. It returns a function
// that can be used to unsubscribe.  Options can be provided to configure
// the subscription and its interactions with published data.
func (s *PubSub) Subscribe(sub Subscription, opts ...SubscribeOption) Unsubscriber {
	c := subscribeConfig{}
	for _, o := range opts {
		o.configure(&c)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	n := s.n
	for _, p := range c.path {
		n = n.AddChild(p)
	}
	id := n.AddSubscription(sub, c.shardID)

	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.cleanupSubscriptionTree(s.n, id, c.path)
	}
}

func (s *PubSub) cleanupSubscriptionTree(n *node.Node, id int64, p []interface{}) {
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
type TreeTraverser func(data interface{}) Paths

// LinearTreeTraverser implements TreeTraverser on behalf of a slice of paths.
// If the data does not traverse multiple paths, then this works well.
func LinearTreeTraverser(a []interface{}) TreeTraverser {
	return func(data interface{}) Paths {
		if len(a) == 0 {
			return FlatPaths(nil)
		}

		return PathsWithTraverser([]interface{}{a[0]}, LinearTreeTraverser(a[1:]))
	}
}

// Paths is returned by a TreeTraverser. It describes how the data is
// both assigned and how to continue to analyze it.
// At will be called with idx ranging from [0, n] where n is the number
// of valid paths. This means that the Paths needs to be prepared
// for an idx that is greater than it has valid data for.
//
// If nextTraverser is nil, then the previous TreeTraverser is used.
type Paths func(idx int, data interface{}) (path interface{}, nextTraverser TreeTraverser, ok bool)

// CombinePaths takes several paths and flattens it into a single path.
func CombinePaths(p ...Paths) Paths {
	var currentStart int
	return Paths(func(idx int, data interface{}) (path interface{}, nextTraverser TreeTraverser, ok bool) {
		for _, pp := range p {
			path, next, ok := pp(idx-currentStart, data)
			if ok {
				return path, next, ok
			}

			if len(p) == 0 {
				break
			}

			currentStart = idx
			p = p[1:]
		}
		return "", nil, false
	})
}

// FlatPaths implements Paths for a slice of paths. It
// returns nil for all nextTraverser meaning to use the given TreeTraverser.
func FlatPaths(p []interface{}) Paths {
	return func(idx int, data interface{}) (interface{}, TreeTraverser, bool) {
		if idx >= len(p) {
			return "", nil, false
		}

		return p[idx], nil, true
	}
}

// PathsWithTraverser implements Paths for both a slice of paths and
// a single TreeTraverser. Each path will return the given TreeTraverser.
func PathsWithTraverser(paths []interface{}, a TreeTraverser) Paths {
	return func(idx int, data interface{}) (interface{}, TreeTraverser, bool) {
		if idx >= len(paths) {
			return "", nil, false
		}

		return paths[idx], a, true
	}
}

// PathAndTraverser is a path and traverser pair.
type PathAndTraverser struct {
	Path      interface{}
	Traverser TreeTraverser
}

// PathsWithTraverser implement Paths and allow a TreeTraverser to have
// multiple paths with multiple traversers.
func PathAndTraversers(t []PathAndTraverser) Paths {
	return func(idx int, data interface{}) (interface{}, TreeTraverser, bool) {
		if idx >= len(t) {
			return "", nil, false
		}

		return t[idx].Path, t[idx].Traverser, true
	}
}

// Publish writes data using the TreeTraverser to the interested subscriptions.
func (s *PubSub) Publish(d interface{}, a TreeTraverser) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.traversePublish(d, d, a, s.n)
}

func (s *PubSub) traversePublish(d, next interface{}, a TreeTraverser, n *node.Node) {
	if n == nil {
		return
	}
	n.ForEachSubscription(func(shardID string, ss []node.SubscriptionEnvelope) {
		if shardID == "" {
			for _, x := range ss {
				x.Subscription(d)
			}
			return
		}

		idx := rand.Intn(len(ss))
		ss[idx].Subscription(d)
	})

	paths := a(next)

	for i := 0; ; i++ {
		child, nextA, ok := paths(i, next)
		if !ok {
			return
		}

		if nextA == nil {
			nextA = a
		}

		c := n.FetchChild(child)

		s.traversePublish(d, next, nextA, c)
	}
}

// rlocker is used to hold either a real sync.RWMutex or a nop lock.
// This is used to turn off locking.
type rlocker interface {
	sync.Locker
	RLock()
	RUnlock()
}

// nopLock is used to turn off locking for the PubSub. It implements the
// rlocker interface.
type nopLock struct{}

// Lock implements rlocker.
func (l nopLock) Lock() {}

// Unlock implements rlocker.
func (l nopLock) Unlock() {}

// RLock implements rlocker.
func (l nopLock) RLock() {}

// RUnlock implements rlocker.
func (l nopLock) RUnlock() {}
