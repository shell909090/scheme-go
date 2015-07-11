package tsfm

import (
	"bitbucket.org/shell909090/scheme-go/scmgo"
)

type Template interface {
	Render(mr *MatchResult) (result *scmgo.SchemeObject, err error)
}

type RootTemplate struct {
	root scmgo.SchemeObject
}

func (t *RootTemplate) Render(mr *MatchResult) (result *scmgo.SchemeObject, err error) {
	return
}
