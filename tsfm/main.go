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

func (t *Transformer) Transform(src scmgo.SchemeObject) (code scmgo.SchemeObject, err error) {
	code = src

	c, ok := code.(*scmgo.Cons)
	if !ok {
		return nil, scmgo.ErrUnknown
	}

	err = Walker(c, func(o *scmgo.Cons) (err error) {
		s, _, err := o.PopSymbol()
		if err != nil { // FIXME: not a symbol is not a error, but maybe.
			err = nil
			return
		}
		syntax, ok := t.syntaxes[s.Name]
		if !ok { // nothing match
			return
		}

		_, err = syntax.Transform(o)
		if err != nil {
			log.Error("%s", err.Error())
			return
		} // FIXME: how to write o back?
		return
	})
	if err != nil {
		return
	}

	return
}
