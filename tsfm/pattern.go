package tsfm

import "bitbucket.org/shell909090/scheme-go/scmgo"

type Pattern interface {
	Match(mr MatchResult, i scmgo.SchemeObject) (yes bool, err error)
}

type PartternAny struct {
}

func (p *PartternAny) Match(mr MatchResult, i scmgo.SchemeObject) (yes bool, err error) {
	return true, nil
}

type PartternVariable struct {
	toName string
}

func (p *PartternVariable) Match(mr MatchResult, i scmgo.SchemeObject) (yes bool, err error) {
	// set x[p.toName] = i
	mr.Add(p.toName, i)
	return true, nil
}

type PartternLiteral struct {
	name string
}

func (p *PartternLiteral) Match(mr MatchResult, i scmgo.SchemeObject) (yes bool, err error) {
	t, ok := i.(*scmgo.Symbol)
	if !ok {
		return false, nil
	}
	return t.Name == p.name, nil
}

type PartternList struct {
	rules []Pattern
}

func (p *PartternList) Match(mr MatchResult, i scmgo.SchemeObject) (yes bool, err error) {
	// o, ok := i.(*scmgo.Cons)
	// if !ok {
	// 	return false, nil
	// }

	// for rule := range p.rules {
	// }
	return
}
