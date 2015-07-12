package tsfm

import (
	"fmt"
	"strings"

	"bitbucket.org/shell909090/scheme-go/scmgo"
)

type MatchResult struct {
	m map[string]scmgo.SchemeObject
}

func CreateMatchResult() (m *MatchResult) {
	m = &MatchResult{
		m: make(map[string]scmgo.SchemeObject),
	}
	return
}

func (m *MatchResult) Add(name string, value scmgo.SchemeObject) {
	m.m[name] = value
}

// FIXME: render with ...
func (m *MatchResult) Copy(obj scmgo.SchemeObject) (result scmgo.SchemeObject, err error) {
	var ok bool
	var i scmgo.SchemeObject
	switch t := obj.(type) {
	case *scmgo.Symbol:
		if t.Name == "..." {

		}
		if o, ok := m.m[t.Name]; ok {
			return o, nil
		}
	case *scmgo.Cons:
		c := &scmgo.Cons{}
		result = c
		for {
			i, err = m.Copy(t.Car)
			if err != nil {
				log.Error("%s", err.Error())
				return
			}
			c.Car = i

			if t.Cdr == scmgo.Onil {
				c.Cdr = scmgo.Onil
				return
			}

			t, ok = t.Cdr.(*scmgo.Cons)
			if !ok { // improper
				i, err = m.Copy(t.Cdr)
				if err != nil {
					log.Error("%s", err.Error())
					return
				}
				c.Cdr = i
			}
			c.Cdr = &scmgo.Cons{}
			c = c.Cdr.(*scmgo.Cons)
		}
		return
	}
	return obj, nil
}

func (m *MatchResult) Format() (r string) {
	var strs []string
	for name, value := range m.m {
		strs = append(strs, fmt.Sprintf("%s = %s", name, value.Format()))
	}
	return strings.Join(strs, "\n")
}
