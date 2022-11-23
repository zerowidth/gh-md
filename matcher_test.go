package main

import (
	"strings"
	"testing"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
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
				assert.True(t, matched, "should have matched")
				assert.Equal(t, test.expected, issue)
			}
		})
	}
}

func TestMatchURLs(t *testing.T) {
	for _, test := range []matcherTest{
		{
			name:     "issue url",
			input:    "https://github.com/owner/repo/issues/123",
			expected: Issue{Owner: "owner", Repo: "repo", Num: "123"},
		},
		{
			name:     "issue url inside markdown",
			input:    "[owner/repo#123: title](https://github.com/owner/repo/issues/123)",
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
				assert.False(t, matched, "should not have matched")
			} else {
				assert.True(t, matched, "should have matched")
				assert.Equal(t, test.expected, reference)
			}
		})
	}
}

func TestMatchURLPrecedence(t *testing.T) {
	ref, matched := match("[owner/repo#123: title](https://github.com/another/repo/issues/456)")
	assert.True(t, matched)
	assert.Equal(t, Issue{Owner: "another", Repo: "repo", Num: "456"}, ref)
}

func TestIssueTitle(t *testing.T) {
	client := testClient(t, "fixtures/issue_title")
	issue := Issue{Owner: "github", Repo: "scientist", Num: "174"}
	title, err := issue.Title(client)
	require.NoError(t, err)
	assert.Equal(t, "Allow setting raise_on_mismatches to base class level for tests", title)
}

func TestPullRequestTitle(t *testing.T) {
	client := testClient(t, "fixtures/pull_request_title")
	pull := Pull{Owner: "github", Repo: "scientist", Num: "175"}
	title, err := pull.Title(client)
	require.NoError(t, err)
	assert.Equal(t, "Add Ruby 3.1 to CI", title)
}

func TestDiscussionTitle(t *testing.T) {
	client := testClient(t, "fixtures/discussion_title")
	discussion := Discussion{Owner: "cli", Repo: "cli", Num: "2673"}
	title, err := discussion.Title(client)
	require.NoError(t, err)
	assert.Equal(t, "Upgrade command", title)
}

func testClient(t *testing.T, fixturePath string) api.GQLClient {
	r, err := recorder.New(fixturePath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = r.Stop() })
	require.Equal(t, recorder.ModeRecordOnce, r.Mode())

	hook := func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		for h := range i.Response.Headers {
			if strings.HasPrefix(h, "X-") {
				delete(i.Response.Headers, h)
			}
		}
		return nil
	}
	r.AddHook(hook, recorder.AfterCaptureHook)

	client, err := gh.GQLClient(&api.ClientOptions{
		EnableCache: false,
		Transport:   r.GetDefaultClient().Transport,
	})
	require.NoError(t, err)

	return client
}
