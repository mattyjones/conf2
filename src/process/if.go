package process
import (
	"fmt"
	"strings"
)

type If struct {
	parent CodeBlock
	Expression Expression
	operations []Operation
}

func (iF *If) Operations() []Operation {
	return iF.operations
}

func (iF *If) SetParent(parent CodeBlock) {
	iF.parent = parent
}

func (iF *If) Parent() CodeBlock {
	return iF.parent
}

func (iF *If) AddOperation(op Operation) {
	if iF.operations == nil {
		iF.operations = make([]Operation, 0, 1)
	}
	iF.operations = append(iF.operations, op)
}

func (iF *If) IsTrue(val interface{}) bool {
	switch x := val.(type) {
	case string:
		clean := strings.TrimSpace(x)
		if len(clean) == 0 || "false" == clean || "0" == clean {
			return false
		}
		return true
	case int:
		return x != 0
	case bool:
		return x
	default:
		s := fmt.Sprintf("%v", val)
		return iF.IsTrue(s)
	}
	return false
}

func (iF *If) Exec(stack *Stack, table Table) error {
	condition, err := iF.Expression.Eval(stack, table)
	if err != nil {
		return err
	}
	if iF.IsTrue(condition) {
		return nil
	}
	return execOperations(stack, iF.operations, table)
}


