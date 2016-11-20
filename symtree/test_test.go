package symtree

func assert(onErr func(string, ...interface{}), cond bool, format string, a ...interface{}) {
	if !cond {
		onErr(format, a...)
	}
}
