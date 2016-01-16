package data

import (
	"strings"
	"bytes"
	"schema"
)

// Immutable otherwise children paths become illegal if parent state changes
type Path struct {
	meta   schema.MetaList
	key    []*Value
	params map[string][]string
	parent *Path
}

func NewRootPath(meta schema.MetaList, params map[string][]string) *Path {
	return &Path{meta:meta, params:params}
}

func NewListItemPath(parent *Path, meta *schema.List, key []*Value) *Path {
	return &Path{parent: parent, meta:meta, key: key}
}

func (path *Path) SetKey(key []*Value) *Path {
	return &Path{parent: path.parent, meta:path.meta, key: key}
}

func NewContainerPath(parent *Path, meta schema.MetaList) *Path {
	return &Path{parent: parent, meta:meta}
}

func (path *Path) Parent() *Path {
	return path.parent
}

func (path *Path) MetaParent() schema.Path {
	if path.parent == nil {
		// subtle difference returning nil and interface reference to nil struct.
		// See http://stackoverflow.com/questions/13476349/check-for-nil-and-nil-interface-in-go
		// by rights in go, all callers should check for interface check for nil and nil interface
		// so this hack some-what contributes to the bad practice of not doing so.
		return nil
	}
	return path.parent
}

func (path *Path) Meta() schema.MetaList {
	return path.meta
}

func (path *Path) Key() []*Value {
	return path.key
}

func (path *Path) Params() map[string][]string {
	p := path
	for p != nil {
		if p.params != nil {
			return p.params
		}
		p = p.parent
	}
	return nil
}

func (seg *Path) String() string {
	strs := make([]string, seg.Len())
	p := seg
	var b bytes.Buffer
	for i := len(strs) - 1; i >= 0; i-- {
		b.Reset()
		p.toBuffer(&b)
		strs[i] = b.String()
		p = p.parent
	}
	return strings.Join(strs, "/")
}

func (seg *Path) toBuffer(b *bytes.Buffer) {
	if seg.meta == nil {
		return
	}
	if b.Len() > 0 {
		b.WriteRune('/')
	}
	b.WriteString(seg.meta.GetIdent())
	if len(seg.key) > 0 {
		b.WriteRune('=')
		for i, k := range seg.key {
			if i != 0 {
				b.WriteRune(',')
			}
			b.WriteString(k.String())
		}
	}
}

func (a *Path) Equal(b *Path) bool {
	if a.Len() != b.Len() {
		return false
	}
	sa := a
	sb := b
	// work up as comparing children are most likely to lead to differences faster
	for sa != nil {
		if ! sa.equalSegment(sb) {
			return false
		}
		sa = sa.parent
		sb = sb.parent
	}
	return true
}

func (path *Path) Len() (len int) {
	p := path
	for p != nil {
		len++
		p = p.parent
	}
	return
}

func (a *Path) equalSegment(b *Path) bool {
	if a.meta == nil {
		if b.meta != nil {
			return false
		}
		if a.meta.GetIdent() != b.meta.GetIdent() {
			return false
		}
	}
	if len(a.key) != len(b.key) {
		return false
	}
	for i, k := range a.key {
		if ! k.Equal(b.key[i]) {
			return false
		}
	}
	return true
}

func (path *Path) Segments() []*Path {
	segs := make([]*Path, path.Len())
	p := path
	for i := len(segs) - 1; i >= 0; i-- {
		segs[i] = p
		p = p.parent
	}
	return segs
}

