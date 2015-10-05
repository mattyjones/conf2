package browse

import (
	"testing"
	"schema"
	"schema/yang"
	"strings"
	"bytes"
)

func LoadSampleModule(t *testing.T) (*schema.Module) {
	m, err:= yang.LoadModule(schema.NewCwdSource(), "../testdata/romancing-the-stone.yang")
	if err != nil {
		t.Error(err.Error())
	}
	return m
}

func TestWalkJson(t *testing.T) {
	config := `{
	"game" : {
		"base-radius" : 14,
		"teams" : [{
		  "team" : {
		  	"color" : "red"
		  }
		}]
	}
}`
	module := LoadSampleModule(t)
	in := &JsonReader{In:strings.NewReader(config), Meta:module}
	var err error
	var actualBuff bytes.Buffer
	out := NewJsonFragmentWriter(&actualBuff)
	if err = Upsert(NewPath(""), in, out); err != nil {
		t.Error(err)
	}
	t.Log(string(actualBuff.Bytes()))
}

func TestWalkYang(t *testing.T) {
	var err error
	module := LoadSampleModule(t)
	var actualBuff bytes.Buffer
	out := NewJsonFragmentWriter(&actualBuff)
	browser := NewSchemaBrowser(module, true)
	if err = Upsert(NewPath(""), browser, out); err != nil {
		t.Error(err)
	} else {
		t.Log(string(actualBuff.Bytes()))
	}
//		actualBuff.Reset()
//
//		var p *Path
//		p, _ = NewPath("module/definitions=game?depth=2")
//		if err = InsertIntoPath(root, out, p); err != nil {
//			t.Error(err)
//		}
//		t.Log(string(actualBuff.Bytes()))
}

