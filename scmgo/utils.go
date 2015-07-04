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
	log          = logging.MustGetLogger("scmgo")
	DefaultNames = make(map[string]SchemeObject)
)

func GetHeadAsSymbol(o *Cons) (s *Symbol, err error) {
	s, ok := o.Car.(*Symbol)
	if !ok {
		return nil, ErrType
	}
	return
}

func StackFormatter(f Frame) (r string) {
	buf := bytes.NewBuffer(nil)
	for c := f; c != nil; c = c.GetParent() {
		if _, ok := c.(*EndFrame); !ok {
			buf.WriteString(c.Format() + "\n")
		}
	}
	return buf.String()
}

func EvalAndReturn(i SchemeObject, e *Environ, p Frame) (next Frame, err error) {
	t, next, err := i.Eval(e, p)
	if err != nil {
		log.Error("%s", err)
		return
	}

	if next != nil {
		return
	}

	next = p
	err = next.Return(t)
	if err != nil {
		log.Error("%s", err)
	}
	return
}

func ReverseList(o *Cons) (result *Cons, err error) {
	var ok bool
	// image a list from left to right.
	l := Onil // that's for left.
	r := o    // and this is right.

	for r != Onil {
		next := r.Cdr        // record the next one of the left.
		r.Cdr = l            // turn right back.
		l = r                // push left forward.
		r, ok = next.(*Cons) // and push right forward, if can.
		if !ok {
			return nil, ErrISNotAList
		}
	}
	return l, nil
}
