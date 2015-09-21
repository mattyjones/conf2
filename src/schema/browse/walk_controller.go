package browse

import (
	"schema"
	"strings"
	"strconv"
	"fmt"
)

func (p *Path) FindTargetController() *FindTargetController {
	return &FindTargetController{path:p}
}

type WalkTargetController struct {
	MaxDepth int
	InitialKey []Value
}

func (p *Path) WalkTargetController() (WalkController, error) {
	wtc, err := NewWalkTargetController(p.query)
	if err != nil {
		return nil, err
	}
	if len(p.Segments) > 0 {
		wtc.InitialKey = p.LastSegment().Key
	}
	return wtc, err
}

func NewWalkTargetController(query string) (*WalkTargetController, error) {
	c := &WalkTargetController{MaxDepth:32}
	c.parseQuery(query)
	return c, nil
}

func (p *WalkTargetController) parseQuery(q string) (err error) {
	params := strings.Split(q, "&")
	for _, param := range params {
		nameValue := strings.Split(param, "=")
		switch nameValue[0] {
		case "depth":
			if p.MaxDepth, err = strconv.Atoi(nameValue[1]); err != nil {
				return
			}
		}
	}
	return
}

func (p *WalkTargetController) CloseSelection(s Selection) error {
	return schema.CloseResource(s)
}

func NewExhaustiveController() WalkController {
	return &WalkTargetController{MaxDepth:32}
}

func (e *WalkTargetController) ListIterator(s Selection, level int, first bool) (hasMore bool, err error) {
	if level >= e.MaxDepth {
		return false, nil
	}
	var key []Value
	if level == 0 {
		key = e.InitialKey
	} else {
		key = NO_KEYS
	}
	return s.Next(key, first)
}

func (e *WalkTargetController) ContainerIterator(s Selection, level int) schema.MetaIterator {
	if level >= e.MaxDepth {
		return schema.EmptyInterator(0)
	}
	return schema.NewMetaListIterator(s.WalkState().Meta, true)
}

type FindTargetController struct {
	path *Path
	target Selection
	resource schema.Resource
}

func newPathController(p *Path) *FindTargetController {
	return &FindTargetController{path:p}
}


func (n *FindTargetController) ListIterator(s Selection, level int, first bool) (hasMore bool, err error) {
	if level == len(n.path.Segments) {
		if len(n.path.Segments[level - 1].Keys) == 0 {
			n.target = s
			return false, nil
		}
		if !first {
			n.target = s
			return false, nil
		}
	}
	if first && level > 0 && level <= len(n.path.Segments) {
		segment := n.path.Segments[level - 1]
		keysAsStrings := segment.Keys
		list, isList := s.WalkState().Meta.(*schema.List)
		if !isList {
			return false, &browseError{Msg:fmt.Sprintf("Key \"%s\" specified when not a list", keysAsStrings)}
		}
		segment.Key, err = CoerseKeys(list, keysAsStrings)
		if err != nil {
			return false, err
		}
		return s.Next(segment.Key, first)
	} else {
		return false, nil
	}
}

func (p *FindTargetController) CloseSelection(s Selection) error {
	if s != p.target {
		return schema.CloseResource(s)
	}
	return nil
}

func (n *FindTargetController) setTarget(s *MySelection) {
	n.target = s
	// we take ownership of resource so it's not released until target is used
	n.resource = s.Resource
	s.Resource = nil
}

func (n *FindTargetController) ContainerIterator(s Selection, level int) schema.MetaIterator {
	//n.path.Key = nil
	if level >= len(n.path.Segments) {
		n.target = s
		return schema.EmptyInterator(0)
	}
	position := schema.FindByIdentExpandChoices(s.WalkState().Meta, n.path.Segments[level].Ident)
	return &schema.SingletonIterator{Meta:position}
}
