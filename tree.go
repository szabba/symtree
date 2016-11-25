//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package symtree implements immutable, s-expression-like trees.
package symtree

// A Tree is either
//
//     - invalid
//     - a symbol
//     - a number
//     - a list of other Trees
//
// Trees are immutable.
//
// Call one of the If* methods to try to access the data.
//
// The == operator doesn't work for Trees.
// Use the Equal function instead.
type Tree struct {
	// A symbolic tree consists of a tag and one field for each valid shape.
	// The tag tells us which shape the value has.

	tag    symTreeTag
	symbol string
	number int
	list   List
}

// Sym creates a Tree that calls the callback passed to IfSymbol.
func Sym(name string) Tree {
	return Tree{tag: symTreeSymbol, symbol: name}
}

// Num creates a Tree that calls the callback passed to IfNumber.
func Num(n int) Tree {
	return Tree{tag: symTreeNumber, number: n}
}

// Lst creates a Tree that calls the callback passed to IfList.
func Lst(elems ...Tree) Tree {
	l := List{elements: make([]Tree, len(elems))}
	copy(l.elements, elems)
	return Tree{tag: symTreeList, list: l}
}

// IfInvalid calls f if the receiver is not a valid Tree.
// There are two ways to get an invalid Tree.
// One is to use the zero value.
// The other is to call List.At with an index that's out of bounds.
func (tree Tree) IfInvalid(f func()) {
	if tree.tag != symTreeInvalid {
		return
	}
	f()
}

// IfSymbol calls f if the receiver is a symbol Tree.
// The argument passed in is the symbol's string representation.
func (tree Tree) IfSymbol(f func(string)) {
	if tree.tag != symTreeSymbol {
		return
	}
	f(tree.symbol)
}

// IfNumber calls f if the receiver is a number Tree.
// The argument passed in is the value of the number.
func (tree Tree) IfNumber(f func(int)) {
	if tree.tag != symTreeNumber {
		return
	}
	f(tree.number)
}

// IfList calls f if the receiver is a list Tree.
// The argument passed in is the List containing all the children of the Tree.
func (tree Tree) IfList(f func(List)) {
	if tree.tag != symTreeList {
		return
	}
	f(tree.list)
}

// A List is an immutable sequence of Trees.
type List struct {
	len      int
	elements []Tree
}

// Len returns the number of elements in the List.
func (l List) Len() int { return len(l.elements) }

// At looks up the i-th element of the List.
// If i is out of bounds, then the returned tree is invalid.
func (l List) At(i int) Tree {
	if i < 0 || i >= len(l.elements) {
		return Tree{}
	}
	return l.elements[i]
}

// Equal compares two Trees for structural equality.
// The == operator doesn't work on Trees.
func Equal(a, b Tree) bool {
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
