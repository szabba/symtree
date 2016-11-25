//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

import "bytes"

// A sexpr wraps a Tree and adds a String method.
// The output is the same as that of WriteSexpr.
type sexpr struct{ Tree }

func (s sexpr) String() string {
	var b bytes.Buffer
	WriteSexpr(&b, s.Tree)
	return b.String()
}
