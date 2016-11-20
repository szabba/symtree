//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

import (
	"bytes"
	"testing"
)

func TestWriteSexpr(t *testing.T) {
	type testcase struct {
		tree     SymTree
		expected string
	}
	cases := map[string]testcase{
		"invalid":    {SymTree{}, "<invalid symtree>"},
		"symbol":     {NewSymbol("+"), "+"},
		"number":     {NewNumber(13), "13"},
		"emptyList":  {NewList(), "()"},
		"flatList":   {NewList(NewSymbol("+"), NewNumber(13), NewNumber(4)), "(+ 13 4)"},
		"nestedList": {NewList(NewList()), "(())"},
	}

	for name, kase := range cases {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer

			WriteSexpr(&b, kase.tree)

			assert(
				t.Errorf, kase.expected == b.String(),
				"expected %q, got %q", kase.expected, b.String(),
			)
		})
	}
}
