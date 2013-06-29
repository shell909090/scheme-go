package scmgo

import (
	"strconv"
)

type SchemeObject interface {
	exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject)
}

type Nil struct {}

var Onil = &Nil{}

func (o *Nil) exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

type Cons struct {
	car SchemeObject
	cdr SchemeObject
}

func (o *Cons) exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

type Symbol struct {
	name string
}

func (o *Symbol) exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

type Quote struct {
	objs SchemeObject
}

func (o *Quote) exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

type Boolean struct {
	b bool
}

func (o *Boolean) exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

var (
	Otrue = &Boolean{b: true}
	Ofalse = &Boolean{b: false}
)

type Integer struct {
	num int
}

func (o *Integer) exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
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

func Float_from_string(s string) (o *Float, err error) {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil { return }
	return &Float{num: i}, nil
}

func (o *Float) exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}

type String struct {
	str string
}

func (o *String) exec(stack *Stack, env *Env, objs SchemeObject) (rslt SchemeObject) {
	return nil
}
