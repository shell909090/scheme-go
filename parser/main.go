package parser

import (
	"errors"
	"io"

	"bitbucket.org/shell909090/scheme-go/scmgo"
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
	log = logging.MustGetLogger("parser")
)

func SourceToAST(source io.ReadCloser) (code scmgo.SchemeObject, err error) {
	p := CreateParser()
	err = Grammar(source, p)
	if err != nil {
		return
	}
	return p.GetCode()
}
