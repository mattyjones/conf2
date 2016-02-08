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

func (self *Tee) Select(sel *Selection, meta schema.MetaList, new bool) (Node, error) {
	var err error
	var child Tee
	if child.Primary, err = self.Primary.Select(sel, meta, new); err != nil {
		return nil, err
	}
	if child.Secondary, err = self.Secondary.Select(sel, meta, new); err != nil {
		return nil, err
	}
	if child.Primary != nil && child.Secondary != nil {
		return &child, nil
	}
	return nil, nil
}

func (self *Tee) Find(sel *Selection, path *Path) (err error) {
	if err = self.Primary.Find(sel, path); err == nil {
		err = self.Secondary.Find(sel, path)
	}
	return
}

func (self *Tee)  Next(sel *Selection, meta *schema.List, new bool, key []*Value, isFirst bool) (Node, error) {
	var err error
	var next Tee
	if next.Primary, err = self.Primary.Next(sel, meta, new, key, isFirst); err != nil {
		return nil, err
	}
	if next.Secondary, err = self.Secondary.Next(sel, meta, new, key, isFirst); err != nil {
		return nil, err
	}
	if next.Primary != nil && next.Secondary != nil {
		return &next, nil
	}
	return nil, nil
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
