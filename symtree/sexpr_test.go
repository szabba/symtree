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
		"oneLetterSymbol":   {"x", Sym("x"), causeEOF},
		"multiLetterSymbol": {"abba", Sym("abba"), causeEOF},
		"firstSymbol":       {"x y z", Sym("x"), noError},
		"emptyList":         {"()", Lst(), noError},
		"unmatchedList":     {"(", SymTree{}, causeUnexpectedEOF},
		"singletonList":     {"(+)", Lst(Sym("+")), noError},
		"nestedEmptyList":   {"(())", Lst(Lst()), noError},
		"nestedList": {
			"(+ x (/ y z))",
			Lst(Sym("+"), Sym("x"), Lst(Sym("/"), Sym("y"), Sym("z"))),
			noError,
		},
		"firstList":        {"(+) (abba u2 rem)", Lst(Sym("+")), noError},
		"digit":            {"7", Num(7), causeEOF},
		"multiDigitNumber": {"13", Num(13), causeEOF},
		"negativeNumber":   {"-9", Num(-9), causeEOF},
		"listWithNumbers":  {"(+ 13 x)", Lst(Sym("+"), Num(13), Sym("x")), noError},
		"initSpaceSkipped": {"   +", Sym("+"), causeEOF},
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
		"symbol":     {Sym("+"), "+"},
		"number":     {Num(13), "13"},
		"emptyList":  {Lst(), "()"},
		"flatList":   {Lst(Sym("+"), Num(13), Num(4)), "(+ 13 4)"},
		"nestedList": {Lst(Lst()), "(())"},
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
