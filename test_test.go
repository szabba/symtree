//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

func assert(onErr func(string, ...interface{}), cond bool, format string, a ...interface{}) {
	if !cond {
		onErr(format, a...)
	}
}
