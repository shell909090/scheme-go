package tsfm

import "bitbucket.org/shell909090/scheme-go/scmgo"

type MatchResult struct {
	m map[string]scmgo.SchemeObject
}

func CreateMatchResult() (m *MatchResult) {
	m = &MatchResult{
		m: make(map[string]scmgo.SchemeObject),
	}
	return
}

func (m *MatchResult) Add(name string, value scmgo.SchemeObject) {
	m.m[name] = value
}
