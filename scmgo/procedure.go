package scmgo

type Procedure interface {
	SchemeObject
	IsApplicativeOrder() bool
	Apply(o *Cons, f Frame) (r SchemeObject, next Frame, err error)
}

type InternalProcedure struct {
	Name        string
	f           func(i *Cons, p Frame) (r SchemeObject, next Frame, err error)
	applicative bool
}

func RegisterInternalProcedure(name string, f func(o *Cons, p Frame) (r SchemeObject, next Frame, err error), applicative bool) {
	DefaultNames[name] = &InternalProcedure{Name: name, f: f, applicative: applicative}
}

func (p *InternalProcedure) IsApplicativeOrder() bool {
	return p.applicative
}

func (p *InternalProcedure) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
	panic("run eval of internal procedure")
}

func (p *InternalProcedure) Apply(o *Cons, f Frame) (r SchemeObject, next Frame, err error) {
	log.Info("apply !%s, argument: %s", p.Name, o.Format())
	r, next, err = p.f(o, f)
	switch {
	case next != nil:
		log.Info("next: %p", next)
	case r != nil:
		log.Info("result: %s", r.Format())
	}
	return
}

func (p *InternalProcedure) Format() (r string) {
	return "!" + p.Name
}

type LambdaProcedure struct {
	Name string
	Env  *Environ
	Args *Cons
	Obj  *Cons
}

func (p *LambdaProcedure) IsApplicativeOrder() bool {
	return true
}

func (p *LambdaProcedure) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
	panic("run eval of lambda procedure")
}

func genNames(p *LambdaProcedure, o *Cons) (r map[string]SchemeObject, err error) {
	var s, s1 *Symbol
	r = make(map[string]SchemeObject)

	pn := p.Args // parameters by name
	pv := o      // parameters by vector

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

func (p *LambdaProcedure) Apply(o *Cons, f Frame) (r SchemeObject, next Frame, err error) {
	names, err := genNames(p, o)
	if err != nil {
		return
	}
	env := p.Env.Fork(names)
	log.Info("apply %s, env:\n%s", p.Format(), env.Format())
	next = CreateBeginFrame(
		p.Obj, env,
		f.GetParent()) // coming from apply, so pass this frame.
	return
}

func (p *LambdaProcedure) Format() (r string) {
	name := p.Name
	if name == "" {
		name = "lambda"
	}
	return "<" + name + ">"
}
