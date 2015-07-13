package scm

type Frame interface {
	GetParent() (p Frame)
	GetEnv() (e *Environ)
	Return(i Obj) (err error)
	Exec() (next Frame, err error)
}

type BaseFrame struct {
	Parent Frame
	Env    *Environ
}

func (f *BaseFrame) GetParent() (p Frame) {
	return f.Parent
}

func (f *BaseFrame) GetEnv() (e *Environ) {
	return f.Env
}

type EndFrame struct {
	result Obj
}

func (f *EndFrame) GetParent() (p Frame) {
	return nil
}

func (f *EndFrame) GetEnv() (e *Environ) {
	return nil
}

func (f *EndFrame) Return(i Obj) (err error) {
	f.result = i
	return
}

func (f *EndFrame) Exec() (next Frame, err error) {
	return nil, ErrUnknown
}

type BeginFrame struct {
	BaseFrame
	Obj *Cons
}

func CreateBeginFrame(o *Cons, e *Environ, p Frame) (f Frame) {
	return &BeginFrame{BaseFrame: BaseFrame{Parent: p, Env: e}, Obj: o}
}

func (f *BeginFrame) Return(i Obj) (err error) {
	return nil
}

func (f *BeginFrame) Exec() (next Frame, err error) {
	var obj Obj
	switch {
	case f.Obj == Onil:
		panic("not make sense")
	case f.Obj.Cdr == Onil: // tail call optimization
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
	BaseFrame
	procedure  Procedure
	Args       *Cons
	EvaledArgs *Cons
}

func CreateApplyFrame(a *Cons, e *Environ, p Frame) (f *ApplyFrame) {
	return &ApplyFrame{BaseFrame: BaseFrame{Parent: p, Env: e},
		Args: a, EvaledArgs: Onil}
}

func (f *ApplyFrame) Return(i Obj) (err error) {
	var ok bool
	if f.procedure != nil {
		f.EvaledArgs = f.EvaledArgs.Push(i)
		return
	}
	f.procedure, ok = i.(Procedure)
	if !ok {
		return ErrNotRunnable
	}
	return
}

func (f *ApplyFrame) Apply(args *Cons) (next Frame, err error) {
	tmp, next, err := f.procedure.Apply(args, f)
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
	if !f.procedure.IsApplicativeOrder() { // normal order
		next, err = f.Apply(f.Args)
		return
	}

	var obj Obj
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
	f.EvaledArgs, err = f.EvaledArgs.Reverse(Onil)
	if err != nil {
		return
	}
	next, err = f.Apply(f.EvaledArgs)
	return
}

type IfFrame struct {
	BaseFrame
	Cond  Obj
	TCase Obj
	ECase Obj
	Hit   Obj
}

func CreateIfFrame(cond, tcase, ecase Obj, e *Environ, p Frame) (f Frame) {
	return &IfFrame{BaseFrame: BaseFrame{Parent: p, Env: e},
		Cond: cond, TCase: tcase, ECase: ecase}
}

func (f *IfFrame) Return(i Obj) (err error) {
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
	}
	return EvalAndReturn(f.Hit, f.Env, f.Parent) // eval case.
}
