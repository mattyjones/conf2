package browse

import (
	"testing"
	"yang"
	"strings"
	"bytes"
)

func LoadSampleModule(t *testing.T) (*yang.Module) {
	f := &yang.FileStreamSource{Root:"../testdata"}
	m, err:= yang.LoadModule(f, "romancing-the-stone.yang")
	if err != nil {
		t.Error(err.Error())
	}
	return m
}

func LoadYangModule(t *testing.T) (*yang.Module) {
	f := &yang.FileStreamSource{Root:"../../../etc"}
	m, err:= yang.LoadModule(f, "yang-1.0.yang")
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
	if root, err := json.GetSelector(module); err != nil {
		t.Error(err)
	} else {
		var err error
		if err = Insert(root, out); err != nil {
			t.Error(err)
		}
		t.Log(string(actualBuff.Bytes()))
	}
}

func TestWalkYang(t *testing.T) {
	module := LoadSampleModule(t)
	yang := LoadYangModule(t)
	var actualBuff bytes.Buffer
	outJson := NewJsonWriter(&actualBuff)
	out, _ := outJson.GetSelector()
	browser := YangBrowser{Meta:yang, Module:module}
	if root, err := browser.RootSelector(); err != nil {
		t.Error(err)
	} else {
		if err = Insert(root, out); err != nil {
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

