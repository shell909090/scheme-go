package tsfm

import (
	"errors"
	"fmt"
	"strings"

	"bitbucket.org/shell909090/scheme-go/scm"
	logging "github.com/op/go-logging"
)

var (
	ErrNotAPattern    = errors.New("not a pattern")
	ErrWrongTpInPtn   = errors.New("wrong type in pattern")
	ErrWrongStruct    = errors.New("wrong struct")
	ErrSyntaxExist    = errors.New("syntax exist")
	ErrElpsAfterNoVar = errors.New("ellipsis after something not a pattern varible")
	ErrNoRule         = errors.New("no rule match in syntax")
)

var (
	log               = logging.MustGetLogger("trans")
	DefineTransformer = &Transformer{syntaxes: make(map[string]*Syntax)}
)

type Literals map[string]int

func ReadLiterals(l *scm.Cons) (literals Literals, err error) {
	var s *scm.Symbol
	literals = make(Literals, 0)
	for l != scm.Onil {
		s, l, err = l.PopSymbol()
		if err != nil {
			return
		}
		literals[s.Name] = 1
	}
	return
}

func (l Literals) CheckLiteral(s string) (yes bool) {
	_, yes = l[s]
	return
}

type MatchResult struct {
	m map[string]scm.Obj
}

func CreateMatchResult() (m *MatchResult) {
	m = &MatchResult{
		m: make(map[string]scm.Obj),
	}
	return
}

func (m *MatchResult) Add(name string, value scm.Obj) {
	m.m[name] = value
}

func (m *MatchResult) Format() (r string) {
	var strs []string
	for name, value := range m.m {
		strs = append(strs, fmt.Sprintf("%s = %s", name, scm.Format(value)))
	}
	return strings.Join(strs, "\n")
}

func isEllipsis(plist *scm.Cons) (yes bool) {
	_, ok := plist.Car.(*scm.Symbol)
	if !ok {
		return false
	}
	next, err := plist.GetN(1)
	if err != nil {
		return false
	}
	next_sym, ok := next.(*scm.Symbol)
	if !ok {
		return false
	}
	return next_sym.Name == "..."
}
