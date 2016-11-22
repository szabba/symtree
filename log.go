//   This Source Code Form is subject to the terms of the Mozilla Public
//   License, v. 2.0. If a copy of the MPL was not distributed with this
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.

package symtree

// SetDebugFn sets a function f the package will use for debug logging.
func SetDebugFn(f func(format string, a ...interface{})) {
	logTo.Debug = f
}

var logTo struct {
	Debug func(format string, a ...interface{})
}

func debug(format string, a ...interface{}) {
	f := logTo.Debug
	if f == nil {
		return
	}
	f(format, a...)
}
