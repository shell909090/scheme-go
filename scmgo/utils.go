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
		buf.WriteString(c.Format() + "\n")
	}
	return buf.String()
}

func EvalMaybeInFrame(p Frame, i SchemeObject, e *Environ) (r SchemeObject, next Frame, err error) {
	_, ok := i.(*Cons)
	if ok {
		next = CreateEvalFrame(p, i, e)
		return
	}

	r, next, err = i.Eval(e, p)
	if next != nil {
		log.Error("fast run a object but not return")
		return nil, nil, ErrUnknown
	}
	return
}
