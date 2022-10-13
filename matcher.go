package main

import (
	"regexp"
)

type Issue struct {
	Owner string
	Repo  string
	Num   string
}

var issueReferencePattern = regexp.MustCompile(`\b([^/]+)/([^/#]+)#(\d+)\b`)

func matchIssue(input string) (Issue, bool) {
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
