package browse

import (
	"testing"
	"yang"
	"strings"
	"bytes"
)

func LoadSampleModule(t *testing.T) (*yang.Module) {
	f := &yang.FileDataSource{Root:"../testdata"}
	m, err:= yang.LoadModule(f, "romancing-the-stone.yang")
	if err != nil {
		t.Error(err.Error())
	}
	return m
}

func LoadYangModule(t *testing.T) (*yang.Module) {
	f := &yang.FileDataSource{Root:"../../../etc"}
	m, err:= yang.LoadModule(f, "yang-1.0.yang")
	if err != nil {
		t.Error(err.Error())
	}
	return m
}

func TestReadControllerDepth(t *testing.T) {
	var p *Path
	var rc *readController
	tests := []struct {
		in string
		maxLevel int
		expected int
	}{
		{"a", 					5, 5},
		{"a?", 					5, 5},
		{"a/b?depth=1",			5, 3},
		{"a/b?depth=2", 		10, 4},
	}
	for _, test := range tests {
		p ,_ = NewPath(test.in)
		rc = &readController{path:p, maxLevel:test.maxLevel}

		for i := 0; i < test.expected - 1; i++ {
			if rc.isMaxLevel() {
				t.Error(test.in, "unexpectedly max-depth'ed at level", i)
			}
			rc = rc.recurse()
		}
		if !rc.isMaxLevel() {
			t.Error(test.in, "expected to max-depth but didn't")
		}
	}
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
	json := JsonTransmitter{strings.NewReader(config)}
	var actualBuff bytes.Buffer
	out := NewJsonReceiver(&actualBuff)
	if root, err := json.GetSelector(module); err != nil {
		t.Error(err)
	} else {
		var err error
		if err = Walk(root, nil, out); err != nil {
			t.Error(err)
		}
		out.Flush()
		t.Log(string(actualBuff.Bytes()))
	}
}

func TestWalkYang(t *testing.T) {
	module := LoadSampleModule(t)
	yang := LoadYangModule(t)
	var actualBuff bytes.Buffer
	out := NewJsonReceiver(&actualBuff)
	browser := MetaTransmitter{meta:yang, module:module}
	if root, err := browser.RootSelector(); err != nil {
		t.Error(err)
	} else {
		var p *Path
		p, _ = NewPath("module/definitions=game?depth=2")
		if err = Walk(root, p, out); err != nil {
			t.Error(err)
		}
		out.Flush()
		t.Log(string(actualBuff.Bytes()))
	}
}

