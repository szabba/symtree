//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package symtree implements immutable, s-expression-like trees.
package symtree

// A SymTree is either
//
//     - invalid
//     - a symbol
//     - a number
//     - a list of other SymTrees
//
// SymTrees are immutable.
//
// Call one of the If* methods to try to access the data.
//
// The == operator doesn't work for SymTrees.
// Use the Equal function instead.
type SymTree struct {
	// A symbolic tree consists of a tag and one field for each valid shape.
	// The tag tells us which shape the value has.

	tag    symTreeTag
	symbol string
	number int
	list   List
}

// Sym creates a SymTree that calls the callback passed to IfSymbol.
func Sym(name string) SymTree {
	return SymTree{tag: symTreeSymbol, symbol: name}
}

// Num creates a SymTree that calls the callback passed to IfNumber.
func Num(n int) SymTree {
	return SymTree{tag: symTreeNumber, number: n}
}

// Lst creates a SymTree that calls the callback passed to IfList.
func Lst(elems ...SymTree) SymTree {
	l := List{elements: make([]SymTree, len(elems))}
	copy(l.elements, elems)
	return SymTree{tag: symTreeList, list: l}
}

// IfInvalid calls f if the receiver is not a valid SymTree.
// There are two ways to get an invalid SymTree.
// One is to use the zero value.
// The other is to call List.At with an index that's out of bounds.
func (tree SymTree) IfInvalid(f func()) {
	if tree.tag != symTreeInvalid {
		return
	}
	f()
}

// IfSymbol calls f if the receiver is a symbol SymTree.
// The argument passed in is the symbol's string representation.
func (tree SymTree) IfSymbol(f func(string)) {
	if tree.tag != symTreeSymbol {
		return
	}
	f(tree.symbol)
}

// IfNumber calls f if the receiver is a number SymTree.
// The argument passed in is the value of the number.
func (tree SymTree) IfNumber(f func(int)) {
	if tree.tag != symTreeNumber {
		return
	}
	f(tree.number)
}

// IfList calls f if the receiver is a list SymTree.
// The argument passed in is the List containing all the children of the SymTree.
func (tree SymTree) IfList(f func(List)) {
	if tree.tag != symTreeList {
		return
	}
	f(tree.list)
}

// A List is an immutable sequence of SymTrees.
type List struct {
	len      int
	elements []SymTree
}

// Len returns the number of elements in the List.
func (l List) Len() int { return len(l.elements) }

// At looks up the i-th element of the List.
// If i is out of bounds, then the returned tree is invalid.
func (l List) At(i int) SymTree {
	if i < 0 || i >= len(l.elements) {
		return SymTree{}
	}
	return l.elements[i]
}

// Equal compares two SymTrees for structural equality.
// The == operator doesn't work on SymTrees.
func Equal(a, b SymTree) bool {
	if a.tag == symTreeInvalid && b.tag == symTreeInvalid {
		return true
	}
	if a.tag == symTreeSymbol && b.tag == symTreeSymbol && a.symbol == b.symbol {
		return true
	}
	if a.tag == symTreeNumber && b.tag == symTreeNumber && a.number == b.number {
		return true
	}
	if a.tag == symTreeList && b.tag == symTreeList {
		return equalLists(a.list, b.list)
	}
	return false
}

func equalLists(a, b List) bool {
	eq := a.Len() == b.Len()
	for i := 0; i < a.Len(); i++ {
		eq = eq && Equal(a.At(i), b.At(i))
	}
	return eq
}

type symTreeTag int

const (
	symTreeInvalid symTreeTag = iota
	symTreeSymbol
	symTreeNumber
	symTreeList
)