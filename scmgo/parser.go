package scmgo

import "strings"

const control_chars = "()' \n\r\t"

func CodeNumber(chunk string) (obj SchemeObject, err error) {
	if strings.Index(chunk, ".") != -1 {
		obj, err = FloatFromString(chunk)
	} else {
		obj, err = IntegerFromString(chunk)
	}
	return
}

func CodeParser(cin chan string) (code SchemeObject, err error) {
	var obj SchemeObject
	var objs []SchemeObject

QUIT:
	for chunk, ok := <-cin; ok; chunk, ok = <-cin {
		switch chunk[0] {
		case '#': // Boolean
			obj, err = BooleanFromString(chunk)
			if err != nil {
				return
			}
		case '"': // String
			obj = String(chunk[1 : len(chunk)-1])
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if chunk[0] == '-' && len(chunk) == 1 {
				// - without number is symbol
				obj = SymbolFromString(chunk)
			} else { // Integer & Float
				obj, err = CodeNumber(chunk)
				if err != nil {
					return
				}
			}
		case '\'': // Quote
			obj = new(Quote)
		case ';': // Comment
			obj = nil
		case '(': // Cons
			obj, err = CodeParser(cin)
			if err != nil {
				return nil, err
			}
		case ')': // return Cons
			break QUIT
		default: // Symbol
			obj = SymbolFromString(chunk)
		}

		if obj == nil { // pass comment
			continue
		}

		// processing Quote
		if len(objs) > 0 {
			o := objs[len(objs)-1]
			if last, ok := o.(*Quote); ok {
				last.objs = obj
				continue
			}
		}

		objs = append(objs, obj)
	}

	if len(objs) > 0 {
		o := objs[len(objs)-1]
		if last, ok := o.(*Quote); ok && last.objs == nil {
			return code, ErrQuoteInEnd
		}
	}

	return ListFromSlice(objs), nil
}
