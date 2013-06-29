package scmgo

import (
	// "fmt"
	"bufio"
	"io"
	"strings"
)

const control_chars = "()' \n\r\t"

func translate_string(dst []rune) ([]rune) {
	s := string(dst)
	s = strings.Replace(s, "\\t", "\t", -1)
	s = strings.Replace(s, "\\r", "\r", -1)
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\\"", "\"", -1)
	s = strings.Replace(s, "\\\\", "\\", -1)
	return []rune(s)
}

// func grammar_parser(source []rune, cout chan []rune) {
func grammar_parser(source io.Reader, cout chan string) {
	var err error
	src := bufio.NewReader(source)
	// defer source.Close()
	defer close(cout)

	// performance
	var symbolbuf []rune
	for c, _, err := src.ReadRune(); err == nil; c, _, err = src.ReadRune() {
		switch c {
		case '(', ')', '\'': // control chars
			if symbolbuf != nil {
				cout <- string(symbolbuf)
				symbolbuf = nil
			}
			cout <- string(c)
		case ' ', '\n', '\r', '\t': // empty chars
			if symbolbuf != nil {
				cout <- string(symbolbuf)
				symbolbuf = nil
			}
		case '"': // string
			if symbolbuf != nil {
				panic("quote in symbol")
			}
			var l string
			line, err := src.ReadString('"')
			for ; line[len(line)-2] == '\\'; {
				l, err = src.ReadString('"')
				switch err {
				case nil:
				case io.EOF: panic("quote not closed")
				default: panic(err)
				}
				line += l
			}
			cout <- line
		case ';': // comment
			if symbolbuf != nil {
				panic("comment in symbol")
			}
			_, err := src.ReadString('\n')
			if err != nil { break }
		default: // symbol
			symbolbuf = append(symbolbuf, c)
		}
	}
	if err == io.EOF { err = nil }
	if err != nil { panic(err) }
}

func code_parser(cin chan string, term bool) (code SchemeObject, err error) {
	var obj SchemeObject
	var objs []SchemeObject
	for chunk, ok := <- cin; ok; chunk, ok = <- cin {
		switch chunk[0]{
		case '(': // Cons
			obj, err = code_parser(cin, true)
			if err != nil { return nil, err }
		case ')': // return Cons
			code = Onil
			for i := len(objs) - 1; i >= 0; i-- {
				code = &Cons{car: objs[i], cdr: code}
			}
			return
		case '#': // Boolean
			if chunk[1] == 't' {
				obj = Otrue
			} else {
				obj = Ofalse
			}
		case '"': // String
			if chunk[len(chunk)-1] != '"' { panic("quote not closed") }
			obj = &String{str: chunk[1:len(chunk)-1]}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// Integer & Float
			if strings.Index(chunk, ".") != -1 {
				obj, err = Float_from_string(chunk)
				if err != nil { panic(err) }
			} else {
				obj, err = Integer_from_string(chunk)
				if err != nil { panic(err) }
			}
		case '\'': // Quote
			obj = new(Quote)
		default: // Symbol
			obj = &Symbol{name: chunk}
		}
		objs = append(objs, obj)
	}
	if term { panic("parenthesis not close") }
	return code, nil
}

// func BuildCode(source []rune) (code SchemeObject, err error) {
func BuildCode(source io.Reader) (code SchemeObject, err error) {
	cpipe := make(chan string)
	go grammar_parser(source, cpipe)
	// for chunk, ok := <- cpipe; ok; chunk, ok = <- cpipe {
	// 	fmt.Println("chunk:", string(chunk))
	// }
	// return nil, nil
	return code_parser(cpipe, false)
}