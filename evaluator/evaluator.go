package evaluator

import (
	"fmt"
	"log"
	"math"
	"yap/ast"
	"yap/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(left, node.Operator, right)
	case *ast.PostfixExpression:
		left := Eval(node.Left)
		return evalPrefixExpression(node.Operator, left)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.TernaryExpression:
		return evalTernaryExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)

		if returnVal, ok := result.(*object.ReturnValue); ok {
			return returnVal.Value
		}
	}
	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalNegativeOperatorExpression(right)
	case "++":
		return evalIncrementOperatorExpression(right)
	case "--":
		return evalDecrementOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right)
	}
}

func evalBangOperatorExpression(obj object.Object) object.Object {
	switch obj {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalNegativeOperatorExpression(obj object.Object) object.Object {
	if obj.Type() != object.INTEGER_OBJ && obj.Type() != object.FLOAT_OBJ {
		return newError("unknown operator: -%s", obj.Type())
	}
	if obj.Type() == object.INTEGER_OBJ {
		val := obj.(*object.Integer).Value
		return &object.Integer{Value: -val}
	} else if obj.Type() == object.FLOAT_OBJ {
		val := obj.(*object.Float).Value
		return &object.Float{Value: -val}
	}
	return NULL
}

func evalInfixIntExpression(left object.Object, operator string, right object.Object) object.Object {
	l_val := left.(*object.Integer).Value
	r_val := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: l_val + r_val}
	case "-":
		return &object.Integer{Value: l_val - r_val}
	case "*":
		return &object.Integer{Value: (l_val * r_val)}
	case "%":
		return &object.Integer{Value: (l_val % r_val)}
	case "/":
		if r_val == 0 {
			log.Fatal("ZERO DIVISION ERROR")
		}
		return &object.Integer{Value: l_val / r_val}
	case "**":
		val := math.Pow(float64(l_val), float64(r_val))
		return &object.Integer{Value: int64(val)}
	case "<":
		return nativeBoolToBooleanObject(l_val < r_val)
	case ">":
		return nativeBoolToBooleanObject(l_val > r_val)
	case "==":
		return nativeBoolToBooleanObject(l_val == r_val)
	case "!=":
		return nativeBoolToBooleanObject(l_val != r_val)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalInfixFloatExpression(left object.Object, operator string, right object.Object) object.Object {
	l_val := left.(*object.Float).Value
	r_val := right.(*object.Float).Value
	switch operator {
	case "+":
		return &object.Float{Value: l_val + r_val}
	case "-":
		return &object.Float{Value: l_val - r_val}
	case "*":
		return &object.Float{Value: (l_val * r_val)}
	case "%":
		return NULL
	case "/":
		if r_val == 0.0 {
			log.Fatal("ZERO DIVISION ERROR")
		}
		return &object.Float{Value: float64(l_val / r_val)}
	case "**":
		val := math.Pow(l_val, r_val)
		return &object.Float{Value: val}
	case "<":
		return nativeBoolToBooleanObject(l_val < r_val)
	case ">":
		return nativeBoolToBooleanObject(l_val > r_val)
	case "==":
		return nativeBoolToBooleanObject(l_val == r_val)
	case "!=":
		return nativeBoolToBooleanObject(l_val != r_val)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}

}

func evalIncrementOperatorExpression(obj object.Object) object.Object {
	if obj.Type() == object.INTEGER_OBJ {
		val := obj.(*object.Integer).Value
		return &object.Integer{Value: val + 1}
	} else if obj.Type() == object.FLOAT_OBJ {
		val := obj.(*object.Float).Value
		return &object.Float{Value: val + 1.0}
	}
	return NULL
}

func evalDecrementOperatorExpression(obj object.Object) object.Object {
	if obj.Type() == object.INTEGER_OBJ {
		val := obj.(*object.Integer).Value
		return &object.Integer{Value: val - 1}
	} else if obj.Type() == object.FLOAT_OBJ {
		val := obj.(*object.Float).Value
		return &object.Float{Value: val - 1.0}
	}
	return NULL
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntExpression(left, operator, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		val := left.(*object.Integer).Value
		left = &object.Float{Value: float64(val)}
		return evalInfixFloatExpression(left, operator, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		val := right.(*object.Integer).Value
		right = &object.Float{Value: float64(val)}
		return evalInfixFloatExpression(left, operator, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalInfixFloatExpression(left, operator, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case (left.Type() != right.Type() &&
		(left.Type() != object.INTEGER_OBJ || left.Type() != object.FLOAT_OBJ) &&
		(right.Type() != object.INTEGER_OBJ || right.Type() != object.FLOAT_OBJ)):
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(exp *ast.IfExpression) object.Object {
	conditions := Eval(exp.Condition)

	if isError(conditions) {
		return conditions
	}
	if isTrue(conditions) {
		return Eval(exp.Consequence)
	} else if exp.Elif != nil {
		condi_len := len(exp.Elif)
		for i := 0; i < condi_len; i++ {
			condis := Eval(exp.Elif[i].Conditions)
			if isTrue(condis) {
				return Eval(exp.Elif[i].Consequences)
			}
		}
	}
	if exp.Alternative != nil {
		return Eval(exp.Alternative)
	} else {
		return NULL
	}
}

func evalTernaryExpression(exp *ast.TernaryExpression) object.Object {
	condition := Eval(exp.Condition)
	if isTrue(condition) {
		return Eval(exp.Consequence)
	} else {
		return Eval(exp.Alternative)
	}
}

func isTrue(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalProgram(obj *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range obj.Statements {
		result = Eval(stmt)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(obj *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range obj.Statements {
		result = Eval(stmt)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
