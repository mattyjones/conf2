package yang
import "testing"

func TestDefList(t *testing.T) {
	g1 := &Grouping{Ident:"G1"}
	g2 := &Grouping{Ident:"G2"}
	c := DefContainer{}
	c.AddDef(g1)
	c.AddDef(g2)
	if c.FirstDef != g1 {
		t.Error("g1 is first child of container")
	}
	if c.LastDef != g2 {
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

func TestDefProxy(t *testing.T) {
	g1 := &Grouping{Ident:"G1"}
	g1a := &Leaf{Ident:"G1A"}
	g1.AddDef(g1a)
	u1 := &Uses{Ident:"G1"}
	groupings := DefContainer{}
	groupings.AddDef(g1)
	u1.grouping = g1
	i := u1.ResolveProxy()
	nextDef := i.NextDef()
	if nextDef == nil {
		t.Error("resolved proxy is nil")
	} else if nextDef != g1a {
		t.Error("expected G1A and got ", nextDef)
	}

	uparent := DefContainer{}
	uparent.AddDef(u1)
	i2 := NewDefListIterator(&uparent, true)
	nextResolvedDef := i2.NextDef()
	if nextResolvedDef != g1a {
		t.Error("resolved in iterator didn't work")
	} else {
		t.Log("AOK")
	}

}
