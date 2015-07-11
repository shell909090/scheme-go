package tsfm

import (
	"bitbucket.org/shell909090/scheme-go/scmgo"
)

type Template struct {
	root scmgo.SchemeObject
}

func (t *Template) Render(mr *MatchResult) (result *scmgo.SchemeObject, err error) {
	return
}
