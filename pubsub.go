package pubsub

import (
	"sync"
	"unsafe"
)

type PubSub struct {
	mu sync.RWMutex
	b  TreeBuilder
	t  TreeTraverser
	n  *node
}

type Data unsafe.Pointer

type TreeBuilder interface {
	PlaceSubscription(sub Subscription, location []string) (key string, ok bool)
}

type TreeTraverser interface {
	Traverse(data Data, location []string) (keys []string)
}

func New(b TreeBuilder, t TreeTraverser) *PubSub {
	return &PubSub{
		b: b,
		t: t,
		n: newNode(),
	}
}

type Subscription interface {
	Write(data Data)
}

func (s *PubSub) Subscribe(subscription Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.traverseSubscribe(subscription, s.n, nil)
}

func (s *PubSub) Publish(d Data) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.traversePublish(d, s.n, nil)
}

type node struct {
	children      map[string]*node
	subscriptions []Subscription
}

func newNode() *node {
	return &node{
		children: make(map[string]*node),
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

	n.subscriptions = append(n.subscriptions, s)
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

func (n *node) listSubscriptions() []Subscription {
	if n == nil {
		return nil
	}

	return n.subscriptions
}

func (s *PubSub) traversePublish(d Data, n *node, l []string) {
	for _, sub := range n.listSubscriptions() {
		sub.Write(d)
	}

	children := s.t.Traverse(d, l)

	for _, child := range children {
		s.traversePublish(d, n.fetchChild(child), append(l, child))
	}
}

func (s *PubSub) traverseSubscribe(ss Subscription, n *node, l []string) {
	child, ok := s.b.PlaceSubscription(ss, l)
	if !ok {
		n.addSubscription(ss)
		return
	}

	s.traverseSubscribe(ss, n.addChild(child), append(l, child))
}
