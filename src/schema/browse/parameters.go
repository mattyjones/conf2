package browse
import (
	"schema"
)

type Parameters struct {
	Given map[string]*Value
	Schema schema.MetaList
}

func NewParameters(schema schema.MetaList) *Parameters {
	return &Parameters{
		Given : make(map[string]*Value, 5),
		Schema : schema,
	}
}

func (p *Parameters) Value(ident string) *Value {
	if v, found := p.Given[ident]; found {
		return v
	} else {
		meta := schema.FindByIdent2(p.Schema, ident)
		if prop, isProp := meta.(schema.HasDataType); isProp {
			t := prop.GetDataType()
			v := &Value{Type:t}
			if len(t.Default) > 0 {
				v.CoerseStrValue(t.Default)
				return v
			}
		}
	}
	return nil
}

func (p *Parameters) Collect() (Selection, error) {
	s := &MySelection{}
	s.OnWrite = func(state *WalkState, meta schema.Meta, op Operation, v *Value) (err error) {
		switch op {
		case UPDATE_VALUE:
			p.Given[meta.GetIdent()] = v
		}
		return nil
	}
	return s, nil
}

func (p *Parameters) Configure(obj interface{}) (err error) {
	i := schema.NewMetaListIterator(p.Schema, true)
	for i.HasNextMeta() {
		m := i.NextMeta()
		if t, hasType := m.(schema.HasDataType); hasType {
			v := p.Value(m.GetIdent())
			if v != nil {
				if err = WriteField(t, obj, v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
