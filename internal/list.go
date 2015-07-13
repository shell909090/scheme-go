package internal

import "bitbucket.org/shell909090/scheme-go/scmgo"

func MakeCons(o *scmgo.Cons, f scmgo.Frame) (value scmgo.Obj, next scmgo.Frame, err error) {
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

	value = &scmgo.Cons{Car: t1, Cdr: t2}
	return
}

func Car(o *scmgo.Cons, f scmgo.Frame) (value scmgo.Obj, next scmgo.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	t, ok := o.Car.(*scmgo.Cons)
	if !ok {
		return nil, nil, scmgo.ErrType
	}

	value = t.Car
	return
}

func Cdr(o *scmgo.Cons, f scmgo.Frame) (value scmgo.Obj, next scmgo.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	t, ok := o.Car.(*scmgo.Cons)
	if !ok {
		return nil, nil, scmgo.ErrType
	}

	value = t.Cdr
	return
}

func IsNull(o *scmgo.Cons, f scmgo.Frame) (value scmgo.Obj, next scmgo.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	t, ok := o.Car.(*scmgo.Cons)
	if !ok {
		return nil, nil, scmgo.ErrType
	}

	if t == scmgo.Onil {
		return scmgo.Otrue, nil, nil
	}
	return scmgo.Ofalse, nil, nil
}

func init() {
	// list cons
	scmgo.RegisterInternalProcedure("cons", MakeCons, true)
	// null? pair?
	scmgo.RegisterInternalProcedure("null?", IsNull, true)
	// car cdr
	scmgo.RegisterInternalProcedure("car", Car, true)
	scmgo.RegisterInternalProcedure("cdr", Cdr, true)
	// caar cadr cdar cddr
	// append
	// map filter reduce

}
