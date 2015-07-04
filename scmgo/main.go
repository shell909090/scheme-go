package scmgo

func Trampoline(f Frame) (result SchemeObject, err error) {
	for {
		log.Debug("stack:\n%s", StackFormatter(f))
		f, err = f.Exec()
		if err != nil {
			log.Error("%s", err)
			return
		}
		if t, ok := f.(*EndFrame); ok {
			return t.result, nil
		}
	}
	return
}

func RunCode(code SchemeObject) (result SchemeObject, err error) {
	progn, ok := code.(*Cons)
	if !ok {
		return nil, ErrType
	}

	env := &Environ{Parent: nil, Names: DefaultNames}
	var f Frame = &EndFrame{Env: env}

	env = &Environ{Parent: env, Names: make(map[string]SchemeObject)}
	f = CreateBeginFrame(progn, env, f)

	result, err = Trampoline(f)
	if result == nil {
		return nil, ErrUnknown
	}
	return
}
