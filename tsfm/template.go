package tsfm

// FIXME: render with ...
func Render(mr *MatchResult, template scmgo.SchemeObject) (result scmgo.SchemeObject, err error) {
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
