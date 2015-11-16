package schema

import "fmt"

type MetaPath struct {
	ParentPath *MetaPath
	Meta       Meta
	Key        string
}

func (p *MetaPath) Parent() MetaList {
	// we know it's a list otherwise it couldn't have a child
	if p.ParentPath == nil {
		return nil
	}
	return p.ParentPath.Meta.(MetaList)
}

func (p *MetaPath) Path() string {
	if p.ParentPath != nil {
		return p.ParentPath.Position()
	}
	if p.Meta == nil {
		return ""
	}
	return p.Meta.GetIdent()
}

func (p *MetaPath) Position() string {
	var s string
	if p.ParentPath == nil {
		if p.Meta == nil {
			return ""
		}
	} else {
		s = p.ParentPath.Position()
	}
	if len(p.Key) > 0 {
		s = fmt.Sprint(s, "=", p.Key)
	}
	if p.Meta == nil {
		return s
	}
	if len(s) == 0 {
		return p.Meta.GetIdent()
	}
	return fmt.Sprint(s, "/", p.Meta.GetIdent())
}

func (p *MetaPath) Root() (root *MetaPath) {
	root = p
	for root.ParentPath != nil {
		root = root.ParentPath
	}
	return
}
