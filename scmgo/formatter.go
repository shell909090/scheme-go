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
			return yes
		}
		i, ok = i.Cdr.(*Cons)
		if !ok {
			return false
		}
	}
	return
}

func formatOneLineList(o *Cons) (r string) {
	strs := make([]string, 0)
	for i := o; i != Onil; {
		strs = append(strs, i.Car.Format())
		t, ok := i.Cdr.(*Cons)
		if !ok {
			strs = append(strs, ".")
			strs = append(strs, i.Cdr.Format())
			break
		}
		i = t
	}
	return "(" + strings.Join(strs, " ") + ")"
}

func formatMultiLineList(s io.Writer, o *Cons, iv int) (rv int, err error) {
	var t SchemeObject
	obj := o
	s.Write([]byte("("))
	rv = iv + 1

	if _, ok := obj.Car.(*Symbol); ok {
		t, obj, err = obj.Pop()
		if err != nil {
			log.Error("%s", err)
			return
		}

		rv, err = PrettyFormat(s, t, rv)
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

	for iv = rv; obj != Onil; {
		rv, err = PrettyFormat(s, obj.Car, rv)
		if err != nil {
			log.Error("%s", err)
			return
		}

		if obj.Cdr != Onil {
			s.Write([]byte("\n"))
			for i := 0; i < iv; i++ {
				s.Write([]byte(" "))
			}
			rv = iv
		}

		o, ok := obj.Cdr.(*Cons)
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
		obj = o
	}

	s.Write([]byte(")"))
	rv += 1
	return
}

func PrettyFormat(s io.Writer, i SchemeObject, iv int) (rv int, err error) {
	var str string

	o, ok := i.(*Cons)
	if !ok { // normal objects
		str = i.Format()
		s.Write([]byte(str))
		return iv + len(str), err
	}

	if o.Car == nil || o.Cdr == nil { // Onil
		_, err = s.Write([]byte("()"))
		return iv + 2, err
	}

	if !AnyList(o) { // all element in one line
		str = formatOneLineList(o)
		s.Write([]byte(str))
		return iv + len(str), err
	}

	return formatMultiLineList(s, o, iv)
}
