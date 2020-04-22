package lexer

import (
	"fmt"
	"io"

	"jakobvarmose/compute/token"
)

type state func(*lexer) state

type lexer struct {
	reader io.Reader
	ch     chan<- token.Token
	state  state
	buf    []byte
	index  int

	lineNum int
	charNum int
}

func Run(r io.Reader) <-chan token.Token {
	ch := make(chan token.Token)
	l := &lexer{
		reader: r,
		ch:     ch,
		state:  outside,

		lineNum: 1,
		charNum: 1,
	}
	go func() {
		for l.state != nil {
			l.state = l.state(l)
		}
	}()
	return ch
}

func (l *lexer) next() {
	if l.index == len(l.buf) {
		b := make([]byte, 1)
		l.reader.Read(b)
		l.buf = append(l.buf, b[0])
	}
	if l.buf[l.index] == '\n' {
		l.lineNum++
		l.charNum = 1
	} else {
		l.charNum++
	}
	l.index++
}

func (l *lexer) peek() uint8 {
	if l.index == len(l.buf) {
		b := make([]byte, 1)
		l.reader.Read(b)
		l.buf = append(l.buf, b[0])
	}
	return l.buf[l.index]
}

func (l *lexer) str() string {
	return string(l.buf[:l.index])
}

func (l *lexer) emit(t token.Type) {
	l.ch <- token.Token{t, l.str()}
	l.buf = l.buf[l.index:]
	l.index = 0
}

func (l *lexer) error() {
	panic(fmt.Sprintf("Error at line %d column %d", l.lineNum, l.charNum))
}

func handleString(l *lexer) state {
	l.next()
	c := l.peek()
	for c != '"' {
		if c == '\\' {
			l.next()
			c = l.peek()
		}
		l.next()
		c = l.peek()
	}
	l.next()
	l.emit(token.TypeString)
	return outside
}

func number(l *lexer) state {
	c := l.peek()
	for c >= '0' && c <= '9' {
		l.next()
		c = l.peek()
	}
	if l.peek() == 'i' {
		l.next()
		l.emit(token.TypeComplex)
		return outside
	}
	l.emit(token.TypeInt)
	return outside
}

func variable(l *lexer) state {
	c := l.peek()
	for c == '@' || c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
		l.next()
		c = l.peek()
	}
	str := l.str()
	if str == "true" || str == "false" {
		l.emit(token.TypeBool)
	} else if str == "or" || str == "and" || str == "xor" {
		l.emit(token.TypeOperator)
	} else {
		l.emit(token.TypeVariable)
	}
	return outside
}

func outside(l *lexer) state {
	c := l.peek()
	if c >= '0' && c <= '9' {
		return number
	} else if c == '@' || c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
		return variable
	} else if c == '"' {
		return handleString
	} else if c == ' ' || c == '\t' {
		l.next()
		l.emit(token.TypeSpace)
	} else if c == '\n' {
		l.next()
		l.emit(token.TypeNewline)
	} else if c == '+' || c == '-' || c == '/' || c == '^' || c == '%' || c == '#' || c == '$' {
		l.next()
		l.emit(token.TypeOperator)
	} else if c == '*' {
		l.next()
		if l.peek() == '*' {
			l.next()
		}
		l.emit(token.TypeOperator)
	} else if c == '<' {
		l.next()
		if l.peek() == '=' || l.peek() == '<' {
			l.next()
		}
		l.emit(token.TypeOperator)
	} else if c == ':' || c == '!' || c == '>' {
		l.next()
		if l.peek() == '=' {
			l.next()
		}
		l.emit(token.TypeOperator)
	} else if c == '=' {
		l.next()
		if l.peek() == '=' || l.peek() == '>' {
			l.next()
		}
		l.emit(token.TypeOperator)
	} else if c == '(' {
		l.next()
		l.emit(token.TypeStart)
	} else if c == ')' {
		l.next()
		l.emit(token.TypeEnd)
	} else if c == '[' {
		l.next()
		l.emit(token.TypeLeftSquare)
	} else if c == ']' {
		l.next()
		l.emit(token.TypeRightSquare)
	} else {
		l.error()
	}
	return outside
}
