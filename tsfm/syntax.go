package tsfm

import (
	"fmt"
	"strings"

	"bitbucket.org/shell909090/scheme-go/scmgo"
)

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

func (m *MatchResult) Format() (r string) {
	var strs []string
	for name, value := range m.m {
		strs = append(strs, fmt.Sprintf("%s = %s", name, value.Format()))
	}
	return strings.Join(strs, "\n")
}

type Rule struct {
	pattern  scmgo.SchemeObject
	template scmgo.SchemeObject
}

func ParseRule(rule *scmgo.Cons) (r *Rule, err error) {
	r = &Rule{}

	r.pattern, rule, err = rule.PopCons()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	log.Info("pattern: %s", r.pattern.Format())

	r.template, rule, err = rule.PopCons()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	log.Info("template: %s", r.template.Format())
	return
}

type Syntax struct {
	Keyword  string
	literals Literals
	rules    []*Rule
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
	log.Info("syntax: %s", s.Keyword)

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
	s.literals, err = ReadLiterals(sliterals)
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
	log.Info("transform: %s", i.Format())
	var yes bool
	for _, rule := range s.rules {
		mr := CreateMatchResult()
		yes, err = Match(rule.pattern, i, s.literals, mr)
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
		if yes {
			log.Info("match result: %s", mr.Format())
			return mr.Copy(rule.t) // render template
		}
	}
	return nil, ErrNoRule
}
