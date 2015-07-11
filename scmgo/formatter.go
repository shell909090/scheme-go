package scmgo

import (
	"errors"
	"io"
	"strings"
)

var ErrQuit = errors.New("quit")

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

func formatMultiLineList(s io.Writer, list *Cons, iv int) (rv int, err error) {
	// iv for Input leVel, rv for Return leVel.
	var tmp SchemeObject
	obj := list
	s.Write([]byte("("))
	rv = iv + 1

	if _, ok := obj.Car.(*Symbol); ok {
		tmp, obj, err = obj.Pop()
		if err != nil {
			log.Error("%s", err)
			return
		}

		rv, err = PrettyFormat(s, tmp, rv)
		if err != nil {
			log.Error("%s", err)
			return
		}

		// this is multi-line, and first element is symbol.
		// so there MUST have other element,
		// which at least one of them is an list.
		s.Write([]byte(" "))
		rv += 1
	}

	for iv = rv; obj != Onil; { // store align rv in iv
		rv, err = PrettyFormat(s, obj.Car, rv)
		if err != nil {
			log.Error("%s", err)
			return
		}

		if obj.Cdr != Onil {
			s.Write([]byte("\n"))
			for i := 0; i < iv; i++ {
				s.Write([]byte(" "))
			} // FIXME: not so good
			rv = iv
		}

		c, ok := obj.Cdr.(*Cons)
		if !ok {
			s.Write([]byte(" . "))
			rv += 3

			rv, err = PrettyFormat(s, obj.Cdr, rv)
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
