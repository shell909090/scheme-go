package parser

import (
	"errors"
	"io"

	logging "github.com/op/go-logging"
	"github.com/shell909090/scheme-go/scm"
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

func SourceToAST(source io.Reader) (code scm.Obj, err error) {
	p := CreateParser()
	err = Grammar(source, p)
	if err != nil {
		return
	}
	return p.GetCode()
}
