package tsfm

import (
	"errors"

	"bitbucket.org/shell909090/scheme-go/scmgo"
)

var (
	ErrNotAPattern    = errors.New("not a pattern")
	ErrWrongTpInPtn   = errors.New("wrong type in pattern")
	ErrWrongStruct    = errors.New("wrong struct")
	ErrSyntaxExist    = errors.New("syntax exist")
	ErrElpsAfterNoVar = errors.New("ellipsis after something not a pattern varible")
	ErrNoRule         = errors.New("no rule match in syntax")
)

var (
	log               = logging.MustGetLogger("trans")
	DefineTransformer = &Transformer{syntaxes: make(map[string]*Syntax)}
)

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
