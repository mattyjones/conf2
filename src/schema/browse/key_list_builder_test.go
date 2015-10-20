package browse

import (
	"testing"
	"strings"
)

func TestKeyListBuilder(t *testing.T) {
	b := NewKeyListBuilder("a/b/c")
	tests := []struct{
		key string
		iskey bool
	} {
		{"a/b/c", false},
		{"no/way", false},
		{"a/b/cc=z/x", false},
		{"a/b/c=a", true},
		{"a/b/c=b/a/b", true},
		{"a/b/c=b/x", true},
		{"a/b/c=c", true},
	}
	for _, test := range tests {
		if b.ParseKey(test.key) != test.iskey {
			t.Errorf("FAIL: %s key? : expected %v", test.key, test.iskey)
		}
	}
	actual := strings.Join(b.List(), "|")
	expected := "a|b|c"
	if actual != expected {
		t.Errorf("\nExpected:%s\n  Actual:%s", expected, actual)
	}
}