package pubsub

import (
	"sync"
)

type PubSub struct {
	mu sync.RWMutex
	e  SubscriptionEnroller
	a  DataAssigner
	n  *node
}

type SubscriptionEnroller interface {
	Enroll(sub Subscription, location []string) (key string, ok bool)
}

type DataAssigner interface {
	Assign(data interface{}, location []string) (keys []string)
}

func New(e SubscriptionEnroller, a DataAssigner) *PubSub {
	return &PubSub{
		e: e,
		a: a,
		n: newNode(),
	}
}

type Subscription interface {
	Write(data interface{})
}

type Unsubscriber func()

func (s *PubSub) Subscribe(subscription Subscription) Unsubscriber {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.traverseSubscribe(subscription, s.n, nil)
}

func (s *PubSub) Publish(d interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.traversePublish(d, s.n, nil)
}

type node struct {
	children      map[string]*node
	subscriptions map[Subscription]struct{}
}

func newNode() *node {
	return &node{
		children:      make(map[string]*node),
		subscriptions: make(map[Subscription]struct{}),
	}
}

func (n *node) addChild(key string) *node {
	if n == nil {
		return nil
	}

	if child, ok := n.children[key]; ok {
		return child
	}

	child := newNode()
	n.children[key] = child
	return child
}

func (n *node) addSubscription(s Subscription) {
	if n == nil {
		return
	}

	// TODO Check for the same subscription twice
	n.subscriptions[s] = struct{}{}
}

func (n *node) fetchChild(key string) *node {
	if n == nil {
		return nil
	}

	if child, ok := n.children[key]; ok {
		return child
	}

	return nil
}

func (n *node) forEachSubscription(f func(s Subscription)) {
	if n == nil {
		return
	}

	for s, _ := range n.subscriptions {
		f(s)
	}
}

func (s *PubSub) traversePublish(d interface{}, n *node, l []string) {
	n.forEachSubscription(func(ss Subscription) {
		ss.Write(d)
	})

	children := s.a.Assign(d, l)

	for _, child := range children {
		s.traversePublish(d, n.fetchChild(child), append(l, child))
	}
}

func (s *PubSub) traverseSubscribe(ss Subscription, n *node, l []string) Unsubscriber {
	child, ok := s.e.Enroll(ss, l)
	if !ok {
		n.addSubscription(ss)
		return func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			current := s.n
			for _, ll := range l {
				current = current.fetchChild(ll)
			}
			delete(current.subscriptions, ss)
		}
	}

	return s.traverseSubscribe(ss, n.addChild(child), append(l, child))
}
