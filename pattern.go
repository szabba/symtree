//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

// A Pattern is something a tree can be matched against or built from.
// A pattern is like a tree with named holes.
// The same hole name can appear in more than one place in a pattern.
type Pattern interface {
	// Match checks whether the candidate matches the pattern.
	// Any subtrees corresponding to named holes will be stored in the match.
	// For a tree to match, every hole with the same name must correspond to an identical subtree.
	//
	// Even when matching fails, the match can be modified.
	Match(candidate Tree, match map[string]Tree) error

	// Substitute creates a new tree based on the pattern.
	// Locations of named holes are replaced with values from match.
	//
	// It fails if a named hole has no value in match.
	Substitute(match map[string]Tree) (Tree, error)
}

// FromExample builds a pattern from a tree and a list of hole names.
// holes contains the names of symbols in the tree that should be considered named hole positions.
func FromExample(holes []string, expr Tree) Pattern {
	var p Pattern = litPattern{expr: expr}
	expr.IfList(fromList(holes, &p))
	expr.IfSymbol(fromSymbol(holes, &p))
	return p
}

func fromList(holes []string, p *Pattern) func(List) {
	return func(l List) {
		children := make(listPattern, l.Len())
		for i := range children {
			children[i] = FromExample(holes, l.At(i))
		}
		*p = children
	}
}

func fromSymbol(holes []string, p *Pattern) func(string) {
	return func(s string) {
		for _, hole := range holes {
			if s == hole {
				*p = holePattern{name: s}
			}
		}
	}
}

type listPattern []Pattern

var _ Pattern = listPattern{}

func (lp listPattern) Match(tree Tree, match map[string]Tree) error {
	var err error
	tree.IfInvalid(func() { err = atomCannotMatchList{} })
	tree.IfSymbol(func(_ string) { err = atomCannotMatchList{} })
	tree.IfNumber(func(_ int) { err = atomCannotMatchList{} })
	tree.IfList(func(example List) {
		if len(lp) != example.Len() {
			err = lenMismatch{expected: len(lp), got: example.Len()}
			return
		}
		for i := 0; err == nil && i < len(lp); i++ {
			err = lp[i].Match(example.At(i), match)
		}
	})
	return err
}

func (lp listPattern) Substitute(match map[string]Tree) (Tree, error) {
	var err error
	elems := make([]Tree, len(lp))
	for i := 0; err == nil && i < len(lp); i++ {
		elems[i], err = lp[i].Substitute(match)
	}
	return Lst(elems...), err
}

type holePattern struct {
	name string
}

var _ Pattern = holePattern{}

func (hp holePattern) Match(tree Tree, match map[string]Tree) error {
	if curr, bound := match[hp.name]; bound && !Equal(curr, tree) {
		return alreadyBound{symbol: hp.name, curr: curr}
	}
	match[hp.name] = tree
	return nil
}

func (hp holePattern) Substitute(match map[string]Tree) (Tree, error) {
	tree, bound := match[hp.name]
	if !bound {
		return Tree{}, notBound(hp.name)
	}
	return tree, nil
}

type litPattern struct {
	expr Tree
}

var _ Pattern = litPattern{}

func (lit litPattern) Match(tree Tree, _ map[string]Tree) error {
	if !Equal(lit.expr, tree) {
		return exactMismatch{expected: lit.expr, got: tree}
	}
	return nil
}

func (lit litPattern) Substitute(_ map[string]Tree) (Tree, error) {
	return lit.expr, nil
}
