package tsfm

import (
	"bitbucket.org/shell909090/scheme-go/scmgo"
	logging "github.com/op/go-logging"
)

var (
	log               = logging.MustGetLogger("trans")
	DefineTransformer = &Transformer{syntaxes: make(map[string]*Syntax)}
)

type Transformer struct {
	syntaxes map[string]*Syntax
}

func (t *Transformer) Parse(obj scmgo.SchemeObject) (err error) {
	code, ok := obj.(*scmgo.Cons)
	if !ok {
		return scmgo.ErrType
	}
	err = code.Iter(func(o scmgo.SchemeObject) (e error) {
		s, e := DefineSyntax(o)
		if e != nil {
			log.Error("%s", e.Error())
			return
		}
		if _, ok := t.syntaxes[s.Keyword]; ok {
			return ErrSyntaxExist
		}
		t.syntaxes[s.Keyword] = s
		return
	}, false)
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	return
}

func Transform(src scmgo.SchemeObject) (code scmgo.SchemeObject, err error) {
	code = src

	p := CreatePatternList()
	p.rule_list = p.rule_list.Push(CreatePatternVariable("dest"))
	p.rule_list = p.rule_list.Push(CreatePatternLiteral("display"))

	c, ok := code.(*scmgo.Cons)
	if !ok {
		return nil, scmgo.ErrUnknown
	}

	err = Walker(c, func(o *scmgo.Cons) (err error) {
		mr := CreateMatchResult()
		yes, err := p.Match(mr, o)
		if err != nil {
			return
		}
		if yes {
			log.Info("%v", mr.m)
			panic("ok")
		}
		return
	})
	if err != nil {
		return
	}

	return
}
