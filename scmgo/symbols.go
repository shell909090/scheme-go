package scmgo

// symbol? string?

// define lambda
// begin compile
// eval apply

// user-init-environment
// current-environment
// import

// let let*
// eq? equal?
// not and or
// cond if when

// display error
// newline
// exit

func Lambda(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	t, o, err := o.Pop()
	if err != nil {
		return
	}

	args, ok := t.(*Cons)
	if ok {
		return nil, nil, ErrType
	}

	r = &LambdaProcedure{
		Env:  p.GetEnv(),
		Args: args,
		Obj:  o,
	}
	return
}
