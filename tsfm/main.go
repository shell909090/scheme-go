package tsfm

import (
	"bitbucket.org/shell909090/scheme-go/scmgo"
	logging "github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("trans")
)

func Transform(src scmgo.SchemeObject) (code scmgo.SchemeObject, err error) {
	code = src
	// c, ok := code.(*scmgo.Cons)
	// if !ok {
	// 	return nil, scmgo.ErrUnknown
	// }
	return
}
