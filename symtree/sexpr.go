//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

import (
	"fmt"
	"io"
)

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
