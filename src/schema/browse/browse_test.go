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
	json := JsonReader{strings.NewReader(config)}
	var actualBuff bytes.Buffer
	outJson := NewJsonWriter(&actualBuff)
	out, _ := outJson.GetSelector()
	if root, err := json.GetSelector(module, false); err != nil {
		t.Error(err)
	} else {
		var err error
		if err = Insert(root, out, NewExhaustiveController()); err != nil {
			t.Error(err)
		}
		t.Log(string(actualBuff.Bytes()))
	}
}

func TestWalkYang(t *testing.T) {
	module := LoadSampleModule(t)
	var actualBuff bytes.Buffer
	outJson := NewJsonWriter(&actualBuff)
	out, _ := outJson.GetSelector()
	browser := NewSchemaBrowser(module, true)
	if root, err := browser.RootSelector(); err != nil {
		t.Error(err)
	} else {
		if err = Insert(root, out, NewExhaustiveController()); err != nil {
			t.Error(err)
		}
		t.Log(string(actualBuff.Bytes()))
//		actualBuff.Reset()
//
//		var p *Path
//		p, _ = NewPath("module/definitions=game?depth=2")
//		if err = InsertIntoPath(root, out, p); err != nil {
//			t.Error(err)
//		}
//		t.Log(string(actualBuff.Bytes()))
	}
}

