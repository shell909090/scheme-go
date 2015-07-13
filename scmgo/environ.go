package scmgo

type Environ struct {
	Parent *Environ
	Names  map[string]Obj
}

func (e *Environ) Fork() (ne *Environ) {
	return &Environ{Parent: e, Names: make(map[string]Obj)}
}

func (e *Environ) Len() (n int) {
	for n = 0; e != nil; e = e.Parent {
		n += 1
	}
	return
}

func (e *Environ) Add(name string, value Obj) {
	log.Debug("add %s in environ length %d", name, e.Len())
	e.Names[name] = value
}

func (e *Environ) Get(name string) (value Obj) {
	var ok bool
	log.Debug("get %s in environ length %d", name, e.Len())
	for ce := e; ce != nil; ce = ce.Parent {
		value, ok = ce.Names[name]
		if ok {
			return
		}
	}
	return nil
}
