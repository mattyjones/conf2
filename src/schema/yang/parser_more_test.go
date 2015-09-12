package yang

import (
	"testing"
	"io/ioutil"
)

func TestSimpleParse(t *testing.T) {
	//yyDebug = 4
	yang, err := ioutil.ReadFile("testdata/simple.yang")
	if err != nil {
		t.Error(err)
	}
	l := lex(string(yang[:]), nil)
	err_code := yyParse(l)
	if err_code != 0 {
		t.Error(l.lastError)
	}
}

func TestStoneParse(t *testing.T) {
	//yyDebug = 4
	yang, err := ioutil.ReadFile("testdata/romancing-the-stone.yang")
	if err != nil {
		t.Error(err)
	}
	l := lex(string(yang[:]), nil)
	err_code := yyParse(l)
	if err_code != 0 {
		t.Error(l.lastError)
	}
}
