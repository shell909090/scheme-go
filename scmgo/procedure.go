package scmgo

import (
	"fmt"
	"io"
)

type Procedure interface {
	IsApplicativeOrder() bool
	Apply(o *Cons, p Frame) (r SchemeObject, next Frame, err error)
}

type InternalProcedure struct {
	Name        string
	f           func(i *Cons, p Frame) (r SchemeObject, next Frame, err error)
	applicative bool
}

func RegisterInternalProcedure(name string, f func(o *Cons, p Frame) (r SchemeObject, next Frame, err error), applicative bool) {
	DefaultNames[name] = &InternalProcedure{Name: name, f: f, applicative: applicative}
}

func (ip *InternalProcedure) IsApplicativeOrder() bool {
	return ip.applicative
}

func (ip *InternalProcedure) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	panic("run eval of internal procedure")
}

func (ip *InternalProcedure) Apply(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	log.Info("apply %s %s", ip.Name, SchemeObjectToString(o))
	r, next, err = ip.f(o, p)
	log.Info("result %s", SchemeObjectToString(r))
	return
}

func (ip *InternalProcedure) Format(s io.Writer, lv int) (rv int, err error) {
	s.Write([]byte("!"))
	rv, err = s.Write([]byte(ip.Name))
	rv += lv + 1
	return
}

type LambdaProcedure struct {
	Name string
	Env  *Environ
	Args *Cons
	Obj  *Cons
}

func (lp *LambdaProcedure) IsApplicativeOrder() bool {
	return true
}

func (lp *LambdaProcedure) Eval(env *Environ, p Frame) (r SchemeObject, next Frame, err error) {
	panic("run eval of lambda procedure")
}

func (lp *LambdaProcedure) GenNames(o *Cons) (r map[string]SchemeObject, err error) {
	var s, s1 *Symbol
	r = make(map[string]SchemeObject)

	pn := lp.Args // parameters by name
	pv := o       // parameters by vector

	for pn != Onil && pv != Onil {
		s, err = GetHeadAsSymbol(pn)
		if err != nil {
			return
		}

		if s.Name == "." {
			_, pn, err = pn.Pop()
			if err != nil {
				return
			}

			s1, err = GetHeadAsSymbol(pn)
			if err != nil {
				return
			}

			r[s1.Name] = pv
			break
		}
		r[s.Name] = pv.Car

		_, pn, err = pn.Pop()
		if err != nil {
			return
		}
		_, pv, err = pv.Pop()
		if err != nil {
			return
		}
	}
	return
}

func (lp *LambdaProcedure) Apply(o *Cons, p Frame) (r SchemeObject, next Frame, err error) {
	names, err := lp.GenNames(o)
	if err != nil {
		return
	}
	// this is called by apply frame, so pass this frame
	next = CreatePrognFrame(p.GetParent(), lp.Obj, lp.Env.Fork(names))
	return
}

func (lp *LambdaProcedure) Format(s io.Writer, lv int) (rv int, err error) {
	name := lp.Name
	if name == "" {
		name = "lambda"
	}
	name = fmt.Sprintf("<%s>", name)
	rv, err = s.Write([]byte(name))
	rv += lv
	return
}