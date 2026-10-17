package bslack

import (
	"io/ioutil"
	"testing"

	"github.com/42wim/matterbridge/bridge"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestExtractTopicOrPurpose(t *testing.T) {
	testcases := map[string]struct {
		input          string
		wantChangeType string
		wantOutput     string
	}{
		"success - topic type":   {"@someone set channel topic: foo bar", "topic", "foo bar"},
		"success - purpose type": {"@someone set channel purpose: foo bar", "purpose", "foo bar"},
		"success - one line":     {"@someone set channel topic: foo bar", "topic", "foo bar"},
		"success - multi-line":   {"@someone set channel topic: foo\nbar", "topic", "foo\nbar"},
		"success - cleared":      {"@someone cleared channel topic", "topic", ""},
		"error - unhandled":      {"some unmatched message", "unknown", ""},
	}

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	cfg := &bridge.Config{Bridge: &bridge.Bridge{Log: logrus.NewEntry(logger)}}
	b := newBridge(cfg)
	for name, tc := range testcases {
		gotChangeType, gotOutput := b.extractTopicOrPurpose(tc.input)

		assert.Equalf(t, tc.wantChangeType, gotChangeType, "This testcase failed: %s", name)
		assert.Equalf(t, tc.wantOutput, gotOutput, "This testcase failed: %s", name)
	}
}

func TestReplaceURL(t *testing.T) {
	testcases := map[string]struct {
		input      string
		wantOutput string
	}{
		"bare link":            {"see <https://example.com/a?b=1&amp;c=2>", "see https://example.com/a?b=1&amp;c=2"},
		"label equals url":     {"<https://example.com|https://example.com>", "https://example.com"},
		"label without scheme": {"<http://example.com|example.com>", "http://example.com"},
		"label trailing slash": {"<http://example.com/|example.com>", "http://example.com/"},
		"label truncated":      {"retro <https://youtu.be/tAJHwUc3trQ?si=JuA0E2Wb5ANaM4sf|youtu.be/tAJHwUc3trQ?si=…>", "retro https://youtu.be/tAJHwUc3trQ?si=JuA0E2Wb5ANaM4sf"},
		"label without www":    {"<https://www.youtube.com/watch?v=tAJHwUc3trQ|youtube.com/watch?v=tAJHwUc3trQ>", "https://www.youtube.com/watch?v=tAJHwUc3trQ"},
		"label with www":       {"<http://www.example.com|www.example.com>", "http://www.example.com"},
		"mailto":               {"<mailto:foo@example.com|foo@example.com>", "foo@example.com"},
		"custom label":         {"<https://example.com|the docs>", "[the docs](https://example.com)"},
		"multiple links":       {"<https://a.com> and <https://b.com|b>", "https://a.com and [b](https://b.com)"},
		"unresolved mention":   {"<@U12345>", "<@U12345>"},
		"unresolved variable":  {"<!here>", "<!here>"},
	}

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	cfg := &bridge.Config{Bridge: &bridge.Bridge{Log: logrus.NewEntry(logger)}}
	b := newBridge(cfg)
	for name, tc := range testcases {
		assert.Equalf(t, tc.wantOutput, b.replaceURL(tc.input), "This testcase failed: %s", name)
	}
}
