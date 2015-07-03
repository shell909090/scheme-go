package scmgo

func StepDebuger(f Frame, r SchemeObject) {
	if r != nil {
		log.Debug("result: %s", SchemeObjectToString(r))
	}
	log.Debug("stack:\n%s", StackFormatter(f))
	return
}

func Trampoline(init_frame Frame, init_obj SchemeObject) (result SchemeObject, err error) {
	f := init_frame
	result = init_obj
	for f != nil {
		// StepDebuger(f, result)
		result, f, err = f.Eval(result)
		if err != nil {
			log.Error("%s", err)
			return
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
	init_frame := CreatePrognFrame(nil, progn, env)

	result, err = Trampoline(init_frame, nil)
	if result == nil {
		return nil, ErrUnknown
	}
	return
}
