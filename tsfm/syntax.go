package tsfm

import "bitbucket.org/shell909090/scheme-go/scmgo"

type Literals map[string]int

func ReadLiterals(l *scmgo.Cons) (literals Literals, err error) {
	var s *scmgo.Symbol
	literals = make(Literals, 0)
	for l != scmgo.Onil {
		s, l, err = l.PopSymbol()
		if err != nil {
			return
		}
		literals[s.Name] = 1
	}
	return
}

func (l Literals) CheckLiteral(s string) (yes bool) {
	_, yes = l[s]
	return
}

type Rule struct {
	p Pattern
	t Template
}

func ParseRule(literals Literals, rule *scmgo.Cons) (r *Rule, err error) {
	r = &Rule{}

	pattern, rule, err := rule.PopCons()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	r.p, err = ParsePattern(literals, pattern)
	if err != nil {
		log.Error("%s", err.Error())
		return
	}

	_, rule, err = rule.PopCons()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	// ParseTemplate

	return
}

type Syntax struct {
	Keyword string
	rules   []*Rule
}

func DefineSyntax(obj scmgo.SchemeObject) (s *Syntax, err error) {
	// get define-syntax and check symbol.
	define, ok := obj.(*scmgo.Cons)
	if !ok {
		return nil, scmgo.ErrType
	}
	sname, define, err := define.PopSymbol()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	if sname.Name != "define-syntax" {
		return // it is ok if not a define-syntax.
	}

	// get syntax keyword
	sname, define, err = define.PopSymbol()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	s = &Syntax{}
	s.Keyword = sname.Name

	syntax, _, err := define.PopCons()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	err = s.Parse(syntax)
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	return
}

func (s *Syntax) Parse(syntax *scmgo.Cons) (err error) {
	// get syntax-rules and check symbol.
	sname, syntax, err := syntax.PopSymbol()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	if sname.Name != "syntax-rules" {
		return ErrWrongStruct
	}

	// get literals.
	sliterals, syntax, err := syntax.PopCons()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	literals, err := ReadLiterals(sliterals)
	if err != nil {
		log.Error("%s", err.Error())
		return
	}

	var rule *scmgo.Cons
	var r *Rule
	for syntax != scmgo.Onil { // get rules.
		rule, syntax, err = syntax.PopCons()
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
		r, err = ParseRule(literals, rule)
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
		s.rules = append(s.rules, r)
	}

	return
}
