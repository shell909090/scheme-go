package impl

import "github.com/shell909090/scheme-go/scm"

func anyFloat(i *scm.Cons) (yes bool, err error) {
	err = i.Iter(func(obj scm.Obj) (e error) {
		switch obj.(type) {
		case scm.Float:
			yes = true
			e = scm.ErrQuit
		case scm.Integer:
		default:
			e = scm.ErrType
		}
		return
	}, false)
	if err == scm.ErrQuit {
		err = nil
	}
	return
}

func ObjToFloat(i scm.Obj) (f float64) {
	switch n := i.(type) {
	case scm.Integer:
		return float64(int(n))
	case scm.Float:
		return float64(n)
	}
	return
}

func IsNumber(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	switch o.Car.(type) {
	case scm.Integer, scm.Float:
		return scm.Otrue, nil, nil
	}
	return scm.Ofalse, nil, nil
}

func IsInteger(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	switch o.Car.(type) {
	case scm.Integer:
		return scm.Otrue, nil, nil
	}
	return scm.Ofalse, nil, nil
}

func Add(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	any, err := anyFloat(o)
	if err != nil {
		return
	}

	if any {
		var s float64
		err = o.Iter(func(obj scm.Obj) (e error) {
			s += ObjToFloat(obj)
			return
		}, false)
		if err != nil {
			return
		}
		value = scm.Float(s)
	} else {
		var s int
		err = o.Iter(func(obj scm.Obj) (e error) {
			s += int(obj.(scm.Integer))
			return
		}, false)
		if err != nil {
			return
		}
		value = scm.Integer(s)
	}
	return
}

func Dec(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	var t scm.Obj
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

		err = o.Iter(func(obj scm.Obj) (e error) {
			s -= ObjToFloat(obj)
			return
		}, false)
		if err != nil {
			return
		}
		value = scm.Float(s)
	} else {
		t, o, err = o.Pop()
		if err != nil {
			return
		}
		s := int(t.(scm.Integer))

		err = o.Iter(func(obj scm.Obj) (e error) {
			s -= int(obj.(scm.Integer))
			return
		}, false)
		if err != nil {
			return
		}
		value = scm.Integer(s)
	}

	return
}

func Mul(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	var t scm.Obj
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

		err = o.Iter(func(obj scm.Obj) (e error) {
			s *= ObjToFloat(obj)
			return
		}, false)
		if err != nil {
			return
		}
		value = scm.Float(s)
	} else {
		t, o, err = o.Pop()
		if err != nil {
			return
		}
		s := int(t.(scm.Integer))

		err = o.Iter(func(obj scm.Obj) (e error) {
			s *= int(obj.(scm.Integer))
			return
		}, false)
		if err != nil {
			return
		}
		value = scm.Integer(s)
	}
	return
}

func Div(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	var t scm.Obj
	t, o, err = o.Pop()
	if err != nil {
		return
	}
	s := ObjToFloat(t)

	err = o.Iter(func(obj scm.Obj) (e error) {
		s /= ObjToFloat(obj)
		return
	}, false)
	if err != nil {
		return
	}
	value = scm.Float(s)
	return
}

func init() {
	scm.RegisterInternalProcedure("number?", IsNumber, true)
	scm.RegisterInternalProcedure("integer?", IsInteger, true)
	// zero? positive? negative?

	scm.RegisterInternalProcedure("+", Add, true)
	scm.RegisterInternalProcedure("-", Dec, true)
	scm.RegisterInternalProcedure("*", Mul, true)
	scm.RegisterInternalProcedure("/", Div, true)

	// = != < > >= <=
	// remainder
	// max min
}
