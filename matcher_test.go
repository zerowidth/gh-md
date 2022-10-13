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
			input:    "[owner/repo#123](owner/repo#123: title)",
			expected: Issue{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:  "no reference",
			input: "whatever / other thing # 123",
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			issue, matched := matchIssue(test.input)
			if (test.expected == Issue{}) {
				assert.False(t, matched)
			} else {
				assert.True(t, matched)
				assert.Equal(t, test.expected, issue)
			}
		})
	}
}
