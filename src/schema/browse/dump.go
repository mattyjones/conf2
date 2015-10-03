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
	s.OnSelect = func(state *WalkState, meta schema.MetaList) (child Selection, err error) {
		nest := created
		created = nil
		return nest, nil
	}
	s.OnWrite = func(state *WalkState, meta schema.Meta, op Operation, v *Value) (err error) {
		switch op {
			case CREATE_CHILD, CREATE_LIST, CREATE_LIST_ITEM:
				created, _ = d.Enter(level + 1)
		}
		d.dumpEditOp(state, op, level)
		d.dumpValue(v, level)
		return
	}
	s.OnNext = func(state *WalkState, meta *schema.List, keys []*Value, first bool) (hasMore bool, err error) {
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

func (d *Dumper) dumpEditOp(state *WalkState, op Operation, level int) {
	line := fmt.Sprintf("%s%s %s\n", Padding[:level], editOps[op], state.String())
	d.out.WriteString(line)
}