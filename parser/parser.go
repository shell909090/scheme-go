package parser

import (
	"strconv"
	"strings"

	"bitbucket.org/shell909090/scheme-go/scmgo"
)

// control_chars
// ()' \n\r\t

func BooleanFromString(s string) (o scmgo.Boolean, err error) {
	if len(s) != 2 || s[0] != '#' {
		return scmgo.Ofalse, ErrBooleanUnknown
	}
	switch s[1] {
	case 't':
		return scmgo.Otrue, nil
	case 'f':
		return scmgo.Ofalse, nil
	}
	return scmgo.Ofalse, ErrBooleanUnknown
}

func CodeNumber(chunk string) (obj scmgo.SchemeObject, err error) {
	if strings.Index(chunk, ".") != -1 {
		var i float64
		i, err = strconv.ParseFloat(chunk, 64)
		obj = scmgo.Float(i)
	} else {
		var i int
		i, err = strconv.Atoi(chunk)
		obj = scmgo.Integer(i)
	}
	return
}

func CodeParser(cin chan string) (code scmgo.SchemeObject, err error) {
	var obj scmgo.SchemeObject
	var objs []scmgo.SchemeObject

QUIT:
	for chunk, ok := <-cin; ok; chunk, ok = <-cin {
		switch chunk[0] {
		case '#': // Boolean
			obj, err = BooleanFromString(chunk)
			if err != nil {
				return
			}
		case '"': // String
			obj = scmgo.String(chunk[1 : len(chunk)-1])
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if chunk[0] == '-' && len(chunk) == 1 {
				// - without number is symbol
				obj = &scmgo.Symbol{Name: chunk}
			} else { // Integer & Float
				obj, err = CodeNumber(chunk)
				if err != nil {
					return
				}
			}
		case '\'': // Quote
			obj = new(scmgo.Quote)
		case ';': // Comment
			obj = nil
		case '(': // Cons
			obj, err = CodeParser(cin)
			if err != nil {
				return nil, err
			}
		case ')': // return Cons
			break QUIT
		default: // Symbol
			obj = &scmgo.Symbol{Name: chunk}
		}

		if obj == nil { // pass comment
			continue
		}

		// processing Quote
		if len(objs) > 0 {
			o := objs[len(objs)-1]
			if last, ok := o.(*scmgo.Quote); ok {
				last.Objs = obj
				continue
			}
		}

		objs = append(objs, obj)
	}

	if len(objs) > 0 {
		o := objs[len(objs)-1]
		if last, ok := o.(*scmgo.Quote); ok && last.Objs == nil {
			return code, ErrQuoteInEnd
		}
	}

	return ListFromSlice(objs), nil
}

func ListFromSlice(s []scmgo.SchemeObject) (o scmgo.SchemeObject) {
	o = scmgo.Onil
	for i := len(s) - 1; i >= 0; i-- {
		o = &scmgo.Cons{Car: s[i], Cdr: o}
	}
	return o
}
