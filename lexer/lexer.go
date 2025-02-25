package lexer

import (
	"strings"
	"yap/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() //using readChar to initialize the Lexer with pos = 0 and readpos = 1
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace()

	switch l.ch {
	case '=': //check for '=='
		if l.nextChar() == '=' {
			char := l.ch
			l.readChar()
			literal := string(char) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '-':
		if l.nextChar() == '-' {
			a := Lexer{input: l.input, position: l.position, readPosition: l.readPosition, ch: l.ch}
			a.readChar()
			if a.nextChar() != '-' {
				a.skipWhiteSpace()
				if isDigital(a.nextChar()) || isLetter(a.nextChar()) || a.nextChar() == ';' {
					char := l.ch
					l.readChar()
					lit := string(char) + string(l.ch)
					tok = token.Token{Type: token.DECREMENT, Literal: lit}
				} else {
					tok = newToken(token.ILLEGAL, l.ch)
				}
			}
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		if l.nextChar() == '+' {
			a := Lexer{input: l.input, position: l.position, readPosition: l.readPosition, ch: l.ch}
			a.readChar()
			if a.nextChar() != '+' {
				a.skipWhiteSpace()
				if isDigital(a.nextChar()) || isLetter(a.nextChar()) || a.nextChar() == ';' {
					char := l.ch
					l.readChar()
					lit := string(char) + string(l.ch)
					tok = token.Token{Type: token.INCREMENT, Literal: lit}
				} else {
					tok = newToken(token.ILLEGAL, l.ch)
				}
			}
		} else {
			tok = newToken(token.PLUS, l.ch)
		}
	case '>':
		tok = newToken(token.GT, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '!': //check for '!='
		if l.nextChar() == '=' {
			char := l.ch
			l.readChar()
			literal := string(char) + string(l.ch)
			tok = token.Token{Type: token.NEQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		if l.nextChar() == '*' {
			char := l.ch
			l.readChar()
			literal := string(char) + string(l.ch)
			tok = token.Token{Type: token.POWER, Literal: literal}
		} else {
			tok = newToken(token.ASTERISK, l.ch)
		}
	case '/':
		tok = newToken(token.DASH, l.ch)
	case '.':
		tok = newToken(token.FLOAT, l.ch)
	case '?':
		tok = newToken(token.TERNARY, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '%':
		tok = newToken(token.MOD, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigital(l.ch) {
			tok = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar() // move the char up every time this is run
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	postion := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[postion:l.position] //group up the whole word
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) nextChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readNumber() token.Token {
	var tok token.Token
	position := l.position
	flo := false
	for isDigital(l.ch) {
		l.readChar()
		if l.ch == '.' {
			tok.Type = token.FLOAT
			l.readChar()
			flo = !flo
			if l.ch == '.' {
				tok.Literal = l.input[position:l.position]
				return tok
			}
		}
	}
	if !flo {
		tok.Type = token.INT
	}
	tok.Literal = l.input[position:l.position]
	return tok
}

func isDigital(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) nextTwoChar() byte {
	if l.readPosition+1 >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition+1]
	}
}

func (l *Lexer) readString() string {
	l.readChar()
	str := []string{}
	str = append(str, string(l.ch))
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
		if l.ch == '\\' {
			if l.nextChar() == 't' {
				l.readChar()
				str = append(str, "\t")
			} else if l.nextChar() == 'n' {
				l.readChar()
				str = append(str, "\n")
			} else {
				str = append(str, string(l.ch))
			}
		} else {
			str = append(str, string(l.ch))
		}
	}
	str = str[:len(str)-1] //removing the last "
	return strings.Join(str, "")
}
