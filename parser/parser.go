package parser

import (
	"strconv"
	"strings"

	"bitbucket.org/shell909090/scheme-go/scmgo"
)

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

func convertDotPair(list *scmgo.Cons) (result *scmgo.Cons) {
	f, c, err := list.Pop(false)
	if err != nil {
		return list
	}

	s, c, err := c.Pop(false)
	if err != nil {
		return list
	}
	if sym, ok := s.(*scmgo.Symbol); !ok || sym.Name != "." {
		return list
	} // secondary element not dot

	t, c, err := c.Pop(false)
	if err != nil {
		return list
	}

	if c != scmgo.Onil {
		return list
	} // list not end
	return &scmgo.Cons{Car: f, Cdr: t} // all matched
}

type Parser struct {
	list  *scmgo.Cons
	stack *scmgo.Cons
}

func CreateParser() (p *Parser) {
	return &Parser{list: scmgo.Onil, stack: scmgo.Onil}
}

func (p *Parser) listToObj() (obj scmgo.SchemeObject, err error) {
	if p.list.Car != nil {
		last, ok := p.list.Car.(*scmgo.Quote)
		if ok && last.Objs == nil {
			return nil, ErrQuoteInEnd
		}
	}
	p.list, err = scmgo.ReverseList(p.list)
	if err != nil {
		log.Error("%s", err)
		return
	}
	return convertDotPair(p.list), nil
}

func (p *Parser) popup() (obj scmgo.SchemeObject, err error) {
	var ok bool
	var t scmgo.SchemeObject

	if p.stack == scmgo.Onil {
		return nil, ErrParenthesisNotClose
	}

	obj, err = p.listToObj()
	if err != nil {
		log.Error("%s", err)
		return
	}

	t, p.stack, err = p.stack.Pop(false)
	if err != nil {
		log.Error("%s", err)
		return
	}

	p.list, ok = t.(*scmgo.Cons)
	if !ok {
		err = scmgo.ErrUnknown
	}
	return
}

func (p *Parser) Write(b []byte) (n int, err error) {
	err = p.ReceiveChunk(string(b))
	if err != nil {
		return
	}
	n = len(b)
	return
}

func (p *Parser) ReceiveChunk(chunk string) (err error) {
	var obj scmgo.SchemeObject

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
	// Comment: no comments now
	// case ';':
	// 	obj = &scmgo.Comment{Content: chunk[1 : len(chunk)-1]}
	// 	return
	case '(': // Cons
		p.stack = p.stack.Push(p.list)
		p.list = scmgo.Onil
		return
	case ')': // return Cons
		obj, err = p.popup()
	default: // Symbol
		obj = &scmgo.Symbol{Name: chunk}
	}

	if err != nil {
		log.Error("%s", err)
		return
	}

	// processing Quote
	if p.list.Car != nil {
		if last, ok := p.list.Car.(*scmgo.Quote); ok {
			last.Objs = obj
			return
		}
	}

	p.list = p.list.Push(obj)
	return
}

func (p *Parser) GetCode() (code scmgo.SchemeObject, err error) {
	if p.stack != scmgo.Onil {
		return nil, ErrParenthesisNotClose
	}
	return p.listToObj()
}
