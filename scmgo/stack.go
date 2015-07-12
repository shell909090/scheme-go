package scmgo

import "fmt"

type Frame interface {
	Formatter
	GetParent() (p Frame)
	GetEnv() (e *Environ)
	Return(i SchemeObject) (err error)
	Exec() (next Frame, err error)
}

type EndFrame struct {
	result SchemeObject
	Env    *Environ
}

func (f *EndFrame) Format() (r string) {
	return "End"
}

func (f *EndFrame) GetParent() (p Frame) {
	return nil
}

func (f *EndFrame) GetEnv() (e *Environ) {
	return f.Env
}

func (f *EndFrame) Return(i SchemeObject) (err error) {
	f.result = i
	return
}

func (f *EndFrame) Exec() (next Frame, err error) {
	return nil, ErrUnknown
}

type BeginFrame struct {
	Parent Frame
	Obj    *Cons
	Env    *Environ
}

func CreateBeginFrame(o *Cons, e *Environ, p Frame) (f Frame) {
	return &BeginFrame{Parent: p, Obj: o, Env: e}
}

func (f *BeginFrame) Format() (r string) {
	n, err := f.Obj.Len(false)
	if err != nil {
		n = 0
	}
	return fmt.Sprintf("Begin: %d", n)
}

func (f *BeginFrame) GetParent() (p Frame) {
	return f.Parent
}

func (f *BeginFrame) GetEnv() (e *Environ) {
	return f.Env
}

func (f *BeginFrame) Return(i SchemeObject) (err error) {
	return nil
}

func (f *BeginFrame) Exec() (next Frame, err error) {
	var obj SchemeObject
	switch {
	case f.Obj == Onil: // FIXME: not make sense
		return f.Parent, nil
	case f.Obj.Cdr == Onil: // jump
		return EvalAndReturn(f.Obj.Car, f.Env, f.Parent)
	default: // eval
		obj, f.Obj, err = f.Obj.Pop()
		if err != nil {
			log.Error("%s", err)
			return
		}
		return EvalAndReturn(obj, f.Env, f)
	}
	return
}

type ApplyFrame struct {
	Parent     Frame
	P          Procedure
	Args       *Cons
	EvaledArgs *Cons
	Env        *Environ
}

func CreateApplyFrame(a *Cons, e *Environ, p Frame) (f *ApplyFrame) {
	return &ApplyFrame{Parent: p, Args: a, EvaledArgs: Onil, Env: e}
}

func (f *ApplyFrame) Format() (r string) {
	return "Apply"
}

func (f *ApplyFrame) GetParent() (p Frame) {
	return f.Parent
}

func (f *ApplyFrame) GetEnv() (e *Environ) {
	return f.Env
}

func (f *ApplyFrame) Return(i SchemeObject) (err error) {
	var ok bool
	if f.P != nil {
		f.EvaledArgs = f.EvaledArgs.Push(i)
		return
	}
	f.P, ok = i.(Procedure)
	if !ok {
		return ErrNotRunnable
	}
	return
}

func (f *ApplyFrame) Apply(args *Cons) (next Frame, err error) {
	tmp, next, err := f.P.Apply(args, f)
	if err != nil {
		log.Error("%s", err)
		return
	}
	if next != nil {
		return
	}
	next = f.Parent
	err = next.Return(tmp)
	return
}

func (f *ApplyFrame) Exec() (next Frame, err error) {
	var obj SchemeObject
	if !f.P.IsApplicativeOrder() { // normal order
		next, err = f.Apply(f.Args)
		return
	}

	for f.Args != Onil {
		// pop up next argument
		obj, f.Args, err = f.Args.Pop()
		if err != nil {
			log.Error("%s", err)
			return
		}

		next, err = EvalAndReturn(obj, f.Env, f)
		if err != nil {
			log.Error("%s", err)
			return
		}
		if next != nil {
			return
		}
	}

	// all args has been evaled
	f.EvaledArgs, err = ReverseList(f.EvaledArgs, Onil)
	if err != nil {
		return
	}
	next, err = f.Apply(f.EvaledArgs)
	return
}

type IfFrame struct {
	Parent Frame
	Env    *Environ
	Cond   SchemeObject
	TCase  SchemeObject
	ECase  SchemeObject
	Hit    SchemeObject
}

func CreateIfFrame(cond, tcase, ecase SchemeObject, e *Environ, p Frame) (f Frame) {
	return &IfFrame{Parent: p, Cond: cond, TCase: tcase, ECase: ecase, Env: e}
}

func (f *IfFrame) Format() (r string) {
	return fmt.Sprintf("If:\n%s", f.Cond.Format())
}

func (f *IfFrame) GetParent() (p Frame) {
	return f.Parent
}

func (f *IfFrame) GetEnv() (e *Environ) {
	return f.Env
}

func (f *IfFrame) Return(i SchemeObject) (err error) {
	b, ok := i.(Boolean)
	if !ok {
		return ErrType
	}
	if bool(b) {
		f.Hit = f.TCase
	} else {
		f.Hit = f.ECase
	}
	return
}

func (f *IfFrame) Exec() (next Frame, err error) {
	switch f.Hit {
	case nil: // eval condition.
		return EvalAndReturn(f.Cond, f.Env, f)
	case Onil: // pass if with no case.
		next = f.Parent
		err = next.Return(Onil)
		return
	default: // eval case.
		return EvalAndReturn(f.Hit, f.Env, f.Parent)
	}
}
