package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type matcherTest struct {
	name     string
	input    string
	expected Match
}

func TestMatchIssues(t *testing.T) {
	for _, test := range []matcherTest{
		{
			name:     "simple reference",
			input:    "owner/repo#123",
			expected: Issue{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:     "reference inside markdown link",
			input:    "[owner/repo#123: title](https://github.com/owner/repo/issues/123)",
			expected: Issue{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:  "no reference",
			input: "whatever / other thing # 123",
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			issue, matched := match(test.input)
			if test.expected == nil {
				assert.False(t, matched)
			} else {
				assert.True(t, matched)
				assert.Equal(t, test.expected, issue)
			}
		})
	}
}

func TestMatchURLs(t *testing.T) {
	for _, test := range []matcherTest{
		{
			name:     "issue url",
			input:    "https://github.com/owner/repo/issue/123",
			expected: Issue{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:     "issue url inside markdown",
			input:    "[owner/repo#123: title](https://github.com/owner/repo/issue/123)",
			expected: Issue{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:     "pull url",
			input:    "https://github.com/owner/repo/pull/123",
			expected: Pull{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:     "pull url inside markdown",
			input:    "[owner/repo#123: title](https://github.com/owner/repo/pull/123)",
			expected: Pull{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:     "repo discussion url",
			input:    "https://github.com/owner/repo/discussions/123",
			expected: Discussion{Owner: "owner", Repo: "repo", Num: "123"},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			reference, matched := match(test.input)
			if (test.expected == Issue{}) {
				assert.False(t, matched)
			} else {
				assert.True(t, matched)
				assert.Equal(t, test.expected, reference)
			}
		})
	}
}

func TestMatchURLPrecedence(t *testing.T) {
	ref, matched := match("[owner/repo#123: title](https://github.com/another/repo/issue/456)")
	assert.True(t, matched)
	assert.Equal(t, Issue{Owner: "another", Repo: "repo", Num: "456"}, ref)
}
