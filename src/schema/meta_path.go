package schema

import (
	"fmt"
	"bytes"
)

type Path interface {
	fmt.Stringer
	Meta() MetaList
	MetaParent() Path
}

type MetaPath struct {
	parent *MetaPath
	meta MetaList
}

func (p *MetaPath) Meta() MetaList {
	return p.meta
}

func (p *MetaPath) MetaParent() Path {
	return p.parent
}

func (p *MetaPath) String() string {
	var b bytes.Buffer
	p.toBuffer(&b)
	return b.String()
}

func (p *MetaPath) toBuffer(buff *bytes.Buffer) {
	if p.parent != nil {
		p.parent.toBuffer(buff)
		buff.WriteRune('/')
	}
	buff.WriteString(p.meta.GetIdent())
}

