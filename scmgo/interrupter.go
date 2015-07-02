package scmgo

type Frame interface {
	SetParent(p Frame)
	GetParent() (p Frame)
	Exec(i SchemeObject) (r SchemeObject, next Frame, err error)
}

// type OFunction struct {
// 	Name   string
// 	Env    Environ
// 	Params *Cons
// 	Objs   SchemeObject
// 	evaled bool
// }

// func GetNAsSymbol(p *Cons, n int) (s *Symbol, err error) {
// 	t, err := p.GetN(n)
// 	if err != nil {
// 		return
// 	}
// 	s, ok := t.(*Symbol)
// 	if !ok {
// 		return nil, ErrRuntimeType
// 	}
// 	return
// }

// func (o *OFunction) Eval(o SchemeObject, stack *Stack, env *Environ) (r SchemeObject, err error) {
// 	r := make(map[string]SchemeObject)
// 	var pn, pv *Cons
// 	pn = o.Params         // parameters by name
// 	pv, ok = objs.(*Cons) // parameters by vector
// 	if !ok {
// 		return ErrRuntimeType
// 	}

// 	for pn != nil && pv != nil {
// 		s, err := GetNAsSymbol(pn, 0)
// 		if err != nil {
// 			return err
// 		}
// 		if s.name == "." {
// 			s1, err := GetNAsSymbol(pn, 1)
// 			if err != nil {
// 				return err
// 			}
// 			r[s1.name] = pv
// 			break
// 		}
// 		r[s.name] = pv.car
// 		pn = pn.cdr
// 		pv = pv.cdr
// 	}
// 	ne := o.Env.Fork(r)
// 	return stack.Jump(&PrognStatus{o.Objs}, ne)
// }

// type PrognStatus struct {
// 	Objs SchemeObject
// }

// func (ps *PrognStatus) Eval(o SchemeObject, stack *Stack, env *Environ) (r SchemeObject, err error) {
// 	if ps.Objs.cdr == nil {
// 		return stack.Jump(ps.Objs.car, envs)
// 	}
// 	t := ps.Objs.car
// 	ps.Objs = ps.Objs.cdr
// 	return stack.Call(t, envs)
// }

// type Frame struct {
// 	Obj    SchemeObject
// 	Env    *Environ
// 	Status map[string]SchemeObject
// }

// func (fm *Frame) Eval(stack *Stack, i SchemeObject) (r SchemeObject, popup bool, err error) {
// 	// Ok, we have obj to eval, env, status, and maybe result
// 	r, err = fm.Obj.Eval(stack, fm.Env)
// 	return
// }

type Stack struct {
	Frames []Frame
}

// func (s *Stack) Eval(o SchemeObject, env *Environ) {
// 	// actually, this is eval. push a frame to stack means to eval it.
// 	// call eval in here, and really do eval in `Frame.Eval` .
// 	f := &Frame{Obj: o, Env: env, Status: make(map[string]SchemeObject)}
// 	s.Frames = append(s.Frames, f)
// }

func (s *Stack) Pop() {
	s.Frames = s.Frames[:len(s.Frames)-2]
}

// func (s *Stack) Exec(o SchemeObject, envs *Environ, objs SchemeObject) (rslt SchemeObject, next bool, err error) {
// 	switch et := e.(type) {
// 	case Symbol:
// 	}
// }

// func (s *Stack) Jump(e SchemeObject, envs *Environ) (f *Frame) {
// 	switch et := e.(type) {
// 	case Symbol:
// 		// return &Frame{Exec: , Envs: }
// 	default:
// 	}
// }

func (s *Stack) Apply(o SchemeObject, p *Cons) {
	// apply a set of parameters to a object
	// if o.IsApplicativeOrder() {

	// } else {

	// }
	return
}

type EvalFrame struct {
	Parent Frame
	Obj    SchemeObject
	Env    *Environ
}

func (ef *EvalFrame) SetParent(p Frame) {
	ef.Parent = p
}

func (ef *EvalFrame) GetParent() (p Frame) {
	return ef.Parent
}

func (ef *EvalFrame) Exec(i SchemeObject) (r SchemeObject, next Frame, err error) {
	r, next, err = ef.Obj.Eval(ef.Env)
	if err != nil {
		return
	}
	if next == nil {
		next = ef.Parent
	} else {
		next.SetParent(ef)
	}
	return
}

type PrognFrame struct {
	Parent Frame
	P      *Cons
	Env    *Environ
}

func (pf *PrognFrame) Exec(i SchemeObject) (r SchemeObject, next Frame, err error) {
	return
}

func (pf *PrognFrame) Eval(stack *Stack, i SchemeObject) (r SchemeObject, popup bool, err error) {
	// var ok bool
	// switch {
	// case pf.P == Onil:
	// 	return Onil, true, nil
	// case pf.P.Cdr == Onil:
	// 	stack.Pop()
	// 	// jump into
	// default:
	// 	obj := pf.P.Car
	// 	pf.P, ok = pf.P.Cdr.(*Cons)
	// 	if !ok {
	// 		return nil, false, ErrISNotAList
	// 	}
	// 	_, err = obj.Eval(pf.Env)
	// 	if err != nil {
	// 		return
	// 	}
	// 	return nil, false, nil
	// }
	// return i, true, nil
	return
}

type ApplyFrame struct {
	P   *Cons
	Env *Environ
}

func (af *ApplyFrame) Eval(stack *Stack, i SchemeObject) (r SchemeObject, popup bool, err error) {
	return
}
