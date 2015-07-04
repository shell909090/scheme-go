package scmgo

type Procedure interface {
	SchemeObject
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
	log.Info("apply !%s, argument: %s", ip.Name, o.Format())
	r, next, err = ip.f(o, p)
	switch {
	case next != nil:
		log.Info("next: %p", next)
	case r != nil:
		log.Info("result: %s", r.Format())
	}
	return
}

func (ip *InternalProcedure) Format() (r string) {
	return "!" + ip.Name
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

func genNames(lp *LambdaProcedure, o *Cons) (r map[string]SchemeObject, err error) {
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
	names, err := genNames(lp, o)
	if err != nil {
		return
	}
	env := lp.Env.Fork(names)
	log.Info("apply %s, env:\n%s", lp.Format(), env.Format())
	next = CreateBeginFrame(
		lp.Obj, env,
		p.GetParent()) // coming from apply, so pass this frame.
	return
}

func (lp *LambdaProcedure) Format() (r string) {
	name := lp.Name
	if name == "" {
		name = "lambda"
	}
	return "<" + name + ">"
}
