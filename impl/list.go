package impl

import "github.com/shell909090/scheme-go/scm"

func List(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	return o, nil, nil
}

func MakeCons(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	err = AssertLen(o, 2)
	if err != nil {
		return
	}

	t1, o, err := o.Pop()
	if err != nil {
		return
	}

	t2, o, err := o.Pop()
	if err != nil {
		return
	}

	value = &scm.Cons{Car: t1, Cdr: t2}
	return
}

func Car(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	t, ok := o.Car.(*scm.Cons)
	if !ok {
		return nil, nil, scm.ErrType
	}

	value = t.Car
	return
}

func Cdr(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	t, ok := o.Car.(*scm.Cons)
	if !ok {
		return nil, nil, scm.ErrType
	}

	value = t.Cdr
	return
}

func IsNull(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	t, ok := o.Car.(*scm.Cons)
	if !ok {
		return nil, nil, scm.ErrType
	}

	if t == scm.Onil {
		return scm.Otrue, nil, nil
	}
	return scm.Ofalse, nil, nil
}

func IsPair(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	if _, ok := o.Car.(*scm.Cons); ok {
		return scm.Otrue, nil, nil
	}
	return scm.Ofalse, nil, nil
}

func init() {
	scm.RegisterInternalProcedure("list", List, true)
	scm.RegisterInternalProcedure("cons", MakeCons, true)
	scm.RegisterInternalProcedure("null?", IsNull, true)
	scm.RegisterInternalProcedure("pair?", IsPair, true)
	scm.RegisterInternalProcedure("car", Car, true)
	scm.RegisterInternalProcedure("cdr", Cdr, true)
}
