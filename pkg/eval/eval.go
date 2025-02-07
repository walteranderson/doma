package eval

import (
	"doma/pkg/parser"
	"fmt"
	"strings"
)

func Eval(expr parser.Expression, env *Env) Object {
	switch expr := expr.(type) {
	case *parser.Program:
		return evalProgram(expr, env)
	case *parser.Number:
		return &Number{Value: expr.Value}
	case *parser.String:
		return &String{Value: expr.Value}
	case *parser.Boolean:
		return &Boolean{Value: expr.Value}
	case *parser.List:
		return &List{Args: expr.Args}
	case *parser.Identifier:
		return evalIdent(expr, env)
	case *parser.MathProcedure:
		return evalMath(expr, env)
	case *parser.EqProcedure:
		return evalEq(expr, env)
	case *parser.DisplayProcedure:
		return evalDisplay(expr, env)
	case *parser.IfProcedure:
		return evalIf(expr, env)
	case *parser.ComparisonProcedure:
		return evalComparison(expr, env)
	case *parser.DefineProcedure:
		return evalDefine(expr, env)
	case *parser.LambdaProcedure:
		return evalLambda(expr, env)
	case *parser.Form:
		return evalForm(expr, env)
	}
	return nil
}

func evalForm(expr *parser.Form, env *Env) Object {
	o := Eval(expr.First, env)
	fn, ok := o.(*Lambda)
	if !ok {
		return newError("expected lambda, got %s", o.Type())
	}
	objs := make([]Object, 0)
	for _, arg := range expr.Rest {
		obj := Eval(arg, env)
		if isError(obj) {
			return obj
		}
		objs = append(objs, obj)
	}
	extendedEnv := extendFnEnv(fn, objs)
	var last Object
	for _, b := range fn.Body {
		last = Eval(b, extendedEnv)
	}
	return last
}

func extendFnEnv(fn *Lambda, args []Object) *Env {
	env := NewEnclosedEnv(fn.Env)
	for idx, param := range fn.Params {
		env.Set(param.Value, args[idx])
	}
	return env
}

func evalLambda(expr *parser.LambdaProcedure, env *Env) Object {
	return &Lambda{
		Params: expr.Params,
		Body:   expr.Body,
		Env:    env,
	}
}

func evalIdent(expr *parser.Identifier, env *Env) Object {
	if ident, ok := env.Get(expr.Value); ok {
		return ident
	}
	return newError("identifier not found: %s", expr.Value)
}

func evalDefine(expr *parser.DefineProcedure, env *Env) Object {
	obj := Eval(expr.Value, env)
	if isError(obj) {
		return obj
	}
	env.Set(expr.Name, obj)
	return obj
}

func evalComparison(expr *parser.ComparisonProcedure, env *Env) Object {
	left := Eval(expr.Left, env)
	if isError(left) {
		return left
	}
	right := Eval(expr.Right, env)
	if isError(right) {
		return right
	}
	if left.Type() != right.Type() {
		return newError("type mismatch - %s and %s", left.Type(), right.Type())
	}
	switch left.(type) {
	case *Number:
		return evalNumberCmp(expr.Operator, left.(*Number), right.(*Number))
	case *String:
		return evalStringCmp(expr.Operator, left.(*String), right.(*String))
	}
	return newError("unsupported type %s", left.Type())
}

func evalStringCmp(op string, left *String, right *String) Object {
	switch op {
	case "<":
		return &Boolean{Value: left.Value < right.Value}
	case ">":
		return &Boolean{Value: left.Value > right.Value}
	case "<=":
		return &Boolean{Value: left.Value <= right.Value}
	case ">=":
		return &Boolean{Value: left.Value >= right.Value}
	}
	return newError("unknown operator: %s", op)
}

func evalNumberCmp(op string, left *Number, right *Number) Object {
	switch op {
	case "<":
		return &Boolean{Value: left.Value < right.Value}
	case ">":
		return &Boolean{Value: left.Value > right.Value}
	case "<=":
		return &Boolean{Value: left.Value <= right.Value}
	case ">=":
		return &Boolean{Value: left.Value >= right.Value}
	}
	return newError("unknown operator: %s", op)
}

func evalIf(expr *parser.IfProcedure, env *Env) Object {
	cond := Eval(expr.Condition, env)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(expr.Consequence, env)
	} else {
		return Eval(expr.Alternative, env)
	}
}

func evalDisplay(expr *parser.DisplayProcedure, env *Env) Object {
	str := make([]string, 0)
	for _, expr := range expr.Args {
		obj := Eval(expr, env)
		if isError(obj) {
			return obj
		}
		str = append(str, obj.Inspect())
	}
	if len(str) > 0 {
		fmt.Println(strings.Join(str, " "))
	}
	return nil
}

func evalEq(expr *parser.EqProcedure, env *Env) Object {
	left := Eval(expr.Left, env)
	if isError(left) {
		return left
	}
	right := Eval(expr.Right, env)
	if isError(right) {
		return right
	}
	if left.Type() != right.Type() {
		return newError("type mismatch - %s and %s", left.Type(), right.Type())
	}
	switch left := left.(type) {
	case *Number:
		return &Boolean{Value: left.Value == right.(*Number).Value}
	case *String:
		return &Boolean{Value: left.Value == right.(*String).Value}
	case *Boolean:
		return &Boolean{Value: left.Value == right.(*Boolean).Value}
	default:
		return newError("eq on invalid type %s", left.Type())
	}
}

func evalMath(expr *parser.MathProcedure, env *Env) Object {
	objs := make([]*Number, 0)
	for _, expr := range expr.Args {
		obj := Eval(expr, env)
		if isError(obj) {
			return obj
		}
		if obj.Type() != NUMBER_OBJ {
			return newError("type mismatch - expected number, got %s", obj.Type())
		}
		objs = append(objs, obj.(*Number))
	}

	if len(objs) == 0 {
		return newError("no arguments")
	}

	val := objs[0].Value
	for i := 1; i < len(objs); i++ {
		obj := objs[i]
		switch expr.Operator {
		case "+":
			val = val + obj.Value
		case "-":
			val = val - obj.Value
		case "*":
			val = val * obj.Value
		case "/":
			val = val / obj.Value
		default:
			return newError("unknown operator: %s", expr.Operator)
		}
	}
	return &Number{Value: val}
}

func evalProgram(program *parser.Program, env *Env) Object {
	var last Object
	for _, expr := range program.Args {
		last = Eval(expr, env)
		if isError(last) {
			return last
		}
	}
	return last
}

func isError(obj Object) bool {
	if obj == nil {
		return false
	}
	return obj.Type() == ERROR_OBJ
}

func isTruthy(obj Object) bool {
	if obj == nil {
		return false
	}
	if obj.Type() == BOOLEAN_OBJ {
		return obj.(*Boolean).Value
	}
	return true
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
