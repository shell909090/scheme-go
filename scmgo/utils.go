package scmgo

import (
	"errors"
	"io"

	"github.com/op/go-logging"
)

var (
	ErrQuotaNotClose       = errors.New("quote not closed")
	ErrQuotaInSymbol       = errors.New("quote in symbol")
	ErrCommentInSymbol     = errors.New("comment in symbol")
	ErrParenthesisNotClose = errors.New("parenthesis not close")
	ErrBooleanUnknown      = errors.New("unknown boolean")
	ErrQuoteInEnd          = errors.New("quote in the end of S-Expression")
	ErrListOutOfIndex      = errors.New("out of index when get list")
	ErrRuntimeType         = errors.New("runtime type error")
	ErrISNotAList          = errors.New("object is not a list")
	ErrRuntimeUnknown      = errors.New("runtime unknown error")
	ErrNameNotFound        = errors.New("name not found")
)

var (
	log = logging.MustGetLogger("scmgo")
)

func Trampoline(init_frame Frame, init_obj SchemeObject) (result SchemeObject, err error) {
	var name string
	f := init_frame
	result = init_obj
	for f != nil {
		name, err = f.Debug()
		if err != nil {
			return
		}
		log.Debug("result: %v, frame: %s %v.", result, name, f)

		result, f, err = f.Exec(result)
		if err != nil {
			log.Error("%s", err)
			return
		}
	}
	return
}

func RunCode(code SchemeObject) (result SchemeObject, err error) {
	progn, ok := code.(*Cons)
	if !ok {
		return nil, ErrRuntimeType
	}
	env := &Environ{Parent: nil, Names: make(map[string]SchemeObject)}
	init_frame := CreatePrognFrame(progn, env)

	result, err = Trampoline(init_frame, nil)
	if result == nil {
		return nil, ErrRuntimeUnknown
	}
	return
}

func BuildCode(source io.ReadCloser) (code SchemeObject, err error) {
	cpipe := make(chan string)
	go GrammarParser(source, cpipe)
	// for chunk, ok := <-cpipe; ok; chunk, ok = <-cpipe {
	// 	fmt.Println("chunk:", string(chunk))
	// }
	// return nil, nil
	return CodeParser(cpipe)
}
