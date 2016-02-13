package data
import (
	"schema"
	"conf2"
)

func (self *Selection) Action(input Node) (*Selection, error) {
conf2.Debug.Printf("Action self=%p", self)
	rpc := self.path.meta.(*schema.Rpc)
	in := Select(rpc.Input, input)
	rpcOutput, rerr := self.node.Action(self, rpc, in)
	if rerr != nil {
		return nil, rerr
	}
	if rpcOutput != nil {
		return Select(rpc.Output, rpcOutput), nil
	}
	return nil, nil
}
