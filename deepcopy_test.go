package deepcopy

import (
	"fmt"
	. "reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleAnything() {
	tests := []interface{}{
		`"Now cut that out!"`,
		39,
		true,
		false,
		2.14,
		[]string{
			"Phil Harris",
			"Rochester van Jones",
			"Mary Livingstone",
			"Dennis Day",
		},
		[2]string{
			"Jell-O",
			"Grape-Nuts",
		},
	}

	for _, expected := range tests {
		actual := MustAnything(expected)
		fmt.Println(actual)
	}
	// Output:
	// "Now cut that out!"
	// 39
	// true
	// false
	// 2.14
	// [Phil Harris Rochester van Jones Mary Livingstone Dennis Day]
	// [Jell-O Grape-Nuts]
}

type Foo struct {
	Foo *Foo
	Bar int
}

func ExampleMap() {
	x := map[string]*Foo{
		"foo": &Foo{Bar: 1},
		"bar": &Foo{Bar: 2},
	}
	y := MustAnything(x).(map[string]*Foo)
	for _, k := range []string{"foo", "bar"} { // to ensure consistent order
		fmt.Printf("x[\"%v\"] = y[\"%v\"]: %v\n", k, k, x[k] == y[k])
		fmt.Printf("x[\"%v\"].Foo = y[\"%v\"].Foo: %v\n", k, k, x[k].Foo == y[k].Foo)
		fmt.Printf("x[\"%v\"].Bar = y[\"%v\"].Bar: %v\n", k, k, x[k].Bar == y[k].Bar)
	}
	// Output:
	// x["foo"] = y["foo"]: false
	// x["foo"].Foo = y["foo"].Foo: false
	// x["foo"].Bar = y["foo"].Bar: true
	// x["bar"] = y["bar"]: false
	// x["bar"].Foo = y["bar"].Foo: false
	// x["bar"].Bar = y["bar"].Bar: true
}

func TestInterface(t *testing.T) {
	x := []interface{}{nil}
	y := MustAnything(x).([]interface{})
	assert.Equal(t, x, y)
	assert.Equal(t, 1, len(y))

	var a interface{}
	b := MustAnything(a)
	assert.Equal(t, a, b)
}

func ExampleAvoidInfiniteLoops() {
	x := &Foo{
		Bar: 4,
	}
	x.Foo = x
	y := MustAnything(x).(*Foo)
	fmt.Printf("x == y: %v\n", x == y)
	fmt.Printf("x == x.Foo: %v\n", x == x.Foo)
	fmt.Printf("y == y.Foo: %v\n", y == y.Foo)
	// Output:
	// x == y: false
	// x == x.Foo: true
	// y == y.Foo: true
}

func TestUnsupportedKind(t *testing.T) {
	x := func() {}

	tests := []interface{}{
		x,
		map[bool]interface{}{true: x},
		[]interface{}{x},
	}

	for _, test := range tests {
		y, err := Anything(test)
		assert.Nil(t, y)
		assert.Error(t, err)
	}
}

func TestUnsupportedKindPanicsOnMust(t *testing.T) {
	x := func() {}
	assert.Panics(t, func() { MustAnything(x) })
}

func TestMismatchedTypesFail(t *testing.T) {
	tests := []struct {
		input interface{}
		kind  Kind
	}{
		{
			map[int]int{1: 2, 2: 4, 3: 8},
			Map,
		},
		{
			[]int{2, 8},
			Slice,
		},
	}
	for _, test := range tests {
		for kind, copier := range copiers {
			if kind == test.kind {
				continue
			}
			actual, err := copier(test.input, nil)
			assert.Nil(t, actual, "%v attempted value %v as %v; should be nil value, got %v", test.kind, test.input, kind, actual)
			assert.Error(t, err, "%v attempted value %v as %v; should have gotten an error", test.kind, test.input, kind)
		}
	}
}
