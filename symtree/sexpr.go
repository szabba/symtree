//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// ReadSexpr reads an s-expression form of a symtree from an io.Reader.
// It is reasonable to expect the
func ReadSexpr(src io.RuneScanner) (SymTree, error) {
	r := reader{src: src}
	return r.parse()
}

type reader struct {
	src io.RuneScanner
	err error
}

func (r *reader) parse() (SymTree, error) {
	r.skipWhile(unicode.IsSpace)
	if r.peek() == '(' {
		return r.parseList()
	}
	return r.parseAtom()
}

func (r *reader) parseList() (SymTree, error) {
	r.accept()
	elems := r.parseListElements()
	if r.peek() != ')' {
		r.eofNotExpected()
	}
	return r.resultIfNoError(Lst(elems...))
}

func (r *reader) resultIfNoError(tree SymTree) (SymTree, error) {
	if r.err != nil {
		return SymTree{}, r.err
	}
	return tree, nil
}

func (r *reader) eofNotExpected() {
	if r.err != io.EOF {
		return
	}
	r.err = io.ErrUnexpectedEOF
}

func (r *reader) parseListElements() []SymTree {
	r.skipWhile(unicode.IsSpace)

	var elems []SymTree
	for r.err == nil && r.peek() != ')' {
		elems = append(elems, r.parseListElem())
	}
	return elems
}

func (r *reader) parseListElem() SymTree {
	elem, _ := r.parse()
	r.skipWhile(unicode.IsSpace)
	return elem
}

func (r *reader) parseAtom() (SymTree, error) {
	atom := r.readWhile(isAtom)

	if atom == "" {
		return SymTree{}, r.err
	}

	if num, err := strconv.Atoi(atom); err == nil {
		return Num(num), r.err
	}

	return Sym(atom), r.err
}

func isAtom(chr rune) bool {
	return !unicode.IsSpace(chr) && chr != ')'
}

func (r *reader) readWhile(f func(rune) bool) string {
	var buf bytes.Buffer
	r.takeWhile(f, &buf)
	return buf.String()
}

func (r *reader) skipWhile(f func(rune) bool) {
	r.takeWhile(f, ioutil.Discard)
}

func (r *reader) takeWhile(f func(rune) bool, dst io.Writer) {
	var buf [4]byte
	for chr := r.peek(); r.err == nil && f(chr); chr = r.peek() {
		r.accept()
		runeLen := utf8.EncodeRune(buf[:], chr)
		dst.Write(buf[:runeLen])
	}
}

func (r *reader) peek() rune {
	chr := r.read()
	r.unread()
	return chr
}

func (r *reader) accept() rune { return r.read() }

func (r *reader) read() rune {
	var chr rune
	if r.err != nil {
		return chr
	}
	chr, _, r.err = r.src.ReadRune()
	return chr
}

func (r *reader) unread() {
	if r.err != nil {
		return
	}
	r.err = r.src.UnreadRune()
}

// WriteSexpr writes the s-expression form of t into dst.
func WriteSexpr(dst io.Writer, t SymTree) (n int, err error) {
	w := writer{dst: dst}
	w.WriteTree(t)
	return w.n, w.err
}

type writer struct {
	dst io.Writer
	n   int
	err error
}

func (w *writer) WriteTree(t SymTree) {
	t.IfInvalid(func() { w.Write("<invalid symtree>") })
	t.IfSymbol(func(s string) { w.Write(s) })
	t.IfNumber(func(n int) { w.Write(n) })
	t.IfList(w.WriteList)
}

func (w *writer) WriteList(list List) {
	w.Write("(")
	for i := 0; i < list.Len(); i++ {
		w.WriteElement(list, i)
	}
	w.Write(")")
}

func (w *writer) WriteElement(list List, i int) {
	w.WriteTree(list.At(i))
	if i+1 < list.Len() {
		w.Write(" ")
	}
}

func (w *writer) Write(v interface{}) {
	if w.err != nil {
		return
	}
	var n int
	n, w.err = fmt.Fprint(w.dst, v)
	w.n += n
}
