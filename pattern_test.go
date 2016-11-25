//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

import (
	"testing"

	"github.com/pkg/errors"
)

func TestLiteralPatternMatchesItsOwnTree(t *testing.T) {
	pattern := FromExample([]string{}, Sym("+"))

	match := map[string]Tree{}
	err := pattern.Match(Sym("+"), match)

	assert(t.Errorf, err == nil, "unexpected error: %s", err)
}

func TestLiteralPatternDoesNotMatchDifferentTree(t *testing.T) {
	pattern := FromExample([]string{}, Sym("+"))
	otherTree := Lst(Sym("-"), Num(17), Sym("y"))

	match := map[string]Tree{}
	err := pattern.Match(otherTree, match)

	assert(t.Errorf, err != nil, "expected an error")

	mismatch, typeOK := errors.Cause(err).(ExactMismatch)
	assert(t.Fatalf, typeOK, "error returned should implement ExactMismatch")
	assert(
		t.Errorf, Equal(Sym("+"), mismatch.Expected()),
		"value expected should be %v, not %v",
		Sym("+"), mismatch.Expected(),
	)
	assert(
		t.Errorf, Equal(otherTree, mismatch.Got()),
		"value got should be %v, not %v",
		otherTree, mismatch.Got(),
	)
}

func TestLiteralPatternProducesItsOwnExample(t *testing.T) {
	pattern := FromExample([]string{}, Sym("+"))

	tree, err := pattern.Substitute(map[string]Tree{})

	assert(t.Errorf, err == nil, "unexpected error: %s", err)
	assert(
		t.Errorf, Equal(Sym("+"), tree),
		"expected %v, got %v", Sym("+"), tree,
	)
}

func TestVariablePatternMatchesATree(t *testing.T) {
	pattern := FromExample([]string{"x"}, Sym("x"))
	tree := Lst(Sym("-"), Num(17), Sym("y"))

	match := map[string]Tree{}
	err := pattern.Match(tree, match)

	assert(t.Errorf, err == nil, "unexpected error: %s", err)
	_, varBound := match["x"]
	assert(t.Errorf, varBound, "variable %q should get bound in the match", "x")
	assert(t.Errorf, Equal(tree, match["x"]), "expected %v, got %v", tree, match["x"])
}

func TestVariablePatternDoesNotRebindWhenMatching(t *testing.T) {
	pattern := FromExample([]string{"x"}, Sym("x"))
	tree := Lst(Sym("-"), Num(17), Sym("y"))

	match := map[string]Tree{"x": Num(19)}
	err := pattern.Match(tree, match)

	assert(t.Fatalf, err != nil, "expected an error")
	alreadyBound, typeOK := err.(SymbolAlreadyBound)
	assert(t.Fatalf, typeOK, "error returned should implement SymbolAlreadyBound")
	assert(t.Errorf, "x" == alreadyBound.Symbol(), "expected symbol %q, not %q", "x", alreadyBound.Symbol())
	assert(
		t.Errorf, Equal(Num(19), alreadyBound.CurrentBinding()),
		"expected current binding %v, not %v", Num(10), alreadyBound.CurrentBinding(),
	)
}

func TestVariablePatternDoesNotFailIfNewBindingWouldBeTheSameAsTheOld(t *testing.T) {
	pattern := FromExample([]string{"x"}, Sym("x"))
	tree := Num(13)

	match := map[string]Tree{"x": Num(13)}
	err := pattern.Match(tree, match)

	assert(t.Errorf, err == nil, "unexpected error: %s", err)
}

func TestVariablePatternSubstitutionWithoutABindingFails(t *testing.T) {
	pattern := FromExample([]string{"x"}, Sym("x"))
	match := map[string]Tree{}

	tree, err := pattern.Substitute(match)

	assert(t.Errorf, Equal(Tree{}, tree), "expected %v, got %v", Tree{}, tree)
	assert(t.Fatalf, err != nil, "expected an error")
	notBound, typeOK := err.(NotBound)
	assert(t.Fatalf, typeOK, "error returned should implement NotBound")
	assert(
		t.Errorf, "x" == notBound.UnboundSymbol(),
		"unbound symbol should be %q, not %q", "x", notBound.UnboundSymbol(),
	)
}

func TestVariablePatternGetsSubstitutedWithBinding(t *testing.T) {
	pattern := FromExample([]string{"x"}, Sym("x"))
	match := map[string]Tree{"x": Num(13)}

	tree, err := pattern.Substitute(match)

	assert(t.Errorf, err == nil, "unexpected error: %s", err)
	assert(t.Errorf, Equal(match["x"], tree), "expected %v, got %v", match["x"], tree)
}

