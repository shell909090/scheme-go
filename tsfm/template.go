package main

import "github.com/shell909090/scheme-go/scm"

func RenderList(mr *MatchResult, template *scm.Cons) (result scm.Obj, err error) {
	var ok bool
	var obj scm.Obj
	c := &scm.Cons{}
	result = c

	for {
		obj, err = Render(mr, template.Car)
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
		c.Car = obj

		if template.Cdr == scm.Onil {
			c.Cdr = scm.Onil
			return
		}

		template, ok = template.Cdr.(*scm.Cons)
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
			obj, ok = mr.m[template.Car.(*scm.Symbol).Name]
			if !ok {
				return nil, ErrElpsAfterNoVar
			}
			c.Cdr = obj
			return
		}

		c.Cdr = &scm.Cons{}
		c = c.Cdr.(*scm.Cons)
	}
	return
}

func Render(mr *MatchResult, template scm.Obj) (result scm.Obj, err error) {
	var ok bool
	switch tmp := template.(type) {
	case *scm.Symbol:
		result, ok = mr.m[tmp.Name]
		if ok {
			return
		}
	case *scm.Cons:
		return RenderList(mr, tmp)
	}
	return template, nil
}
