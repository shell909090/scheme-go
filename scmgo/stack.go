package scmgo

import "fmt"

type Frame interface {
	Debug() (r string)
	SetParent(p Frame)
	GetEnv() (e *Environ)
	Exec(i SchemeObject) (r SchemeObject, next Frame, err error)
}

type EvalFrame struct {
	Parent Frame
	Obj    SchemeObject
	Env    *Environ
}

func CreateEvalFrame(p Frame, o SchemeObject, e *Environ) (f Frame) {
	return &EvalFrame{Parent: p, Obj: o, Env: e}
}

func (ef *EvalFrame) Debug() (r string) {
	return fmt.Sprintf("EvalFrame: {%s}", SchemeObjectToString(ef.Obj))
}

func (ef *EvalFrame) SetParent(p Frame) {
	ef.Parent = p
}

func (ef *EvalFrame) GetEnv() (e *Environ) {
	return ef.Env
}

func (ef *EvalFrame) Exec(i SchemeObject) (r SchemeObject, next Frame, err error) {
	r, next, err = ef.Obj.Eval(ef.Env, ef.Parent)
	if next == nil {
		next = ef.Parent
	}
	return
}

type PrognFrame struct {
	Parent Frame
	Obj    *Cons
	Env    *Environ
}

func CreatePrognFrame(p Frame, o *Cons, e *Environ) (f Frame) {
	return &PrognFrame{Parent: p, Obj: o, Env: e}
}

func (pf *PrognFrame) Debug() (r string) {
	return "PrognFrame"
}

func (pf *PrognFrame) SetParent(p Frame) {
	pf.Parent = p
}

func (pf *PrognFrame) GetEnv() (e *Environ) {
	return pf.Env
}

func (pf *PrognFrame) Exec(i SchemeObject) (r SchemeObject, next Frame, err error) {
	var obj SchemeObject

	switch {
	case pf.Obj == Onil:
		return Onil, pf.Parent, nil
	case pf.Obj.Cdr == Onil: // jump
		obj := pf.Obj.Car
		next = CreateEvalFrame(pf.Parent, obj, pf.Env)
	default:
		obj, pf.Obj, err = pf.Obj.Pop()
		if err != nil {
			return
		}

		next = CreateEvalFrame(pf, obj, pf.Env)
	}
	return
}

type ApplyFrame struct {
	Parent     Frame
	P          Procedure
	Args       *Cons
	EvaledArgs *Cons
	EvaledTail *Cons
	Env        *Environ
}

func CreateApplyFrame(p Frame, o *Cons, e *Environ) (f Frame) {
	return &ApplyFrame{Parent: p, Args: o, Env: e}
}

func (af *ApplyFrame) Debug() (r string) {
	return "ApplyFrame"
}

func (af *ApplyFrame) SetParent(p Frame) {
	af.Parent = p
}

func (af *ApplyFrame) GetEnv() (e *Environ) {
	return af.Env
}

func (af *ApplyFrame) Exec(i SchemeObject) (r SchemeObject, next Frame, err error) {
	var ok bool
	var obj SchemeObject

	// accept argument
	if af.P == nil {
		af.P, ok = i.(Procedure)
		if !ok {
			return nil, nil, ErrNotRunnable
		}

		if !af.P.IsApplicativeOrder() {
			// TODO: normal order
			return
		}

		af.EvaledArgs = Onil
	} else {
		t := &Cons{Car: i, Cdr: Onil}
		if af.EvaledArgs == Onil {
			af.EvaledArgs = t
			af.EvaledTail = t
		} else {
			af.EvaledTail.Cdr = t
			af.EvaledTail = t
		}
	}

	if af.Args == Onil { // all args has been evaled
		r, next, err = af.P.Apply(af.EvaledArgs, af)
		if err != nil {
			log.Error("%s", err)
			return
		}
		if next == nil {
			next = af.Parent
		}
		return
	}

	obj, af.Args, err = af.Args.Pop()
	if err != nil {
		return
	}

	next = CreateEvalFrame(af, obj, af.Env)
	return
}
