package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type issueTest struct {
	name     string
	input    string
	expected Issue
}

func TestMatchIssue(t *testing.T) {
	for _, test := range []issueTest{
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
			issue, matched := matchIssueReference(test.input)
			if (test.expected == Issue{}) {
				assert.False(t, matched)
			} else {
				assert.True(t, matched)
				assert.Equal(t, test.expected, issue)
			}
		})
	}
}

type urlTest struct {
	name     string
	url      string
	expected Reference
}

func TestMatchURL(t *testing.T) {
	for _, test := range []urlTest{
		{
			name:     "issue url",
			url:      "https://github.com/owner/repo/issue/123",
			expected: Issue{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:     "issue url inside markdown",
			url:      "[owner/repo#123: title](https://github.com/owner/repo/issue/123)",
			expected: Issue{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:     "pull url",
			url:      "https://github.com/owner/repo/pull/123",
			expected: Pull{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:     "pull url inside markdown",
			url:      "[owner/repo#123: title](https://github.com/owner/repo/pull/123)",
			expected: Pull{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:     "repo discussion url",
			url:      "https://github.com/owner/repo/discussions/123",
			expected: Discussion{Owner: "owner", Repo: "repo", Num: "123"},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			reference, matched := matchURL(test.url)
			if (test.expected == Issue{}) {
				assert.False(t, matched)
			} else {
				assert.True(t, matched)
				assert.Equal(t, test.expected, reference)
			}
		})
	}
}
