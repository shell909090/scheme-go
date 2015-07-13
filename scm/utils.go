package scm

import (
	"errors"

	logging "github.com/op/go-logging"
)

var (
	ErrQuit           = errors.New("quit")
	ErrListOutOfIndex = errors.New("out of index when get list")
	ErrType           = errors.New("runtime type error")
	ErrISNotAList     = errors.New("object is not a list")
	ErrUnknown        = errors.New("unknown error")
	ErrNameNotFound   = errors.New("name not found")
	ErrNotRunnable    = errors.New("object not runnable")
)

var (
	log                   = logging.MustGetLogger("scm")
	DefaultNames          = make(map[string]Obj)
	DefaultEnv   *Environ = &Environ{Parent: nil, Names: DefaultNames}
)

type Formatter interface {
	Format() (r string)
}

func Eval(o Obj, env *Environ, f Frame) (value Obj, next Frame, err error) {
	switch t := o.(type) {
	case *Symbol:
		value = env.Get(t.Name)
		if value == nil {
			return nil, nil, ErrNameNotFound
		}
	case *Quote:
		value = t.Objs
	case *Cons:
		var procedure Obj
		procedure, t, err = t.Pop()
		if err != nil {
			return
		}

		next = CreateApplyFrame(t, env, f) // not sure about procedure yet.
		// get a result now, or get a frame which can return in future.
		next, err = EvalAndReturn(procedure, env, next)
	case *InternalProcedure, *LambdaProcedure:
		panic("run eval in procedure")
	default:
		value = o
	}
	return
}

func EvalAndReturn(i Obj, e *Environ, f Frame) (next Frame, err error) {
	t, next, err := Eval(i, e, f)
	if err != nil {
		log.Error("%s", err)
		return
	}
	if next != nil {
		return
	}

	next = f
	err = next.Return(t)
	if err != nil {
		log.Error("%s", err)
	}
	return
}

func Trampoline(f Frame) (result Obj, err error) {
	for {
		log.Debug("stack:\n%s", StackFormat(f))
		f, err = f.Exec()
		if err != nil {
			log.Error("%s", err)
			return
		}
		if t, ok := f.(*EndFrame); ok {
			if t.result == nil {
				return nil, ErrUnknown
			}
			return t.result, nil
		}
	}
	return nil, ErrUnknown
}

func RunCode(code Obj) (result Obj, env *Environ, err error) {
	list, ok := code.(*Cons)
	if !ok {
		return nil, nil, ErrType
	}
	env = DefaultEnv.Fork()
	result, err = Trampoline(CreateBeginFrame(list, env, &EndFrame{}))
	return
}

func Apply(p Procedure, args *Cons) (result Obj, err error) {
	return Trampoline(&ApplyFrame{BaseFrame: BaseFrame{
		Parent: &EndFrame{}, Env: DefaultEnv.Fork()},
		procedure: p, Args: args, EvaledArgs: Onil})
}
