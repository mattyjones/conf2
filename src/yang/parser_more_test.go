package yang

import (
	"testing"
	"io/ioutil"
)

func TestSimpleParse(t *testing.T) {
	//yyDebug = 4
	yang, err := ioutil.ReadFile("test_data/simple.test.yang")
	if err != nil {
		t.Error(err)
	}
	l := lex(string(yang[:]))
	err_code := yyParse(l)
	if err_code != 0 {
		t.Fail()
	}
}

func TestStoneParse(t *testing.T) {
	//yyDebug = 4
	yang, err := ioutil.ReadFile("test_data/romancing-the-stone.yang")
	if err != nil {
		t.Error(err)
	}
	l := lex(string(yang[:]))
	err_code := yyParse(l)
	if err_code != 0 {
		t.Fail()
	}
}
