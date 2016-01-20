package data
import (
	"schema"
)

func (self *Selection) Action(input Node) (*Selection, error) {
	rpc := self.path.meta.(*schema.Rpc)
	in := NewSelection(rpc.Input, input)
	rpcOutput, rerr := self.node.Action(self, rpc, in)
	if rerr != nil {
		return nil, rerr
	}
	if rpcOutput != nil {
		return NewSelection(rpc.Output, rpcOutput), nil
	}
	return nil, nil
}
