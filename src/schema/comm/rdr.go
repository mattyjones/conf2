package comm

import (
	"errors"
	"schema"
	"schema/browse"
)

type Reader struct {

	// Not using a bytes.Reader and instead reading from []byte directly because
	// doesn't seem to add a lot of value in this case
	data []byte

	pos int
}

func NewReader(data []byte) *Reader {
	return &Reader{data: data}
}

func (c *Reader) ReadValue(typ *schema.DataType) (v *browse.Value, err error) {
	var formatCode int
	if formatCode, err = c.ReadInt(); err != nil {
		return
	}
	format := schema.DataFormat(formatCode)
	v = &browse.Value{Type: typ}
	switch format {
	case schema.FMT_INT32:
		v.Int, err = c.ReadInt()
	case schema.FMT_ENUMERATION:
		var enumId int
		enumId, err = c.ReadInt()
		v.SetEnum(enumId)
	case schema.FMT_BOOLEAN:
		v.Bool, err = c.ReadBool()
	}

	if schema.IsListFormat(format) {
		var listLen int
		if listLen, err = c.ReadInt(); err != nil {
			return
		}
		switch format {
		case schema.FMT_INT32_LIST, schema.FMT_ENUMERATION_LIST:
			v.Intlist = make([]int, listLen)
			for i := 0; i < listLen; i++ {
				if v.Intlist[i], err = c.ReadInt(); err != nil {
					return nil, err
				}
			}
			if format == schema.FMT_ENUMERATION_LIST {
				v.SetEnumList(v.Intlist)
			}
		case schema.FMT_BOOLEAN_LIST:
			v.Boollist = make([]bool, listLen)
			for i := 0; i < listLen; i++ {
				if v.Boollist[i], err = c.ReadBool(); err != nil {
					return nil, err
				}
			}
		case schema.FMT_STRING_LIST:
			v.Strlist = make([]string, listLen)
			for i := 0; i < listLen; i++ {
				if v.Strlist[i], err = c.ReadString(); err != nil {
					return nil, err
				}
			}
		}
	}
	return
}

var EOF = errors.New("Reached end of data buffer before reading expected contents")

func (c *Reader) ReadInt() (v int, err error) {
	if c.pos+4 > len(c.data) {
		return 0, EOF
	}
	v = int(c.data[c.pos])
	v += int(c.data[c.pos+1]) << 8
	v += int(c.data[c.pos+2]) << 16
	v += int(c.data[c.pos+3]) << 24
	c.pos += 4
	return
}

func (c *Reader) ReadString() (v string, err error) {
	start := c.pos
	for ; ; c.pos++ {
		if c.pos >= len(c.data) {
			break
		}
		if c.data[c.pos] == 0 {
			c.pos++
			return string(c.data[start : c.pos-1]), nil
		}
	}
	return "", EOF
}

func (c *Reader) ReadBool() (v bool, err error) {
	if c.pos >= len(c.data) {
		return false, EOF
	}
	c.pos++
	return (c.data[c.pos-1] > 0), nil
}
