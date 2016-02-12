package data
import (
	"schema"
	"errors"
)

func Config(operational Node, config Node) Node {
	n := &MyNode{
		Label: "Persist",
		OnSelect: func(sel *Selection, meta schema.MetaList, new bool) (Node, error) {
			operChild, err := operational.Select(sel, meta, new)
			if err != nil || operChild == nil {
				return nil, err
			}
			if ! sel.IsConfig(meta) {
				return operChild, nil
			}
			configChild, storeErr := config.Select(sel, meta, new)
			if storeErr != nil {
				return nil, err
			}
			if configChild == nil && ! new {
				configChild, storeErr = config.Select(sel, meta, true)
				if storeErr != nil {
					return nil, err
				}
				if configChild == nil {
					return nil, errors.New("Could not create storage node for " + sel.String())
				}
			}
			return Config(operChild, configChild), nil
		},
		OnNext: func(sel *Selection, meta *schema.List, new bool, key []*Value, first bool) (next Node, err error) {
			operChild, err := operational.Next(sel, meta, new, key, first)
			if err != nil || operChild == nil {
				return nil, err
			}
			if ! sel.IsConfig(meta) {
				return operChild, nil
			}
			configChild, storeErr := config.Next(sel, meta, new, key, first)
			if storeErr != nil {
				return nil, err
			}
			if configChild == nil && ! new {
				configChild, storeErr = config.Next(sel, meta, true, key, first)
				if storeErr != nil {
					return nil, err
				}
				if configChild == nil {
					return nil, errors.New("Could not create storage node for " + sel.String())
				}
			}
			return Config(operChild, configChild), nil
		},
		OnWrite: func(sel *Selection, meta schema.HasDataType, val *Value) error {
			if err := operational.Write(sel, meta, val); err != nil {
				return err
			}
			return config.Write(sel, meta, val)
		},
		OnEvent: func(sel *Selection, e Event) error {
			if err := operational.Event(sel, e); err != nil {
				return err
			}
			if err := config.Event(sel, e); err != nil {
				return err
			}
			return nil
		},
		OnRead : func(sel *Selection, meta schema.HasDataType) (*Value, error) {
			if meta.(schema.HasDetails).Details().Config(sel.Path()) {
				return config.Read(sel, meta)
			}
			return operational.Read(sel, meta)
		},
		OnChoose : operational.Choose,
		OnAction : operational.Action,
		OnPeek: operational.Peek,
	}
	if storeAware, ok := operational.(ChangeAwareNode); ok {
		storeAware.DirectChanges(n)
	}
	return n
}
