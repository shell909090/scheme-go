package tsfm

import "bitbucket.org/shell909090/scheme-go/scmgo"

type PatternObject struct {
}

func (p *PatternObject) Eval(env *scmgo.Environ, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
	panic("run eval of partten object")
}

type Pattern interface {
	Eval(env *scmgo.Environ, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error)
	Format() (r string)
	Match(mr *MatchResult, i scmgo.SchemeObject) (yes bool, err error)
}

type PatternAny struct {
	PatternObject
}

func (p *PatternAny) Format() (r string) {
	return "_"
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

func (p *PatternVariable) Format() (r string) {
	return p.toName
}

func (p *PatternVariable) Match(mr *MatchResult, i scmgo.SchemeObject) (yes bool, err error) {
	mr.Add(p.toName, i)
	return true, nil
}

type PatternEllipses struct {
	PatternObject
	toName string
}

func CreatePatternEllipses(pv *PatternVariable) (p *PatternEllipses) {
	p = &PatternEllipses{toName: pv.toName}
	return
}

func (p *PatternEllipses) Format() (r string) {
	return p.toName + " ..."
}

func (p *PatternEllipses) Match(mr *MatchResult, i scmgo.SchemeObject) (yes bool, err error) {
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

func (p *PatternLiteral) Format() (r string) {
	return "'" + p.name
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

func (p *PatternList) Format() (r string) {
	return p.rule_list.Format()
}

func (p *PatternList) Match(mr *MatchResult, i scmgo.SchemeObject) (yes bool, err error) {
	olist, ok := i.(*scmgo.Cons)
	if !ok {
		return false, nil
	}
	rlist := p.rule_list

	for olist != scmgo.Onil || rlist != scmgo.Onil {
		log.Info("now match: %s %s", olist.Format(), rlist.Format())

		obj := olist.Car
		subp, ok := rlist.Car.(Pattern)
		if !ok {
			return false, ErrNotAPattern
		}

		if _, ok := subp.(*PatternEllipses); ok {
			yes, err = subp.Match(mr, olist)
			if err != nil {
				log.Error("%s", err.Error())
				return
			}
			return yes, nil
		}

		yes, err = subp.Match(mr, obj)
		if err != nil {
			log.Error("%s", err.Error())
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

func ParsePattern(literals Literals, pattern *scmgo.Cons) (p Pattern, err error) {
	var ok bool
	var pv *PatternVariable
	pl := CreatePatternList()
	o := pattern
LOOP:
	for o != scmgo.Onil {
		switch ttmp := o.Car.(type) {
		case *scmgo.Symbol:
			switch {
			case ttmp.Name == "_":
				p = &PatternAny{}
			case ttmp.Name == "...":
				pv, ok = pl.rule_list.Car.(*PatternVariable)
				if !ok {
					return nil, ErrElpsAfterNoVar
				}
				pl.rule_list.Car = CreatePatternEllipses(pv)
				// that's the end
				break LOOP
			case literals.CheckLiteral(ttmp.Name):
				p = CreatePatternLiteral(ttmp.Name)
			default:
				p = CreatePatternVariable(ttmp.Name)
			}
		case *scmgo.Cons:
			p, err = ParsePattern(literals, ttmp)
		default:
			panic("not support yet")
		}

		if err != nil {
			log.Error("%s", err.Error())
			return
		}
		pl.rule_list = pl.rule_list.Push(p)

		o, ok = o.Cdr.(*scmgo.Cons)
		if !ok {
			pl.rule_list, err = scmgo.ReverseList(pl.rule_list, o.Cdr)
			return pl, err
		}
	}

	pl.rule_list, err = scmgo.ReverseList(pl.rule_list, scmgo.Onil)
	return pl, err
}
