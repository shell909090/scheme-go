package tsfm

import "bitbucket.org/shell909090/scheme-go/scmgo"

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

	err = Walker(c, func(o *scmgo.Cons) (result scmgo.SchemeObject, err error) {
		s, ok := o.Car.(*scmgo.Symbol)
		if !ok { // not a symbol is not a error.
			return
		}
		syntax, ok := t.syntaxes[s.Name]
		if !ok { // nothing match
			return
		}

		result, err = syntax.Transform(o)
		if err != nil {
			log.Error("%s", err.Error())
			return
		}

		log.Warning("render result: %s", result.Format())
		return
	})
	if err != nil {
		return
	}

	return
}

func Walker(o *scmgo.Cons, f func(o *scmgo.Cons) (scmgo.SchemeObject, error)) (err error) {
	var ok bool
	var tmplist *scmgo.Cons
	var tmp scmgo.SchemeObject
	for n := o; n != scmgo.Onil; {
		tmplist, ok = n.Car.(*scmgo.Cons)
		if ok {
			tmp, err = f(tmplist)
			if err != nil {
				log.Error("%s", err.Error())
				return
			}
			if tmp != nil {
				n.Car = tmp
			}

			tmplist, ok = n.Car.(*scmgo.Cons)
			if ok {
				err = Walker(tmplist, f)
				if err != nil {
					log.Error("%s", err.Error())
					return
				}
			}
		}

		n, ok = n.Cdr.(*scmgo.Cons)
		if !ok { // improper
			return
		}
	}
	return
}
