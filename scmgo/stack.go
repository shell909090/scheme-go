package scmgo

import "fmt"

type Frame interface {
	Format() (r string)
	GetParent() (f Frame)
	GetEnv() (e *Environ)
	Return(i SchemeObject) (err error)
	Eval() (next Frame, err error)
}

type EndFrame struct {
	result SchemeObject
	Env    *Environ
}

func (ef *EndFrame) Format() (r string) {
	return "End"
}

func (ef *EndFrame) GetParent() (f Frame) {
	return nil
}

func (ef *EndFrame) GetEnv() (e *Environ) {
	return ef.Env
}

func (ef *EndFrame) Return(i SchemeObject) (err error) {
	ef.result = i
	return
}

func (ef *EndFrame) Eval() (next Frame, err error) {
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

func (bf *BeginFrame) Format() (r string) {
	n, err := bf.Obj.Len()
	if err != nil {
		n = 0
	}
	return fmt.Sprintf("Begin: %d", n)
}

func (bf *BeginFrame) GetParent() (f Frame) {
	return bf.Parent
}

func (bf *BeginFrame) GetEnv() (e *Environ) {
	return bf.Env
}

func (bf *BeginFrame) Return(i SchemeObject) (err error) {
	return nil
}

func (bf *BeginFrame) Eval() (next Frame, err error) {
	var obj SchemeObject

	for {
		switch {
		case bf.Obj == Onil: // FIXME: not make sense
			return bf.Parent, nil
		case bf.Obj.Cdr == Onil: // jump
			obj := bf.Obj.Car

			next, err = EvalAndReturn(obj, bf.Env, bf.Parent)
			return
		default: // eval
			obj, bf.Obj, err = bf.Obj.Pop()
			if err != nil {
				return
			}

			_, next, err = obj.Eval(bf.Env, bf)
			if err != nil {
				log.Error("%s", err)
				return
			}
			if next != nil {
				return
			}
		}
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

func (af *ApplyFrame) Format() (r string) {
	return "Apply"
}

func (af *ApplyFrame) GetParent() (f Frame) {
	return af.Parent
}

func (af *ApplyFrame) GetEnv() (e *Environ) {
	return af.Env
}

func (af *ApplyFrame) Return(i SchemeObject) (err error) {
	var ok bool
	if af.P != nil {
		af.EvaledArgs = af.EvaledArgs.Push(i)
		return
	}
	af.P, ok = i.(Procedure)
	if !ok {
		return ErrNotRunnable
	}
	return
}

func (af *ApplyFrame) Apply(o *Cons) (next Frame, err error) {
	t, next, err := af.P.Apply(o, af)
	if err != nil {
		log.Error("%s", err)
		return
	}
	if next != nil {
		return
	}
	next = af.Parent
	err = next.Return(t)
	return
}

func (af *ApplyFrame) Eval() (next Frame, err error) {
	var t, obj SchemeObject

	if !af.P.IsApplicativeOrder() { // normal order
		next, err = af.Apply(af.Args)
		return
	}

	for af.Args != Onil {
		// pop up next argument
		obj, af.Args, err = af.Args.Pop()
		if err != nil {
			return
		}

		t, next, err = obj.Eval(af.Env, af)
		if err != nil {
			log.Error("%s", err)
			return
		}
		if next != nil { // if had to call next frame, quit to call
			return
		}
		// append to evaled args.
		af.EvaledArgs = af.EvaledArgs.Push(t)
	}

	// all args has been evaled
	af.EvaledArgs, err = ReverseList(af.EvaledArgs)
	if err != nil {
		return
	}
	next, err = af.Apply(af.EvaledArgs)
	return
}

type CondFrame struct {
	Parent Frame
	Obj    *Cons
	Env    *Environ
	Hit    SchemeObject
}

func CreateCondFrame(o *Cons, e *Environ, p Frame) (f Frame) {
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

func (cf *CondFrame) Return(i SchemeObject) (err error) {
	var t SchemeObject
	b, ok := i.(Boolean)
	if !ok {
		return ErrType
	}

	t, cf.Obj, err = cf.Obj.Pop()
	if err != nil {
		return
	}

	cond, ok := t.(*Cons)
	if !ok {
		return ErrType
	}

	t, err = cond.GetN(1)
	if err != nil {
		return
	}

	if bool(b) {
		cf.Hit = t
	}
	return
}

func (cf *CondFrame) Eval() (next Frame, err error) {
	var ok bool
	var cond *Cons
	var t SchemeObject

	if cf.Hit != nil { // finally, we matched a condition.
		next, err = EvalAndReturn(cf.Hit, cf.Env, cf.Parent)
		return
	}

	t = cf.Obj.Car

	cond, ok = t.(*Cons)
	if !ok {
		return nil, ErrType
	}

	t = cond.Car
	if n, ok := t.(*Symbol); ok && n.Name == "else" {
		log.Debug("hit else")
		t, err = cond.GetN(1)
		if err != nil {
			log.Error("%s", err)
			return
		}
		cf.Hit = t
		return cf, nil
	}

	// actually eval a condition.

	next, err = EvalAndReturn(t, cf.Env, cf)
	return
}
