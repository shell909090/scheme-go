package scmgo

type Procedure interface {
	SchemeObject
	IsApplicativeOrder() bool
	Apply(args *Cons, f Frame) (value SchemeObject, next Frame, err error)
}

type InternalProcedure struct {
	Name        string
	procedure   func(i *Cons, f Frame) (value SchemeObject, next Frame, err error)
	applicative bool
}

func RegisterInternalProcedure(name string, procedure func(o *Cons, f Frame) (value SchemeObject, next Frame, err error), applicative bool) {
	DefaultNames[name] = &InternalProcedure{
		Name:        name,
		procedure:   procedure,
		applicative: applicative,
	}
}

func (p *InternalProcedure) IsApplicativeOrder() bool {
	return p.applicative
}

// func (p *InternalProcedure) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
// 	panic("run eval of internal procedure")
// }

func (p *InternalProcedure) Apply(args *Cons, f Frame) (value SchemeObject, next Frame, err error) {
	log.Info("apply !%s, argument: %s", p.Name, args.Format())
	value, next, err = p.procedure(args, f)
	switch {
	case next != nil:
		log.Info("next: %p", next)
	case value != nil:
		log.Info("result: %s", value.Format())
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

// func (p *LambdaProcedure) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
// 	panic("run eval of lambda procedure")
// }

func setup_args(env *Environ, p *LambdaProcedure, o *Cons) (err error) {
	var t SchemeObject
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

func (p *LambdaProcedure) Apply(args *Cons, f Frame) (value SchemeObject, next Frame, err error) {
	env := p.Env.Fork()
	err = setup_args(env, p, args)
	if err != nil {
		return
	}
	log.Info("apply %s, env:\n%s", p.Format(), env.Format())
	next = CreateBeginFrame(
		p.Obj, env, f.GetParent()) // coming from apply, so pass this frame.
	return
}

func (p *LambdaProcedure) Format() (r string) {
	name := p.Name
	if name == "" {
		name = "lambda"
	}
	return "<" + name + ">"
}
