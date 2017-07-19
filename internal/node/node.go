package node

import "math/rand"

type Subscription interface {
	Write(data interface{})
}

type Node struct {
	children      map[string]*Node
	subscriptions map[int64]Subscription
}

func New() *Node {
	return &Node{
		children:      make(map[string]*Node),
		subscriptions: make(map[int64]Subscription),
	}
}

func (n *Node) AddChild(key string) *Node {
	if n == nil {
		return nil
	}

	if child, ok := n.children[key]; ok {
		return child
	}

	child := New()
	n.children[key] = child
	return child
}

func (n *Node) FetchChild(key string) *Node {
	if n == nil {
		return nil
	}

	if child, ok := n.children[key]; ok {
		return child
	}

	return nil
}

func (n *Node) DeleteChild(key string) {
	if n == nil {
		return
	}

	delete(n.children, key)
}

func (n *Node) ChildLen() int {
	return len(n.children)
}

func (n *Node) AddSubscription(s Subscription) int64 {
	if n == nil {
		return 0
	}

	var id int64
	for {
		id = rand.Int63()
		if _, ok := n.subscriptions[id]; !ok {
			break
		}
	}

	n.subscriptions[id] = s
	return id
}

func (n *Node) DeleteSubscription(id int64) {
	if n == nil {
		return
	}

	delete(n.subscriptions, id)
}

func (n *Node) SubscriptionLen() int {
	return len(n.subscriptions)
}

func (n *Node) ForEachSubscription(f func(s Subscription)) {
	if n == nil {
		return
	}

	for _, s := range n.subscriptions {
		f(s)
	}
}
