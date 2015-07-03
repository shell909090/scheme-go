package scmgo

import "fmt"

func Define(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	t, o, err := o.Pop()
	if err != nil {
		return
	}

	args, ok := t.(*Cons)
	if !ok {
		return nil, nil, ErrType
	}

	name, err := GetHeadAsSymbol(args)
	if err != nil {
		return
	}
	_, args, err = args.Pop()
	if err != nil {
		return
	}

	r = &LambdaProcedure{
		Name: name.Name,
		Env:  p.GetEnv(),
		Args: args,
		Obj:  o,
	}

	p.GetEnv().Add(name.Name, r)
	return
}

func Lambda(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	t, o, err := o.Pop()
	if err != nil {
		return
	}

	args, ok := t.(*Cons)
	if !ok {
		return nil, nil, ErrType
	}

	r = &LambdaProcedure{
		Env:  p.GetEnv(),
		Args: args,
		Obj:  o,
	}
	return
}

func Cond(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	// this is called by apply frame, pass it.
	next = CreateCondFrame(p.GetParent(), o, p.GetEnv())
	return
}

func If(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	// TODO:
	return
}

func Display(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	n, err := o.Len()
	if err != nil {
		return
	}
	if n != 1 {
		return nil, nil, ErrArguments
	}

	fmt.Printf("%s", SchemeObjectToString(o.Car))
	return Onil, nil, nil
}

func Newline(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	fmt.Printf("\n")
	return Onil, nil, nil
}

func init() {
	// symbol? string?

	RegisterInternalProcedure("define", Define, false)
	RegisterInternalProcedure("lambda", Lambda, false)
	// begin compile
	// eval apply

	// user-init-environment
	// current-environment
	// import

	// let let*
	// eq? equal?
	// not and or
	// cond if when
	RegisterInternalProcedure("cond", Cond, false)

	RegisterInternalProcedure("display", Display, true)
	RegisterInternalProcedure("error", Display, true)
	RegisterInternalProcedure("newline", Newline, true)
	// exit
}
