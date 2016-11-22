//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

import "bytes"

// A sexpr wraps a SymTree and adds a String method.
// The output is the same as that of WriteSexpr.
type sexpr struct{ SymTree }

func (s sexpr) String() string {
	var b bytes.Buffer
	WriteSexpr(&b, s.SymTree)
	return b.String()
}
