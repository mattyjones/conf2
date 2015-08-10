package browse
import (
	"strings"
	"strconv"
	"yang"
)

type Path struct {
	Segments []PathSegment
	URL string
	Depth int
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
		if err = p.parseQuery(path[qmark + 1:]); err != nil {
			return nil, err
		}
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

func (p *Path) parseQuery(q string) (err error) {
	params := strings.Split(q, "&")
	for _, param := range params {
		nameValue := strings.Split(param, "=")
		switch nameValue[0] {
		case "depth":
			if p.Depth, err = strconv.Atoi(nameValue[1]); err != nil {
				return
			}
		}
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

type pathWalkController struct {
	path *Path
	target *Selection
}

func newPathController(p *Path) *pathWalkController {
	return &pathWalkController{path:p}
}

func (n *pathWalkController) ListIterator(s *Selection, level int, first bool) (hasMore bool, err error) {
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
		return s.Iterate(n.path.Segments[level - 1].Keys, first)
	} else {
		return false, nil
	}
}

func (n *pathWalkController) ContainerIterator(s *Selection, level int) yang.MetaIterator {
	if level >= len(n.path.Segments) {
		n.target = s
		return yang.EmptyInterator(0)
	}
	position := yang.FindByIdent2(s.Meta, n.path.Segments[level].Ident)
	return &yang.SingletonIterator{Meta:position}
}
