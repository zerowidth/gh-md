package main

import (
	"regexp"
	"strconv"

	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type Match interface {
	// Reference returns a string to reference this match by
	Reference() string
	// Title fetches the title of a match from the API
	Title(client api.GQLClient) (string, error)
}

type Issue struct {
	Owner string
	Repo  string
	Num   string
}

type Pull struct {
	Owner string
	Repo  string
	Num   string
}

type Discussion struct {
	Owner string
	Repo  string
	Num   string
}

func (i Issue) Reference() string {
	return i.Owner + "/" + i.Repo + "#" + i.Num
}

func (i Issue) Title(client api.GQLClient) (string, error) {
	var query struct {
		Repository struct {
			IssueOrPullRequest struct {
				Issue struct {
					Title string
				} `graphql:"... on Issue"`
				PullRequest struct {
					Title string
				} `graphql:"... on PullRequest"`
			} `graphql:"issueOrPullRequest(number: $number)"`
			Discussion struct {
				Title string
			} `graphql:"discussion(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	num, err := strconv.Atoi(i.Num)
	if err != nil {
		return "", err
	}
	variables := map[string]interface{}{
		"owner":  graphql.String(i.Owner),
		"name":   graphql.String(i.Repo),
		"number": graphql.Int(num),
	}
	err = client.Query("IssueTitle", &query, variables)
	// ignore error unless no title was found - a "not found" is expected for either
	// the issue-ish or the discussion.
	if query.Repository.IssueOrPullRequest.Issue.Title != "" {
		return query.Repository.IssueOrPullRequest.Issue.Title, nil
	}
	if query.Repository.IssueOrPullRequest.PullRequest.Title != "" {
		return query.Repository.IssueOrPullRequest.PullRequest.Title, nil
	}
	if query.Repository.Discussion.Title != "" {
		return query.Repository.Discussion.Title, nil
	}
	return "", err
}

func (p Pull) Reference() string {
	return p.Owner + "/" + p.Repo + "#" + p.Num
}

func (p Pull) Title(client api.GQLClient) (string, error) {
	var query struct {
		Repository struct {
			IssueOrPullRequest struct {
				Issue struct {
					Title string
				} `graphql:"... on Issue"`
				PullRequest struct {
					Title string
				} `graphql:"... on PullRequest"`
			} `graphql:"issueOrPullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	num, err := strconv.Atoi(p.Num)
	if err != nil {
		return "", err
	}
	variables := map[string]interface{}{
		"owner":  graphql.String(p.Owner),
		"name":   graphql.String(p.Repo),
		"number": graphql.Int(num),
	}
	err = client.Query("IssueTitle", &query, variables)
	if query.Repository.IssueOrPullRequest.Issue.Title != "" {
		return query.Repository.IssueOrPullRequest.Issue.Title, nil
	}
	if query.Repository.IssueOrPullRequest.PullRequest.Title != "" {
		return query.Repository.IssueOrPullRequest.PullRequest.Title, nil
	}
	return "", err
}

func (d Discussion) Reference() string {
	return d.Owner + "/" + d.Repo + "#" + d.Num
}

func (d Discussion) Title(client api.GQLClient) (string, error) {
	var query struct {
		Repository struct {
			Discussion struct {
				Title string
			} `graphql:"discussion(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	num, err := strconv.Atoi(d.Num)
	if err != nil {
		return "", err
	}
	variables := map[string]interface{}{
		"owner":  graphql.String(d.Owner),
		"name":   graphql.String(d.Repo),
		"number": graphql.Int(num),
	}
	err = client.Query("DiscussionTitle", &query, variables)
	if query.Repository.Discussion.Title != "" {
		return query.Repository.Discussion.Title, nil
	}
	return "", err
}

var urlReferencePattern = regexp.MustCompile(`https://github.com/([^/]+)/([^/]+)/(issues|pull|discussions)/(\d+)(.*)`)
var issueReferencePattern = regexp.MustCompile(`\b([^/]+)/([^/#]+)#(\d+)\b`)

// match attempts to find a Match (issue, PR, discussion) in the given input.
func match(input string) (Match, bool) {
	if ref, ok := matchURL(input); ok {
		return ref, true
	}

	if issue, ok := matchIssueReference(input); ok {
		return issue, true
	}

	return nil, false
}

func matchIssueReference(input string) (Issue, bool) {
	matches := issueReferencePattern.FindStringSubmatch(input)
	if matches != nil {
		return Issue{
			Owner: matches[1],
			Repo:  matches[2],
			Num:   matches[3],
		}, true
	}

	return Issue{}, false
}

func matchURL(url string) (Match, bool) {
	matches := urlReferencePattern.FindStringSubmatch(url)
	if matches == nil {
		return nil, false
	}

	switch matches[3] {
	case "issues":
		return Issue{
			Owner: matches[1],
			Repo:  matches[2],
			Num:   matches[4],
		}, true
	case "pull":
		return Pull{
			Owner: matches[1],
			Repo:  matches[2],
			Num:   matches[4],
		}, true
	case "discussions":
		return Discussion{
			Owner: matches[1],
			Repo:  matches[2],
			Num:   matches[4],
		}, true
	}

	return nil, false
}
