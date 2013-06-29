package scmgo

import (
	"io"
	"strconv"
)

type SchemeObject interface {
	Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject)
	Format(s io.Writer) (err error)
}

type Nil struct {}

var Onil = &Nil{}

func (o *Nil) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Nil) Format(s io.Writer) (err error) {
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

func (o *Cons) Format(s io.Writer) (err error) {
	obj := o
	s.Write([]byte("("))
	obj.car.Format(s)
	for {
		switch u := obj.cdr.(type) {
		case *Nil:
			s.Write([]byte(")"))
			return
		case *Cons:
			obj = u
			s.Write([]byte(" "))
			obj.car.Format(s)
		default:
			s.Write([]byte(" . "))
			obj.car.Format(s)
			return
		}
	}
	return
}

type Symbol struct {
	name string
}

func (o *Symbol) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Symbol) Format(s io.Writer) (err error) {
	s.Write([]byte(o.name))
	return nil
}

type Quote struct {
	objs SchemeObject
}

func (o *Quote) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Quote) Format(s io.Writer) (err error) {
	s.Write([]byte("'"))
	o.objs.Format(s)
	return nil
}

type Boolean struct {
	b bool
}

func (o *Boolean) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Boolean) Format(s io.Writer) (err error) {
	if o.b {
		s.Write([]byte("#t"))
	} else {
		s.Write([]byte("#f"))
	}
	return nil
}

var (
	Otrue = &Boolean{b: true}
	Ofalse = &Boolean{b: false}
)

type Integer struct {
	num int
}

func (o *Integer) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Integer) Format(s io.Writer) (err error) {
	s.Write([]byte(strconv.FormatInt(int64(o.num), 10)))
	return nil
}

func Integer_from_string(s string) (o *Integer, err error) {
	i, err := strconv.Atoi(string(s))
	if err != nil { return }
	return &Integer{num: i}, nil
}

type Float struct {
	num float64
}

func (o *Float) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *Float) Format(s io.Writer) (err error) {
	s.Write([]byte(strconv.FormatFloat(o.num, 'f', 2, 64)))
	return nil
}

func Float_from_string(s string) (o *Float, err error) {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil { return }
	return &Float{num: i}, nil
}

type String struct {
	str string
}

func (o *String) Exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

func (o *String) Format(s io.Writer) (err error) {
	s.Write([]byte("\""))
	s.Write([]byte(o.str))
	s.Write([]byte("\""))
	return nil
}
