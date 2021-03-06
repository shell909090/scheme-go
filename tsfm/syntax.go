package main

import (
	"github.com/shell909090/scheme-go/impl"
	"github.com/shell909090/scheme-go/scm"
)

type Rule struct {
	pattern  scm.Obj
	template scm.Obj
}

func ParseRule(rule *scm.Cons) (r *Rule, err error) {
	var pattern, template *scm.Cons
	rule, err = impl.ParseParameters(rule, &pattern, &template)
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	r = &Rule{pattern: pattern, template: template}
	log.Debug("pattern: %s", scm.Format(r.pattern))
	log.Debug("template: %s", scm.Format(r.template))
	return
}

type Syntax struct {
	Keyword  string
	literals Literals
	rules    []*Rule
}

func DefineSyntax(obj scm.Obj) (s *Syntax, err error) {
	// get define-syntax and check symbol.
	list, ok := obj.(*scm.Cons)
	if !ok {
		return nil, scm.ErrType
	}
	var define, keyword string
	list, err = impl.ParseParameters(list, &define, &keyword)
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

	syntax, ok := list.Car.(*scm.Cons)
	if !ok {
		return nil, scm.ErrType
	}
	err = s.Parse(syntax)
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	return
}

func (s *Syntax) Parse(syntax *scm.Cons) (err error) {
	// get syntax-rules and check symbol, and get literals.
	var syntax_rules string
	var literals *scm.Cons
	syntax, err = impl.ParseParameters(syntax, &syntax_rules, &literals)
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

	var rule *scm.Cons
	var r *Rule
	for syntax != scm.Onil { // get rules.
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

func (s *Syntax) Transform(i scm.Obj) (result scm.Obj, err error) {
	log.Info("transform: %s", scm.Format(i))
	var yes bool
	for _, rule := range s.rules {
		mr := CreateMatchResult()
		log.Debug("match: %s", scm.Format(rule.pattern))
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
