package browse
import (
	"strings"
)

type Path struct {
	Segments []*PathSegment
	URL string
	query string
}

type PathSegment struct {
	Path *Path
	Index int
	Ident string
	Keys []string
	Key []Value
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
	p.Segments = make([]*PathSegment, len(segments))
	for i, segment := range segments {
		p.Segments[i] = &PathSegment{Path:p, Index:i}
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

func (p *Path) LastSegment() *PathSegment {
	return p.Segments[len(p.Segments) - 1]
}

