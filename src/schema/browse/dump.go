package browse
import (
	"io"
	"bufio"
	"fmt"
	"strconv"
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
	s.OnSelect = func() (child Selection, err error) {
		return d.Enter(level + 1)
	}
	s.OnWrite = func(op Operation, v *Value) (err error) {
		d.dumpEditOp(s, op, level)
		d.dumpValue(v, level)
		return
	}
	s.OnNext = func(keys []interface{}, first bool) (hasMore bool, err error) {
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
	if v.IsList {
		switch v.Type.Ident {
		case "string":
			s = fmt.Sprintf("%v", v.Strlist)
		case "int32":
			s = fmt.Sprintf("%v", v.Intlist)
		case "boolean":
			s = fmt.Sprintf("%v", v.Boollist)
		}
	} else {
		switch v.Type.Ident {
		case "string":
			s = v.Str
		case "int32":
			s = strconv.Itoa(v.Int)
		case "boolean":
			if v.Bool {
				s = "true"
			} else {
				s = "false"
			}
		}
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