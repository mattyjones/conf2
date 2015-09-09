package yang
import (
	"testing"
	"io/ioutil"
	"container/list"
	"fmt"
)

func TestSimpleExample(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/simple.yang")

	if err != nil {
		t.Fatalf(err.Error())
	}
	l := lex(string(data), nil)
	tokens := list.New()
	for  {
		token, err := l.nextToken()
		if err != nil {
			t.Errorf(err.Error())
		}
		if token.typ == ParseEof {
			break
		}
		tokens.PushBack(token)
	}
	if tokens.Len() != 281 {
		for e := tokens.Front() ; e != nil; e  = e.Next() {
			fmt.Println(e.Value)
		}
		LogTokens(l)
		t.Fatalf("wrong num tokens %d", tokens.Len())
	}
}


func TestStoneLex(t *testing.T) {
	stone, err := ioutil.ReadFile("testdata/romancing-the-stone.yang")
	if err != nil {
		t.Errorf("could not load file %s", err)
	}
	lex(string(stone), nil)
}