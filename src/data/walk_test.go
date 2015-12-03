package data

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
	m := LoadSampleModule(t)
	rdr := NewJsonReader(strings.NewReader(config)).Node()
	var actualBuff bytes.Buffer
	wtr := NewJsonWriter(&actualBuff).Node()
	if err := NodeToNode(rdr, wtr, m).Upsert(); err != nil {
		t.Error(err)
	}
	t.Log(string(actualBuff.Bytes()))
}

func TestWalkYang(t *testing.T) {
	var err error
	module := LoadSampleModule(t)
	var actualBuff bytes.Buffer
	wtr := NewJsonWriter(&actualBuff).Node()
	browser := NewSchemaData(module, true)
	if err = NodeToNode(browser.Node(), wtr, browser.Schema()).Upsert(); err != nil {
		t.Error(err)
	} else {
		t.Log(string(actualBuff.Bytes()))
	}
}
