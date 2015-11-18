package browse

import (
	"bytes"
	"schema"
	"schema/yang"
	"strings"
	"testing"
)

func LoadSampleModule(t *testing.T) *schema.Module {
	m, err := yang.LoadModuleFromByteArray([]byte(yang.TestDataRomancingTheStone), nil)
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
  		  "color" : "red",
		  "team" : {
		    "members" : ["joe","mary"]
		  }
		}]
	}
}`
	module := LoadSampleModule(t)
	rdr, err := NewJsonReader(strings.NewReader(config)).Node()
	if err != nil {
		t.Fatal(err)
	}
	in := NewSelection(rdr, module)
	var actualBuff bytes.Buffer
	out := NewSelection(NewJsonWriter(&actualBuff).Container(), module)
	if err = Upsert(in, out); err != nil {
		t.Error(err)
	}
	t.Log(string(actualBuff.Bytes()))
}

func TestWalkYang(t *testing.T) {
	var err error
	module := LoadSampleModule(t)
	var actualBuff bytes.Buffer
	out := NewSelection(NewJsonWriter(&actualBuff).Container(), module)
	browser := NewSchemaData(module, true)
	var in *Selection
	in, err = browser.Selector(NewPath(""))
	if err = Upsert(in, out); err != nil {
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
