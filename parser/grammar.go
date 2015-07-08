package parser

import (
	"bufio"
	"io"
)

func readString(src *bufio.Reader) (s string, err error) {
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

func Grammar(source io.ReadCloser, p *Parser) (err error) {
	src := bufio.NewReader(source)
	defer source.Close()

	// performance
	var symbolbuf []byte
	var str string
	var c byte
	for c, err = src.ReadByte(); err == nil; c, err = src.ReadByte() {
		switch c {
		case '(', ')', '\'': // control chars
			if symbolbuf != nil {
				_, err = p.Write(symbolbuf)
				if err != nil {
					return
				}
				symbolbuf = nil
			}
			_, err = p.Write([]byte{c})
		case ' ', '\n', '\r', '\t': // empty chars
			if symbolbuf != nil {
				_, err = p.Write(symbolbuf)
				if err != nil {
					return
				}
				symbolbuf = nil
			}
		case '"': // string
			if symbolbuf != nil {
				return ErrQuotaInSymbol
			}
			str, err = readString(src)
			if err != nil {
				return
			}
			err = p.ReceiveChunk(str)
		case ';': // comment
			if symbolbuf != nil {
				return ErrCommentInSymbol
			}
			str, err = src.ReadString('\n')
			if err != nil {
				return
			}
			err = p.ReceiveChunk(";" + str)
		default: // symbol
			symbolbuf = append(symbolbuf, c)
		}
		if err != nil {
			log.Error("%s", err)
			return
		}
	}

	if err == io.EOF {
		err = nil
	}
	return err
}
