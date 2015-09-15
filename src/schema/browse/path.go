package browse
import (
	"strings"
	"strconv"
	"schema"
)

type Path struct {
	Segments []PathSegment
	URL string
	query string
}

type PathSegment struct {
	Path *Path
	Index int
	Ident string
	Keys []string
}

func NewPath(path string) (p *Path, err error) {
	p = &Path{}

	if path == "" {
		return
	}

	qmark := strings.Index(path, "?")
	if qmark >= 0 {
		p.URL = path[:qmark]
		p.query = path[qmark + 1:]
	} else {
		p.URL = path
	}
	segments := strings.Split(p.URL, "/")
	p.Segments = make([]PathSegment, len(segments))
	for i, segment := range segments {
		p.Segments[i] = PathSegment{Path:p, Index:i}
		p.Segments[i].parseSegment(segment)
	}
	return
}

func (ps *PathSegment) parseSegment(segment string) {
	equalsMark := strings.Index(segment, "=")
	if equalsMark >= 0 {
		ps.Ident = segment[:equalsMark]
		ps.Keys = strings.Split(segment[equalsMark + 1:], ",")
	} else {
		ps.Ident = segment
	}
}

func (p *Path) FindTargetController() *FindTargetController {
	return &FindTargetController{path:p}
}

func (p *Path) WalkTargetController() (WalkController, error) {
	if len(p.query) > 0 {
		return NewWalkTargetController(p.query)
	}
	return NewExhaustiveController(), nil
}

type WalkTargetController struct {
	MaxDepth int
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
	return s.Next(NO_KEYS, first)
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
		keysAsStrings := n.path.Segments[level - 1].Keys
		list, isList := s.WalkState().Meta.(*schema.List)
		if !isList {

		}
		keys, err := CoerseKeys(list, keysAsStrings)
		if err != nil {
			return false, err
		}
		return s.Next(keys, first)
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
	if level >= len(n.path.Segments) {
		n.target = s
		return schema.EmptyInterator(0)
	}
	position := schema.FindByIdent2(s.WalkState().Meta, n.path.Segments[level].Ident)
	return &schema.SingletonIterator{Meta:position}
}
