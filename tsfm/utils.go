package tsfm

import (
	"errors"

	"bitbucket.org/shell909090/scheme-go/scmgo"
)

var (
	ErrNotAPattern = errors.New("not a pattern")
	ErrWrongStruct = errors.New("wrong struct")
	ErrSyntaxExist = errors.New("syntax exist")
)

func Walker(o *scmgo.Cons, f func(o *scmgo.Cons) (err error)) (err error) {
	err = o.Iter(func(i scmgo.SchemeObject) (err error) {
		c, ok := i.(*scmgo.Cons)
		if !ok {
			return
		}
		err = f(c)
		if err != nil {
			log.Error("%s", err)
			return
		}
		err = Walker(c, f)
		return
	}, false)
	if err != nil {
		log.Error("%s", err)
		return
	}
	return
}

func FilterList(o *scmgo.Cons, f func(i scmgo.SchemeObject) (yes bool, err error)) (err error) {
	var ok, yes bool
	for c := o; c != scmgo.Onil; {
		yes, err = f(c.Car)
		if err != nil {
			log.Error("%s", err)
			return
		}

		if yes {
			n, ok := c.Cdr.(*scmgo.Cons)
			if !ok {
				return scmgo.ErrISNotAList
			}
			c.Car = n.Car
			c.Cdr = n.Cdr
		} else {
			c, ok = c.Cdr.(*scmgo.Cons)
			if !ok {
				return scmgo.ErrISNotAList
			}
		}
	}
	return
}
