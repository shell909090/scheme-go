package scmgo

import (
	"bytes"
	"errors"
	"io"

	logging "github.com/op/go-logging"
)

var (
	ErrQuotaNotClose       = errors.New("quote not closed")
	ErrQuotaInSymbol       = errors.New("quote in symbol")
	ErrCommentInSymbol     = errors.New("comment in symbol")
	ErrParenthesisNotClose = errors.New("parenthesis not close")
	ErrBooleanUnknown      = errors.New("unknown boolean")
	ErrQuoteInEnd          = errors.New("quote in the end of S-Expression")
)

var (
	ErrListOutOfIndex = errors.New("out of index when get list")
	ErrType           = errors.New("runtime type error")
	ErrISNotAList     = errors.New("object is not a list")
	ErrUnknown        = errors.New("unknown error")
	ErrNameNotFound   = errors.New("name not found")
	ErrNotRunnable    = errors.New("object not runnable")
	ErrArguments      = errors.New("wrong arguments")
)

var (
	log          = logging.MustGetLogger("scmgo")
	DefaultNames = make(map[string]SchemeObject)
)

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

func Trampoline(init_frame Frame, init_obj SchemeObject) (result SchemeObject, err error) {
	f := init_frame
	result = init_obj
	for f != nil {
		log.Debug("frame: %s, result: %s.",
			f.Debug(), SchemeObjectToString(result))

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
		return nil, ErrType
	}
	env := &Environ{Parent: nil, Names: DefaultNames}
	init_frame := CreatePrognFrame(nil, progn, env)

	result, err = Trampoline(init_frame, nil)
	if result == nil {
		return nil, ErrUnknown
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
