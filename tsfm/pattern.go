package tsfm

import "bitbucket.org/shell909090/scheme-go/scm"

func MatchList(plist, olist *scm.Cons, literals Literals, mr *MatchResult) (yes bool, err error) {
	for plist != scm.Onil && olist != scm.Onil {
		if isEllipsis(plist) {
			// ellipsis, capture rest into varible and return.
			mr.Add(plist.Car.(*scm.Symbol).Name, olist)
			return true, nil
		}

		yes, err = Match(plist.Car, olist.Car, literals, mr)
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
		if !yes {
			return false, nil
		}

		pnext, ok := plist.Cdr.(*scm.Cons)
		if !ok { // improper
			if _, ok = olist.Cdr.(*scm.Cons); ok {
				return false, nil
			}
			return Match(plist.Cdr, olist.Cdr, literals, mr)
		}

		plist = pnext
		olist, ok = olist.Cdr.(*scm.Cons)
		if !ok {
			return false, nil
		}
	}
	if isEllipsis(plist) {
		mr.Add(plist.Car.(*scm.Symbol).Name, olist)
		return true, nil
	}
	return olist == scm.Onil && plist == scm.Onil, nil
}

func Match(pattern, obj scm.Obj, literals Literals, mr *MatchResult) (yes bool, err error) {
	switch tmp := pattern.(type) {
	case *scm.Symbol:
		switch {
		case tmp.Name == "_":
			return true, nil
		case literals.CheckLiteral(tmp.Name):
			tmpo, ok := obj.(*scm.Symbol)
			if !ok {
				return false, nil
			}
			return tmpo.Name == tmp.Name, nil
		default:
			mr.Add(tmp.Name, obj)
			return true, nil
		}
	case *scm.Cons:
		tmpo, ok := obj.(*scm.Cons)
		if !ok {
			return false, nil
		}
		return MatchList(tmp, tmpo, literals, mr)
	}
	return false, ErrWrongTpInPtn
}
