package main

import (
	"regexp"
)

type Match interface {
	// Reference returns a string to reference this match by
	Reference() string
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

func (p Pull) Reference() string {
	return p.Owner + "/" + p.Repo + "#" + p.Num
}

func (d Discussion) Reference() string {
	return d.Owner + "/" + d.Repo + "#" + d.Num
}

var nwoReferencePattern = regexp.MustCompile(`https://github.com/([^/]+)/([^/]+)/(issue|pull|discussions)/(\d+)(.*)`)
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
	matches := nwoReferencePattern.FindStringSubmatch(url)
	if matches == nil {
		return nil, false
	}

	switch matches[3] {
	case "issue":
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
