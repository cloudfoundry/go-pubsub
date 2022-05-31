# onpar
[![GoDoc][go-doc-badge]][go-doc] [![travis][travis-badge]][travis]

Parallel testing framework for Go

### Assertions
OnPar provides built-in assertion mechanisms using `Expect` statement to make expectations. Here is some more information about `Expect` and some of the matchers that come with OnPar.

- [Expect](expect/README.md)
- [Matchers](matchers/README.md)

### Specs
Test assertions are done within a `Spec()` function. Each `Spec` has a name and a function. The function takes a `testing.T` as an argument and any output from a `BeforeEach()`. Each `Spec` is run in parallel (`t.Parallel()` is invoked for each spec before calling the given function).

```go
func TestSpecs(t *testing.T) {
    o := onpar.New()
    defer o.Run(t)

    o.BeforeEach(func(t *testing.T) (*testing.T, int, float64) {
        return t, 99, 101.0
    })

    o.AfterEach(func(t *testing.T, a int, b float64) {
            // ...
    })

    o.Spec("something informative", func(t *testing.T, a int, b float64) {
        if a != 99 {
            t.Errorf("%d != 99", a)
        }
    })
}
```

### Grouping
`Group`s are used to keep `Spec`s in logical place. The intention is to gather each `Spec` in a reasonable place. Each `Group` can have a `BeforeEach()` and a `AfterEach()` but are not required to.


```go
func TestGrouping(t *testing.T) {
    o := onpar.New()
    defer o.Run(t)

    o.BeforeEach(func(t *testing.T) (*testing.T, int, float64) {
        return t, 99, 101.0
    })

    o.Group("some-group", func() {
        o.BeforeEach(func(t *teting.T, a int, b float64) (*testing.T, string) {
            return t, "foo"
        })

        o.AfterEach(func(t *testing.T, s string) {
            // ...
        })

        o.Spec("something informative", func(t *testing.T, s string) {
            // ...
        })
    })
}
```

### Run Order
Each `BeforeEach()` runs before any `Spec` in the same `Group`. It will also run before any sub-group `Spec`s and their `BeforeEach`es. Any `AfterEach()` will run after the `Spec` and before parent `AfterEach`es.

``` go
func TestRunOrder(t *testing.T) {
    o := onpar.New()
    defer o.Run(t)

    o.BeforeEach(func(t *testing.T) (*testing.T, int, string) {
        // Spec "A": Order = 1
        // Spec "B": Order = 1
        // Spec "C": Order = 1
        return t, 99, "foo"
    })

    o.AfterEach(func(t *testing.T, i int, s string) {
        // Spec "A": Order = 4
        // Spec "B": Order = 6
        // Spec "C": Order = 6
    })

    o.Group("DA", func() {
        o.AfterEach(func(t *testing.T, i int, s string) {
            // Spec "A": Order = 3
            // Spec "B": Order = 5
            // Spec "C": Order = 5
        })

        o.Spec("A", func(t *testing.T, i int, s string) {
            // Spec "A": Order = 2
        })

        o.Group("DB", func() {
            o.BeforeEach(func(t *testing.T, i int, s string) (*testing.T, float64) {
                // Spec "B": Order = 2
                // Spec "C": Order = 2
                return t, 101
            })

            o.AfterEach(func(t *testing.T, f float64) {
                // Spec "B": Order = 4
                // Spec "C": Order = 4
            })

            o.Spec("B", func(t *testing.T, f float64) {
                // Spec "B": Order = 3
            })

            o.Spec("C", func(t *testing.T, f float64) {
                // Spec "C": Order = 3
            })
        })

        o.Group("DC", func() {
            o.BeforeEach(func(t *testing.T, i int, s string) *testing.T {
                // Will not be invoked
            })

            o.AfterEach(func(t *testing.T) {
                // Will not be invoked
            })
        })
    })
}
```

### Avoiding Closure
Why bother with returning values from a `BeforeEach`? To avoid closure of course! When running `Spec`s in parallel (which they always do), each variable needs a new instance to avoid race conditions. If you use closure, then this gets tough. So onpar will pass the arguments to the given function returned by the `BeforeEach`.

The `BeforeEach` is a gatekeeper for arguments. The returned values from `BeforeEach` are required for the following `Spec`s. Child `Group`s are also passed what their direct parent `BeforeEach` returns.

[go-doc-badge]:             https://godoc.org/github.com/poy/onpar?status.svg
[go-doc]:                   https://godoc.org/github.com/poy/onpar
[travis-badge]:             https://travis-ci.org/poy/onpar.svg?branch=master
[travis]:                   https://travis-ci.org/poy/onpar?branch=master
