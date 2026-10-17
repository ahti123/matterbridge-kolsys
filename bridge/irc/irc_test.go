package birc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceMarkdownLinks(t *testing.T) {
	testcases := map[string]struct {
		input      string
		wantOutput string
	}{
		"plain url":    {"see https://example.com", "see https://example.com"},
		"labeled link": {"read [the docs](https://example.com/a?b=1) now", "read the docs <https://example.com/a?b=1> now"},
		"two links":    {"[a](https://a.com) [b](http://b.com)", "a <https://a.com> b <http://b.com>"},
		"not a link":   {"array[0](x)", "array[0](x)"},
		"no text":      {"nothing here", "nothing here"},
	}
	for name, tc := range testcases {
		assert.Equalf(t, tc.wantOutput, replaceMarkdownLinks(tc.input), "This testcase failed: %s", name)
	}
}
