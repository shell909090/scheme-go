package impl

import (
	"fmt"

	"bitbucket.org/shell909090/scheme-go/scm"
)

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
	// symbol? string?

	scm.RegisterInternalProcedure("define", Define, false)
	scm.RegisterInternalProcedure("lambda", Lambda, false)
	// begin compile
	// eval apply

	// user-init-environment
	// current-environment
	// import

	// let let*
	// eq? equal?
	// not and or
	// cond if when
	scm.RegisterInternalProcedure("if", If, false)

	scm.RegisterInternalProcedure("display", Display, true)
	scm.RegisterInternalProcedure("newline", Newline, true)
	// exit
}
