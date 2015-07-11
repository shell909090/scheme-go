package scmgo

import (
	"io"
	"strings"
)

var EmptySpace = []byte("                                        ")

func AnyList(o *Cons) (yes bool) {
	ok := true
	for i := o; i != Onil; {
		if _, yes = i.Car.(*Cons); yes {
			return true
		}
		i, ok = i.Cdr.(*Cons)
		if !ok {
			return false
		}
	}
	return false
}

func formatOneLineList(list *Cons) (r string) {
	strs := make([]string, 0)
	for c := list; c != Onil; {
		strs = append(strs, c.Car.Format())
		tmp, ok := c.Cdr.(*Cons)
		if !ok {
			strs = append(strs, ".")
			strs = append(strs, c.Cdr.Format())
			break
		}
		c = tmp
	}
	return "(" + strings.Join(strs, " ") + ")"
}

func fullSpace(s io.Writer, iv int) {
	for i := 0; i < iv; i++ {
		if iv-i > len(EmptySpace) {
			s.Write(EmptySpace)
			i += len(EmptySpace)
		} else {
			s.Write(EmptySpace[:iv-i])
			i = iv
		}
	}
}

func formatMultiLineList(s io.Writer, list *Cons, iv int) (rv int, err error) {
	// iv for Input leVel, rv for Return leVel.
	obj := list
	s.Write([]byte("("))
	rv = iv + 1
	iv += 3

	for obj != Onil { // store align rv in iv
		rv, err = PrettyFormat(s, obj.Car, rv)
		if err != nil {
			log.Error("%s", err)
			return
		}

		if obj.Cdr != Onil {
			s.Write([]byte("\n"))
			fullSpace(s, iv)
			rv = iv
		}

		c, ok := obj.Cdr.(*Cons)
		if !ok {
			s.Write([]byte(" . "))
			rv, err = PrettyFormat(s, obj.Cdr, rv+3)
			if err != nil {
				log.Error("%s", err)
				return
			}
			break
		}
		obj = c
	}

	s.Write([]byte(")"))
	rv += 1
	return
}

func PrettyFormat(s io.Writer, o SchemeObject, iv int) (rv int, err error) {
	var str string

	c, ok := o.(*Cons)
	if !ok { // normal objects
		str = o.Format()
		s.Write([]byte(str))
		return iv + len(str), err
	}

	if c.Car == nil || c.Cdr == nil { // Onil
		// actually, Onil can be format right in formatOneLineList.
		// we do this for short path.
		_, err = s.Write([]byte("()"))
		return iv + 2, err
	}

	if AnyList(c) {
		return formatMultiLineList(s, c, iv)
	}

	// all element in one line
	str = formatOneLineList(c)
	s.Write([]byte(str))
	return iv + len(str), err
}
