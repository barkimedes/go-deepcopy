package deepcopy

import (
	. "reflect"
	"testing"
)

func TestExampleAnything(t *testing.T) {
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
		if !DeepEqual(expected, actual) {
			t.Errorf("want '%v', got '%v'", expected, actual)
		}
	}
}

type Foo struct {
	Foo *Foo
	Bar int
}

func TestExampleMap(t *testing.T) {
	x := map[string]*Foo{
		"foo": {Bar: 1},
		"bar": {Bar: 2},
	}

	y := MustAnything(x).(map[string]*Foo)
	for _, k := range []string{"foo", "bar"} { // to ensure consistent order
		if x[k] == y[k] {
			t.Errorf("x[\"%v\"] == y[\"%v\"]: want '%v' got '%v'", k, k, false, x[k] == y[k])
		}
		if x[k].Foo == y[k].Foo {
			t.Errorf("x[\"%v\"].Foo == y[\"%v\"].Foo: want '%v' got '%v'", k, k, false, x[k].Foo == y[k].Foo)
		}
		if x[k].Bar != y[k].Bar {
			t.Errorf("x[\"%v\"].Bar == y[\"%v\"].Bar: want '%v' got '%v'", k, k, true, x[k].Bar == y[k].Bar)
		}
	}
}

func TestInterface(t *testing.T) {
	x := []interface{}{nil}
	y := MustAnything(x).([]interface{})
	if !DeepEqual(x, y) || len(y) != 1 {
		t.Errorf("expect %v == %v; y had length %v (expected 1)", x, y, len(y))
	}
	var a interface{}
	b := MustAnything(a)
	if a != b {
		t.Errorf("expected %v == %v", a, b)
	}
}

func TestExampleAvoidInfiniteLoops(t *testing.T) {
	x := &Foo{
		Bar: 4,
	}
	x.Foo = x
	y := MustAnything(x).(*Foo)

	if x == y {
		t.Errorf("x == y: want '%v' got '%v'", false, x == y)
	}
	if x != x.Foo {
		t.Errorf("x == x.Foo: want '%v' got '%v'", true, x == x.Foo)
	}
	if y != y.Foo {
		t.Errorf("y == y.Foo: want '%v' got '%v'", true, y == y.Foo)
	}
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
		if y != nil {
			t.Errorf("expected %v to be nil", y)
		}
		if err == nil {
			t.Errorf("expected err to not be nil")
		}
	}
}

func TestUnsupportedKindPanicsOnMust(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected a panic; didn't get one")
		}
	}()
	x := func() {}
	MustAnything(x)
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
			if actual != nil {

				t.Errorf("%v attempted value %v as %v; should be nil value, got %v", test.kind, test.input, kind, actual)
			}
			if err == nil {
				t.Errorf("%v attempted value %v as %v; should have gotten an error", test.kind, test.input, kind)
			}
		}
	}
}
