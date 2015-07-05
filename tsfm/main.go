package tsfm

import (
	"bitbucket.org/shell909090/scheme-go/scmgo"
	logging "github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("trans")
)

func RemoveComment(o *scmgo.Cons) (err error) {
	err = FilterList(o, func(i scmgo.SchemeObject) (yes bool, err error) {
		if _, ok := i.(*scmgo.Comment); ok {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		log.Error("%s", err)
		return
	}

	return
}

func Transform(src scmgo.SchemeObject) (code scmgo.SchemeObject, err error) {
	code = src
	c, ok := code.(*scmgo.Cons)
	if !ok {
		return nil, scmgo.ErrUnknown
	}
	err = Walker(c, RemoveComment)
	if err != nil {
		log.Error("%s", err)
		return
	}
	return
}
