package ast

import (
	"testing"
	"yap/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&SayStatement{
				Token: token.Token{Type: token.LET, Literal: "say"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "say_var"},
					Value: "say_var",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
			&ConstStaement{
				Token: token.Token{Type: token.CONST, Literal: "ackchyually"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "constVar"},
					Value: "constVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "constVar2"},
					Value: "constVar2",
				},
			},
			&GlobalStatement{
				Token: token.Token{Type: token.GLOBAL, Literal: "worldwide"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "globalVar"},
					Value: "globalVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "globalVar2"},
					Value: "globalVar2",
				},
			},
		},
	}

	if program.Statements[0].String() != "say say_var = anotherVar;" {
		t.Errorf("program.String() wrong, got= %q", program.String())
	}
}
