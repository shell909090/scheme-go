package scmgo

// symbol? string?

// define lambda
// begin compile
// eval apply

// user-init-environment
// current-environment
// import

// let let*
// eq? equal?
// not and or
// cond if when

// display error
// newline
// exit

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
