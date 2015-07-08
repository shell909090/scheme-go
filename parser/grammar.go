package parser

import (
	"bufio"
	"io"
)

func readString(src *bufio.Reader) (buf []byte, err error) {
	var c byte
	buf = append(buf, '"')

	for c, err = src.ReadByte(); err == nil; c, err = src.ReadByte() {
		switch c {
		case '"':
			return
		case '\\':
			c, err := src.ReadByte()
			switch err {
			case nil:
			case io.EOF:
				return nil, ErrQuotaNotClose
			default:
				return nil, err
			}
			switch c {
			case '\\':
				buf = append(buf, c)
			case 't':
				buf = append(buf, '\t')
			case 'r':
				buf = append(buf, '\r')
			case 'n':
				buf = append(buf, '\n')
			case '"':
				buf = append(buf, '"')
				// case 'x'?
			}
		default:
			buf = append(buf, c)
		}
	}
	return
}

func Grammar(source io.ReadCloser, dst io.Writer) (err error) {
	src := bufio.NewReader(source)
	defer source.Close()

	// performance
	var c byte
	var buf []byte
	for c, err = src.ReadByte(); err == nil; c, err = src.ReadByte() {
		switch c {
		case '(', ')', '\'': // control chars
			if buf != nil {
				_, err = dst.Write(buf)
				if err != nil {
					return
				}
				buf = nil
			}
			_, err = dst.Write([]byte{c})
		case ' ', '\n', '\r', '\t': // empty chars
			if buf != nil {
				_, err = dst.Write(buf)
				if err != nil {
					return
				}
				buf = nil
			}
		case '"': // string
			if buf != nil {
				return ErrQuotaInSymbol
			}
			buf, err = readString(src)
			if err != nil {
				return
			}
			_, err = dst.Write(buf)
			buf = nil
		case ';': // comment
			if buf != nil {
				return ErrCommentInSymbol
			}
			_, err = src.ReadSlice('\n')
		default: // symbol
			buf = append(buf, c)
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
