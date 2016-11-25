//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

import "fmt"

// A SymbolAlreadyBound error means an attempt to override a binding has happened.
type SymbolAlreadyBound interface {
	error
	Symbol() string
	CurrentBinding() Tree
}

type alreadyBound struct {
	symbol string
	curr   Tree
}

var _ SymbolAlreadyBound = alreadyBound{}

func (crs alreadyBound) Error() string {
	return fmt.Sprintf(
		"cannot rebind %s, already bound to %v",
		crs.Symbol(), crs.CurrentBinding(),
	)
}

func (crs alreadyBound) Symbol() string { return crs.symbol }

func (crs alreadyBound) CurrentBinding() Tree { return crs.curr }

// A NotBound error means a symbol that should be bound was not.
type NotBound interface {
	error
	UnboundSymbol() string
}

type notBound string

func (nb notBound) Error() string {
	return fmt.Sprintf("unbound symbol %q", nb.UnboundSymbol())
}

func (nb notBound) UnboundSymbol() string { return string(nb) }

// An ExactMismatch error means that a subtree didn't match a hole-free subpattern.
type ExactMismatch interface {
	error
	Expected() Tree
	Got() Tree
}

type exactMismatch struct {
	expected, got Tree
}

var _ ExactMismatch = exactMismatch{}

func (m exactMismatch) Error() string {
	return fmt.Sprintf("expected %v, got %v", m.Expected(), m.Got())
}

func (m exactMismatch) Expected() Tree { return m.expected }
func (m exactMismatch) Got() Tree      { return m.got }

// An AtomCannotMatchList error means that a tree was an atom and a pattern was for a list.
type AtomCannotMatchList interface {
	error
	AtomCannotMatchList()
}

type atomCannotMatchList struct{}

var _ AtomCannotMatchList = atomCannotMatchList{}

func (_ atomCannotMatchList) Error() string { return "atom cannot match list" }

func (_ atomCannotMatchList) AtomCannotMatchList() {}

// A LengthMismatch error means that the list has a different length than the pattern.
type LengthMismatch interface {
	error
	ExpectedLen() int
	GotLen() int
}

type lenMismatch struct {
	expected, got int
}

var _ LengthMismatch = lenMismatch{}

func (lm lenMismatch) Error() string {
	return fmt.Sprintf("expected a list of length %d, not %d", lm.ExpectedLen(), lm.GotLen())
}

func (lm lenMismatch) ExpectedLen() int { return lm.expected }
func (lm lenMismatch) GotLen() int      { return lm.got }
