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
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`{"name": "Monkey"}[func(x){x}];`,
			"unusable as hash key: FUNCTION",
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

func TestSayStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"propose a = 5; a;", 5},
		{"propose a = 5 * 5; a;", 25},
		{"propose a = 5; propose b = a; b;", 5},
		{"propose a =5; propose b = a; propose c = a + b +5; c;", 15},
	}

	for _, test := range tests {
		testIntegerObject(t, testEval(test.input), test.expected)
	}
}

func TestFucntionObject(t *testing.T) {
	input := "func(x) {x + 2;};"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not a function object, got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has the wrong parameters, Paramenters= %+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("Error: body is not %q, got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"propose identity = func(x) {x;}; identity(5);", 5},
		{"propose identity = func(x) { sayless x;}; identity(5);", 5},
		{"propose double = func(x) {x * 2;} double(5);", 10},
		{"propose add = func(x, y) { x + y;}; add(5, 5);", 10},
		{"propose add = func(x , y) {x + y;}; add(5, add(5, 5));", 15},
		{"func(x) {x;}(5)", 5},
	}

	for _, test := range tests {
		testIntegerObject(t, testEval(test.input), test.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"propose a = \"hello\"; a;", "hello"},
		{"propose a = \"2\"; a;", "2"},
		{"propose a = \" \"; a;", " "},
	}

	for _, test := range tests {
		testStringObject(t, testEval(test.input), test.expected)
	}
}

func TestStringOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`propose a = "hello" + " " + "there"; a`, "hello there"},
	}

	for _, test := range tests {
		testStringObject(t, testEval(test.input), test.expected)
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("hello")`, 5},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments, expect=1, got=2"},
	}
	for _, test := range tests {
		eval := testEval(test.input)

		switch expected := test.expected.(type) {
		case int:
			testIntegerObject(t, eval, int64(expected))
		case string:
			objErr, ok := eval.(*object.Error)
			if !ok {
				t.Errorf("object is not Error, got=%T", eval)
			}
			if objErr.Message != expected {
				t.Errorf("wrong error message. expect=%s, got=%s", expected, objErr.Message)
			}
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	input := "[1,2 * 2, 3 +3]"

	eval := testEval(input)
	result, ok := eval.(*object.Array)
	if !ok {
		t.Fatalf("Object error: expect= object.Array, got=%T", eval)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("Array length error: expected=3, got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestForLoopOpperation(t *testing.T) {
	input := `
        propose a = 4;
        for (propose b = 3; b <6; ++b){
            a = a + b;
        }
        a;
    `
	eval := testEval(input)

	result, ok := eval.(*object.Integer)
	if !ok {
		t.Fatalf("Unexpect object error: expect= object.For, got=%T", eval)
	}

	testIntegerObject(t, result, 16)

}

func TestForLoopReturn(t *testing.T) {
	tests := []struct {
		input     string
		expectVal any
	}{
		{
			`
            for (propose a = 5; a < 10; ++a){
                perhaps (a == 7){
                    sayless a;
                }
            }
            `,
			7,
		},
		{
			`
            propose a = 6;
            for (nocap){
                ++a;
                perhaps ( a == 10){
                    sayless "hello";
                }
            }
            `,
			"hello",
		},
	}

	for _, test := range tests {
		result := testEval(test.input)

		switch obj := result.(type) {
		case *object.Integer:
			num := int64(test.expectVal.(int))
			testIntegerObject(t, obj, num)
		case *object.String:
			str := test.expectVal.(string)
			testStringObject(t, obj, str)
		default:
			t.Fatalf("Unexpected object type of %T", obj)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParserProgram()
	env := object.NewEnviroment()

	return Eval(program, env)
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1,2,3,4][0]",
			1,
		},
		{
			"[1,2,3][1]",
			2,
		},
		{
			"[1,2,3][2]",
			3,
		},
		{
			"[1,2,3][1+1]",
			3,
		},
		{
			"propose arr = [1,2,3]; arr[1]",
			2,
		},
		{
			"propose arr = [1,2,3]; arr[0] + arr[1] +arr[2]",
			6,
		},
		{
			"propose arr = [1,2,3]; propose i = arr[1]; arr[i]",
			3,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1,2,3][-1]",
			nil,
		},
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

func TestHashLiteral(t *testing.T) {
	input := `
    propose two = "two";
    {
        "one": 10 - 9,
        "two": 1 + 1,
        "thr" + "ee" :6 / 2,
        4: 4,
        true: 5,
        false: 6
    }
    `

	eval := testEval(input)
	result, ok := eval.(*object.Hash)
	if !ok {
		t.Fatalf("Object Error: expect=object.Hash, got=%T", eval)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs, got=%d", len(result.Pairs))
	}

	for expectedKey, expectedVal := range expected {
		pair, ok := result.Pairs[expectedKey]

		if !ok {
			t.Errorf("No pair for give key in pairs")
		}
		testIntegerObject(t, pair.Value, expectedVal)
	}
}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5 :5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
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

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	str, ok := obj.(*object.String)

	if !ok {
		t.Errorf("Object error: expect=object.String, got =%T", obj)
		return false
	}

	if str.Value != expected {
		t.Errorf("String error: mismatch string, expect=%s, got=%s",
			expected, str.Value)
		return false
	}
	return true
}
