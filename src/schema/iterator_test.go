package schema
import "testing"

func TestEmptyIterator(t *testing.T) {
	i := EmptyInterator(0)
	if i.HasNextMeta() {
		t.Fail()
	}
}

func TestSingletonIterator(t *testing.T) {
	leaf := &Leaf{Ident:"L"}
	i := &SingletonIterator{leaf}
	if ! i.HasNextMeta() {
		t.Fail()
	}
	if i.NextMeta() != leaf {
		t.Fail()
	}
	if i.HasNextMeta() {
		t.Fail()
	}
}