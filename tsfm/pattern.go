package tsfm

import (
	"bitbucket.org/shell909090/scheme-go/scmgo"
)

type PatternObject struct {
}

func (p *PatternObject) Eval(env *scmgo.Environ, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
	panic("run eval of partten object")
}

func (p *PatternObject) Format() (r string) {
	return "parttern"
}

type Pattern interface {
	Match(mr *MatchResult, i scmgo.SchemeObject) (yes bool, err error)
}

type PatternAny struct {
	PatternObject
}

func (p *PatternAny) Match(mr *MatchResult, i scmgo.SchemeObject) (yes bool, err error) {
	return true, nil
}

type PatternVariable struct {
	PatternObject
	toName string
}

func CreatePatternVariable(toname string) (p *PatternVariable) {
	p = &PatternVariable{toName: toname}
	return
}

func (p *PatternVariable) Match(mr *MatchResult, i scmgo.SchemeObject) (yes bool, err error) {
	mr.Add(p.toName, i)
	return true, nil
}

type PatternLiteral struct {
	PatternObject
	name string
}

func CreatePatternLiteral(name string) (p *PatternLiteral) {
	p = &PatternLiteral{name: name}
	return
}

func (p *PatternLiteral) Match(mr *MatchResult, i scmgo.SchemeObject) (yes bool, err error) {
	t, ok := i.(*scmgo.Symbol)
	if !ok {
		return false, nil
	}
	return t.Name == p.name, nil
}

type PatternList struct {
	PatternObject
	rule_list *scmgo.Cons
}

func CreatePatternList() (p *PatternList) {
	p = &PatternList{rule_list: scmgo.Onil}
	return
}

func (p *PatternList) Match(mr *MatchResult, i scmgo.SchemeObject) (yes bool, err error) {
	olist, ok := i.(*scmgo.Cons)
	if !ok {
		return false, nil
	}
	rlist := p.rule_list

	for olist != scmgo.Onil || rlist != scmgo.Onil {
		log.Info("%s %s", olist.Format(), rlist.Format())

		obj := olist.Car
		subp, ok := rlist.Car.(Pattern)
		if !ok {
			return false, ErrNotAPattern
		}

		yes, err = subp.Match(mr, obj)
		if err != nil {
			return
		}
		if !yes {
			return false, nil
		}

		switch t := olist.Cdr.(type) {
		case *scmgo.Cons: // continue on list.
			olist = t
			rlist, ok = rlist.Cdr.(*scmgo.Cons)
			if !ok {
				return false, nil
			}
		default: // pair in the end.
			subp, ok = rlist.Cdr.(Pattern)
			if !ok { // they are just not match
				return false, nil
			}

			yes, err = subp.Match(mr, olist.Cdr)
			return // match or not, they are final.
		}
	}

	if olist == scmgo.Onil && rlist == scmgo.Onil {
		return true, nil
	}
	return false, nil
}
