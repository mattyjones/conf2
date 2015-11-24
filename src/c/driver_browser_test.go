package c

import (
	"testing"
	"schema/browse"
	"schema"
)

func TestLeafListInt(t *testing.T) {
	v := &browse.Value{Type:&schema.TYPE_INT32,Intlist:[]int{ 1, 2, 3 }}
	if c_val, err := leafListValue(v); err != nil {
		t.Error(err)
	} else {
		if c_val.datalen != 12 {
			t.Fail()
		}
		if c_val.islist == 0 {
			t.Fail()
		}
	}
}

func TestLeafListStr(t *testing.T) {
	v := &browse.Value{Type:&schema.TYPE_STRING,Strlist:[]string{ "a", "bb", "ccc" }}
	if c_val, err := leafListValue(v); err != nil {
		t.Error(err)
	} else {
		// 6 runes + 3 null terminators
		if c_val.datalen != 9 {
			t.Fail()
		}
	}
}
