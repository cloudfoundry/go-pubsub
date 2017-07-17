package pubsub

import (
	"sync"

	"github.com/apoydence/pubsub/internal/node"
)

type PubSub struct {
	mu sync.RWMutex
	e  SubscriptionEnroller
	a  DataAssigner
	n  *node.Node
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
		n: node.New(),
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
