package yang
import "testing"

func TestMetaList(t *testing.T) {
	g1 := &Grouping{Ident:"G1"}
	g2 := &Grouping{Ident:"G2"}
	c := MetaContainer{}
	c.AddMeta(g1)
	c.AddMeta(g2)
	if c.FirstMeta != g1 {
		t.Error("g1 is first child of container")
	}
	if c.LastMeta != g2 {
		t.Error("g2 is last child of container")
	}
	if g1.GetParent() != &c {
		t.Error("g1 parent is not container")
	}
	if g2.GetParent() != &c {
		t.Error("g2 parent is not container")
	}
	if g1.Sibling != g2 {
		t.Error("g1 is not linked to g2")
	}
	if g2.Sibling != nil {
		t.Error("g2 sibling should be nil")
	}
}

func TestMetaProxy(t *testing.T) {
	g1 := &Grouping{Ident:"G1"}
	g1a := &Leaf{Ident:"G1A"}
	g1.AddMeta(g1a)
	u1 := &Uses{Ident:"G1"}
	groupings := MetaContainer{}
	groupings.AddMeta(g1)
	u1.grouping = g1
	i := u1.ResolveProxy()
	nextMeta := i.NextMeta()
	if nextMeta == nil {
		t.Error("resolved proxy is nil")
	} else if nextMeta != g1a {
		t.Error("expected G1A and got ", nextMeta)
	}

	uparent := MetaContainer{}
	uparent.AddMeta(u1)
	i2 := NewMetaListIterator(&uparent, true)
	nextResolvedMeta := i2.NextMeta()
	if nextResolvedMeta != g1a {
		t.Error("resolved in iterator didn't work")
	} else {
		t.Log("AOK")
	}

}
