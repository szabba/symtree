//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestReadSexpr(t *testing.T) {
	type testcase struct {
		input      string
		tree       SymTree
		checkError func(*testing.T, error)
	}
	noError := func(t *testing.T, err error) {
		assert(t.Errorf, nil == err, "expected no error, got %q", err)
	}

	causeEOF := func(t *testing.T, err error) {
		cause := errors.Cause(err)
		assert(t.Errorf, io.EOF == cause, "expected error cause %q, got %q", io.EOF, cause)
	}

	causeUnexpectedEOF := func(t *testing.T, err error) {
		cause := errors.Cause(err)
		assert(t.Errorf, io.ErrUnexpectedEOF == cause, "expected error cause %q, got %q", io.ErrUnexpectedEOF, cause)
	}

	cases := map[string]testcase{
		"noInput":           {"", SymTree{}, causeEOF},
		"oneLetterSymbol":   {"x", NewSymbol("x"), causeEOF},
		"multiLetterSymbol": {"abba", NewSymbol("abba"), causeEOF},
		"firstSymbol":       {"x y z", NewSymbol("x"), noError},
		"emptyList":         {"()", NewList(), noError},
		"unmatchedList":     {"(", SymTree{}, causeUnexpectedEOF},
		"singletonList":     {"(+)", NewList(NewSymbol("+")), noError},
		"nestedEmptyList":   {"(())", NewList(NewList()), noError},
		"nestedList": {
			"(+ x (/ y z))",
			NewList(NewSymbol("+"), NewSymbol("x"), NewList(NewSymbol("/"), NewSymbol("y"), NewSymbol("z"))),
			noError,
		},
		"firstList":        {"(+) (abba u2 rem)", NewList(NewSymbol("+")), noError},
		"digit":            {"7", NewNumber(7), causeEOF},
		"multiDigitNumber": {"13", NewNumber(13), causeEOF},
		"negativeNumber":   {"-9", NewNumber(-9), causeEOF},
		"listWithNumbers":  {"(+ 13 x)", NewList(NewSymbol("+"), NewNumber(13), NewSymbol("x")), noError},
		"initSpaceSkipped": {"   +", NewSymbol("+"), causeEOF},
	}

	for name, kase := range cases {
		t.Run(name, func(t *testing.T) {

			tree, err := ReadSexpr(strings.NewReader(kase.input))

			assert(t.Errorf, Equal(kase.tree, tree), "expected tree %v, got %v", sexpr{kase.tree}, sexpr{tree})
			kase.checkError(t, err)
		})
	}
}

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
