# pubsub
[![GoDoc][go-doc-badge]][go-doc]

A tree based PubSub library for data in Go.

PubSub publishes data to subscriptions. However it can do so much more that just push some data to a subscription.  Each subscription is placed in a tree. When data is published, it traverses the tree and finds each interested subscription. This allows for sophisticated filters and routing.

Simple Example:
```go
ps := pubsub.New()
subscription := func(name string) pubsub.SubscriptionFunc {
	return func(data interface{}) {
		fmt.Printf("%s -> %v\n", name, data)
	}
}

ps.Subscribe(subscription("sub-1"), []string{"a", "b", "c"})
ps.Subscribe(subscription("sub-2"), []string{"a", "b", "d"})
ps.Subscribe(subscription("sub-3"), []string{"a", "b", "e"})

ps.Publish("data-1", pubsub.LinearDataAssigner([]string{"a", "b"}))
ps.Publish("data-2", pubsub.LinearDataAssigner([]string{"a", "b", "c"}))
ps.Publish("data-3", pubsub.LinearDataAssigner([]string{"a", "b", "d"}))
ps.Publish("data-3", pubsub.LinearDataAssigner([]string{"x", "y"}))

// Output:
// sub-1 -> data-2
// sub-2 -> data-3
```

In this example the `LinearDataAssigner` is used to traverse the tree of subscriptions. When an interested subscription is found (in this case `sub-1` and `sub-2` for `data-2` and `data-3` respectively), the subscription is handed the data.

More complex examples can be found in the [examples](https://github.com/apoydence/pubsub/tree/master/examples) directory.

### DataAssigners
A `DataAssigner` is used to traverse the subscription tree and find what subscriptions should have the data published to them. There are a few implementations provided, however it is likely a user will need to implement their own to suit their data.

When creating a `DataAssigner` it is important to note how the data is structured. A `DataAssigner` must be deterministic and ideally stateless. The order the data is parsed and returned (via `Assign()`) must align with the given path of `Subscribe()`.
This means if the `DataAssigner` intends to look at field A, then B, and then finally C, then the subscription path must be A, B and then C (and not B, A, C or something).

### Subscriptions
A `Subscription` is used when publishing data. The given path is used to determine it's placement in the subscription tree.

[go-doc-badge]:             https://godoc.org/github.com/apoydence?status.svg
[go-doc]:                   https://godoc.org/github.com/apoydence
