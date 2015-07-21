package browser

import (
	"yang"
)

type browserError struct {
	Msg string
}

func (err *browserError) Error() string {
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

type Editor interface {
	// Receiver
	StartTransaction()
	NewObject(yang.MetaList)
	NewList(*yang.List)
	PutIntLeaf(*yang.Leaf, int)
	PutStringLeaf(*yang.Leaf, string)
	PutStringLeafList(*yang.LeafList, []string)
	PutIntLeafList(*yang.LeafList, []int)
	ExitObject(yang.MetaList)
	ExitList(*yang.List)
	EndTransaction()
    // Editor
	DeleteObject(*yang.Container)
	DeleteList(*yang.List)
}

type Transmitter interface {
	Transmit() error
}
