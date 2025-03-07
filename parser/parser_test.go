package parser

import (
	"fmt"
	"testing"
	"yap/ast"
	"yap/lexer"
)

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestSayStatment(t *testing.T) {
	tests := []struct {
		input               string
		expectedIndentifier string
		expectedValue       interface{}
	}{
		{"propose a = 4;", "a", 4},
		{"propose a = 5.5;", "a", 5.5},
		{"propose x = 5;", "x", 5},
		{"propose y = true;", "y", true},
		{"propose foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Statement length error: expect=1, got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, test.expectedIndentifier) {
			return
		}

		exp := stmt.(*ast.SayStatement).Value
		if !testLiteralExpression(t, exp, test.expectedValue) {
			return
		}
	}
}

func TestPrefixExpression(t *testing.T) {
	prefixTest := []struct {
		input  string
		prefix string
		value  interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, test := range prefixTest {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserError(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements error, expected=1, got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement error, expected=*ast.ExpressionStatement, got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expression error, expected *ast.PrefixExpression, got=%T", exp)
		}
		if exp.Operator != test.prefix {
			t.Fatalf("operator mismatch, expected=%q, got=%q", test.prefix, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, test.value) {
			return
		}
	}
}

func TestParsingPostfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     interface{}
		operator string
	}{
		{"5++;", 5, "++"},
		{"5--;", 5, "--"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Statement length error: expect=1, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expression error: expect= *ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PostfixExpression)
		if !ok {
			t.Fatalf("Expression Error: expect= *ast.PostfixExpression, got=%T", stmt.Expression)
		}

		if !testLiteralExpression(t, exp.Left, test.left) {
			return
		}

		if test.operator != exp.Operator {
			t.Fatalf("Mistaching operator: expect=%s, got=%s", test.operator, exp.Operator)
		}

	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true == false", true, "==", false},
		{"false == false", false, "==", false},
		{"false != true", false, "!=", true},
	}

	for _, test := range infixTests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("statements len error: expected 1, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("statement error, statement is not *ast.ExpressionStatement")
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expression is not a type of *ast.InfixExpression")
		}

		if !testInfixExpression(t, stmt.Expression, test.left, test.operator, test.right) {
			return
		}

		if test.operator != exp.Operator {
			t.Fatalf("Operator error: expected= %s, got=%s", test.operator, exp.Operator)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"(a + b) + c",
			"((a + b) + c)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 > 4 != 3 > 4",
			"((5 > 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) +5",
			"((1 + (2 + 3)) + 5)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserError(t, p)

		actual := program.String()
		if actual != test.expected {
			t.Fatalf("Parsing error: expect=%s, got=%s", test.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue bool
	}{
		{"cap", false},
		{"nocap", true},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Unexpected amount of program statements: expect=1, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("unexpected statement error: expecte ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("the expression type is not a boolean, got=*%T", stmt.Expression)
		}

		if boolean.Value != test.expectedValue {
			t.Fatalf("unexpected value error: expected=%t, got=%t", test.expectedValue, boolean.Value)
		}
	}
}
func TestIfExpression(t *testing.T) {
	input := "perhaps (a == b) { a }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Statement length error: expect=1, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement error, expecting ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expression Error: expect ast.IfExpression, got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "a", "==", "b") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("Consequence statement error: expect='3', got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence statement error: expect= ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifer(t, consequence.Expression, "a") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statments was not nil, got=%+v", exp.Alternative)
	}

}

func TestIfElseExpression(t *testing.T) {
	input := "perhaps (a == b) { a } otherwise { b }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Statement length error: expect=1, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expression Statement error: expect *ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expression error: expected ast.IfExpression, got=%T", stmt.Expression)
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("Length Consequence error: expect=1, got=%d", len(exp.Consequence.Statements))
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Fatalf("Length Alternative error: expec=1, got=%d", len(exp.Alternative.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expression error: could not convert to ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	if testIdentifer(t, consequence.Expression, "a") {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expression error: could not convert to ast.ExpressionStatement, got=%T", exp.Alternative.Statements[0])
	}

	if testIdentifer(t, alternative.Expression, "b") {
		return
	}
}

func TestElifStatement(t *testing.T) {
	input := `
	perhaps (a == b)
		{ a }
	perchance (a > b)
		{ b }
	perchance (a > b)
		{ a }
	otherwise
		{ c }
	`
	l := lexer.New(input)
	p := New(l)
	programp := p.ParserProgram()
	checkParserError(t, p)
	if len(programp.Statements) != 1 {
		t.Fatalf("Statement Legnth error: expect=1, got=%d, %s", len(programp.Statements), programp.String())
	}

	stmt, ok := programp.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expression error: expect ast.ExpressionStatement, got=%T", programp.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expression error: expect ast.IfExpression, got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "a", "==", "b") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("Consequence length error: expect 1, got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expression Error: expect ast.ExpressionStatment, got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifer(t, consequence.Expression, "a") {
		return
	}

	if len(exp.Elif) != 2 {
		t.Fatalf("Elif error: expect 2, got=%d", len(exp.Elif))
	}

	for i, elif := range exp.Elif {
		if !testInfixExpression(t, elif.Conditions, "a", ">", "b") {
			return
		}
		consequence, ok := elif.Consequences.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expression error: expect ast.ExpressionStatement, got=%T", elif.Consequences.Statements)
		}
		if i == 0 {
			if !testIdentifer(t, consequence.Expression, "b") {
				return
			}
		} else {
			if !testIdentifer(t, consequence.Expression, "a") {
				return
			}
		}
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Fatalf("Expression Alternative length error: expect 1, got=%d", len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative ExpressionStatement error: expect ast.ExpressionStatement, got=%T", exp.Alternative.Statements[0])
	}

	if testIdentifer(t, alternative.Expression, "c") {
		return
	}

}

func TestTernaryStatement(t *testing.T) {
	input := `(a == b)? a : b`

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Statement length error: expect=1, got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expression error: expect ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.TernaryExpression)
	if !ok {
		t.Fatalf("Expression error: expect= ast.TernaryExpression, got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "a", "==", "b") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("Cosequence Statement length error: expect=1, got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence Expression error: expect= ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifer(t, consequence.Expression, "a") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Fatalf("Alternative Expression length error: expect=1, got=%d", len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative Expression error: expect= ast.ExpressionStatement, got=%T", exp.Alternative.Statements[0])
	}
	if !testIdentifer(t, alternative.Expression, "b") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `
	func (a , b){
		a * b;
		c + e;
	}
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Statement Length error: expect=1, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expression error: expect= ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("Expression error: expect= ast.FunctionalExpression, got=%T", stmt.Expression)
	}

	if len(exp.Parameters) != 2 {
		t.Fatalf("Parameter length error: expect=2, got=%d", len(exp.Parameters))
	}

	if !testIdentifer(t, exp.Parameters[0], "a") {
		return
	}
	if !testIdentifer(t, exp.Parameters[1], "b") {
		return
	}

	if len(exp.Body.Statements) != 2 {
		t.Fatalf("Body Statement length error: expect=1, got=%d", len(exp.Body.Statements))
	}

	body, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expression error: expect= ast.ExpressionStatement, got=%T", exp.Body.Statements[0])
	}

	if !testInfixExpression(t, body.Expression, "a", "*", "b") {
		return
	}

}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{
			input:          "func() {}",
			expectedParams: []string{},
		},
		{
			input:          "func (x) { return x }",
			expectedParams: []string{"x"},
		},
		{
			input: `func (x , y, z, v) {
				x + y;
				return c
			}`,
			expectedParams: []string{"x", "y", "z", "v"},
		},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserError(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp := stmt.Expression.(*ast.FunctionExpression)

		if len(exp.Parameters) != len(test.expectedParams) {
			t.Fatalf("Length mismatching error: expect=%d, got=%d", len(test.expectedParams), len(exp.Parameters))
		}

		for i, arg := range test.expectedParams {
			if !testIdentifer(t, exp.Parameters[i], arg) {
				t.Fatalf("Mismatching parameter: expect=%s, got=%s", arg, exp.Parameters[i])
			}
		}
	}
}

func TestForExpression(t *testing.T) {
	tests := []struct {
		input         string
		condition_len int
		BSTlen        int
	}{
		{
			`
            for (propose a = 3; a < 5; ++a){
                b = b + a;
            }
            `,
			2,
			1,
		},
		{
			`
            for (a < 5){
                ++a;
            }
            `,
			1,
			1,
		},
		{
			`
            for (true){
                propose a = [1,2,3,4];
                b = b + a[3];
                a = append(a, b);
            }
            `,
			1,
			3,
		},
		{
			`
            for (a < 5; ++a){}
            `,
			2,
			0,
		},
	}

	for _, test := range tests {

		l := lexer.New(test.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Statement length error, expect=1, got=%d", len(program.Statements))
		}

		exp, ok := program.Statements[0].(*ast.ForExpression)
		if !ok {
			t.Fatalf("Expression error: expect= ast.Expression, got=%T", program.Statements[0])
		}

		if len(exp.Conditions) != test.condition_len {
			t.Fatalf("Condition length error: expect=3, got= %d", len(exp.Conditions))
		}

		if len(exp.Statements.Statements) != test.BSTlen {
			t.Fatalf("Block Expression lenght error: expect=1, got=%d", len(exp.Statements.Statements))
		}

	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 + 3, 4 * 5)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Statement length error: expect=1. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expression Statement error: expect= ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Expression error: expect= ast.CallExpression, got=%T", stmt.Expression)
	}
	if !testIdentifer(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("Argument length error: expect=3, got=%d", len(exp.Arguments))
	}
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "+", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "*", 5)
}

func TestArrayLiteral(t *testing.T) {
	inputs := [...]string{
		"propose a = [1, 2, 3, 5];",
		"propose a = [false, true, true];",
		"propose a = [\"true\", \"hello\", \"there\"];",
		"propose a = [[], [], []];",
	}

	for _, input := range inputs {
		t.Log(input)
		l := lexer.New(input)
		p := New(l)
		program := p.ParserProgram()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Unexpected length: expect=1, got=%d", len(program.Statements))
		}
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "arr[1+2];"

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserError(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expression error: expect=ast.Expression, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("Expression error: expect=ast.IndexExpression, got=%T", stmt.Expression)
	}

	if !testIdentifer(t, exp.Left, "arr") {
		return
	}

	if !testInfixExpression(t, exp.Index, 1, "+", 2) {
		return
	}
}

func TestHashLiteral(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserError(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("Expression error: expect=as.HashLiteral, got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("Hash length error: expect=3, got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, val := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not a StringLiteral, got=%T", key)
		}

		expectedVal := expected[literal.String()]
		testIntegerLiteral(t, val, expectedVal)
	}
}

func TestParsingEmptuHash(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserError(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("Expresison error: expect=HashLiteral, got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Fatalf("Hash length error: expect=0, got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralExpression(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10-8, "three": 15/5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserError(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("Expression error: expect=HashLiteral, got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("Hash length error: expect=3, got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, val := range hash.Pairs {
		lit, ok := key.(*ast.StringLiteral)

		if !ok {
			t.Errorf("Key is not ast.StringLiteral, gpt=%T", key)
			continue
		}

		testFunc, ok := tests[lit.String()]
		if !ok {
			t.Errorf("No test function for key %q found", lit.String())
		}
		testFunc(val)
	}
}

func checkParserError(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parse has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, intLit ast.Expression, value int64) bool {
	intex, ok := intLit.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("intLit is not ast.IntegerLiteral, got=%T", intex)
		return false
	}
	if intex.Value != value {
		t.Errorf("Mismatching value error, expected=%d, got=%d", value, intex.Value)
		return false
	}

	if intex.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("intex.TokenLiteral not %d, got=%s", value, intex.TokenLiteral())
		return false
	}
	return true
}

func testIdentifer(t *testing.T, exp ast.Expression, value string) bool {
	indent, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Could not casted expression as ast.Identifier")
		return false
	}
	if indent.Value != value {
		t.Errorf("Mismatching value error: expect=%s, got=%s", value, indent.Value)
		return false
	}
	if indent.TokenLiteral() != value {
		t.Errorf("Mismatching token literal error: expect=%s, got=%s", value, indent.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifer(t, exp, v)
	case bool:
		return testBoolLiteral(t, exp, v)
	case float64:
		return testFloatLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operation string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got=%T(%s)", exp, exp)
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operation {
		t.Errorf("exp.Operator is not '%s', got=%q", operation, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testBoolLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("expression is not ast.Boolean")
		return false
	}

	if bo.Value != value {
		t.Errorf("unexpected value error: expected=%t, got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t, got=%s", value, bo.TokenLiteral())
		return false
	}
	return true
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "propose" {
		t.Errorf("stmt TokenLiteral error: expect='say', got=%s", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.SayStatement)
	if !ok {
		t.Errorf("Statement error: expect ast.SayStatment, got=%T", stmt)
	}

	if letStmt.Name.Value != name {
		t.Errorf("Statement value error: expect=%s, got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("Statement Token Literal error: expect=%s, got=%s", name, letStmt.TokenLiteral())
		return false
	}

	return true

}

func testFloatLiteral(t *testing.T, exp ast.Expression, val float64) bool {
	flo, ok := exp.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("Expression Error: expect= ast.FloatLiteral, got=%T", exp)
		return false
	}

	if flo.Value != val {
		t.Errorf("Value error: expect=%f, got=%f", val, flo.Value)
		return false
	}

	if flo.TokenLiteral() != fmt.Sprintf("%g", val) {
		t.Errorf("Token Literal error: expect=%g, got=%s", val, flo.TokenLiteral())
		return false
	}

	return true
}
