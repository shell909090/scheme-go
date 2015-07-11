package scmgo

import "strings"

type Environ struct {
	Formatter
	Parent *Environ
	Names  map[string]SchemeObject
}

func (e *Environ) Format() (r string) {
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

func (e *Environ) Fork() (ne *Environ) {
	return &Environ{Parent: e, Names: make(map[string]SchemeObject)}
}

func (e *Environ) Add(name string, value SchemeObject) {
	log.Info("add %s in environ %p", name, e)
	e.Names[name] = value
}

func (e *Environ) Get(name string) (value SchemeObject) {
	var ok bool
	log.Info("get %s in environ %p", name, e)
	for ce := e; ce != nil; ce = ce.Parent {
		value, ok = ce.Names[name]
		if ok {
			return
		}
	}
	return nil
}
