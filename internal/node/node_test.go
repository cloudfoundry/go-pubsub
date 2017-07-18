package node_test

import (
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"github.com/apoydence/pubsub/internal/node"
)

type TN struct {
	*testing.T
	n *node.Node
}

func TestNode(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) TN {
		return TN{
			T: t,
			n: node.New(),
		}
	})

	o.Spec("returns nil for nil node", func(t TN) {
		var nilNode *node.Node
		n := nilNode.FetchChild("invalid")
		Expect(t, n == nil).To(BeTrue())
	})

	o.Spec("returns nil for unknown child", func(t TN) {
		n := t.n.FetchChild("invalid")
		Expect(t, n == nil).To(BeTrue())
	})

	o.Spec("returns child", func(t TN) {
		n1 := t.n.AddChild("a")
		n2 := t.n.FetchChild("a")
		Expect(t, n1).To(Equal(n2))
	})

	o.Spec("returns all subscriptions", func(t TN) {
		s1 := spySubscription{id: "a"}
		s2 := spySubscription{id: "b"}
		s3 := spySubscription{id: "c"}
		id1 := t.n.AddSubscription(s1)

		t.n.AddSubscription(s2)
		t.n.AddSubscription(s3)
		t.n.DeleteSubscription(id1)

		var ss []node.Subscription
		t.n.ForEachSubscription(func(s node.Subscription) {
			ss = append(ss, s)
		})
		Expect(t, ss).To(HaveLen(2))
		Expect(t, ss).To(Contain(s2))
		Expect(t, ss).To(Contain(s3))
	})
}

type spySubscription struct {
	node.Subscription
	id string
}
