package data

import (
	"schema"
	"fmt"
)


// when writing values, splits output into two nodes.
// when reading, reads from secondary only when primary returns nil
type Tee struct {
	Primary Node
	Secondary Node
}

func (self *Tee) String() string {
	return fmt.Sprintf("Tee(%s,%s)", self.Primary.String(), self.Secondary.String())
}

func (self *Tee) Select(sel *Selection, r ContainerRequest) (Node, error) {
	var err error
	var child Tee
	if child.Primary, err = self.Primary.Select(sel, r); err != nil {
		return nil, err
	}
	if child.Secondary, err = self.Secondary.Select(sel, r); err != nil {
		return nil, err
	}
	if child.Primary != nil && child.Secondary != nil {
		return &child, nil
	}
	return nil, nil
}

func (self *Tee)  Next(sel *Selection, r ListRequest) (Node, []*Value, error) {
	var err error
	var next Tee
	key := r.Key
	if next.Primary, key, err = self.Primary.Next(sel, r); err != nil {
		return nil, nil, err
	}
	if next.Secondary, _, err = self.Secondary.Next(sel, r); err != nil {
		return nil, nil, err
	}
	if next.Primary != nil && next.Secondary != nil {
		return &next, key, nil
	}
	return nil, nil, nil
}

func (self *Tee) Read(sel *Selection, meta schema.HasDataType) (*Value, error) {
	// merging results, prefer first
	if v, err := self.Primary.Read(sel, meta); err != nil {
		return nil, err
	} else if v != nil {
		return v, nil
	}
	return self.Secondary.Read(sel, meta)
}

func (self *Tee) Write(sel *Selection, meta schema.HasDataType, val *Value) (err error) {
	if err = self.Primary.Write(sel, meta, val); err == nil {
		err = self.Secondary.Write(sel, meta, val)
	}
	return
}

func (self *Tee) Choose(sel *Selection, choice *schema.Choice) (m schema.Meta, err error) {
	return self.Primary.Choose(sel, choice)
}

func (self *Tee) Event(sel *Selection, e Event) (err error) {
	if err = self.Primary.Event(sel, e); err == nil {
		err = self.Secondary.Event(sel, e)
	}
	return
}

func (self *Tee) Action(sel *Selection, meta *schema.Rpc, input *Selection) (output Node, err error) {
	return self.Primary.Action(sel, meta, input)
}

func (self *Tee) Peek(sel *Selection, peekId string) interface{} {
	if v := self.Primary.Peek(sel, peekId); v != nil {
		return v
	}
	return self.Secondary.Peek(sel, peekId)
}
