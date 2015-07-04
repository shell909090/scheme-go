package scmgo

import (
	"errors"
	"io"
	"strings"
)

var ErrQuit = errors.New("quit")

func AnyList(o *Cons) (yes bool, err error) {
	err = o.Iter(func(obj SchemeObject) (e error) {
		_, yes = obj.(*Cons)
		if yes {
			e = ErrQuit
		}
		return
	})
	if err == ErrQuit {
		err = nil
	}
	return
}

func formatOneLineList(o *Cons) (r string, err error) {
	strs := make([]string, 0)
	err = o.Iter(func(obj SchemeObject) (e error) {
		strs = append(strs, obj.Format())
		return
	})
	if err != nil {
		log.Error("%s", err)
		return
	}
	return "(" + strings.Join(strs, " ") + ")", nil
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

		if obj != Onil {
			s.Write([]byte(" "))
			rv += 1
		}
	}

	iv = rv
	for ok := true; obj != Onil; obj, ok = obj.Cdr.(*Cons) {
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

	anycons, err := AnyList(o)
	if err != nil {
		log.Error("%s", err)
		return
	}
	if !anycons { // all element in one line
		str, err = formatOneLineList(o)
		if err != nil {
			log.Error("%s", err)
			return
		}
		s.Write([]byte(str))
		return iv + len(str), err
	}

	return formatMultiLineList(s, o, iv)
}
