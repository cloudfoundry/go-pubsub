package diff

import (
	"fmt"
	"reflect"
	"strings"
)

// Opt is an option type that can be passed to New.
//
// Most of the time, you'll want to use at least one
// of Actual or Expected, to differentiate the two
// in your output.
type Opt func(Differ) Differ

// WithFormat returns an Opt that wraps up differences
// using a format string.  The format should contain
// one '%s' to add the difference string in.
func WithFormat(format string) Opt {
	return func(d Differ) Differ {
		d.wrappers = append(d.wrappers, func(v string) string {
			return fmt.Sprintf(format, v)
		})
		return d
	}
}

// Sprinter is any type which can print a string.
type Sprinter interface {
	Sprint(...interface{}) string
}

// WithSprinter returns an Opt that wraps up differences
// using a Sprinter.
func WithSprinter(s Sprinter) Opt {
	return func(d Differ) Differ {
		d.wrappers = append(d.wrappers, func(v string) string {
			return s.Sprint(v)
		})
		return d
	}
}

func applyOpts(o *Differ, opts ...Opt) {
	for _, opt := range opts {
		*o = opt(*o)
	}
}

// Actual returns an Opt that only applies other Opt values
// to the actual value.
func Actual(opts ...Opt) Opt {
	return func(d Differ) Differ {
		if d.actual == nil {
			d.actual = &Differ{}
		}
		applyOpts(d.actual, opts...)
		return d
	}
}

// Expected returns an Opt that only applies other Opt values
// to the expected value.
func Expected(opts ...Opt) Opt {
	return func(d Differ) Differ {
		if d.expected == nil {
			d.expected = &Differ{}
		}
		applyOpts(d.expected, opts...)
		return d
	}
}

// Differ is a type that can diff values.  It keeps its own
// diffing style.
type Differ struct {
	wrappers []func(string) string

	actual   *Differ
	expected *Differ
}

// New creates a Differ, using the passed in opts to manipulate
// its diffing behavior and output.
//
// If opts is empty, the default options used will be:
// [ WithFormat(">%s<"), Actual(WithFormat("%s!=")) ]
//
// opts will be applied to the text in the order they
// are passed in, so you can do things like color a value
// and then wrap the colored text up in custom formatting.
//
// See the examples on the different Opt types for more
// detail.
func New(opts ...Opt) *Differ {
	d := Differ{}
	if len(opts) == 0 {
		opts = []Opt{WithFormat(">%s<"), Actual(WithFormat("%s!="))}
	}
	for _, opt := range opts {
		d = opt(d)
	}
	return &d
}

// format is used to format a string using the wrapper functions.
func (d Differ) format(v string) string {
	for _, w := range d.wrappers {
		v = w(v)
	}
	return v
}

// Diff takes two values and returns a string showing a
// diff of them.
func (d *Differ) Diff(actual, expected interface{}) string {
	return d.diff(reflect.ValueOf(actual), reflect.ValueOf(expected))
}

func (d *Differ) genDiff(format string, actual, expected interface{}) string {
	afmt := fmt.Sprintf(format, actual)
	if d.actual != nil {
		afmt = d.actual.format(afmt)
	}
	efmt := fmt.Sprintf(format, expected)
	if d.expected != nil {
		efmt = d.expected.format(efmt)
	}
	return d.format(afmt + efmt)
}

