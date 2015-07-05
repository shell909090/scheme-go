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
