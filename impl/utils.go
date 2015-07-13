package impl

import (
	"errors"

	"bitbucket.org/shell909090/scheme-go/scm"
	logging "github.com/op/go-logging"
)

var (
	ErrArguments = errors.New("wrong arguments")
)

var (
	log = logging.MustGetLogger("impl")
)

func AssertLen(o *scm.Cons, length int) (err error) {
	n, err := o.Len(false)
	if err != nil {
		return
	}
	if n != length {
		return ErrArguments
	}
	return
}

func ParseParameters(list *scm.Cons, v ...interface{}) (next *scm.Cons, err error) {
	for _, a := range v {
		switch arg := a.(type) {
		case *string:
			switch first := list.Car.(type) {
			case *scm.Symbol:
				*arg = first.Name
			case *scm.String:
				*arg = string(*first)
			default:
				return nil, ErrArguments
			}
		case **scm.Cons:
			first, ok := list.Car.(*scm.Cons)
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
