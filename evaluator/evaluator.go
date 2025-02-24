package evaluator

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"yap/ast"
	"yap/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Enviroment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.FunctionExpression:
		params := node.Parameters
		body := node.Body
		if node.Name != nil {
			fu := &object.Function{Parameters: params, Env: env, Body: body}
			env.Set(node.Name.Value, fu)
		} else {
			return &object.Function{Parameters: params, Env: env, Body: body}
		}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.SayStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.PotentialStatement:
		if env.Exist(node.Name.Value) {
			val := Eval(node.Value, env)
			if isError(val) {
				return val
			}
			if !env.TypeComp(node.Name.Value, val.Type()) {
				return newError("type mismatch error: could not set %s into '%s' variable (Type = %s)",
					val.Type(), node.Name.String(), env.GetType(node.Name.Value).Type())
			}
			env.Set(node.Name.Value, val)
		} else {
			return newError("valariable %s does not exist, (perhaps not yet declare?)", node.Name.String())
		}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		objType := elements[0].Type()
		for i := 1; i < len(elements); i++ {
			if elements[i].Type() != objType {
				return newError("Type mismatch, cannot have an array of %s and %s",
					objType, elements[i].Type())
			}
		}
		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.Right.TokenLiteral(), env)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(left, node.Operator, right, env)
	case *ast.PostfixExpression:
		left := Eval(node.Left, env)
		return evalPrefixExpression(node.Operator, left, node.Left.TokenLiteral(), env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.TernaryExpression:
		return evalTernaryExpression(node, env)
	case *ast.StringLiteral:
		return &object.String{Value: node.Literal}
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}
	return nil
}

func evalStatements(stmts []ast.Statement, env *object.Enviroment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)

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

func evalPrefixExpression(operator string, right object.Object, name string, env *object.Enviroment) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalNegativeOperatorExpression(right)
	case "++":
		return evalIncrementOperatorExpression(right, env, name)
	case "--":
		return evalDecrementOperatorExpression(right, env, name)
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

func evalInfixIntExpression(left object.Object, operator string,
	right object.Object, env *object.Enviroment) object.Object {
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

func evalInfixFloatExpression(left object.Object, operator string,
	right object.Object) object.Object {
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

func evalStringInfixExpression(left object.Object, operator string,
	right object.Object) object.Object {

	lVal := left.(*object.String).Value
	rVal := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: (lVal + rVal)}
	case "*":
		num, err := strconv.Atoi(lVal)
		if err == nil {
			str := rVal
			for i := 0; i < num; i++ {
				str += str
			}
			return &object.String{Value: str}
		}
		num, err = strconv.Atoi(rVal)
		if err == nil {
			str := lVal
			for i := 0; i < num; i++ {
				str += lVal
			}
			return &object.String{Value: str}
		}
		return newError("Cannot do a multiplication operator on %s and %s",
			lVal, rVal)
	default:
		return newError("Operator '%s' is not supported for string operation", operator)
	}
}

func evalIncrementOperatorExpression(obj object.Object, env *object.Enviroment, name string) object.Object {
	if obj.Type() == object.INTEGER_OBJ {
		val := obj.(*object.Integer).Value
		env.Set(name, &object.Integer{Value: val + 1})
		return &object.Integer{Value: val + 1}
	} else if obj.Type() == object.FLOAT_OBJ {
		val := obj.(*object.Float).Value
		env.Set(name, &object.Float{Value: val + 1.0})
		return &object.Float{Value: val + 1.0}
	}
	return NULL
}

func evalDecrementOperatorExpression(obj object.Object, env *object.Enviroment, name string) object.Object {
	log.Println(name)
	if obj.Type() == object.INTEGER_OBJ {
		val := obj.(*object.Integer).Value
		env.Set(name, &object.Integer{Value: val - 1})
		return &object.Integer{Value: val - 1}
	} else if obj.Type() == object.FLOAT_OBJ {
		val := obj.(*object.Float).Value
		env.Set(name, &object.Float{Value: val - 1.0})
		return &object.Float{Value: val - 1.0}
	}
	return NULL
}

func evalInfixExpression(left object.Object, operator string,
	right object.Object, env *object.Enviroment) object.Object {

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntExpression(left, operator, right, env)
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
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(left, operator, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.INTEGER_OBJ:
		val := int(right.(*object.Integer).Value)
		strObj := &object.String{Value: strconv.Itoa(val)}
		return evalStringInfixExpression(left, operator, strObj)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.STRING_OBJ:
		val := int(left.(*object.Integer).Value)
		strObj := &object.String{Value: strconv.Itoa(val)}
		return evalStringInfixExpression(strObj, operator, right)
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

func evalIfExpression(exp *ast.IfExpression, env *object.Enviroment) object.Object {
	conditions := Eval(exp.Condition, env)

	if isError(conditions) {
		return conditions
	}
	if isTrue(conditions) {
		return Eval(exp.Consequence, env)
	} else if exp.Elif != nil {
		condi_len := len(exp.Elif)
		for i := 0; i < condi_len; i++ {
			condis := Eval(exp.Elif[i].Conditions, env)
			if isTrue(condis) {
				return Eval(exp.Elif[i].Consequences, env)
			}
		}
	}
	if exp.Alternative != nil {
		return Eval(exp.Alternative, env)
	} else {
		return NULL
	}
}

func evalTernaryExpression(exp *ast.TernaryExpression, env *object.Enviroment) object.Object {
	condition := Eval(exp.Condition, env)
	if isTrue(condition) {
		return Eval(exp.Consequence, env)
	} else {
		return Eval(exp.Alternative, env)
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

func evalProgram(obj *ast.Program, env *object.Enviroment) object.Object {
	var result object.Object

	for _, stmt := range obj.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(obj *ast.BlockStatement, env *object.Enviroment) object.Object {
	var result object.Object

	for _, stmt := range obj.Statements {
		result = Eval(stmt, env)

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

func evalIdentifier(obj *ast.Identifier, env *object.Enviroment) object.Object {
	if val, ok := env.Get(obj.Value); ok {
		return val
	}

	if builtins, ok := builtins[obj.Value]; ok {
		return builtins
	}

	return newError("identifier not found: " + obj.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Enviroment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch function := fn.(type) {
	case *object.Function:
		extendedEvn := extendFunctionEnv(function, args)
		evaluated := Eval(function.Body, extendedEvn)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return function.Fn(args...)
	}

	return newError("not a function: %s", fn.Type())
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Enviroment {
	env := object.NewEncloseEnviroment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("Index operator not support %s", left.Type())
	}
}

func evalArrayIndexExpression(arr, index object.Object) object.Object {
	arrObj := arr.(*object.Array)
	idex := index.(*object.Integer).Value
	max := int64(len(arrObj.Elements) - 1)

	if idex < 0 || idex > max {
		return NULL
	}

	return arrObj.Elements[idex]
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Enviroment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	keys := []object.Object{}
	for keyNode, valNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		keys = append(keys, key)
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable hash key %s", key.Type())
		}

		val := Eval(valNode, env)
		if isError(val) {
			return val
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: val}
	}

	return &object.Hash{Pairs: pairs, Keys: keys}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}
