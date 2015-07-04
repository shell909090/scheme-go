package scmgo

import (
	"bytes"
	"strconv"
)

type Evalor interface {
	Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error)
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

func (o *Symbol) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = env.Get(o.Name)
	if r == nil {
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

func (o *Quote) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = o.Objs
	return
}

func (o *Quote) Format() (r string) {
	return "'" + o.Objs.Format()
}

type Comment struct {
	Content string
}

func (c *Comment) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	return nil, nil, ErrUnknown
}

func (c *Comment) Format() (r string) {
	return ";" + c.Content
}

type Boolean bool

func (o Boolean) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = o
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

func (o Integer) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = o
	return
}

func (o Integer) Format() (r string) {
	return strconv.FormatInt(int64(o), 10)
}

type Float float64

func (o Float) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = o
	return
}

func (o Float) Format() (r string) {
	return strconv.FormatFloat(float64(o), 'f', 2, 64)
}

type String string

func (o String) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	r = o
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

func (o *Cons) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	var procedure SchemeObject
	procedure, o, err = o.Pop()
	if err != nil {
		return
	}

	next = CreateApplyFrame(o, env, p) // not sure about procedure yet.
	p = next

	// get a result now, or get a frame which can return in future.
	procedure, next, err = procedure.Eval(env, next)
	if err != nil {
		return
	}
	if next != nil {
		return
	}
	// get return now
	next = p
	err = next.Return(procedure)
	return
}

func (o *Cons) Format() (r string) {
	buf := bytes.NewBuffer(nil)
	_, err := PrettyFormat(buf, o, 0)
	if err != nil {
		log.Error("%s", err)
		return ""
	}
	return buf.String()
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
	if o == Onil {
		return nil, nil, ErrListOutOfIndex
	}
	r = o.Car
	next, ok := o.Cdr.(*Cons)
	if !ok {
		return nil, nil, ErrISNotAList
	}
	return
}

func (o *Cons) Push(i SchemeObject) (next *Cons) {
	return &Cons{Car: i, Cdr: o}
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
