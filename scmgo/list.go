package scmgo

func MakeCons(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
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

	r = &Cons{Car: t1, Cdr: t2}
	return
}

func Car(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	n, err := o.Len()
	if err != nil {
		return
	}
	if n != 1 {
		return nil, nil, ErrArguments
	}

	t, ok := o.Car.(*Cons)
	if !ok {
		return nil, nil, ErrType
	}

	r = t.Car
	return
}

func Cdr(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	n, err := o.Len()
	if err != nil {
		return
	}
	if n != 1 {
		return nil, nil, ErrArguments
	}

	t, ok := o.Car.(*Cons)
	if !ok {
		return nil, nil, ErrType
	}

	r = t.Cdr
	return
}

func IsNull(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	n, err := o.Len()
	if err != nil {
		return
	}
	if n != 1 {
		return nil, nil, ErrArguments
	}

	t, ok := o.Car.(*Cons)
	if !ok {
		return nil, nil, ErrType
	}

	if t == Onil {
		return Otrue, nil, nil
	}
	return Ofalse, nil, nil
}

func init() {
	// list cons
	RegisterInternalProcedure("cons", MakeCons, true)
	// null? pair?
	RegisterInternalProcedure("null?", IsNull, true)
	// car cdr
	RegisterInternalProcedure("car", Car, true)
	RegisterInternalProcedure("cdr", Cdr, true)
	// caar cadr cdar cddr
	// append
	// map filter reduce

}
