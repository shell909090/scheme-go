package scmgo

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var EmptySpace = []byte("                                        ")

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
		strs = append(strs, Format(c.Car))
		tmp, ok := c.Cdr.(*Cons)
		if !ok {
			strs = append(strs, ".")
			strs = append(strs, Format(c.Cdr))
			break
		}
		c = tmp
	}
	return "(" + strings.Join(strs, " ") + ")"
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

func PrettyFormat(s io.Writer, o Obj, iv int) (rv int, err error) {
	var str string

	c, ok := o.(*Cons)
	if !ok { // normal objects
		str = Format(o)
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

func Format(o Obj) (r string) {
	switch t := o.(type) {
	case *Symbol:
		return t.Name
	case *Quote:
		return "'" + Format(t.Objs)
	case *Cons:
		buf := bytes.NewBuffer(nil)
		if _, err := PrettyFormat(buf, t, 0); err != nil {
			log.Error("%s", err)
			return ""
		}
		return buf.String()
	case Boolean:
		if t {
			return "#t"
		} else {
			return "#f"
		}
	case Integer:
		return strconv.FormatInt(int64(t), 10)
	case Float:
		return strconv.FormatFloat(float64(t), 'f', 2, 64)
	case String:
		return fmt.Sprintf("\"%s\"", t)
	case *InternalProcedure:
		return "!" + t.Name
	case *LambdaProcedure:
		name := t.Name
		if name == "" {
			name = "lambda"
		}
		return "<" + name + ">"
	}
	return fmt.Sprintf("%v", o)
}

func EnvFormat(e *Environ) (r string) {
	str := make([]string, 0)
	for ce := e; ce != nil; ce = ce.Parent {
		strname := make([]string, 0)
		for k, _ := range ce.Names {
			strname = append(strname, k)
		}
		str = append(str, strings.Join(strname, " "))
	}
	return strings.Join(str, "\n")
}

func StackFormat(f Frame) (r string) {
	buf := bytes.NewBuffer(nil)
	for ; f != nil; f = f.GetParent() {
		if _, ok := f.(*EndFrame); !ok {
			buf.WriteString(FrameFormat(f) + "\n")
		}
	}
	return buf.String()
}

func FrameFormat(f Frame) (r string) {
	switch t := f.(type) {
	case *EndFrame:
		return "End"
	case *BeginFrame:
		n, err := t.Obj.Len(false)
		if err != nil {
			n = 0
		}
		return fmt.Sprintf("Begin: %d", n)
	case *ApplyFrame:
		return "Apply"
	case *IfFrame:
		return fmt.Sprintf("If:\n%s", Format(t.Cond))
	}
	return "Unknown Frame"
}
