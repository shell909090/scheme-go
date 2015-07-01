package scmgo

import (
	"io"
	"strconv"
)

type Evalor interface {
	Eval(stack *Stack, env *Environ) (r SchemeObject, err error)
}

type Formatter interface {
	Format(s io.Writer, lv int) (err error)
}

type SchemeObject interface {
	Evalor
	Formatter
}

type Nil struct{}

var Onil = &Nil{}

func (o *Nil) Eval(stack *Stack, env *Environ) (r SchemeObject, err error) {
	return nil, nil
}

func (o *Nil) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte("()"))
	return nil
}

type Cons struct {
	Car SchemeObject
	Cdr SchemeObject
}

func (o *Cons) Eval(stack *Stack, env *Environ) (r SchemeObject, err error) {
	// f, err := stack.Eval(o.Car, env)
	// if err != nil {
	// 	return
	// }
	// Apply(f, o.Cdr)
	return nil, nil
}

func (o *Cons) Format(s io.Writer, lv int) (err error) {
	anycons := o.anyCons()

	obj := o
	s.Write([]byte("("))
	obj.Car.Format(s, lv+1)

	if u, ok := o.Car.(*Symbol); anycons && ok {
		switch u.Name {
		case "if":
			lv += 3
			obj, ok = obj.Cdr.(*Cons)
			if !ok {
				panic("")
			}
			s.Write([]byte(" "))
			obj.Car.Format(s, lv+4)
		case "define", "lambda":
			lv += 1
			obj, ok = obj.Cdr.(*Cons)
			if !ok {
				panic("")
			}
			s.Write([]byte(" "))
			obj.Car.Format(s, lv+6)
		}
	}

	for {
		switch u := obj.Cdr.(type) {
		case *Nil:
			s.Write([]byte(")"))
			return
		case *Cons:
			obj = u
			if anycons {
				s.Write([]byte("\n"))
				for i := 0; i < lv; i++ {
					s.Write([]byte(" "))
				}
			} else {
				s.Write([]byte(" "))
			}
			obj.Car.Format(s, lv+1)
		default:
			s.Write([]byte(" . "))
			obj.Car.Format(s, lv+1)
			return
		}
	}
	return
}

func (o *Cons) Iter(f func(obj SchemeObject) bool) {
	for i := o; ; {
		if f(i.Car) {
			return
		}
		switch t := i.Cdr.(type) {
		case *Cons:
			i = t
		case *Nil:
			return
		default:
			f(t)
			return
		}
	}
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

func (o *Cons) anyCons() (yes bool) {
	o.Iter(func(obj SchemeObject) bool {
		if _, yes = obj.(*Cons); yes {
			return true
		}
		return false
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

type Symbol struct {
	Name string
}

func (o *Symbol) Eval(stack *Stack, env *Environ) (r SchemeObject, err error) {
	return nil, nil
}

func (o *Symbol) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte(o.Name))
	return nil
}

func SymbolFromString(s string) (o *Symbol) {
	return &Symbol{Name: s}
}

type Quote struct {
	objs SchemeObject
}

func (o *Quote) Eval(stack *Stack, env *Environ) (r SchemeObject, err error) {
	return nil, nil
}

func (o *Quote) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte("'"))
	o.objs.Format(s, lv)
	return nil
}

type Boolean bool

func (o Boolean) Eval(stack *Stack, env *Environ) (r SchemeObject, err error) {
	return nil, nil
}

func (o Boolean) Format(s io.Writer, lv int) (err error) {
	if o {
		s.Write([]byte("#t"))
	} else {
		s.Write([]byte("#f"))
	}
	return nil
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

func (o Integer) Eval(stack *Stack, env *Environ) (r SchemeObject, err error) {
	return nil, nil
}

func (o Integer) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte(strconv.FormatInt(int64(o), 10)))
	return nil
}

func IntegerFromString(s string) (o Integer, err error) {
	i, err := strconv.Atoi(s)
	return Integer(i), err
}

type Float float64

func (o Float) Eval(stack *Stack, env *Environ) (r SchemeObject, err error) {
	return nil, nil
}

func (o Float) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte(strconv.FormatFloat(float64(o), 'f', 2, 64)))
	return nil
}

func FloatFromString(s string) (o Float, err error) {
	i, err := strconv.ParseFloat(s, 64)
	return Float(i), err
}

type String string

func (o String) Eval(stack *Stack, env *Environ) (r SchemeObject, err error) {
	return nil, nil
}

func (o String) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte("\""))
	s.Write([]byte(o))
	s.Write([]byte("\""))
	return nil
}
