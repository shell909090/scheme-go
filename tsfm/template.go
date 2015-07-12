package tsfm

import "bitbucket.org/shell909090/scheme-go/scmgo"

func RenderList(mr *MatchResult, template *scmgo.Cons) (result scmgo.SchemeObject, err error) {
	var ok bool
	var obj scmgo.SchemeObject
	c := &scmgo.Cons{}
	result = c

	for {
		obj, err = Render(mr, template.Car)
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
		c.Car = obj

		if template.Cdr == scmgo.Onil {
			c.Cdr = scmgo.Onil
			return
		}

		template, ok = template.Cdr.(*scmgo.Cons)
		if !ok { // improper
			obj, err = Render(mr, template.Cdr)
			if err != nil {
				log.Error("%s", err.Error())
				return
			}
			c.Cdr = obj
			return
		}

		if isEllipsis(template) {
			obj, ok = mr.m[template.Car.(*scmgo.Symbol).Name]
			if !ok {
				return nil, ErrElpsAfterNoVar
			}
			c.Cdr = obj
			return
		}

		c.Cdr = &scmgo.Cons{}
		c = c.Cdr.(*scmgo.Cons)
	}
	return
}

func Render(mr *MatchResult, template scmgo.SchemeObject) (result scmgo.SchemeObject, err error) {
	var ok bool
	switch tmp := template.(type) {
	case *scmgo.Symbol:
		result, ok = mr.m[tmp.Name]
		if ok {
			return
		}
	case *scmgo.Cons:
		return RenderList(mr, tmp)
	}
	return template, nil
}