func (d *Differ) diff(av, ev reflect.Value) string {
	if !av.IsValid() {
		if !ev.IsValid() {
			return "<nil>"
		}
		if ev.Kind() == reflect.Ptr {
			return d.diff(av, ev.Elem())
		}
		return d.genDiff("%v", "<nil>", ev.Interface())
	}
	if !ev.IsValid() {
		if av.Kind() == reflect.Ptr {
			return d.diff(av.Elem(), ev)
		}
		return d.genDiff("%v", av.Interface(), "<nil>")
	}

	if av.Kind() != ev.Kind() {
		return d.genDiff("%s", av.Type(), ev.Type())
	}

	if av.CanInterface() {
		switch av.Interface().(type) {
		case []rune, []byte, string:
			return d.strDiff(av, ev)
		}
	}

	switch av.Kind() {
	case reflect.Ptr, reflect.Interface:
		return d.diff(av.Elem(), ev.Elem())
	case reflect.Slice, reflect.Array, reflect.String:
		if av.Len() != ev.Len() {
			// TODO: do a more thorough diff of values
			return d.genDiff(fmt.Sprintf("%s(len %%d)", av.Type()), av.Len(), ev.Len())
		}
		var elems []string
		for i := 0; i < av.Len(); i++ {
			elems = append(elems, d.diff(av.Index(i), ev.Index(i)))
		}
		return "[ " + strings.Join(elems, ", ") + " ]"
	case reflect.Map:
		var parts []string
		for _, kv := range ev.MapKeys() {
			k := kv.Interface()
			bmv := ev.MapIndex(kv)
			amv := av.MapIndex(kv)
			if !amv.IsValid() {
				parts = append(parts, d.genDiff("%s", fmt.Sprintf("missing key %v", k), fmt.Sprintf("%v: %v", k, bmv.Interface())))
				continue
			}
			parts = append(parts, fmt.Sprintf("%v: %s", k, d.diff(amv, bmv)))
		}
		for _, kv := range av.MapKeys() {
			// We've already compared all keys that exist in both maps; now we're
			// just looking for keys that only exist in a.
			if !ev.MapIndex(kv).IsValid() {
				k := kv.Interface()
				parts = append(parts, d.genDiff("%s", fmt.Sprintf("extra key %v: %v", k, av.MapIndex(kv).Interface()), fmt.Sprintf("%v: nil", k)))
				continue
			}
		}
		return "{" + strings.Join(parts, ", ") + "}"
	case reflect.Struct:
		if av.Type().Name() != ev.Type().Name() {
			return d.genDiff("%s", av.Type().Name(), ev.Type().Name()) + "(mismatched types)"
		}
		var parts []string
		for i := 0; i < ev.NumField(); i++ {
			f := ev.Type().Field(i)
			if f.PkgPath != "" && !f.Anonymous {
				// unexported
				continue
			}
			name := f.Name
			bfv := ev.Field(i)
			afv := av.Field(i)
			parts = append(parts, fmt.Sprintf("%s: %s", name, d.diff(afv, bfv)))
		}
		return fmt.Sprintf("%s{%s}", av.Type(), strings.Join(parts, ", "))
	default:
		if av.Type().Comparable() {
			a, b := av.Interface(), ev.Interface()
			if a != b {
				return d.genDiff("%#v", a, b)
			}
			return fmt.Sprintf("%#v", a)
		}
		return d.format(fmt.Sprintf("UNSUPPORTED: could not compare values of type %s", av.Type()))
	}
}

// strDiff helps us generate a diff between two strings.
//
// TODO: make this maybe less naive?  string diffs are hard.
func (d *Differ) strDiff(av, ev reflect.Value) string {
	strTyp := reflect.TypeOf("")
	var out string

	diffs := smallestDiff(av, ev, 0, 0)
	i := 0
	for _, diff := range diffs {
		out += av.Slice(i, diff.aStart).Convert(strTyp).Interface().(string)
		curra := av.Slice(diff.aStart, diff.aEnd).Convert(strTyp).Interface().(string)
		curre := ev.Slice(diff.eStart, diff.eEnd).Convert(strTyp).Interface().(string)
		out += d.genDiff("%s", curra, curre)
		i = diff.aEnd
	}
	out += av.Slice(i, av.Len()).Convert(strTyp).Interface().(string)
	return out
}

type diffs []difference

func (d diffs) cost() int {
	cost := 0
	for _, diff := range d {
		cost += diff.cost()
	}
	return cost
}

type difference struct {
	aStart, aEnd, eStart, eEnd int
}

func (d difference) cost() int {
	return greater(d.aEnd-d.aStart, d.eEnd-d.eStart)
}

func greater(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type diffIdx struct {
	aStart, eStart int
}

func smallestDiff(av, ev reflect.Value, aStart, eStart int) []difference {
	cache := make(map[diffIdx][]difference)
	return smallestCachingDiff(cache, av, ev, aStart, eStart)
}

func smallestCachingDiff(cache map[diffIdx][]difference, av, ev reflect.Value, aStart, eStart int) []difference {
	for aStart < av.Len() && eStart < ev.Len() && av.Index(aStart).Interface() == ev.Index(eStart).Interface() {
		aStart++
		eStart++
	}
	if d, ok := cache[diffIdx{aStart: aStart, eStart: eStart}]; ok {
		return d
	}
	if aStart == av.Len() && eStart == ev.Len() {
		return nil
	}

	shortest := diffs{{aStart: aStart, aEnd: av.Len(), eStart: eStart, eEnd: ev.Len()}}
	if aStart == av.Len() || eStart == ev.Len() {
		return shortest
	}
	for i := aStart; i < av.Len(); i++ {
		for j := eStart; j < ev.Len(); j++ {
			if av.Index(i).Interface() != ev.Index(j).Interface() {
				continue
			}
			currDiffs := append(diffs{{
				aStart: aStart,
				aEnd:   i,
				eStart: eStart,
				eEnd:   j,
			}}, smallestCachingDiff(cache, av, ev, i+1, j+1)...)
			if currDiffs.cost() < shortest.cost() {
				shortest = currDiffs
			}
		}
	}
	cache[diffIdx{aStart: aStart, eStart: eStart}] = shortest
	return shortest
}
