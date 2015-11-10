package schema

import "testing"

func TestMetaPathString(t *testing.T) {
	p1 := MetaPath{Meta: &Container{Ident: "p1"}}
	if p1.String() != "p1" {
		t.Error(p1.String())
	}
	p2 := MetaPath{ParentPath: &p1, Meta: &List{Ident: "p2"}}
	if p2.String() != "p1/p2" {
		t.Error(p2.String())
	}
	p2.Key = "x"
	if p2.String() != "p1=x/p2" {
		t.Error(p2.String())
	}
}
