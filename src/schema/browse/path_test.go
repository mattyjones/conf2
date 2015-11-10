package browse

import (
	"fmt"
	"testing"
)

func TestNullPath(t *testing.T) {
	p, _ := ParsePath("")
	if len(p.Segments) > 0 {
		t.Error("expected zero segments")
	}
}

func TestPathSegment(t *testing.T) {
	tests := []struct {
		in       string
		expected []string
	}{
		{"a/b", []string{"a", "b"}},
		{"a/b=c", []string{"a", "b"}},
		{"a/b=c,q,aaa/d", []string{"a", "b", "d"}},
		{"a/b?foo=1", []string{"a", "b"}},
	}
	for _, test := range tests {
		p, _ := ParsePath(test.in)
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
		in       string
		expected [][]string
	}{
		{"a/b/c", [][]string{none, none, none}},
		{"a/b=c/d", [][]string{none, []string{"c"}, none}},
		{"a=c,q,aaa/b/d", [][]string{[]string{"c", "q", "aaa"}, none, none}},
		{"a/b/d=x", [][]string{none, none, []string{"x"}}},
	}
	for _, test := range tests {
		p, _ := ParsePath(test.in)
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
		in       string
		expected int
	}{
		{"", 32},
		{"depth=1", 1},
		{"depth=99", 99},
	}
	for _, test := range tests {
		p := LimitedWalk(test.in)
		if p.MaxDepth != test.expected {
			desc := fmt.Sprintf("\"%s\" - expected depth \"%d\" - got \"%d\"",
				test.in, test.expected, p.MaxDepth)
			t.Error(desc)
		}
	}
}
