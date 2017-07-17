package node

type Subscription interface {
	Write(data interface{})
}

type Node struct {
	children      map[string]*Node
	subscriptions map[Subscription]struct{}
}

func New() *Node {
	return &Node{
		children:      make(map[string]*Node),
		subscriptions: make(map[Subscription]struct{}),
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

func (n *Node) AddSubscription(s Subscription) {
	if n == nil {
		return
	}

	// TODO Check for the same subscription twice
	n.subscriptions[s] = struct{}{}
}

func (n *Node) DeleteSubscription(s Subscription) {
	if n == nil {
		return
	}

	delete(n.subscriptions, s)
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

func (n *Node) ForEachSubscription(f func(s Subscription)) {
	if n == nil {
		return
	}

	for s, _ := range n.subscriptions {
		f(s)
	}
}
