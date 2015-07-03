package internal

import "bitbucket.org/shell909090/scheme-go/scmgo"

func MakeCons(o *scmgo.Cons, p scmgo.Frame) (r scmgo.SchemeObject, next scmgo.Frame, err error) {
	n, err := o.Len()
	if err != nil {
		return
	}
	if n != 2 {
		return nil, nil, ErrArguments
	}

	t1, o, err := o.Pop()
	if err != nil {
		return
	}

	t2, o, err := o.Pop()
	if err != nil {
		return
	}

	r = &scmgo.Cons{Car: t1, Cdr: t2}
	return
}

func Car(o *scmgo.Cons, p scmgo.Frame) (r scmgo.SchemeObject, next scmgo.Frame, err error) {
	n, err := o.Len()
	if err != nil {
		return
	}
	if n != 1 {
		return nil, nil, ErrArguments
	}

	t, ok := o.Car.(*scmgo.Cons)
	if !ok {
		return nil, nil, scmgo.ErrType
	}

	r = t.Car
	return
}

func Cdr(o *scmgo.Cons, p scmgo.Frame) (r scmgo.SchemeObject, next scmgo.Frame, err error) {
	n, err := o.Len()
	if err != nil {
		return
	}
	if n != 1 {
		return nil, nil, ErrArguments
	}

	t, ok := o.Car.(*scmgo.Cons)
	if !ok {
		return nil, nil, scmgo.ErrType
	}

	r = t.Cdr
	return
}

func IsNull(o *scmgo.Cons, p scmgo.Frame) (r scmgo.SchemeObject, next scmgo.Frame, err error) {
	n, err := o.Len()
	if err != nil {
		return
	}
	if n != 1 {
		return nil, nil, ErrArguments
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
