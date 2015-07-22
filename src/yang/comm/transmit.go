package comm

import (
	"yang"
)

type commError struct {
	Msg string
}

func (err *commError) Error() string {
	return err.Msg
}

type Visitor interface {
	Visit(*yang.Meta, interface{})
}

type Receiver interface {
	StartTransaction()
	NewObject(yang.MetaList)
	NewList(*yang.List)
	NewListItem(yang.Meta)
	PutIntLeaf(*yang.Leaf, int)
	PutStringLeaf(*yang.Leaf, string)
	PutStringLeafList(*yang.LeafList, []string)
	PutIntLeafList(*yang.LeafList, []int)
	ExitList(*yang.List)
	ExitListItem(yang.Meta)
	ExitObject(yang.MetaList)
	EndTransaction()
}

type Transmitter interface {
	Transmit() error
}
