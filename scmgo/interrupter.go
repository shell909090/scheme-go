package scmgo

import "fmt"

type Frame interface {
	Debug() (r string)
	SetParent(p Frame)
	// GetParent() (p Frame)
	Exec(i SchemeObject) (r SchemeObject, next Frame, err error)
}

type Procedure interface {
	IsApplicativeOrder() bool
	Apply(i SchemeObject, p Frame) (r SchemeObject, next Frame, err error)
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
// 		return nil, ErrType
// 	}
// 	return
// }

// func (o *OFunction) Eval(o SchemeObject, stack *Stack, env *Environ) (r SchemeObject, err error) {
// 	r := make(map[string]SchemeObject)
// 	var pn, pv *Cons
// 	pn = o.Params         // parameters by name
// 	pv, ok = objs.(*Cons) // parameters by vector
// 	if !ok {
// 		return ErrType
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

func (ef *EvalFrame) Debug() (r string) {
	return fmt.Sprintf("EvalFrame: {%s}", SchemeObjectToString(ef.Obj))
}

func (ef *EvalFrame) SetParent(p Frame) {
	ef.Parent = p
}

// func (ef *EvalFrame) GetParent() (p Frame) {
// 	return ef.Parent
// }

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

func CreatePrognFrame(o *Cons, e *Environ) (f Frame) {
	return &PrognFrame{Obj: o, Env: e}
}

func (pf *PrognFrame) Debug() (r string) {
	return "PrognFrame"
}

func (pf *PrognFrame) SetParent(p Frame) {
	pf.Parent = p
}

// func (pf *PrognFrame) GetParent() (p Frame) {
// 	return pf.Parent
// }

func (pf *PrognFrame) Exec(i SchemeObject) (r SchemeObject, next Frame, err error) {
	var obj SchemeObject

	switch {
	case pf.Obj == Onil:
		return Onil, pf.Parent, nil
	case pf.Obj.Cdr == Onil: // jump
		obj := pf.Obj.Car
		next = CreateEvalFrame(obj, pf.Env)
		next.SetParent(pf.Parent)
	default:
		obj, pf.Obj, err = pf.Obj.Pop()
		if err != nil {
			return
		}

		next = CreateEvalFrame(obj, pf.Env)
		next.SetParent(pf)
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

func CreateApplyFrame(o *Cons, e *Environ) (f Frame) {
	return &ApplyFrame{Args: o, Env: e}
}

func (af *ApplyFrame) Debug() (r string) {
	return "ApplyFrame"
}

func (af *ApplyFrame) SetParent(p Frame) {
	af.Parent = p
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
		r, next, err = af.P.Apply(af.EvaledArgs, af.Parent)
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

	next = CreateEvalFrame(obj, af.Env)
	next.SetParent(af)
	return
}
