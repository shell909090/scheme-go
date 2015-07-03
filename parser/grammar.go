package parser

import (
	"bufio"
	"io"
)

func GrammarReadString(src *bufio.Reader) (s string, err error) {
	var strbuf []rune
	strbuf = append(strbuf, '"')
	for c, _, err := src.ReadRune(); err == nil; c, _, err = src.ReadRune() {
		switch c {
		case '"':
			return string(strbuf), nil
		case '\\':
			c, _, err := src.ReadRune()
			switch err {
			case nil:
			case io.EOF:
				return s, ErrQuotaNotClose
			default:
				return s, err
			}
			switch c {
			case '\\':
				strbuf = append(strbuf, c)
			case 't':
				strbuf = append(strbuf, '\t')
			case 'r':
				strbuf = append(strbuf, '\r')
			case 'n':
				strbuf = append(strbuf, '\n')
			case '"':
				strbuf = append(strbuf, '"')
				// case 'x'?
			}
		default:
			strbuf = append(strbuf, c)
		}
	}
	return
}

func GrammarParser(source io.ReadCloser, cout chan string) (err error) {
	src := bufio.NewReader(source)
	defer source.Close()
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
				return ErrQuotaInSymbol
			}
			line, err := GrammarReadString(src)
			if err != nil {
				return err
			}
			cout <- line
		case ';': // comment
			if symbolbuf != nil {
				return ErrCommentInSymbol
			}
			s, err := src.ReadString('\n')
			if err != nil {
				return err
			}
			cout <- ";" + s
		default: // symbol
			symbolbuf = append(symbolbuf, c)
		}
	}

	if err == io.EOF {
		err = nil
	}
	return err
}
