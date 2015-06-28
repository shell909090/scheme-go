package scmgo

import (
	"io"
	"strconv"
)

type SchemeObject interface {
	Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject)
	Format(s io.Writer, lv int) (err error)
}

type Nil struct{}

var Onil = &Nil{}

func (o *Nil) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Nil) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte("()"))
	return nil
}

type Cons struct {
	car SchemeObject
	cdr SchemeObject
}

func (o *Cons) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Cons) Format(s io.Writer, lv int) (err error) {
	anycons := o.anyCons()

	obj := o
	s.Write([]byte("("))
	obj.car.Format(s, lv+1)

	if u, ok := o.car.(*Symbol); anycons && ok {
		switch u.name {
		case "if":
			lv += 3
			obj, ok = obj.cdr.(*Cons)
			if !ok {
				panic("")
			}
			s.Write([]byte(" "))
			obj.car.Format(s, lv+4)
		case "define", "lambda":
			lv += 1
			obj, ok = obj.cdr.(*Cons)
			if !ok {
				panic("")
			}
			s.Write([]byte(" "))
			obj.car.Format(s, lv+6)
		}
	}

	for {
		switch u := obj.cdr.(type) {
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
			obj.car.Format(s, lv+1)
		default:
			s.Write([]byte(" . "))
			obj.car.Format(s, lv+1)
			return
		}
	}
	return
}

func (o *Cons) Iter(f func(obj SchemeObject) bool) {
	for i := o; ; {
		if f(i.car) {
			return
		}
		switch t := i.cdr.(type) {
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
		o = &Cons{car: s[i], cdr: o}
	}
	return o
}

type Symbol struct {
	name string
}

func (o *Symbol) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Symbol) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte(o.name))
	return nil
}

func SymbolFromString(s string) (o *Symbol) {
	return &Symbol{name: s}
}

type Quote struct {
	objs SchemeObject
}

func (o *Quote) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Quote) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte("'"))
	o.objs.Format(s, lv)
	return nil
}

type Boolean struct {
	b bool
}

func (o *Boolean) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Boolean) Format(s io.Writer, lv int) (err error) {
	if o.b {
		s.Write([]byte("#t"))
	} else {
		s.Write([]byte("#f"))
	}
	return nil
}

var (
	Otrue  = &Boolean{b: true}
	Ofalse = &Boolean{b: false}
)

func BooleanFromString(s string) (o *Boolean, err error) {
	// FIXME: not so good
	switch s[1] {
	case 't':
		return Otrue, nil
	case 'f':
		return Ofalse, nil
	default:
		return nil, ErrBooleanUnknown
	}
}

type Integer struct {
	num int
}

func (o *Integer) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Integer) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte(strconv.FormatInt(int64(o.num), 10)))
	return nil
}

func IntegerFromString(s string) (o *Integer, err error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return
	}
	return &Integer{num: i}, nil
}

type Float struct {
	num float64
}

func (o *Float) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Float) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte(strconv.FormatFloat(o.num, 'f', 2, 64)))
	return nil
}

func FloatFromString(s string) (o *Float, err error) {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return
	}
	return &Float{num: i}, nil
}

type String struct {
	str string
}

func (o *String) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *String) Format(s io.Writer, lv int) (err error) {
	s.Write([]byte("\""))
	s.Write([]byte(o.str))
	s.Write([]byte("\""))
	return nil
}

func StringFromString(s string) (o *String) {
	return &String{str: s}
}
