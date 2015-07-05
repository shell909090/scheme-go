package tsfm

import "bitbucket.org/shell909090/scheme-go/scmgo"

type MatchResult interface {
	Add(name string, value scmgo.SchemeObject)
}
