package main

import (
	"regexp"
)

type Reference interface {
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

var nwoReferencePattern = regexp.MustCompile(`https://github.com/([^/]+)/([^/]+)/(issue|pull|discussions)/(\d+)(.*)`)
var issueReferencePattern = regexp.MustCompile(`\b([^/]+)/([^/#]+)#(\d+)\b`)

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

func matchURL(url string) (Reference, bool) {
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
