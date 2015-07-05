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
	list, ok := code.(*Cons)
	if !ok {
		return nil, ErrType
	}

	env := &Environ{Parent: DefaultEnv, Names: make(map[string]SchemeObject)}
	f := CreateBeginFrame(list, env, &EndFrame{Env: DefaultEnv})

	result, err = Trampoline(f)
	if result == nil {
		return nil, ErrUnknown
	}
	return
}
