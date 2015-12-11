package data
import (
	"schema"
)

func PathAction(data Data, path string, input Node, output Node) (error) {
	p, err := schema.ParsePath(path, data.Schema())
	if err != nil {
		return  err
	}
	s, serr := WalkPath(NewSelection(data.Node(), data.Schema()), p)
	if serr != nil {
		return serr
	}
	return SelectionAction(s, input, output)
}

func SelectionAction(sel *Selection, input Node, output Node) (error) {
	rpc := sel.State.Position().(*schema.Rpc)
	rpcOutput, rerr := sel.Node.Action(sel, rpc, input)
	if rerr != nil {
		return rerr
	}
	if rpc.Output != nil && output != nil {
		return NodeToNode(rpcOutput, output, rpc.Output).Insert()
	}
	return nil
}
