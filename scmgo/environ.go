package scmgo

import "strings"

type Environ struct {
	Formatter
	Parent *Environ
	Names  map[string]SchemeObject
}

func formatNames(names map[string]SchemeObject) (r string) {
	str := make([]string, 0)
	for k, _ := range names {
		str = append(str, k)
	}
	return strings.Join(str, " ")
}

func (e *Environ) Format() (r string) {
	str := make([]string, 0)
	for ce := e; ce != nil; ce = ce.Parent {
		str = append(str, formatNames(ce.Names))
	}
	return strings.Join(str, "\n")
}

func (e *Environ) Fork(r map[string]SchemeObject) (ne *Environ) {
	if r == nil {
		r = make(map[string]SchemeObject)
	}
	ne = &Environ{
		Parent: e,
		Names:  r,
	}
	return ne
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
