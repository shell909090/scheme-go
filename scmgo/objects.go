package scmgo

import (
	"bytes"
	"strconv"
)

type Evalor interface {
	Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error)
}

type Formatter interface {
	Format() (r string)
}

type SchemeObject interface {
	Evalor
	Formatter
}

type Symbol struct {
	Name string
}

func (o *Symbol) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
	value = env.Get(o.Name)
	if value == nil {
		return nil, nil, ErrNameNotFound
	}
	return
}

func (o *Symbol) Format() (r string) {
	return o.Name
}

type Quote struct {
	Objs SchemeObject
}

func (o *Quote) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
	value = o.Objs
	return
}

func (o *Quote) Format() (r string) {
	return "'" + o.Objs.Format()
}

type Boolean bool

func (o Boolean) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
	value = o
	return
}

func (o Boolean) Format() (r string) {
	if o {
		return "#t"
	} else {
		return "#f"
	}
}

const (
	Otrue  = Boolean(true)
	Ofalse = Boolean(false)
)

type Integer int

func (o Integer) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
	value = o
	return
}

func (o Integer) Format() (r string) {
	return strconv.FormatInt(int64(o), 10)
}

type Float float64

func (o Float) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
	value = o
	return
}

func (o Float) Format() (r string) {
	return strconv.FormatFloat(float64(o), 'f', 2, 64)
}

type String string

func (o String) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
	value = o
	return
}

func (o String) Format() (r string) {
	return "\"" + string(o) + "\""
}

type Cons struct {
	Car SchemeObject
	Cdr SchemeObject
}

var Onil = &Cons{}

func (o *Cons) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
	var procedure SchemeObject
	procedure, o, err = o.Pop()
	if err != nil {
		return
	}

	next = CreateApplyFrame(o, env, f) // not sure about procedure yet.
	f = next

	// get a result now, or get a frame which can return in future.
	procedure, next, err = procedure.Eval(env, next)
	if err != nil {
		return
	}
	if next != nil {
		return
	}
	// get return immediately
	next = f
	err = next.Return(procedure)
	return
}

func (o *Cons) Format() (r string) {
	buf := bytes.NewBuffer(nil)
	if _, err := PrettyFormat(buf, o, 0); err != nil {
		log.Error("%s", err)
		return ""
	}
	return buf.String()
}

func (o *Cons) Pop() (r SchemeObject, next *Cons, err error) {
	if o == Onil {
		return nil, nil, ErrListOutOfIndex
	}
	r = o.Car
	next, ok := o.Cdr.(*Cons)
	if !ok {
		return nil, nil, ErrISNotAList
	}
	return
} // O(1)

func (o *Cons) Push(i SchemeObject) (next *Cons) {
	return &Cons{Car: i, Cdr: o}
} // O(1)

func (o *Cons) IsImproper() bool {
	ok := true
	for i := o; i != Onil; i, ok = i.Cdr.(*Cons) {
		if !ok {
			return true
		}
	}
	return false
} // O(n)

func (o *Cons) Iter(f func(obj SchemeObject) (e error), improper bool) (err error) {
	ok := true
	for i := o; i != Onil; {
		err = f(i.Car)
		if err != nil {
			return
		}
		i, ok = i.Cdr.(*Cons)
		if !ok {
			if !improper {
				return ErrISNotAList
			}
			return f(i.Cdr)
		}
	}
	return
} // O(n)

func (o *Cons) Len(improper bool) (n int, err error) {
	err = o.Iter(func(obj SchemeObject) (e error) {
		n += 1
		return
	}, improper)
	return
} // O(n)

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
} // O(n)

func (o *Cons) PopSymbol() (s *Symbol, next *Cons, err error) {
	t, next, err := o.Pop()
	if err != nil {
		return
	}
	s, ok := t.(*Symbol)
	if !ok {
		return nil, nil, ErrType
	}
	return
}

func (o *Cons) PopCons() (s *Cons, next *Cons, err error) {
	t, next, err := o.Pop()
	if err != nil {
		return
	}
	s, ok := t.(*Cons)
	if !ok {
		return nil, nil, ErrType
	}
	return
}
