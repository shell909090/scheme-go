package scmgo

import (
	"io"
	"strconv"
)

type Evalor interface {
	Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error)
}

type Formatter interface {
	Format(s io.Writer, lv int) (rv int, err error)
}

type SchemeObject interface {
	Evalor
	Formatter
}

type Symbol struct {
	Name string
}

func (o *Symbol) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = env.Get(o.Name)
	if r == nil {
		return nil, nil, ErrNameNotFound
	}
	return
}

func (o *Symbol) Format(s io.Writer, lv int) (rv int, err error) {
	s.Write([]byte(o.Name))
	return lv + len(o.Name), nil
}

type Quote struct {
	objs SchemeObject
}

func (o *Quote) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = o.objs
	return
}

func (o *Quote) Format(s io.Writer, lv int) (rv int, err error) {
	s.Write([]byte("'"))
	rv, err = o.objs.Format(s, lv+1)
	return
}

type Boolean bool

func (o Boolean) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = o
	return
}

func (o Boolean) Format(s io.Writer, lv int) (rv int, err error) {
	if o {
		s.Write([]byte("#t"))
	} else {
		s.Write([]byte("#f"))
	}
	return lv + 2, nil
}

const (
	Otrue  = Boolean(true)
	Ofalse = Boolean(false)
)

type Integer int

func (o Integer) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = o
	return
}

func (o Integer) Format(s io.Writer, lv int) (rv int, err error) {
	rv, err = s.Write([]byte(strconv.FormatInt(int64(o), 10)))
	rv += lv
	return
}

type Float float64

func (o Float) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = o
	return
}

func (o Float) Format(s io.Writer, lv int) (rv int, err error) {
	rv, err = s.Write([]byte(strconv.FormatFloat(float64(o), 'f', 2, 64)))
	rv += lv
	return
}

type String string

func (o String) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = o
	return
}

func (o String) Format(s io.Writer, lv int) (rv int, err error) {
	s.Write([]byte("\""))
	rv, err = s.Write([]byte(o))
	s.Write([]byte("\""))
	rv += lv + 2
	return
}

type Cons struct {
	Car SchemeObject
	Cdr SchemeObject
}

var Onil = &Cons{}

func (o *Cons) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	var procedure SchemeObject
	procedure, o, err = o.Pop()
	if err != nil {
		return
	}

	af := CreateApplyFrame(p, o, env)
	next = CreateEvalFrame(af, procedure, env)
	return
}

func (o *Cons) Format(s io.Writer, lv int) (rv int, err error) {
	if o.Car == nil || o.Cdr == nil {
		_, err = s.Write([]byte("()"))
		return lv + 2, err
	}

	anycons, err := o.anyCons()
	if err != nil {
		return
	}
	if anycons {
		return o.PrettyFormat(s, lv)
	}

	s.Write([]byte("("))
	rv = lv + 1

	ok := true
	for i := o; i != Onil; i, ok = i.Cdr.(*Cons) {
		if !ok {
			return rv, ErrISNotAList
		}
		rv, err = i.Car.Format(s, rv)
		if err != nil {
			return
		}
		rv += 1
		if i.Cdr != Onil { // not last one
			s.Write([]byte(" "))
		}
	}
	s.Write([]byte(")")) // last one here
	return
}

func (o *Cons) PrettyFormat(s io.Writer, lv int) (rv int, err error) {
	obj := o
	s.Write([]byte("("))
	lv += 1

	if _, ok := obj.Car.(*Symbol); ok {
		rv, err = obj.Car.Format(s, lv)
		if err != nil {
			return
		}

		if obj.Cdr != Onil {
			s.Write([]byte(" "))
			lv = rv + 1
		}
		obj, ok = obj.Cdr.(*Cons)
	} else {
		ok = true
	}

	ok := true
	for ; obj != Onil; obj, ok = obj.Cdr.(*Cons) {
		if !ok {
			s.Write([]byte(" . "))
			lv += 3

			rv, err = obj.Cdr.Format(s, lv)
			if err != nil {
				return
			}

			s.Write([]byte(")"))
			rv += 1
			return
		}

		rv, err = obj.Car.Format(s, lv)
		if err != nil {
			return
		}
		if obj.Cdr != Onil {
			s.Write([]byte("\n"))
			for i := 0; i < lv; i++ {
				s.Write([]byte(" "))
			}
		}
	}

	s.Write([]byte(")"))
	rv = lv + 1
	return
}

func (o *Cons) Iter(f func(obj SchemeObject) (e error)) (err error) {
	ok := true
	for i := o; i != Onil; i, ok = i.Cdr.(*Cons) {
		if !ok {
			return ErrISNotAList
		}
		err = f(i.Car)
		if err != nil {
			return
		}
	}
	return
}

func (o *Cons) Pop() (r SchemeObject, next *Cons, err error) {
	r = o.Car
	next, ok := o.Cdr.(*Cons)
	if !ok {
		return nil, nil, ErrISNotAList
	}
	return
}

func (o *Cons) Len() (n int, err error) {
	err = o.Iter(func(obj SchemeObject) (e error) {
		n += 1
		return
	})
	return
}

func (o *Cons) GetN(n int) (r SchemeObject, err error) {
	var ok bool
	c := o
	for i := 0; i < n; i++ {
		switch c.Cdr {
		case nil:
			return nil, ErrUnknown
		case Onil:
			return nil, ErrListOutOfIndex
		}
		c, ok = o.Cdr.(*Cons)
		if !ok {
			return nil, ErrISNotAList
		}
	}
	return c.Car, nil
}

func (o *Cons) anyCons() (yes bool, err error) {
	err = o.Iter(func(obj SchemeObject) (e error) {
		_, yes = obj.(*Cons)
		if yes {
			e = ErrUnknown
		}
		return
	})
	if err == ErrUnknown {
		err = nil
	}
	return
}
