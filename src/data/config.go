package data
import (
	"schema"
	"errors"
)

func Config(operational Node, config Node) Node {
	n := &MyNode{
		Label: "Persist",
		OnSelect: func(sel *Selection, r ContainerRequest) (Node, error) {
			operChild, err := operational.Select(sel, r)
			if err != nil || operChild == nil {
				return nil, err
			}
			if ! sel.IsConfig(r.Meta) {
				return operChild, nil
			}
			configChild, storeErr := config.Select(sel, r)
			if storeErr != nil {
				return nil, err
			}
			if configChild == nil && ! r.New {
				r.New = true
				configChild, storeErr = config.Select(sel, r)
				if storeErr != nil {
					return nil, err
				}
				if configChild == nil {
					return nil, errors.New("Could not create storage node for " + sel.String())
				}
			}
			return Config(operChild, configChild), nil
		},
		OnNext: func(sel *Selection, r ListRequest) (next Node, key []*Value, err error) {
			var operChild, configChild Node
			if operChild, key, err = operational.Next(sel, r); err != nil || operChild == nil {
				return
			}
			if ! sel.IsConfig(r.Meta) {
				return operChild, key, nil
			}
			r.Key = key
			if configChild, _, err = config.Next(sel, r); err != nil {
				return
			}
			if configChild == nil && ! r.New {
				r.New = true
				if configChild, _, err = config.Next(sel, r); err != nil {
					return
				}
				if configChild == nil {
					return nil, nil, errors.New("Could not create storage node for " + sel.String())
				}
			}
			return Config(operChild, configChild), key, nil
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
