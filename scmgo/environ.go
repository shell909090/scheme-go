package scmgo

type Environ struct {
	Parent *Environ
	Names  map[string]SchemeObject
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
