package internal

import (
	"fmt"

	"bitbucket.org/shell909090/scheme-go/scmgo"
)

func Define(o *scmgo.Cons, p scmgo.Frame) (r scmgo.SchemeObject, next scmgo.Frame, err error) {
	t, o, err := o.Pop()
	if err != nil {
		return
	}

	args, ok := t.(*scmgo.Cons)
	if !ok {
		return nil, nil, scmgo.ErrType
	}

	name, err := scmgo.GetHeadAsSymbol(args)
	if err != nil {
		return
	}
	_, args, err = args.Pop()
	if err != nil {
		return
	}

	r = &scmgo.LambdaProcedure{
		Name: name.Name,
		Env:  p.GetEnv(),
		Args: args,
		Obj:  o,
	}

	p.GetEnv().Add(name.Name, r)
	return
}

func Lambda(o *scmgo.Cons, p scmgo.Frame) (r scmgo.SchemeObject, next scmgo.Frame, err error) {
	t, o, err := o.Pop()
	if err != nil {
		return
	}

	args, ok := t.(*scmgo.Cons)
	if !ok {
		return nil, nil, scmgo.ErrType
	}

	r = &scmgo.LambdaProcedure{
		Env:  p.GetEnv(),
		Args: args,
		Obj:  o,
	}
	return
}

func Cond(o *scmgo.Cons, p scmgo.Frame) (r scmgo.SchemeObject, next scmgo.Frame, err error) {
	// coming from apply, so pass this frame.
	next = scmgo.CreateCondFrame(o, p.GetEnv(), p.GetParent())
	return
}

func If(o *scmgo.Cons, p scmgo.Frame) (r scmgo.SchemeObject, next scmgo.Frame, err error) {
	// TODO:
	return
}

func Display(o *scmgo.Cons, p scmgo.Frame) (r scmgo.SchemeObject, next scmgo.Frame, err error) {
	n, err := o.Len()
	if err != nil {
		return
	}
	if n != 1 {
		return nil, nil, ErrArguments
	}

	fmt.Printf("%s", o.Car.Format())
	return scmgo.Onil, nil, nil
}

func Newline(o *scmgo.Cons, p scmgo.Frame) (r scmgo.SchemeObject, next scmgo.Frame, err error) {
	fmt.Printf("\n")
	return scmgo.Onil, nil, nil
}

func init() {
	// symbol? string?

	scmgo.RegisterInternalProcedure("define", Define, false)
	scmgo.RegisterInternalProcedure("lambda", Lambda, false)
	// begin compile
	// eval apply

	// user-init-environment
	// current-environment
	// import

	// let let*
	// eq? equal?
	// not and or
	// cond if when
	scmgo.RegisterInternalProcedure("cond", Cond, false)

	scmgo.RegisterInternalProcedure("display", Display, true)
	scmgo.RegisterInternalProcedure("error", Display, true)
	scmgo.RegisterInternalProcedure("newline", Newline, true)
	// exit
}