func TestListPatternFailsToMatchASymbol(t *testing.T) {
	pattern := FromExample([]string{}, Lst())
	tree := Sym("+")

	match := map[string]Tree{}
	err := pattern.Match(tree, match)

	assert(t.Fatalf, err != nil, "expected an error")
	_, typeOK := err.(AtomCannotMatchList)
	assert(t.Errorf, typeOK, "error returned should implement AtomCannotMatchList")
}

func TestListPatternFailsToMatchANumber(t *testing.T) {
	pattern := FromExample([]string{}, Lst())
	tree := Num(13)

	match := map[string]Tree{}
	err := pattern.Match(tree, match)

	assert(t.Fatalf, err != nil, "expected an error")
	_, typeOK := err.(AtomCannotMatchList)
	assert(t.Errorf, typeOK, "error returned should implement AtomCannotMatchList")
}

func TestListPatternFailsToMatchAnInvalidTree(t *testing.T) {
	pattern := FromExample([]string{}, Lst())
	tree := Tree{}

	match := map[string]Tree{}
	err := pattern.Match(tree, match)

	assert(t.Fatalf, err != nil, "expected an error")
	_, typeOK := err.(AtomCannotMatchList)
	assert(t.Errorf, typeOK, "error returned should implement AtomCannotMatchList")
}

func TestListPatternFailsToMatchAListTreeOfDifferentLength(t *testing.T) {
	pattern := FromExample([]string{}, Lst(Sym("+"), Num(13)))
	otherTree := Lst(Sym("+"), Num(13), Num(7))

	match := map[string]Tree{}
	err := pattern.Match(otherTree, match)

	assert(t.Fatalf, err != nil, "expected an error")
	lenMismatch, typeOK := err.(LengthMismatch)
	assert(t.Fatalf, typeOK, "error returned should implement LengthMismatch")
	assert(
		t.Errorf, 2 == lenMismatch.ExpectedLen(),
		"length expected should be %d, not %d", 2, lenMismatch.ExpectedLen(),
	)
	assert(
		t.Errorf, 3 == lenMismatch.GotLen(),
		"length got should be %d, not %d", 3, lenMismatch.GotLen(),
	)
}

func TestListPatternOfLiteralPatternsMatchesItsDefinition(t *testing.T) {
	def := Lst(Sym("+"), Num(13))
	pattern := FromExample([]string{}, def)

	match := map[string]Tree{}
	err := pattern.Match(def, match)

	assert(t.Errorf, err == nil, "unexpected error: %s", err)
}

func TestListPatternWithVariableMatch(t *testing.T) {
	pattern := FromExample(
		[]string{"x"},
		Lst(Sym("+"), Sym("x"), Num(13)),
	)
	tree := Lst(Sym("+"), Num(7), Num(13))

	match := map[string]Tree{}
	err := pattern.Match(tree, match)

	assert(t.Errorf, err == nil, "unexpected error: %s", err)
	assert(
		t.Errorf, Equal(Num(7), match["x"]),
		"expected %v, got %v", Num(7), match["x"],
	)
}

func TestNestedListPatternMatch(t *testing.T) {
	pattern := FromExample(
		[]string{"x", "y"},
		Lst(
			Sym("+"),
			Lst(Sym("sin"), Sym("x")),
			Lst(Sym("cos"), Sym("y")),
		),
	)
	tree := Lst(
		Sym("+"),
		Lst(Sym("sin"), Sym("π")),
		Lst(Sym("cos"), Sym("τ")),
	)

	match := map[string]Tree{}
	err := pattern.Match(tree, match)

	assert(t.Errorf, err == nil, "unexpected error: %s", err)
	x, xBound := match["x"]
	y, yBound := match["y"]
	assert(t.Errorf, xBound, "variable %q unbound", "x")
	assert(
		t.Errorf, !xBound || Equal(Sym("π"), x),
		"expected %v, got %v", Sym("π"), x,
	)
	assert(t.Errorf, yBound, "variable %q unbound", "y")
	assert(
		t.Errorf, !yBound || Equal(Sym("τ"), y),
		"expected %v, got %v", Sym("τ"), y,
	)
}

func TestListPatternSubstitute(t *testing.T) {
	pattern := FromExample(
		[]string{"x"}, Lst(Sym("sin"), Sym("x")),
	)

	match := map[string]Tree{"x": Sym("π")}
	tree, err := pattern.Substitute(match)

	assert(t.Errorf, err == nil, "unexpected error: %s", err)
	expectedTree := Lst(Sym("sin"), Sym("π"))
	assert(
		t.Errorf, Equal(expectedTree, tree),
		"expected %v, got %v", expectedTree, tree,
	)
}
