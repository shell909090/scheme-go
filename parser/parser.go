package parser

import (
	"strconv"
	"strings"

	"bitbucket.org/shell909090/scheme-go/scm"
)

func StringToBoolean(b []byte) (o scm.Boolean, err error) {
	if len(b) != 2 || b[0] != '#' {
		return scm.Ofalse, ErrBooleanUnknown
	}
	switch b[1] {
	case 't':
		return scm.Otrue, nil
	case 'f':
		return scm.Ofalse, nil
	}
	return scm.Ofalse, ErrBooleanUnknown
}

func StringToNumber(b []byte) (obj scm.Obj, err error) {
	chunk := string(b)
	if strings.Index(chunk, ".") != -1 {
		var i float64
		i, err = strconv.ParseFloat(chunk, 64)
		obj = scm.Float(i)
	} else {
		var i int
		i, err = strconv.Atoi(chunk)
		obj = scm.Integer(i)
	}
	return
}

func convertDotPair(list *scm.Cons) (result *scm.Cons) {
	f, c, err := list.Pop()
	if err != nil {
		return list
	}

	s, c, err := c.Pop()
	if err != nil {
		return list
	}
	if sym, ok := s.(*scm.Symbol); !ok || sym.Name != "." {
		return list
	} // secondary element not dot

	t, c, err := c.Pop()
	if err != nil {
		return list
	}

	if c != scm.Onil {
		return list
	} // list not end
	return &scm.Cons{Car: f, Cdr: t} // all matched
}

type Parser struct {
	list  *scm.Cons
	stack *scm.Cons
}

func CreateParser() (p *Parser) {
	return &Parser{list: scm.Onil, stack: scm.Onil}
}

func (p *Parser) listToObj() (obj scm.Obj, err error) {
	if p.list.Car != nil {
		last, ok := p.list.Car.(*scm.Quote)
		if ok && last.Objs == nil {
			return nil, ErrQuoteInEnd
		}
	}
	p.list, err = p.list.Reverse(scm.Onil)
	if err != nil {
		log.Error("%s", err)
		return
	}
	return convertDotPair(p.list), nil
}

func (p *Parser) Write(chunk []byte) (n int, err error) {
	var obj scm.Obj
	n = len(chunk)

	switch chunk[0] {
	case '#': // Boolean
		obj, err = StringToBoolean(chunk)
	case '"': // String
		obj = scm.String(string(chunk[1 : len(chunk)-1]))
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if chunk[0] == '-' && len(chunk) == 1 {
			// - without number is symbol
			obj = &scm.Symbol{Name: string(chunk)}
		} else { // Integer & Float
			obj, err = StringToNumber(chunk)
		}
	case '\'': // Quote
		obj = new(scm.Quote)
	case '(': // Cons
		p.stack = p.stack.Push(p.list)
		p.list = scm.Onil
		return
	case ')': // return Cons
		obj, err = p.listToObj()
		if err != nil {
			log.Error("%s", err)
			return
		}
		p.list, p.stack, err = p.stack.PopCons()
	default: // Symbol
		obj = &scm.Symbol{Name: string(chunk)}
	}

	if err != nil {
		log.Error("%s", err)
		return
	}

	// processing Quote
	if p.list.Car != nil {
		if last, ok := p.list.Car.(*scm.Quote); ok {
			last.Objs = obj
			return
		}
	}

	p.list = p.list.Push(obj)
	return
}

func (p *Parser) GetCode() (code scm.Obj, err error) {
	if p.stack != scm.Onil {
		log.Error("%s", scm.Format(p.stack))
		return nil, ErrParenthesisNotClose
	}
	return p.listToObj()
}
