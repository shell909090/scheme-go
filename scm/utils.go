package scm

import (
	"errors"

	logging "github.com/op/go-logging"
)

var (
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

func ReverseList(head *Cons, tail Obj) (result *Cons, err error) {
	// image a list from left to right.
	var ok bool
	left := tail
	right := head
	for right != Onil {
		next := right.Cdr        // record the next one of the left.
		right.Cdr = left         // turn right back.
		left = right             // push left forward.
		right, ok = next.(*Cons) // and push right forward, if can.
		if !ok {                 // improper
			return nil, ErrISNotAList
		}
	}
	return left.(*Cons), nil
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
		f = next

		// get a result now, or get a frame which can return in future.
		procedure, next, err = Eval(procedure, env, next)
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
		if next != nil {
			return
		}
		// get return immediately
		next = f
		err = next.Return(procedure)
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
			return t.result, nil
		}
	}
	return
}

func RunCode(code Obj) (result Obj, err error) {
	list, ok := code.(*Cons)
	if !ok {
		return nil, ErrType
	}

	env := &Environ{Parent: DefaultEnv, Names: make(map[string]Obj)}
	f := CreateBeginFrame(list, env, &EndFrame{Env: DefaultEnv})

	result, err = Trampoline(f)
	if result == nil {
		return nil, ErrUnknown
	}
	return
}
