package main

import "github.com/shell909090/scheme-go/scm"

type Transformer struct {
	syntaxes map[string]*Syntax
}

func (t *Transformer) Parse(obj scm.Obj) (err error) {
	code, ok := obj.(*scm.Cons)
	if !ok {
		return scm.ErrType
	}
	err = code.Iter(func(o scm.Obj) (e error) {
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

func (t *Transformer) Transform(src scm.Obj) (code scm.Obj, err error) {
	code = src

	c, ok := code.(*scm.Cons)
	if !ok {
		return nil, scm.ErrUnknown
	}

	err = Walker(c, func(o *scm.Cons) (result scm.Obj, err error) {
		s, ok := o.Car.(*scm.Symbol)
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

		log.Info("render result: %s", scm.Format(result))
		return
	})
	if err != nil {
		return
	}
	return
}

func Walker(o *scm.Cons, f func(o *scm.Cons) (scm.Obj, error)) (err error) {
	var ok bool
	var tmplist *scm.Cons
	var tmp scm.Obj
	for n := o; n != scm.Onil; {
		tmplist, ok = n.Car.(*scm.Cons)
		if ok {
			tmp, err = f(tmplist)
			if err != nil {
				log.Error("%s", err.Error())
				return
			}
			if tmp != nil {
				n.Car = tmp
			}

			tmplist, ok = n.Car.(*scm.Cons)
			if ok {
				err = Walker(tmplist, f)
				if err != nil {
					log.Error("%s", err.Error())
					return
				}
			}
		}

		n, ok = n.Cdr.(*scm.Cons)
		if !ok { // improper
			return
		}
	}
	return
}
