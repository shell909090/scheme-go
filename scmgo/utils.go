package scmgo

import (
	"bytes"
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
	log                   = logging.MustGetLogger("scmgo")
	DefaultNames          = make(map[string]SchemeObject)
	DefaultEnv   *Environ = &Environ{Parent: nil, Names: DefaultNames}
)

func StackFormatter(f Frame) (r string) {
	buf := bytes.NewBuffer(nil)
	for c := f; c != nil; c = c.GetParent() {
		if _, ok := c.(*EndFrame); !ok {
			buf.WriteString(c.Format() + "\n")
		}
	}
	return buf.String()
}

func EvalAndReturn(i SchemeObject, e *Environ, f Frame) (next Frame, err error) {
	t, next, err := i.Eval(e, f)
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

func ReverseList(head *Cons, tail SchemeObject) (result *Cons, err error) {
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

func Trampoline(f Frame) (result SchemeObject, err error) {
	for {
		log.Debug("stack:\n%s", StackFormatter(f))
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

func RunCode(code SchemeObject) (result SchemeObject, err error) {
	list, ok := code.(*Cons)
	if !ok {
		return nil, ErrType
	}

	env := &Environ{Parent: DefaultEnv, Names: make(map[string]SchemeObject)}
	f := CreateBeginFrame(list, env, &EndFrame{Env: DefaultEnv})

	result, err = Trampoline(f)
	if result == nil {
		return nil, ErrUnknown
	}
	return
}
