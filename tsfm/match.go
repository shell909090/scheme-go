package tsfm

import (
	"fmt"
	"strings"

	"bitbucket.org/shell909090/scheme-go/scmgo"
)

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

// func (m *MatchResult) Copy(obj scmgo.SchemeObject) (result scmgo.SchemeObject, err error) {
// 	switch t := obj.(type) {
// 	case *scmgo.Symbol:
// 		o, ok = m.m[t.Name]
// 		if ok {
// 		}
// 	}
// }

func (m *MatchResult) Format() (r string) {
	var strs []string
	for name, value := range m.m {
		strs = append(strs, fmt.Sprintf("%s = %s", name, value.Format()))
	}
	return strings.Join(strs, "\n")
}
