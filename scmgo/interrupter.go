package scmgo

type Frame interface {
	Debug() (r string, err error)
	SetParent(p Frame)
	// GetParent() (p Frame)
	Exec(i SchemeObject) (r SchemeObject, next Frame, err error)
}

type Applicable interface {
	IsApplicativeOrder() bool
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

type EvalFrame struct {
	Parent Frame
	Obj    SchemeObject
	Env    *Environ
}

func CreateEvalFrame(o SchemeObject, e *Environ) (f Frame) {
	return &EvalFrame{Obj: o, Env: e}
}

func (ef *EvalFrame) Debug() (r string, err error) {
	return "EvalFrame", nil
}

func (ef *EvalFrame) SetParent(p Frame) {
	ef.Parent = p
}

// func (ef *EvalFrame) GetParent() (p Frame) {
// 	return ef.Parent
// }

func (ef *EvalFrame) Exec(i SchemeObject) (r SchemeObject, next Frame, err error) {
	r, next, err = ef.Obj.Eval(ef.Env, ef)
	return
}

type PrognFrame struct {
	Parent Frame
	Obj    *Cons
	Env    *Environ
}

func CreatePrognFrame(o *Cons, e *Environ) (f Frame) {
	return &PrognFrame{Obj: o, Env: e}
}

func (pf *PrognFrame) Debug() (r string, err error) {
	return "PrognFrame", nil
}

func (pf *PrognFrame) SetParent(p Frame) {
	pf.Parent = p
}

// func (pf *PrognFrame) GetParent() (p Frame) {
// 	return pf.Parent
// }

func (pf *PrognFrame) Exec(i SchemeObject) (r SchemeObject, next Frame, err error) {
	var ok bool
	switch {
	case pf.Obj == Onil:
		return Onil, pf.Parent, nil
	case pf.Obj.Cdr == Onil: // jump
		obj := pf.Obj.Car
		next = CreateEvalFrame(obj, pf.Env)
		next.SetParent(pf.Parent)
	default:
		obj := pf.Obj.Car
		pf.Obj, ok = pf.Obj.Cdr.(*Cons)
		if !ok {
			err = ErrISNotAList
			return
		}

		next = CreateEvalFrame(obj, pf.Env)
		next.SetParent(pf)
	}
	return
}

type ApplyFrame struct {
	Parent Frame
	P      Procedure
	Args   *Cons
	Env    *Environ
}

func CreateApplyFrame(o *Cons, e *Environ) (f Frame) {
	return &ApplyFrame{Args: o, Env: e}
}

func (af *ApplyFrame) Debug() (r string, err error) {
	return "ApplyFrame", nil
}

func (af *ApplyFrame) SetParent(p Frame) {
	af.Parent = p
}

func (af *ApplyFrame) Exec(i SchemeObject) (r SchemeObject, next Frame, err error) {
	if af.P == nil {
		af.P = i
	}

	if af.Args == Onil { // all args has been evaled

	}
	return
}
