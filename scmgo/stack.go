package scmgo

import "fmt"

type Frame interface {
	Format() (r string)
	GetParent() (f Frame)
	GetEnv() (e *Environ)
	Eval(i SchemeObject) (r SchemeObject, next Frame, err error)
}

type EvalFrame struct {
	Parent Frame
	Obj    SchemeObject
	Env    *Environ
}

func CreateEvalFrame(p Frame, o SchemeObject, e *Environ) (f Frame) {
	return &EvalFrame{Parent: p, Obj: o, Env: e}
}

func (ef *EvalFrame) Format() (r string) {
	return fmt.Sprintf("Eval: {%s}", SchemeObjectToString(ef.Obj))
}

func (ef *EvalFrame) GetParent() (f Frame) {
	return ef.Parent
}

func (ef *EvalFrame) GetEnv() (e *Environ) {
	return ef.Env
}

func (ef *EvalFrame) Eval(i SchemeObject) (r SchemeObject, next Frame, err error) {
	log.Info("eval: {%s}", SchemeObjectToString(ef.Obj))
	r, next, err = ef.Obj.Eval(ef.Env, ef.Parent)
	if next == nil {
		next = ef.Parent
	}
	if err != nil {
		log.Error("%s", err)
		return
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

func (pf *PrognFrame) Format() (r string) {
	n, err := pf.Obj.Len()
	if err != nil {
		n = 0
	}
	return fmt.Sprintf("Progn: %d", n)
}

func (pf *PrognFrame) GetParent() (f Frame) {
	return pf.Parent
}

func (pf *PrognFrame) GetEnv() (e *Environ) {
	return pf.Env
}

func (pf *PrognFrame) Eval(i SchemeObject) (r SchemeObject, next Frame, err error) {
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
	return &ApplyFrame{Parent: p, Args: o, EvaledArgs: Onil, Env: e}
}

func (af *ApplyFrame) Format() (r string) {
	return "Apply"
}

func (af *ApplyFrame) GetParent() (f Frame) {
	return af.Parent
}

func (af *ApplyFrame) GetEnv() (e *Environ) {
	return af.Env
}

func (af *ApplyFrame) AppendEvaled(i SchemeObject) {
	t := &Cons{Car: i, Cdr: Onil}
	if af.EvaledArgs == Onil {
		af.EvaledArgs = t
		af.EvaledTail = t
	} else {
		af.EvaledTail.Cdr = t
		af.EvaledTail = t
	}
}

func (af *ApplyFrame) Apply(o *Cons) (r SchemeObject, next Frame, err error) {
	r, next, err = af.P.Apply(o, af)
	if err != nil {
		log.Error("%s", err)
		return
	}
	if next == nil {
		next = af.Parent
	}
	return
}

func (af *ApplyFrame) Eval(i SchemeObject) (r SchemeObject, next Frame, err error) {
	var ok bool
	var obj SchemeObject

	// accept argument
	if af.P == nil {
		af.P, ok = i.(Procedure)
		if !ok {
			return nil, nil, ErrNotRunnable
		}

		if !af.P.IsApplicativeOrder() { // normal order
			r, next, err = af.Apply(af.Args)
			return
		}
	} else {
		af.AppendEvaled(i)
	}

	if af.Args == Onil { // all args has been evaled
		r, next, err = af.Apply(af.EvaledArgs)
		return
	}

	// eval next argument
	obj, af.Args, err = af.Args.Pop()
	if err != nil {
		return
	}

	next = CreateEvalFrame(af, obj, af.Env)
	return
}

type CondFrame struct {
	Parent Frame
	Obj    *Cons
	Env    *Environ
}

func CreateCondFrame(p Frame, o *Cons, e *Environ) (f Frame) {
	return &CondFrame{Parent: p, Obj: o, Env: e}
}

func (cf *CondFrame) Format() (r string) {
	return fmt.Sprintf("Cond:\n%s", SchemeObjectToString(cf.Obj))
}

func (cf *CondFrame) GetParent() (f Frame) {
	return cf.Parent
}

func (cf *CondFrame) GetEnv() (e *Environ) {
	return cf.Env
}

func (cf *CondFrame) Eval(i SchemeObject) (r SchemeObject, next Frame, err error) {
	var ok bool
	var t SchemeObject
	var cond *Cons
	var b Boolean

	if i != nil {
		b, ok = i.(Boolean)
		if !ok {
			return nil, nil, ErrType
		}

		t, cf.Obj, err = cf.Obj.Pop()
		if err != nil {
			return
		}

		if bool(b) {
			cond, ok = t.(*Cons)
			if !ok {
				return nil, nil, ErrType
			}

			t, err = cond.GetN(1)
			if err != nil {
				return
			}

			next = CreateEvalFrame(cf.Parent, t, cf.Env)
			return
		}
	}

	t = cf.Obj.Car

	cond, ok = t.(*Cons)
	if !ok {
		return nil, nil, ErrType
	}

	t = cond.Car
	if n, ok := t.(*Symbol); ok && n.Name == "else" {
		log.Debug("hit else")
		t, err = cond.GetN(1)
		if err != nil {
			return
		}
		next = CreateEvalFrame(cf.Parent, t, cf.Env)
		return
	}
	next = CreateEvalFrame(cf, t, cf.Env)
	return
}
