package browse

import (
	"strings"
)

type Path struct {
	Segments []*PathSegment
	URL      string
	Query    string
}

type PathSegment struct {
	Path  *Path
	Index int
	Ident string
	Keys  []string
}

func NewPath(path string) (p *Path) {
	var err error
	if p, err = ParsePath(path); err != nil {
		if err != nil {
			panic(err.Error())
		}
	}
	return p
}

func ParsePath(path string) (p *Path, err error) {
	p = &Path{}

	if path == "" || path == "/" {
		return
	}
	// UNSUPPORTED OPINION: little reason not to just fix absolute positioned url
	if path[0] == '/' {
		path = path[1:]
	}

	qmark := strings.Index(path, "?")
	if qmark >= 0 {
		p.URL = path[:qmark]
		p.SetQuery(path[qmark+1:])
	} else {
		p.URL = path
	}
	segments := strings.Split(p.URL, "/")
	p.Segments = make([]*PathSegment, len(segments))
	for i, segment := range segments {
		p.Segments[i] = &PathSegment{Path: p, Index: i}
		p.Segments[i].parseSegment(segment)
	}
	return
}

func (ps *PathSegment) parseSegment(segment string) {
	equalsMark := strings.Index(segment, "=")
	if equalsMark >= 0 {
		ps.Ident = segment[:equalsMark]
		ps.Keys = strings.Split(segment[equalsMark+1:], ",")
	} else {
		ps.Ident = segment
	}
}

func (p *Path) LastSegment() *PathSegment {
	if len(p.Segments) == 0 {
		return nil
	}
	return p.Segments[len(p.Segments)-1]
}

func (p *Path) SetQuery(query string) {
	p.Query = query
}
