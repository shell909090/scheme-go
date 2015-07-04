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

func SchemeObjectToString(o SchemeObject) (s string) {
	if o == nil {
		return ""
	}

	buf := bytes.NewBuffer(nil)
	_, err := o.Format(buf, 0)
	if err != nil {
		log.Error("%s", err)
		return "<unknown>"
	}
	return buf.String()
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

func ReverseList(o *Cons) (r *Cons, err error) {
	if o == Onil {
		return o, nil
	}

	var ok bool
	r = Onil // that's for right
	l := o   // and this is left
	for l != Onil {
		next := l.Cdr // next is the left of left
		l.Cdr = r
		r = l
		l, ok = next.(*Cons)
		if !ok {
			return nil, ErrISNotAList
		}
	}
	return
}
