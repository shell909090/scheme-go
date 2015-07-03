package scmgo

func anyFloat(i *Cons) (yes bool, err error) {
	err = i.Iter(func(obj SchemeObject) (e error) {
		switch obj.(type) {
		case Float:
			yes = true
			e = ErrUnknown
		case Integer:
		default:
			e = ErrType
		}
		return
	})
	if err == ErrUnknown {
		err = nil
	}
	return
}

func ObjToFloat(i SchemeObject) (f float64) {
	switch n := i.(type) {
	case Integer:
		return float64(int(n))
	case Float:
		return float64(n)
	}
	return
}

func IsNumber(o *Cons, p *ApplyFrame) (r SchemeObject, next Frame, err error) {
	n, err := o.Len()
	if err != nil {
		return
	}
	if n != 1 {
		return nil, nil, ErrArguments
	}
	switch o.Car.(type) {
	case Integer, Float:
		return Otrue, nil, nil
	}
	return Ofalse, nil, nil
}

func IsInteger(o *Cons, p *ApplyFrame) (r SchemeObject, next Frame, err error) {
	n, err := o.Len()
	if err != nil {
		return
	}
	if n != 1 {
		return nil, nil, ErrArguments
	}
	switch o.Car.(type) {
	case Integer:
		return Otrue, nil, nil
	}
	return Ofalse, nil, nil
}

func Add(o *Cons, p *ApplyFrame) (r SchemeObject, next Frame, err error) {
	f, err := anyFloat(o)
	if err != nil {
		return
	}

	if f {
		var s float64
		err = o.Iter(func(obj SchemeObject) (e error) {
			s += ObjToFloat(obj)
			return
		})
		if err != nil {
			return
		}
		r = Float(s)
	} else {
		var s int
		err = o.Iter(func(obj SchemeObject) (e error) {
			s += int(obj.(Integer))
			return
		})
		if err != nil {
			return
		}
		r = Integer(s)
	}
	return
}

func Dec(o *Cons, p *ApplyFrame) (r SchemeObject, next Frame, err error) {
	var t SchemeObject
	f, err := anyFloat(o)
	if err != nil {
		return
	}

	if f {
		t, o, err = o.Pop()
		if err != nil {
			return
		}
		s := ObjToFloat(t)

		err = o.Iter(func(obj SchemeObject) (e error) {
			s -= ObjToFloat(obj)
			return
		})
		if err != nil {
			return
		}
		r = Float(s)
	} else {
		t, o, err = o.Pop()
		if err != nil {
			return
		}
		s := int(t.(Integer))

		err = o.Iter(func(obj SchemeObject) (e error) {
			s -= int(obj.(Integer))
			return
		})
		if err != nil {
			return
		}
		r = Integer(s)
	}

	return
}

func Mul(o *Cons, p *ApplyFrame) (r SchemeObject, next Frame, err error) {
	var t SchemeObject
	f, err := anyFloat(o)
	if err != nil {
		return
	}

	if f {
		t, o, err = o.Pop()
		if err != nil {
			return
		}
		s := ObjToFloat(t)

		err = o.Iter(func(obj SchemeObject) (e error) {
			s *= ObjToFloat(obj)
			return
		})
		if err != nil {
			return
		}
		r = Float(s)
	} else {
		t, o, err = o.Pop()
		if err != nil {
			return
		}
		s := int(t.(Integer))

		err = o.Iter(func(obj SchemeObject) (e error) {
			s *= int(obj.(Integer))
			return
		})
		if err != nil {
			return
		}
		r = Integer(s)
	}
	return
}

func Div(o *Cons, p *ApplyFrame) (r SchemeObject, next Frame, err error) {
	var t SchemeObject
	t, o, err = o.Pop()
	if err != nil {
		return
	}
	s := ObjToFloat(t)

	err = o.Iter(func(obj SchemeObject) (e error) {
		s /= ObjToFloat(obj)
		return
	})
	if err != nil {
		return
	}
	r = Float(s)
	return
}

func init() {
	RegisterInternalProcedure("number?", IsNumber, true)
	RegisterInternalProcedure("integer?", IsInteger, true)
	// zero? positive? negative?

	RegisterInternalProcedure("+", Add, true)
	RegisterInternalProcedure("-", Dec, true)
	RegisterInternalProcedure("*", Mul, true)
	RegisterInternalProcedure("/", Div, true)

	// = != < > >= <=
	// remainder
	// max min
}
