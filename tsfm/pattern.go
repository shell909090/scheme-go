package tsfm

import "bitbucket.org/shell909090/scheme-go/scmgo"

func isEllipsis(plist *scmgo.Cons) (yes bool) {
	_, ok := plist.Car.(*scmgo.Symbol)
	if !ok {
		return false
	}
	next, err := plist.GetN(1)
	if err != nil {
		return false
	}
	next_sym, ok := next.(*scmgo.Symbol)
	if !ok {
		return false
	}
	return next_sym.Name == "..."
}

func MatchList(plist, olist *scmgo.Cons, literals Literals, mr *MatchResult) (yes bool, err error) {
	var ok bool
	for plist != scmgo.Onil && olist != scmgo.Onil {
		log.Debug("now match: %s %s", plist.Format(), olist.Format())

		if isEllipsis(plist) {
			// ellipsis, capture rest into varible and return.
			mr.Add(plist.Car.(*scmgo.Symbol).Name, olist)
			return
		}

		yes, err = Match(plist.Car, olist.Car, literals, mr)
		if err != nil {
			log.Error("%s", err.Error())
			return
		}
		if !yes {
			return false, nil
		}

		switch pnext := plist.Cdr.(type) {
		case *scmgo.Cons: // continue on list.
			plist = pnext
			olist, ok = olist.Cdr.(*scmgo.Cons)
			if !ok {
				return false, nil
			}
		default: // pair in the end.
			if _, ok = olist.Cdr.(*scmgo.Cons); ok {
				return false, nil
			}
			return Match(plist.Cdr, olist.Cdr, literals, mr)
		}
	}
	return olist == scmgo.Onil && plist == scmgo.Onil, nil
}

func Match(pattern, obj scmgo.SchemeObject, literals Literals, mr *MatchResult) (yes bool, err error) {
	switch tmp := pattern.(type) {
	case *scmgo.Symbol:
		switch {
		case tmp.Name == "_":
			return true, nil
		case literals.CheckLiteral(tmp.Name):
			tmpo, ok := obj.(*scmgo.Symbol)
			if !ok {
				return false, nil
			}
			return tmpo.Name == tmp.Name, nil
		default:
			mr.Add(tmp.Name, obj)
			return true, nil
		}
	case *scmgo.Cons:
		tmpo, ok := obj.(*scmgo.Cons)
		if !ok {
			return false, nil
		}
		return MatchList(tmp, tmpo, literals, mr)
	}
	return false, ErrWrongTpInPtn
}
