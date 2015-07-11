package internal

import "bitbucket.org/shell909090/scheme-go/scmgo"

func anyFloat(i *scmgo.Cons) (yes bool, err error) {
	err = i.Iter(func(obj scmgo.SchemeObject) (e error) {
		switch obj.(type) {
		case scmgo.Float:
			yes = true
			e = ErrQuit
		case scmgo.Integer:
		default:
			e = scmgo.ErrType
		}
		return
	}, false)
	if err == ErrQuit {
		err = nil
	}
	return
}

func ObjToFloat(i scmgo.SchemeObject) (f float64) {
	switch n := i.(type) {
	case scmgo.Integer:
		return float64(int(n))
	case scmgo.Float:
		return float64(n)
	}
	return
}

func IsNumber(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	switch o.Car.(type) {
	case scmgo.Integer, scmgo.Float:
		return scmgo.Otrue, nil, nil
	}
	return scmgo.Ofalse, nil, nil
}

func IsInteger(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	switch o.Car.(type) {
	case scmgo.Integer:
		return scmgo.Otrue, nil, nil
	}
	return scmgo.Ofalse, nil, nil
}

func Add(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
	any, err := anyFloat(o)
	if err != nil {
		return
	}

	if any {
		var s float64
		err = o.Iter(func(obj scmgo.SchemeObject) (e error) {
			s += ObjToFloat(obj)
			return
		}, false)
		if err != nil {
			return
		}
		value = scmgo.Float(s)
	} else {
		var s int
		err = o.Iter(func(obj scmgo.SchemeObject) (e error) {
			s += int(obj.(scmgo.Integer))
			return
		}, false)
		if err != nil {
			return
		}
		value = scmgo.Integer(s)
	}
	return
}

func Dec(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
	var t scmgo.SchemeObject
	any, err := anyFloat(o)
	if err != nil {
		return
	}

	if any {
		t, o, err = o.Pop()
		if err != nil {
			return
		}
		s := ObjToFloat(t)

		err = o.Iter(func(obj scmgo.SchemeObject) (e error) {
			s -= ObjToFloat(obj)
			return
		}, false)
		if err != nil {
			return
		}
		value = scmgo.Float(s)
	} else {
		t, o, err = o.Pop()
		if err != nil {
			return
		}
		s := int(t.(scmgo.Integer))

		err = o.Iter(func(obj scmgo.SchemeObject) (e error) {
			s -= int(obj.(scmgo.Integer))
			return
		}, false)
		if err != nil {
			return
		}
		value = scmgo.Integer(s)
	}

	return
}

func Mul(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
	var t scmgo.SchemeObject
	any, err := anyFloat(o)
	if err != nil {
		return
	}

	if any {
		t, o, err = o.Pop()
		if err != nil {
			return
		}
		s := ObjToFloat(t)

		err = o.Iter(func(obj scmgo.SchemeObject) (e error) {
			s *= ObjToFloat(obj)
			return
		}, false)
		if err != nil {
			return
		}
		value = scmgo.Float(s)
	} else {
		t, o, err = o.Pop()
		if err != nil {
			return
		}
		s := int(t.(scmgo.Integer))

		err = o.Iter(func(obj scmgo.SchemeObject) (e error) {
			s *= int(obj.(scmgo.Integer))
			return
		}, false)
		if err != nil {
			return
		}
		value = scmgo.Integer(s)
	}
	return
}

func Div(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
	var t scmgo.SchemeObject
	t, o, err = o.Pop()
	if err != nil {
		return
	}
	s := ObjToFloat(t)

	err = o.Iter(func(obj scmgo.SchemeObject) (e error) {
		s /= ObjToFloat(obj)
		return
	}, false)
	if err != nil {
		return
	}
	value = scmgo.Float(s)
	return
}

func init() {
	scmgo.RegisterInternalProcedure("number?", IsNumber, true)
	scmgo.RegisterInternalProcedure("integer?", IsInteger, true)
	// zero? positive? negative?

	scmgo.RegisterInternalProcedure("+", Add, true)
	scmgo.RegisterInternalProcedure("-", Dec, true)
	scmgo.RegisterInternalProcedure("*", Mul, true)
	scmgo.RegisterInternalProcedure("/", Div, true)

	// = != < > >= <=
	// remainder
	// max min
}
