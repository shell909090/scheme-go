package tsfm

import (
	"bitbucket.org/shell909090/scheme-go/internal"
	"bitbucket.org/shell909090/scheme-go/scmgo"
)

type Rule struct {
	pattern  scmgo.SchemeObject
	template scmgo.SchemeObject
}

func ParseRule(rule *scmgo.Cons) (r *Rule, err error) {
	var pattern, template *scmgo.Cons
	rule, err = internal.ParseParameters(rule, &pattern, &template)
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	r = &Rule{pattern: pattern, template: template}
	log.Debug("pattern: %s", scmgo.Format(r.pattern))
	log.Debug("template: %s", scmgo.Format(r.template))
	return
}

type Syntax struct {
	Keyword  string
	literals Literals
	rules    []*Rule
}

func DefineSyntax(obj scmgo.SchemeObject) (s *Syntax, err error) {
	// get define-syntax and check symbol.
	list, ok := obj.(*scmgo.Cons)
	if !ok {
		return nil, scmgo.ErrType
	}
	var define, keyword string
	list, err = internal.ParseParameters(list, &define, &keyword)
	if err != nil { // is it ok if obj is totally not a define-syntax?
		err = nil // even not a list, or first two elements not symbols.
		return
	}
	if define != "define-syntax" {
		return // it is ok if not a define-syntax.
	}

	// get syntax keyword
	s = &Syntax{Keyword: keyword}
	log.Info("syntax: %s", s.Keyword)

	syntax, ok := list.Car.(*scmgo.Cons)
	if !ok {
		return nil, scmgo.ErrType
	}
	err = s.Parse(syntax)
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	return
}

func (s *Syntax) Parse(syntax *scmgo.Cons) (err error) {
	// get syntax-rules and check symbol, and get literals.
	var syntax_rules string
	var literals *scmgo.Cons
	syntax, err = internal.ParseParameters(syntax, &syntax_rules, &literals)
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	if syntax_rules != "syntax-rules" {
		return ErrWrongStruct
	}
	s.literals, err = ReadLiterals(literals)
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
		r, err = ParseRule(rule)
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
		s.rules = append(s.rules, r)
	}

	return
}

func (s *Syntax) Transform(i scmgo.SchemeObject) (result scmgo.SchemeObject, err error) {
	log.Info("transform: %s", scmgo.Format(i))
	var yes bool
	for _, rule := range s.rules {
		mr := CreateMatchResult()
		log.Debug("match: %s", scmgo.Format(rule.pattern))
		yes, err = Match(rule.pattern, i, s.literals, mr)
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
		if yes {
			log.Debug("result: %s", mr.Format())
			return Render(mr, rule.template)
		}
	}
	return nil, ErrNoRule
}
