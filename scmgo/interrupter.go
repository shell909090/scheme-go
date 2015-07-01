package scmgo

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

type Environ struct {
	Parent *Environ
	Names  map[string]SchemeObject
	// Fast   map[string]SchemeObject
}

// func (e *Env) GenFast() {
// }

func (e *Environ) Fork(r map[string]SchemeObject) (ne *Environ) {
	if r == nil {
		r = make(map[string]SchemeObject)
	}
	ne = &Environ{
		Parent: e,
		Names:  r,
	}
	return ne
}

func (e *Environ) Add(name string, value SchemeObject) {
	e.Names[name] = value
	// e.Fast[name] = value
}

func (e *Environ) Get(name string) (value SchemeObject, ok bool) {
	for ce := e; ce != nil; ce = e.Parent {
		value, ok = e.Names[name]
		if ok {
			return
		}
	}
	return nil, false
}

type Frame struct {
	Obj    SchemeObject
	Env    *Environ
	Status map[string]SchemeObject
}

func (fm *Frame) Eval(stack *Stack, i SchemeObject) (r SchemeObject, popup bool, err error) {
	// Ok, we have obj to eval, env, status, and maybe result
	r, err = fm.Obj.Eval(stack, fm.Env)
	return
}

type Stack struct {
	Frames []*Frame
}

func (s *Stack) Push(o SchemeObject, env *Environ) {
	f := &Frame{Obj: o, Env: env, Status: make(map[string]SchemeObject)}
	s.Frames = append(s.Frames, f)
}

func (s *Stack) Pop() {
	s.Frames = s.Frames[:len(s.Frames)-2]
}

func (s *Stack) Eval(o SchemeObject, envs *Environ) (r SchemeObject, err error) {
	// call eval in here, and really do eval in Frame.Eval
	return
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

func (s *Stack) Trampoline() (result SchemeObject, err error) {
	var popup bool

	for len(s.Frames) > 0 {
		f := s.Frames[len(s.Frames)-1]

		result, popup, err = f.Eval(s, result)
		if err != nil {
			return
		}

		if popup {
			s.Pop()
		}
	}
	return
}
