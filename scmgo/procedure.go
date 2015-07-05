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

func (p *InternalProcedure) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
	panic("run eval of internal procedure")
}

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

func (p *LambdaProcedure) Eval(env *Environ, f Frame) (value SchemeObject, next Frame, err error) {
	panic("run eval of lambda procedure")
}

func genNames(p *LambdaProcedure, o *Cons) (names map[string]SchemeObject, err error) {
	var s, s1 *Symbol
	names = make(map[string]SchemeObject)

	pn := p.Args // parameters by name
	pv := o      // parameters by vector

	for pn != Onil && pv != Onil {
		s, err = GetHeadAsSymbol(pn)
		if err != nil {
			return
		}

		if s.Name == "." {
			_, pn, err = pn.Pop(false)
			if err != nil {
				return
			}

			s1, err = GetHeadAsSymbol(pn)
			if err != nil {
				return
			}

			names[s1.Name] = pv
			break
		}
		names[s.Name] = pv.Car

		_, pn, err = pn.Pop(false)
		if err != nil {
			return
		}
		_, pv, err = pv.Pop(false)
		if err != nil {
			return
		}
	}
	return
}

func (p *LambdaProcedure) Apply(args *Cons, f Frame) (value SchemeObject, next Frame, err error) {
	names, err := genNames(p, args)
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
