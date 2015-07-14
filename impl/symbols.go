package impl

import (
	"fmt"

	"bitbucket.org/shell909090/scheme-go/scm"
)

func IsSymbol(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	_, _, err = o.PopSymbol()
	if err == scm.ErrType {
		return scm.Ofalse, nil, nil
	}
	return scm.Otrue, nil, nil
}

func IsString(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	t, _, err := o.Pop()
	if err != nil {
		return
	}
	if _, ok := t.(scm.String); !ok {
		return scm.Ofalse, nil, nil
	}
	return scm.Otrue, nil, nil
}

func Define(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	args, o, err := o.PopCons()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}

	name, args, err := args.PopSymbol()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}

	value = &scm.LambdaProcedure{
		Name: name.Name,
		Env:  f.GetEnv(),
		Args: args,
		Obj:  o,
	}

	f.GetEnv().Add(name.Name, value)
	return
}

func Lambda(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	args, o, err := o.PopCons()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}

	value = &scm.LambdaProcedure{
		Env:  f.GetEnv(),
		Args: args,
		Obj:  o,
	}
	return
}

func Begin(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	return nil, scm.CreateBeginFrame(o, f.GetEnv(), f.GetParent()), nil
}

func Eval(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	a, o, err := o.Pop()
	if err != nil {
		return
	}
	next, err = scm.EvalAndReturn(a, f.GetEnv(), f.GetParent())
	return
}

func Apply(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	p, args, err := o.Pop()
	if err != nil {
		return
	}
	procedure, ok := p.(scm.Procedure)
	if !ok {
		return nil, nil, scm.ErrNotRunnable
	}
	return procedure.Apply(args, f) // apply will realized this is an Apply frame.
}

func If(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	cond, o, err := o.Pop()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	tcase, o, err := o.Pop()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}
	var ecase scm.Obj = scm.Onil
	if o != scm.Onil {
		ecase, o, err = o.Pop()
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
	}
	next = scm.CreateIfFrame(cond, tcase, ecase, f.GetEnv(), f.GetParent())
	return
}

func Display(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	fmt.Printf("%s", scm.Format(o.Car))
	return scm.Onil, nil, nil
}

func Newline(o *scm.Cons, f scm.Frame) (value scm.Obj, next scm.Frame, err error) {
	fmt.Printf("\n")
	return scm.Onil, nil, nil
}

func init() {
	scm.RegisterInternalProcedure("symbol?", IsSymbol, true)
	scm.RegisterInternalProcedure("string?", IsString, true)

	scm.RegisterInternalProcedure("define", Define, false)
	scm.RegisterInternalProcedure("lambda", Lambda, false)
	// begin compile
	scm.RegisterInternalProcedure("begin", Begin, false)
	scm.RegisterInternalProcedure("eval", Eval, true)
	scm.RegisterInternalProcedure("apply", Apply, true)

	// user-init-environment
	// current-environment
	// import

	// eq? equal?
	scm.RegisterInternalProcedure("if", If, false)

	scm.RegisterInternalProcedure("display", Display, true)
	scm.RegisterInternalProcedure("newline", Newline, true)
	// exit
}
