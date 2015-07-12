package internal

import (
	"fmt"

	"bitbucket.org/shell909090/scheme-go/scmgo"
)

func Define(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
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

	value = &scmgo.LambdaProcedure{
		Name: name.Name,
		Env:  f.GetEnv(),
		Args: args,
		Obj:  o,
	}

	f.GetEnv().Add(name.Name, value)
	return
}

func Lambda(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
	args, o, err := o.PopCons()
	if err != nil {
		log.Error("%s", err.Error())
		return
	}

	value = &scmgo.LambdaProcedure{
		Env:  f.GetEnv(),
		Args: args,
		Obj:  o,
	}
	return
}

// func Cond(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
// 	// coming from apply, so pass this frame.
// 	next = scmgo.CreateCondFrame(o, f.GetEnv(), f.GetParent())
// 	return
// }

func If(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
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
	var ecase scmgo.SchemeObject = scmgo.Onil
	if o != scmgo.Onil {
		ecase, o, err = o.Pop()
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
	}
	next = scmgo.CreateIfFrame(cond, tcase, ecase, f.GetEnv(), f.GetParent())
	return
}

func Display(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
	err = AssertLen(o, 1)
	if err != nil {
		return
	}

	fmt.Printf("%s", o.Car.Format())
	return scmgo.Onil, nil, nil
}

func Newline(o *scmgo.Cons, f scmgo.Frame) (value scmgo.SchemeObject, next scmgo.Frame, err error) {
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
	// scmgo.RegisterInternalProcedure("cond", Cond, false)
	scmgo.RegisterInternalProcedure("if", If, false)

	scmgo.RegisterInternalProcedure("display", Display, true)
	scmgo.RegisterInternalProcedure("error", Display, true)
	scmgo.RegisterInternalProcedure("newline", Newline, true)
	// exit
}
