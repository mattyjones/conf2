package browse
import (
	"io"
	"bufio"
	"fmt"
	"strconv"
	"schema"
)

var editOps = map[Operation]string {
	CREATE_CHILD : "CREATE_CHILD",
	POST_CREATE_CHILD : "POST_CREATE_CHILD",
	CREATE_LIST : "CREATE_LIST",
	POST_CREATE_LIST : "POST_CREATE_LIST",
	UPDATE_VALUE : "UPDATE_VALUE",
	DELETE_CHILD : "DELETE_CHILD",
	DELETE_LIST : "DELETE_LIST",
	BEGIN_EDIT : "BEGIN_EDIT",
	END_EDIT : "END_EDIT",
	CREATE_LIST_ITEM : "CREATE_LIST_ITEM",
	POST_CREATE_LIST_ITEM : "POST_CREATE_LIST_ITEM",
}
const Padding = "                                                                                       "

type Dumper struct {
	out *bufio.Writer
}

func NewDumper(out io.Writer) *Dumper {
	return &Dumper{
		out:bufio.NewWriter(out),
	}
}

func (d *Dumper) GetSelector() (Selection, error) {
	return d.Enter(0)
}

func (d *Dumper) Enter(level int) (Selection, error) {
	row := 0
	s := &MySelection{}
	var created Selection
	s.OnSelect = func(meta schema.MetaList) (child Selection, err error) {
		nest := created
		created = nil
		return nest, nil
	}
	s.OnWrite = func(meta schema.Meta, op Operation, v *Value) (err error) {
		switch op {
			case CREATE_CHILD, CREATE_LIST, CREATE_LIST_ITEM:
				created, _ = d.Enter(level + 1)
		}
		d.dumpEditOp(s, op, level)
		d.dumpValue(v, level)
		return
	}
	s.OnNext = func(keys []*Value, first bool) (hasMore bool, err error) {
		d.out.WriteString(fmt.Sprintf("%sITERATE row=%d, first=%v\n", Padding[:level], row, first))
		row++
		return false, nil
	}
	return s, nil
}

func (d *Dumper) dumpValue(v *Value, level int) {
	if v == nil {
		return
	}
	s := "?"
	t := v.Type.Ident
	switch v.Type.Format {
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

func (d *Dumper) dumpEditOp(s *MySelection, op Operation, level int) {
	ident := ""
	if s.State.Position != nil {
		ident = s.State.Position.GetIdent()
	}
	line := fmt.Sprintf("%s%s %s\n", Padding[:level], editOps[op], ident)
	d.out.WriteString(line)
}