package evaluator

import (
	"testing"
	"yap/lexer"
	"yap/object"
	"yap/parser"
)

func TestEvalIntegralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"10", 10},
		{"10", 10},
		{"-5", -5},
		{"-40", -40},
		{"5 + 6 + 7 + 6 + 4", 28},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"2 * (5 + 10)", 30},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"50 / 2 * 2 + 10", 60},
		{"35 / 5 + 5", 12},
		{"2 ** 2", 4},
		{"3 ** 2", 9},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 != 2", true},
		{"1 == 2", false},
		{"true == true", true},
		{"true == false", false},
		{"true != false", true},
		{"false == false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(2 > 1) == true", true},
		{"(2 > 1) == false", false},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input     string
		expeceted bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expeceted)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"perhaps (true) { 10 }", 10},
		{"perhaps (false) { 10 }", nil},
		{"perhaps (1) { 10 }", 10},
		{"perhaps (1 < 2) {10}", 10},
		{"perhaps (1 > 2) {10}", nil},
		{"perhaps (1 > 2) {10} otherwise {20}", 20},
		{"perhaps (1 < 2) {10} otherwise {20}", 10},
		{"perhaps (1 > 2) {30} perchance ( 2 == 4) {3} perchance (2**2 == 4) {23}", 23},
		{"perhaps (2 > 3) {3} perchance (3 < 2) {3} otherwise {4}", 4},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestTernaryExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"(2 == 4)? 4 : 5", 5},
		{"(true)? 25 : 4", 25},
		{"(3 == 3) ? 5 : 1", 5},
		{"(false) ? 7 : 3", 3},
		{"(!false) ? 7 : 3", 7},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		integer, ok := test.expected.(int)
		if ok {
			testIntegerObject(t, eval, int64(integer))
		} else {
			testNullObject(t, eval)
		}
	}
}

func TestReturnStatments(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"sayless 10;", 10},
		{"sayless 10; 9;", 10},
		{"sayless 2 * 4;", 8},
		{"10; sayless 2 * 5 + 9; 8;", 19},
		{"perhaps (10 > 1) {perhaps (10 > 1) { sayless 10; } sayless 1;}", 10},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		testIntegerObject(t, eval, test.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input       string
		expectedMsg string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"perhaps (10 > 1) { true + false;}",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"perhaps (10 > 1) { perhaps (10 > 1) { sayless true + false;} sayless 1;}",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)

		errObj, ok := eval.(*object.Error)
		if !ok {
			t.Errorf("no error object return, got=%T (%+v)", eval, eval)
			continue
		}

		if errObj.Message != test.expectedMsg {
			t.Errorf("wrong error message, expected=%s, got=%s",
				test.expectedMsg, errObj.Message)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParserProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Error: object is not an Integer, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Value error: expect=%d, got=%d", expected, result.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Object error: expect= object.Boolean, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Value error: expect=%t, got=%t", expected, result.Value)
		return false
	}
	return true
}
func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL, got=%T", obj)
		return false
	}
	return true
}
