package node_test

import (
	"math/rand"
	"testing"

	"code.cloudfoundry.org/go-pubsub/internal/node"
	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
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
			n: node.New(rand.Int63n),
		}
	})

	o.Spec("returns nil for nil node", func(t TN) {
		var nilNode *node.Node
		n := nilNode.FetchChild(99)
		Expect(t, n == nil).To(BeTrue())
	})

	o.Spec("returns nil for unknown child", func(t TN) {
		n := t.n.FetchChild(99)
		Expect(t, n == nil).To(BeTrue())
	})

	o.Spec("returns child", func(t TN) {
		n1 := t.n.AddChild(1)
		n2 := t.n.FetchChild(1)
		Expect(t, n1).To(Equal(n2))
		Expect(t, t.n.ChildLen()).To(Equal(1))

		// Removes child upon deletion
		t.n.DeleteChild(1)
		Expect(t, t.n.FetchChild(1) == nil).To(BeTrue())
	})

	o.Spec("returns all subscriptions", func(t TN) {
		id1 := t.n.AddSubscription(func(interface{}) {}, "")

		t.n.AddSubscription(func(interface{}) {}, "")
		t.n.AddSubscription(func(interface{}) {}, "")
		t.n.DeleteSubscription(id1)

		var ss []func(interface{})
		t.n.ForEachSubscription(func(id string, s []node.SubscriptionEnvelope) {
			for _, x := range s {
				ss = append(ss, x.Subscription)
			}
		})
		Expect(t, ss).To(HaveLen(2))
		Expect(t, t.n.SubscriptionLen()).To(Equal(2))
	})

	o.Spec("it handles ID collisions", func(t TN) {
		n := node.New(func(int64) int64 { return 0 })
		id1 := n.AddSubscription(func(interface{}) {}, "")
		id2 := n.AddSubscription(func(interface{}) {}, "")

		Expect(t, id1).To(Not(Equal(id2)))
	})
}

type spySubscription struct {
	subscription func(interface{})
	id           string
}
