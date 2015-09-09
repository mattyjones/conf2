package browse

import (
	"testing"
	"fmt"
	"yang"
)

func TestNullPath(t  *testing.T) {
	p, _ := NewPath("")
	if len(p.Segments) > 0 {
		t.Error("expected zero segments")
	}
}


func TestPathSegment(t *testing.T) {
	tests := []struct {
		in string
		expected []string
	}{
		{"a/b", []string{"a", "b"}},
		{"a/b=c", []string{"a", "b"}},
		{"a/b=c,q,aaa/d", []string{"a", "b", "d"}},
		{"a/b?foo=1", []string{"a", "b"}},
	}
	for _, test := range tests {
		p, _ := NewPath(test.in)
		if len(test.expected) != len(p.Segments) {
			t.Error("wrong number of expected segments for", test.in)
		}
		for i, e := range test.expected {
			if e != p.Segments[i].Ident {
				msg := "expected to find \"%s\" as segment number %d in \"%s\" but got \"%s\""
				t.Error(fmt.Sprintf(msg, e, i, test.in, p.Segments[i].Ident))
			}
		}
	}
}

func TestPathSegmentKeys(t *testing.T) {
	none := []string{}
	tests := []struct {
		in string
		expected [][]string
	}{
		{"a/b/c", 			[][]string{none, none, none}},
		{"a/b=c/d", 		[][]string{none, []string{"c"}, none}},
		{"a=c,q,aaa/b/d", 	[][]string{[]string{"c", "q", "aaa"}, none, none}},
		{"a/b/d=x", 		[][]string{none, none, []string{"x"}}},
	}
	for _, test := range tests {
		p, _ := NewPath(test.in)
		if len(test.expected) != len(p.Segments) {
			t.Error("wrong number of expected segments for", test.in)
		}
		for i, segment := range test.expected {
			for j, key := range segment {
				if p.Segments[i].Keys[j] != key {
					desc := fmt.Sprintf("\"%s\" segment \"%s\" - expected \"%s\" - got \"%s\"",
						test.in, p.Segments[i].Ident, key, p.Segments[i].Keys[j])
					t.Error(desc)
				}
			}
		}
	}
}

func TestPathParams(t *testing.T) {
	tests := []struct {
		in string
		expected int
	}{
		{"", 				32},
		{"depth=1", 		1},
		{"depth=99", 		99},
	}
	for _, test := range tests {
		p, _ := NewWalkTargetController(test.in)
		if p.MaxDepth != test.expected {
			desc := fmt.Sprintf("\"%s\" - expected depth \"%d\" - got \"%d\"",
				test.in, test.expected, p.MaxDepth)
			t.Error(desc)
		}
	}
}

func TestPathControllerContainerIterator(t *testing.T) {
	var p *Path
	tests := []struct {
		in string
		expected int
	}{
		{"a", 					1},
		{"a?", 					1},
		{"a/b?depth=1",			2},
		{"a/b=key",     		2},
		{"a=key/b",		        2},
	}
	selection := &Selection{}
	selection.Meta = &yang.Container{Ident:"root"}
	selection.Meta.AddMeta(&yang.Container{Ident:"a"})
	selection.Meta.AddMeta(&yang.Container{Ident:"b"})
	for _, test := range tests {
		p ,_ = NewPath(test.in)
		rc := newPathController(p)

		for i := 0; i < test.expected; i++ {
			if ! rc.ContainerIterator(selection, i).HasNextMeta() {
				t.Error(test.in, "- unexpectedly max-depth'ed at level", i)
			}
		}
		if rc.ContainerIterator(selection, test.expected).HasNextMeta() {
			t.Error(test.in, "- expected to max-depth but didn't")
		} else if rc.target != selection {
			t.Error(test.in, "- target incorrect")
		}
	}
}

func TestPathControllerListIterator(t *testing.T) {
	var p *Path
	tests := []struct {
		in string
		expected int
		last bool
	}{
		{"a", 					1, false},
		{"a?", 					1, false},
		{"a/b?depth=1",			2, false},
		{"a/b=key",     		2, true},
		{"a=key/b",		        2, false},
	}
	selection := &Selection{}
	selection.Iterate = func([]string, bool) (bool, error) {
		return true, nil
	}
	var more bool
	for _, test := range tests {
		p ,_ = NewPath(test.in)
		rc := newPathController(p)

		for i := 1; i < test.expected; i++ {
			more, _ = rc.ListIterator(selection, i, true)
			if ! more {
				t.Error(test.in, "- unexpectedly max-depth'ed at level", i)
			}
			more, _ = rc.ListIterator(selection, i, false)
			if more {
				t.Error(test.in, "- unexpectedly found 2nd item", i)
			}
		}
		more, _ = rc.ListIterator(selection, test.expected, true)
		if more != test.last {
			t.Error(test.in, "- expected to max-depth but didn't")
		} else {
			more, _ = rc.ListIterator(selection, test.expected, false)
		    if rc.target != selection {
				t.Error(test.in, "- target incorrect")
			}
			more, _ = rc.ListIterator(selection, test.expected + 1, true)
			if more {
				t.Error(test.in, "- expected to max-depth but didn't")
			}
		}
	}
}


