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
	tree := Sym(name)

	tree.IfSymbol(func(namePassedIn string) {
		assert(t.Errorf, name == namePassedIn, "expected %q, got %q", name, namePassedIn)
	})
}

func TestNumberCallsCallbackWithRighInt(t *testing.T) {
	number := 13
	tree := Num(number)

	tree.IfNumber(func(numberPassedIn int) {
		assert(t.Errorf, number == numberPassedIn, "expected %d, got %d", number, numberPassedIn)
	})
}

func TestTreesAreEqualToThemselves(t *testing.T) {
	trees := map[string]Tree{
		"invalid":    Tree{},
		"symbol":     Sym("+"),
		"number":     Num(13),
		"emptyList":  Lst(),
		"flatList":   Lst(Sym("sin"), Num(13)),
		"nestedList": Lst(Lst()),
	}
	for caseName, tree := range trees {
		t.Run(caseName, func(t *testing.T) {
			assert(t.Errorf, Equal(tree, tree), "tree %v should be equal to itself", sexpr{tree})
		})
	}
}

func TestDifferentTreesAreNotEqual(t *testing.T) {
	trees := map[string]Tree{
		"invalid":    Tree{},
		"symbol":     Sym("+"),
		"number":     Num(13),
		"emptyList":  Lst(),
		"flatList":   Lst(Sym("sin"), Num(13)),
		"nestedList": Lst(Lst()),
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
	args := []Tree{Sym("+"), Num(13), Sym("x")}
	tree := Lst(args...)

	tree.IfList(func(list List) {
		assert(
			t.Errorf, len(args) == list.Len(),
			"the list should have %d elements, not %d",
			len(args), list.Len(),
		)
	})
}

func TestListElementsAreEqualToEachOther(t *testing.T) {
	args := []Tree{Sym("+"), Num(13), Sym("x")}
	tree := Lst(args...)

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
	tree := Lst()

	invalid := false
	tree.IfList(func(list List) {
		list.At(0).IfInvalid(func() {
			invalid = true
		})
	})

	assert(t.Errorf, invalid, "an out-of-range list element should be invalid")
}
