package scm

type Procedure interface {
	Obj
	IsApplicativeOrder() bool
	Apply(args *Cons, f Frame) (value Obj, next Frame, err error)
}

type InternalProcedure struct {
	Name        string
	procedure   func(i *Cons, f Frame) (value Obj, next Frame, err error)
	applicative bool
}

func RegisterInternalProcedure(name string, procedure func(o *Cons, f Frame) (value Obj, next Frame, err error), applicative bool) {
	DefaultNames[name] = &InternalProcedure{
		Name:        name,
		procedure:   procedure,
		applicative: applicative,
	}
}

func (p *InternalProcedure) IsApplicativeOrder() bool {
	return p.applicative
}

func (p *InternalProcedure) Apply(args *Cons, f Frame) (value Obj, next Frame, err error) {
	log.Info("internal %s", p.Name)
	// internal procedure has to remember, we now in apply frame, will be abandon next.
	value, next, err = p.procedure(args, f)
	switch {
	case next != nil:
		log.Info("next: %p", next)
	case value != nil:
		log.Info("result: %s", Format(value))
	}
	return
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

func setup_args(env *Environ, p *LambdaProcedure, o *Cons) (err error) {
	var t Obj
	var s *Symbol

	pn := p.Args // parameters by name
	pv := o      // parameters by vector

	for pn != Onil && pv != Onil {
		s, pn, err = pn.PopSymbol()
		if err != nil {
			log.Error("%s", err.Error())
			return
		}

		if s.Name == "." {
			s, pn, err = pn.PopSymbol()
			if err != nil {
				log.Error("%s", err.Error())
				return
			}
			env.Add(s.Name, pv)
			break
		}

		t, pv, err = pv.Pop()
		if err != nil {
			return
		}
		env.Add(s.Name, t)
	}
	return
}

func (p *LambdaProcedure) Apply(args *Cons, f Frame) (value Obj, next Frame, err error) {
	env := p.Env.Fork()
	err = setup_args(env, p, args)
	if err != nil {
		return
	}
	log.Info("lambda %s", Format(p))
	next = CreateBeginFrame(
		p.Obj, env, f.GetParent()) // coming from apply, so pass this frame.
	return
}
