package data
import (
	"testing"
	"fmt"
)

func TestPathParams(t *testing.T) {
	tests := []struct {
		in       map[string][]string
		expected int
	}{
		{
			map[string][]string{},
			32,
		},
		{
			map[string][]string{
				"depth" : []string{ "99" },
			},
			99,
		},
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
