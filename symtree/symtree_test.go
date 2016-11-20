//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

import (
	"fmt"
	"testing"
)

func TestSymbolCallsCallbackWithRightString(t *testing.T) {
	name := "a-symbol"
	tree := NewSymbol(name)

	tree.IfSymbol(func(namePassedIn string) {
		assert(t.Errorf, name == namePassedIn, "expected %q, got %q", name, namePassedIn)
	})
}

func TestNumberCallsCallbackWithRighInt(t *testing.T) {
	number := 13
	tree := NewNumber(number)

	tree.IfNumber(func(numberPassedIn int) {
		assert(t.Errorf, number == numberPassedIn, "expected %d, got %d", number, numberPassedIn)
	})
}

func TestSymTreesAreEqualToThemselves(t *testing.T) {
	trees := map[string]SymTree{
		"invalid":    SymTree{},
		"symbol":     NewSymbol("+"),
		"number":     NewNumber(13),
		"emptyList":  NewList(),
		"flatList":   NewList(NewSymbol("sin"), NewNumber(13)),
		"nestedList": NewList(NewList()),
	}
	for caseName, tree := range trees {
		t.Run(caseName, func(t *testing.T) {
			assert(t.Errorf, Equal(tree, tree), "tree %v should be equal to itself", sexpr{tree})
		})
	}
}

func TestDifferentSymTreesAreNotEqual(t *testing.T) {
	trees := map[string]SymTree{
		"invalid":    SymTree{},
		"symbol":     NewSymbol("+"),
		"number":     NewNumber(13),
		"emptyList":  NewList(),
		"flatList":   NewList(NewSymbol("sin"), NewNumber(13)),
		"nestedList": NewList(NewList()),
	}
	for leftName, left := range trees {
		for rightName, right := range trees {

			caseName := fmt.Sprintf("%s-%s", leftName, rightName)
			if leftName == rightName {
				continue
			}

			t.Run(caseName, func(t *testing.T) {
				assert(
					t.Errorf, !Equal(left, right),
					"%v should not equal %v", sexpr{left}, sexpr{right},
				)
			})
		}
	}
}

func TestListLengthIsTheNumberOfElementsItWasCreatedWith(t *testing.T) {
	args := []SymTree{NewSymbol("+"), NewNumber(13), NewSymbol("x")}
	tree := NewList(args...)

	tree.IfList(func(list List) {
		assert(
			t.Errorf, len(args) == list.Len(),
			"the list should have %d elements, not %d",
			len(args), list.Len(),
		)
	})
}

func TestListElementsAreEqualToEachOther(t *testing.T) {
	args := []SymTree{NewSymbol("+"), NewNumber(13), NewSymbol("x")}
	tree := NewList(args...)

	tree.IfList(func(list List) {
		for i, arg := range args {
			assert(
				t.Errorf, Equal(arg, list.At(i)),
				"%d-th element: expected %v, got %v",
				i, sexpr{arg}, sexpr{list.At(i)},
			)
		}
	})
}

func TestListElementOutOfRangeIsInvalid(t *testing.T) {
	tree := NewList()

	invalid := false
	tree.IfList(func(list List) {
		list.At(0).IfInvalid(func() {
			invalid = true
		})
	})

	assert(t.Errorf, invalid, "an out-of-range list element should be invalid")
}
