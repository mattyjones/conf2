package comm

import (
	"C"
	"bufio"
	"bytes"
	"errors"
	"schema"
	"data"
)

type Writer struct {
	buffer bytes.Buffer
	out    *bufio.Writer
	temp   [4]byte
	pos    int
}

func NewWriter() *Writer {
	w := &Writer{}
	w.out = bufio.NewWriter(&w.buffer)
	return w
}

func (c *Writer) Data() (message []byte) {
	c.out.Flush()
	message = c.buffer.Bytes()
	c.buffer.Reset()
	return
}

func (c *Writer) WriteValues(vals []*data.Value) (err error) {
	c.WriteInt(len(vals))
	for _, val := range vals {
		if err = c.WriteValue(val); err != nil {
			return
		}
	}
	return
}

// Format of value message
// ===================
//  format_code size:4
//
//  If format_code implies fixed len data (e.g. boolean, int, ...)
//     [data]     size: implied by format type
//
//  If format_code implies variable len data (e.g. string, list of ints, ...)
//     data_len   size 4
//     [data]
//
//  [data] format
// ==================
//   If format_code implies list
//     list_len    size 4
//     [data]      size = data_len - 4
//
//  If format code implies non-list variable, then each data type is formatted independantly
//  and detailed in each put method
//
func (c *Writer) WriteValue(val *data.Value) (err error) {
	c.WriteInt(int(val.Type.Format))
	// Implied Fix-Length data types
	switch val.Type.Format {
	case schema.FMT_INT32, schema.FMT_ENUMERATION:
		c.WriteInt(val.Int)
	case schema.FMT_BOOLEAN:
		c.WriteBool(val.Bool)
	case schema.FMT_STRING:
		c.WriteString(val.Str)
	case schema.FMT_STRING_LIST:
		c.WriteInt(len(val.Strlist))
		for _, s := range val.Strlist {
			c.WriteString(s)
		}
	case schema.FMT_INT32_LIST, schema.FMT_ENUMERATION_LIST:
		c.WriteInt(len(val.Intlist))
		for _, i := range val.Intlist {
			c.WriteInt(i)
		}
	case schema.FMT_BOOLEAN_LIST:
		c.WriteInt(len(val.Boollist))
		for _, b := range val.Boollist {
			c.WriteBool(b)
		}
	default:
		return errors.New("Unsupported type")
	}

	return
}

func (c *Writer) WriteInt(i int) {
	// x86 is little endian - TODO: detect and support others
	c.temp[0] = byte(i)
	c.temp[1] = byte(i >> 8)
	c.temp[2] = byte(i >> 16)
	c.temp[3] = byte(i >> 24)
	c.out.Write(c.temp[:4])
}

func (c *Writer) WriteShort(i int16) {
	// x86 is little endian - TODO: detect and support others
	c.temp[0] = byte(i)
	c.temp[1] = byte(i >> 8)
	c.out.Write(c.temp[:2])
}

func (c *Writer) WriteBool(b bool) {
	// TODO: Performance - could make this smaller using bit field
	if b {
		c.out.WriteByte(TRUE_BYTE)
	} else {
		c.out.WriteByte(FALSE_BYTE)
	}
}

// string format in UTF-8**
// ======================
//   byte, byte, NULL
//
// ** This implies strings cannot have NULL in byte array which is legal in most
//   computer languages except C.  If we'd like to change this, there shouldn't be
//   too many assumptions built around this format
func (c *Writer) WriteString(s string) (err error) {
	if _, err = c.out.WriteString(s); err == nil {
		err = c.out.WriteByte(CSTR_TERM)
	}
	return
}
