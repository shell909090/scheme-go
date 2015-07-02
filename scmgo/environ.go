package scmgo

type Environ struct {
	Parent *Environ
	Names  map[string]SchemeObject
	// Fast   map[string]SchemeObject
}

// func (e *Env) GenFast() {
// }

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
	e.Names[name] = value
	// e.Fast[name] = value
}

func (e *Environ) Get(name string) (value SchemeObject, ok bool) {
	for ce := e; ce != nil; ce = e.Parent {
		value, ok = e.Names[name]
		if ok {
			return
		}
	}
	return nil, false
}
