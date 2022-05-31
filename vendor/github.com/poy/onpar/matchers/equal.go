package matchers

import (
	"fmt"
	"reflect"
)

// EqualMatcher performs a DeepEqual between the actual and expected.
type EqualMatcher struct {
	expected interface{}
	differ   Differ
}

// Equal returns an EqualMatcher with the expected value
func Equal(expected interface{}) *EqualMatcher {
	return &EqualMatcher{
		expected: expected,
	}
}

func (m *EqualMatcher) UseDiffer(d Differ) {
	m.differ = d
}

func (m EqualMatcher) Match(actual interface{}) (interface{}, error) {
	if !reflect.DeepEqual(actual, m.expected) {
		if m.differ == nil {
			return nil, fmt.Errorf("%+v (%[1]T) to equal %+v (%[2]T)", actual, m.expected)
		}
		return nil, fmt.Errorf("expected %v to equal %v\ndiff: %s", actual, m.expected, m.differ.Diff(actual, m.expected))
	}

	return actual, nil
}
