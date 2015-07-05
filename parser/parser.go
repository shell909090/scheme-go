package parser

import (
	"strconv"
	"strings"

	"bitbucket.org/shell909090/scheme-go/scmgo"
)

// control_chars
// ()' \n\r\t

func StringToBoolean(s string) (o scmgo.Boolean, err error) {
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

func StringToNumber(chunk string) (obj scmgo.SchemeObject, err error) {
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

func listToObj(list *scmgo.Cons) (obj scmgo.SchemeObject, err error) {
	if list.Car != nil {
		last, ok := list.Car.(*scmgo.Quote)
		if ok && last.Objs == nil {
			return nil, ErrQuoteInEnd
		}
	}
	return scmgo.ReverseList(list)
}

func popup(ilist, istack *scmgo.Cons) (obj scmgo.SchemeObject, list, stack *scmgo.Cons, err error) {
	if stack == scmgo.Onil {
		return nil, nil, nil, ErrParenthesisNotClose
	}

	obj, err = listToObj(ilist)
	if err != nil {
		log.Error("%s", err)
		return
	}

	t, stack, err := istack.Pop()
	if err != nil {
		log.Error("%s", err)
		return
	}

	list, ok := t.(*scmgo.Cons)
	if !ok {
		err = scmgo.ErrUnknown
	}
	return
}

func Code(cin chan string) (code scmgo.SchemeObject, err error) {
	var obj scmgo.SchemeObject
	list := scmgo.Onil
	stack := scmgo.Onil

	for chunk, ok := <-cin; ok; chunk, ok = <-cin {
		switch chunk[0] {
		case '#': // Boolean
			obj, err = StringToBoolean(chunk)
		case '"': // String
			obj = scmgo.String(chunk[1 : len(chunk)-1])
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if chunk[0] == '-' && len(chunk) == 1 {
				// - without number is symbol
				obj = &scmgo.Symbol{Name: chunk}
			} else { // Integer & Float
				obj, err = StringToNumber(chunk)
			}
		case '\'': // Quote
			obj = new(scmgo.Quote)
		case ';': // Comment
			obj = &scmgo.Comment{Content: chunk[1 : len(chunk)-1]}
		case '(': // Cons
			stack = stack.Push(list)
			list = scmgo.Onil
			continue
		case ')': // return Cons
			obj, list, stack, err = popup(list, stack)
		default: // Symbol
			obj = &scmgo.Symbol{Name: chunk}
		}

		if err != nil {
			log.Error("%s", err)
			return
		}

		// processing Quote
		if list.Car != nil {
			if last, ok := list.Car.(*scmgo.Quote); ok {
				last.Objs = obj
				continue
			}
		}

		list = list.Push(obj)
	}

	if stack != scmgo.Onil {
		return nil, ErrParenthesisNotClose
	}
	return listToObj(list)
}
