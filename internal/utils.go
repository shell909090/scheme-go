package internal

import (
	"errors"

	"bitbucket.org/shell909090/scheme-go/scmgo"
	logging "github.com/op/go-logging"
)

var (
	ErrQuit      = errors.New("quit")
	ErrArguments = errors.New("wrong arguments")
)

var (
	log = logging.MustGetLogger("internal")
)

func AssertLen(o *scmgo.Cons, length int) (err error) {
	n, err := o.Len(false)
	if err != nil {
		return
	}
	if n != length {
		return ErrArguments
	}
	return
}

func ParseParameters(list *scmgo.Cons, v ...interface{}) (next *scmgo.Cons, err error) {
	for _, a := range v {
		switch arg := a.(type) {
		case *string:
			switch first := list.Car.(type) {
			case *scmgo.Symbol:
				*arg = first.Name
			case *scmgo.String:
				*arg = string(*first)
			default:
				return nil, ErrArguments
			}
		case **scmgo.Cons:
			first, ok := list.Car.(*scmgo.Cons)
			if !ok {
				return nil, ErrArguments
			}
			*arg = first
		default:
			return nil, ErrArguments
		}
		_, list, err = list.Pop()
		if err != nil {
			return
		}
	}
	return list, nil
}
