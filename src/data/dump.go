package data

import (
	"bufio"
	"fmt"
	"io"
	"schema"
	"strconv"
)

const Padding = "                                                                                       "

type Dumper struct {
	out *bufio.Writer
}

func NewDumper(out io.Writer) *Dumper {
	return &Dumper{
		out: bufio.NewWriter(out),
	}
}

func (self *Dumper) Select(meta schema.MetaList) *Selection {
	return Select(meta, self.enter(0))
}

func (d *Dumper) enter(level int) Node {
	row := 0
	s := &MyNode{}
	s.OnSelect = func(state *Selection, r ContainerRequest) (child Node, err error) {
		if ! r.New {
			return
		}
		return d.enter(level + 1), nil
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, v *Value) (err error) {
		d.dumpValue(v, level)
		return
	}
	s.OnNext = func(state *Selection, r ListRequest) (next Node, key []*Value, err error) {
		if ! r.New {
			return
		}
		d.out.WriteString(fmt.Sprintf("%sITERATE row=%d, first=%v\n", Padding[:level], row, r.First))
		row++
		return s, r.Key, nil
	}
	return s
}

func (d *Dumper) dumpValue(v *Value, level int) {
	if v == nil {
		return
	}
	s := "?"
	t := v.Type.Ident
	switch v.Type.Format() {
	case schema.FMT_STRING:
		s = v.Str
	case schema.FMT_STRING_LIST:
		s = fmt.Sprintf("%v", v.Strlist)
	case schema.FMT_INT32:
		s = strconv.Itoa(v.Int)
	case schema.FMT_INT32_LIST:
		s = fmt.Sprintf("%v", v.Intlist)
	case schema.FMT_BOOLEAN:
		if v.Bool {
			s = "true"
		} else {
			s = "false"
		}
	case schema.FMT_BOOLEAN_LIST:
		s = fmt.Sprintf("%v", v.Boollist)
	}
	line := fmt.Sprintf("%s-> \"%s\" type=%s\n", Padding[:level], s, t)
	d.out.WriteString(line)
}
