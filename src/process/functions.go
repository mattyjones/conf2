package process
import (
	"fmt"
	"bytes"
)

var builtins map[string]interface{}
func init() {
	builtins = map[string]interface{} {
		"concat" : ConcatFunction,
		"contains" : ContainsFunction,
		"equals" : func(stack *Stack, table Table, a interface{}, b interface{}) bool {
			return ArgAsString(a) == ArgAsString(b)
		},
		"length" : func(stack *Stack, table Table, arg interface{}) int {
			return len(ArgAsString(arg))
		},
		"not" : func(stack *Stack, table Table, arg bool) bool {
			return !arg
		},
	}
}

func ContainsFunction(stack *Stack, table Table, ident string) bool {
	x, err := table.Get(ident)
	if x != nil && err != nil {
		return true
	}
	x, err = table.Select(ident, false)
	if x != nil && err != nil {
		return true
	}
	return false
}

func ArgAsString(arg interface{}) string {
	if tostr, ok := arg.(fmt.Stringer); ok {
		return tostr.String()
	}

	return fmt.Sprintf("%v", arg)
}

func ConcatFunction(stack *Stack, table Table, args ...interface{}) string {
	var buff bytes.Buffer
	for _, arg := range args {
		buff.WriteString(ArgAsString(arg))
	}
	return buff.String()
}
