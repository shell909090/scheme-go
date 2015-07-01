package scmgo

import (
	"errors"
	"io"
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
)

func BuildCode(source io.ReadCloser) (code SchemeObject, err error) {
	cpipe := make(chan string)
	go GrammarParser(source, cpipe)
	// for chunk, ok := <-cpipe; ok; chunk, ok = <-cpipe {
	// 	fmt.Println("chunk:", string(chunk))
	// }
	// return nil, nil
	return CodeParser(cpipe)
}
