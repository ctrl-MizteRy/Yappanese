package evaluator

import (
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
		return evalStatements(node.Statements)

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
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return intOrfloat(left, node.Operator, right)
	case *ast.PostfixExpression:
		left := Eval(node.Left)
		return evalPrefixExpression(node.Operator, left)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)
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
		return NULL
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
		return &object.Float{Value: float64(l_val) / float64(r_val)}
	case "**":
		val := math.Pow(float64(l_val), float64(r_val))
		return &object.Float{Value: val}
	default:
		return NULL
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
	default:
		return NULL
	}

}

func intOrfloat(left object.Object, operator string, right object.Object) object.Object {
	l, err := left.(*object.Integer)
	if !err {
		r, err := right.(*object.Integer)
		if !err {
			return evalInfixFloatExpression(left, operator, right)
		}
		r_float := float64(r.Value)
		righ := &object.Float{Value: r_float}
		return evalInfixFloatExpression(left, operator, righ)
	}
	r, err := right.(*object.Integer)
	if !err {
		lef := &object.Float{Value: float64(l.Value)}
		return evalInfixFloatExpression(lef, operator, right)
	}
	return evalInfixIntExpression(left, operator, r)
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
