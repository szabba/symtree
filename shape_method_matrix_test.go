//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

import (
	"fmt"
	"testing"
)

// Each SymTree is one of four possibl "shapes": invalid, a symbol, a number or a list.
// No If* method works the same for all shapes.
// Naturally, we want to test each method with each shape.
//
// Writing all 16 cases by hand would be error prone.
// Adding a new shape later on would require writing 9 new tests.
// More generally, given n shapes, adding a new one requires adding 2n - 1 tests.
//
// Here's the trick to making the effort O(1).
// The If* methods all take callbacks.
// For most shape-method pairs, the callbacks should not be called.
// The callbacks should be called only when the shape and method match up.
//
// Each test case needs a tree and a way to check if the method calls the callback.
// We can specify each n times, to get n^2 test cases.
// Adding a new shape only requires constant work.
//
// We actually package those two things up with the shape and method name.
// That's a technicality though.

func TestTreeShapeMethodMatrix(t *testing.T) {
	cases := []testCaseParts{
		{
			SymTree{}, "invalid", "IfInvalid", func(t SymTree, wasCalled *bool) {
				t.IfInvalid(func() { *wasCalled = true })
			},
		}, {
			Sym("a-symbol"), "symbol", "IfSymbol", func(t SymTree, wasCalled *bool) {
				t.IfSymbol(func(_ string) { *wasCalled = true })
			},
		}, {
			Num(13), "number", "IfNumber", func(t SymTree, wasCalled *bool) {
				t.IfNumber(func(_ int) { *wasCalled = true })
			},
		}, {
			Lst(), "list", "IfList", func(t SymTree, wasCalled *bool) {
				t.IfList(func(_ List) { *wasCalled = true })
			},
		},
	}
	for _, tree := range cases {
		for _, method := range cases {
			mc := testCase{tree: tree, method: method}
			t.Run(mc.Name(), mc.Test)
		}
	}
}

type testCase struct {
	tree   testCaseParts
	method testCaseParts
}

type testCaseParts struct {
	Tree       SymTree
	ShapeName  string
	MethodName string
	Check      func(SymTree, *bool)
}

func (mc *testCase) Name() string {
	return fmt.Sprintf("%sCall%s", mc.tree.ShapeName, mc.method.MethodName)
}

func (mc *testCase) Test(t *testing.T) {
	wasCalled := false
	mc.method.Check(mc.tree.Tree, &wasCalled)
	if mc.tree.ShapeName != mc.method.ShapeName {
		assert(
			t.Errorf, !wasCalled,
			"a %s should not call the function passed to it's %s method",
			mc.tree.ShapeName, mc.method.MethodName,
		)
	} else {
		assert(
			t.Errorf, wasCalled,
			"a %s should call the function passed to it's %s method",
			mc.tree.ShapeName, mc.method.MethodName,
		)
	}
}
