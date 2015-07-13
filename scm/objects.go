package scm

type Obj interface{}

type Symbol struct {
	Name string
}

type Quote struct {
	Objs Obj
}

type Boolean bool

const (
	Otrue  = Boolean(true)
	Ofalse = Boolean(false)
)

type Integer int
type Float float64
type String string

type Cons struct {
	Car Obj
	Cdr Obj
}

var Onil = &Cons{}

func (o *Cons) Pop() (r Obj, next *Cons, err error) {
	if o == Onil {
		return nil, nil, ErrListOutOfIndex
	}
	r = o.Car
	next, ok := o.Cdr.(*Cons)
	if !ok {
		return nil, nil, ErrISNotAList
	}
	return
} // O(1)

func (o *Cons) Push(i Obj) (next *Cons) {
	return &Cons{Car: i, Cdr: o}
} // O(1)

func (o *Cons) IsImproper() bool {
	ok := true
	for i := o; i != Onil; i, ok = i.Cdr.(*Cons) {
		if !ok {
			return true
		}
	}
	return false
} // O(n)

func (o *Cons) Iter(f func(obj Obj) (e error), improper bool) (err error) {
	ok := true
	for i := o; i != Onil; {
		err = f(i.Car)
		if err != nil {
			return
		}
		i, ok = i.Cdr.(*Cons)
		if !ok {
			if !improper {
				return ErrISNotAList
			}
			return f(i.Cdr)
		}
	}
	return
} // O(n)

func (o *Cons) Len(improper bool) (n int, err error) {
	err = o.Iter(func(obj Obj) (e error) {
		n += 1
		return
	}, improper)
	return
} // O(n)

func (o *Cons) GetN(n int) (r Obj, err error) {
	var ok bool
	c := o
	for i := 0; i < n; i++ {
		switch c.Cdr {
		case nil:
			return nil, ErrUnknown
		case Onil:
			return nil, ErrListOutOfIndex
		}
		c, ok = o.Cdr.(*Cons)
		if !ok {
			return nil, ErrISNotAList
		}
	}
	return c.Car, nil
} // O(n)

func (o *Cons) PopSymbol() (s *Symbol, next *Cons, err error) {
	t, next, err := o.Pop()
	if err != nil {
		return
	}
	s, ok := t.(*Symbol)
	if !ok {
		return nil, nil, ErrType
	}
	return
} // O(1)

func (o *Cons) PopCons() (s *Cons, next *Cons, err error) {
	t, next, err := o.Pop()
	if err != nil {
		return
	}
	s, ok := t.(*Cons)
	if !ok {
		return nil, nil, ErrType
	}
	return
} // O(1)
