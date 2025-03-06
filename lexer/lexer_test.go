package lexer

import (
	"testing"
	"yap/token"
)

func TestNextToken(t *testing.T) {
	input := `propose five = 5;
	propose ten = 10;

    ++a;
    a++;
    --a;
    a--;

	propose add = func(x,y){
	x + y;
	};

	propose result = add(five, ten);
	!/*5;
	5 < 10 > 5; 

	perhaps (a) {
		sayless true;
	} perchance (b) {
		sayless false; 
	} otherwise {
		sayless ten; 
	}
	5 == 5;
	5 != 6;

    [1, 2];
    for();
	`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "propose"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "propose"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INCREMENT, "++"},
		{token.IDENT, "a"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "a"},
		{token.INCREMENT, "++"},
		{token.SEMICOLON, ";"},
		{token.DECREMENT, "--"},
		{token.IDENT, "a"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "a"},
		{token.DECREMENT, "--"},
		{token.SEMICOLON, ";"},
		{token.LET, "propose"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "propose"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.DASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "perhaps"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "sayless"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELIF, "perchance"},
		{token.LPAREN, "("},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "sayless"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "otherwise"},
		{token.LBRACE, "{"},
		{token.RETURN, "sayless"},
		{token.IDENT, "ten"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "5"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.NEQ, "!="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.FOR, "for"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d], expected=%q, got=%q", i, tt.expectedType, tok.Type)
		} else if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d], expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestString(t *testing.T) {
	input := `"This is a test"`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STRING, "This is a test"},
	}

	l := New(input)
	for _, test := range tests {
		tok := l.NextToken()
		if tok.Type != test.expectedType {
			t.Fatalf("ExpectedType error: expect=%T, got=%T", test.expectedType, tok.Type)
		}
		if tok.Literal != test.expectedLiteral {
			t.Fatalf("ExpectedLiteral error: expect=%s, got=%s", test.expectedLiteral, tok.Literal)
		}
	}
}
