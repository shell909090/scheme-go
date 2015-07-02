package scmgo

import (
	"io"
	"strconv"
)

type Evalor interface {
	Eval(env *Environ, p Frame) (r SchemeObject, f Frame, err error)
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

func (o *Symbol) Eval(env *Environ, p Frame) (r SchemeObject, f Frame, err error) {
	r, ok := env.Get(o.Name)
	if !ok {
		return nil, nil, ErrNameNotFound
	}
	return
}

func (o *Symbol) Format(s io.Writer, lv int) (rv int, err error) {
	s.Write([]byte(o.Name))
	return lv + len(o.Name), nil
}

func SymbolFromString(s string) (o *Symbol) {
	return &Symbol{Name: s}
}

type Quote struct {
	objs SchemeObject
}

func (o *Quote) Eval(env *Environ, p Frame) (r SchemeObject, f Frame, err error) {
	r = o.objs
	return
}

func (o *Quote) Format(s io.Writer, lv int) (rv int, err error) {
	s.Write([]byte("'"))
	rv, err = o.objs.Format(s, lv+1)
	return
}

type Boolean bool

func (o Boolean) Eval(env *Environ, p Frame) (r SchemeObject, f Frame, err error) {
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
	Otrue  = true
	Ofalse = false
)

func BooleanFromString(s string) (o Boolean, err error) {
	// FIXME: not so good
	switch s[1] {
	case 't':
		return Otrue, nil
	case 'f':
		return Ofalse, nil
	default:
		return Otrue, ErrBooleanUnknown
	}
}

type Integer int

func (o Integer) Eval(env *Environ, p Frame) (r SchemeObject, f Frame, err error) {
	r = o
	return
}

func (o Integer) Format(s io.Writer, lv int) (rv int, err error) {
	rv, err = s.Write([]byte(strconv.FormatInt(int64(o), 10)))
	rv += lv
	return
}

func IntegerFromString(s string) (o Integer, err error) {
	i, err := strconv.Atoi(s)
	return Integer(i), err
}

type Float float64

func (o Float) Eval(env *Environ, p Frame) (r SchemeObject, f Frame, err error) {
	r = o
	return
}

func (o Float) Format(s io.Writer, lv int) (rv int, err error) {
	rv, err = s.Write([]byte(strconv.FormatFloat(float64(o), 'f', 2, 64)))
	rv += lv
	return
}

func FloatFromString(s string) (o Float, err error) {
	i, err := strconv.ParseFloat(s, 64)
	return Float(i), err
}

type String string

func (o String) Eval(env *Environ, p Frame) (r SchemeObject, f Frame, err error) {
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

func (o *Cons) Eval(env *Environ, p Frame) (r SchemeObject, f Frame, err error) {
	procedure := o.Car
	o, ok := o.Cdr.(*Cons)
	if !ok {
		return nil, nil, ErrISNotAList
	}

	af := ApplyFrame(o, env)
	af.SetParent(p)
	f = CreateEvalFrame(procedure, env)
	f.SetParent(af)
	return
}

func (o *Cons) Format(s io.Writer, lv int) (rv int, err error) {
	var ok bool

	if o.Car == nil || o.Cdr == nil {
		_, err = s.Write([]byte("()"))
		return lv + 2, err
	}

	anycons, err := o.anyCons()
	if err != nil {
		return
	}
	if !anycons {
		s.Write([]byte("("))
		rv = lv + 1
		ok = true
		for i := o; i != Onil; i, ok = i.Cdr.(*Cons) {
			if !ok {
				return rv, ErrISNotAList
			}
			rv, err = i.Car.Format(s, rv)
			if err != nil {
				return
			}
			rv += 1
			if i.Cdr != Onil {
				s.Write([]byte(" "))
			}
		}
		s.Write([]byte(")"))
		return
	}

	obj := o
	s.Write([]byte("("))
	lv += 1

	if _, ok = obj.Car.(*Symbol); ok {
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

	for ; ; obj, ok = obj.Cdr.(*Cons) {
		switch {
		case !ok:
			s.Write([]byte(" . "))
			lv += 3

			rv, err = obj.Cdr.Format(s, lv)
			if err != nil {
				return
			}

			s.Write([]byte(")"))
			rv += 1
			return
		case obj == Onil:
			s.Write([]byte(")"))
			rv = lv + 1
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

	return
}

func (o *Cons) Iter(f func(obj SchemeObject) bool) (err error) {
	ok := true
	for i := o; i != Onil; i, ok = i.Cdr.(*Cons) {
		if !ok {
			return ErrISNotAList
		}
		if !f(i.Car) {
			return
		}
	}
	return
}

func (o *Cons) GetN(n int) (r SchemeObject, err error) {
	var ok bool
	c := o
	for i := 0; i < n; i++ {
		switch c.Cdr {
		case nil:
			return nil, ErrRuntimeUnknown
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
	err = o.Iter(func(obj SchemeObject) bool {
		_, yes = obj.(*Cons)
		return !yes
	})
	return
}

func ListFromSlice(s []SchemeObject) (o SchemeObject) {
	o = Onil
	for i := len(s) - 1; i >= 0; i-- {
		o = &Cons{Car: s[i], Cdr: o}
	}
	return o
}

// type Function struct {
// }

// func (f *Function) IsApplicativeOrder() bool {
// 	return true
// }

// func (f *Function) Eval(env *Environ, p Frame) (r SchemeObject, f Frame, err error) {

// }

// func (f *Function) Format(s io.Writer, lv int) (rv int, err error) {

// }
