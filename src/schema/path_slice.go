package schema

import (
	"strings"
	"errors"
	"bytes"
	"net/url"
)

type PathSlice struct {
	Head   *Path
	Tail   *Path
}

func (slice *PathSlice) Params() map[string][]string {
	if slice.Head != nil {
		return slice.Head.Params()
	}
	return nil
}

func NewPathSlice(path string, meta MetaList) (p *PathSlice) {
	var err error
	if p, err = ParsePath(path, meta); err != nil {
		if err != nil {
			panic(err.Error())
		}
	}
	return p
}

func ParsePath(path string, meta MetaList) (*PathSlice, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	p := NewRootPath(meta, map[string][]string(u.Query()))
	slice := &PathSlice{
		Head: p,
		Tail: p,
	}
	segments := strings.Split(u.EscapedPath(), "/")
	for _, segment := range segments {
		if segment == "" {
			break
		}
		seg := &Path{parent:p}
		equalsMark := strings.Index(segment, "=")
		var ident string
		var keyStrs []string
		if equalsMark >= 0 {
			if ident, err =  url.QueryUnescape(segment[:equalsMark]); err != nil {
				return nil, err
			}
			keyStrs = strings.Split(segment[equalsMark+1:], ",")
			for i, escapedKeystr := range keyStrs {
				if keyStrs[i], err = url.QueryUnescape(escapedKeystr); err != nil {
					return nil, err
				}
			}
		} else {
			if ident, err =  url.QueryUnescape(segment); err != nil {
				return nil, err
			}
		}
		m := FindByIdentExpandChoices(p.meta, ident)
		var notLeaf bool
		if m == nil {
			return nil, errors.New(ident + " not found in " + p.meta.GetIdent())
		}
		if seg.meta, notLeaf = m.(MetaList); ! notLeaf {
			return nil, errors.New("paths cannot contain leaf types:" + ident)
		}
		if len(keyStrs) > 0 {
			if seg.key, err = CoerseKeys(seg.meta.(*List), keyStrs); err != nil {
				return nil, err
			}
		}
		slice.AppendPath(seg)
		p = seg
	}
	return slice, nil
}

func (path *PathSlice) Empty() bool {
	return path.Tail == nil
}

func (path *PathSlice) Equal(bPath *PathSlice) bool {
	if path.Len() != bPath.Len() {
		return false
	}
	a := path.Tail
	b := bPath.Tail
	for a != nil {
		if a.meta != b.meta {
			return false
		}
		if len(a.key) != len(b.key) {
			return false
		}
		for i, k := range a.key {
			if ! k.Equal(b.key[i]) {
				return false
			}
		}
		a = a.parent
		b = b.parent
	}
	return true
}

func (slice *PathSlice) NextAfter(path *Path) (p *Path) {
	if path == slice.Tail {
		return nil
	}
	candidate := slice.Tail
	for candidate != nil {
		if candidate.parent == path {
			return candidate
		}
		candidate = candidate.parent
	}
	return nil
}

func (path *PathSlice) AppendPath(child *Path) {
	if path.Tail == nil {
		path.Head = child
		path.Tail = child
	} else {
		child.parent = path.Tail
		path.Tail = child
	}
}

func (path *PathSlice) Len() (len int) {
	p := path.Tail
	for p != nil {
		len++
		p = p.parent
	}
	return
}

func (path *PathSlice) String() string {
	var b bytes.Buffer
	for _, segment := range path.Segments() {
		segment.toBuffer(&b)
	}
	return b.String()
}

func (path *PathSlice) Segments() []*Path {
	segments := make([]*Path, path.Len())
	p := path.Tail
	for i := len(segments) - 1; i >= 0; i-- {
		segments[i] = p
		p = p.parent
	}
	return segments
}

