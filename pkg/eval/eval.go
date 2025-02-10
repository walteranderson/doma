package eval

import (
	"doma/pkg/lexer"
	"doma/pkg/parser"
	"fmt"
	"strings"
)

func Eval(expr parser.Expression, env *Env) Object {
	switch expr := expr.(type) {
	case *parser.Number:
		return &Number{Value: expr.Value}
	case *parser.String:
		return &String{Value: expr.Value}
	case *parser.Boolean:
		return &Boolean{Value: expr.Value}
	case *parser.List:
		return &List{Args: expr.Args}
	case *parser.BuiltinIdentifier:
		return &Builtin{Value: expr.Token.Type}
	case *parser.Symbol:
		return &Symbol{Value: expr.Value}
	case *parser.Program:
		var last Object
		for _, expr := range expr.Args {
			last = Eval(expr, env)
			if isError(last) {
				return last
			}
		}
		return last
	case *parser.Identifier:
		if ident, ok := env.Get(expr.Value); ok {
			return ident
		}
		return newError("identifier not found: %s", expr.Value)
	case *parser.Form:
		obj := Eval(expr.First, env)
		switch obj := obj.(type) {
		case *Builtin:
			return applyBuiltin(obj, expr, env)
		case *Procedure:
			return applyLambda(obj.Value, expr, env)
		case *Lambda:
			return applyLambda(obj, expr, env)
		default:
			return newError("unknown procedure: %s", expr.First)
		}
	}
	return nil
}

func applyBuiltin(ident *Builtin, expr *parser.Form, env *Env) Object {
	switch ident.Value {
	case lexer.PLUS,
		lexer.MINUS,
		lexer.ASTERISK,
		lexer.SLASH:
		return evalMath(ident, expr, env)
	case lexer.EQ:
		return evalEq(expr, env)
	case lexer.DISPLAY:
		return evalDisplay(expr, env)
	case lexer.DEFINE:
		return evalDefine(expr, env)
	case lexer.STRINGREF:
		return evalStringRef(expr, env)
	case lexer.IF:
		return evalIf(expr, env)
	case lexer.LAMBDA:
		return evalLambda(expr, env)
	case lexer.LT,
		lexer.GT,
		lexer.LTE,
		lexer.GTE:
		return evalComparison(ident, expr, env)
	case lexer.FIRST:
		return evalFirst(expr, env)
	case lexer.REST:
		return evalRest(expr, env)
	default:
		return newError("unknown identifier: %s", ident.Value)
	}
}

func evalFirst(expr *parser.Form, env *Env) Object {
	if len(expr.Rest) != 1 {
		return newError("first expects 1 argument, got %d", len(expr.Rest))
	}
	obj := Eval(expr.Rest[0], env)
	if isError(obj) {
		return obj
	}
	lst, ok := obj.(*List)
	if !ok {
		return newError("first expects a list, received %s", obj.Type())
	}
	ident, ok := lst.Args[0].(*parser.Identifier)
	if ok {
		return &Symbol{Value: ident.Value}
	}
	return Eval(lst.Args[0], env)
}

func evalRest(expr *parser.Form, env *Env) Object {
	if len(expr.Rest) != 1 {
		return newError("first expects 1 argument, got %d", len(expr.Rest))
	}
	obj := Eval(expr.Rest[0], env)
	if isError(obj) {
		return obj
	}
	lst, ok := obj.(*List)
	if !ok {
		return newError("first expects a list, received %s", obj.Type())
	}
	lst.Args = lst.Args[1:]
	return lst
}

func evalStringRef(expr *parser.Form, env *Env) Object {
	if len(expr.Rest) != 2 {
		return newError("string-ref expects 2 arguments, got %d", len(expr.Rest))
	}
	strObj := Eval(expr.Rest[0], env)
	if isError(strObj) {
		return strObj
	}
	str, ok := strObj.(*String)
	if !ok {
		return newError("string-ref expects first argument to be a string, got %s", expr.Rest[0].String())
	}
	idxObj := Eval(expr.Rest[1], env)
	if isError(idxObj) {
		return idxObj
	}
	idx, ok := idxObj.(*Number)
	if !ok {
		return newError("string-ref expects second argument to be a number, got %s", expr.Rest[1].String())
	}
	val := str.Value[idx.Value]
	return &String{Value: string(val)}
}

func applyLambda(fn *Lambda, expr *parser.Form, env *Env) Object {
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

func evalLambda(expr *parser.Form, env *Env) Object {
	if len(expr.Rest) < 2 {
		return newError("lambda expects at least 2 arguments, got %d", len(expr.Rest))
	}
	lst, ok := expr.Rest[0].(*parser.List)
	if !ok {
		return newError("lambda expects first argument to be a list, got %s", expr.First.TokenLiteral())
	}
	params := make([]*parser.Identifier, 0)
	for _, arg := range lst.Args {
		ident, ok := arg.(*parser.Identifier)
		if !ok {
			return newError("lambda args expect to be all parameters to be identifiers, got %s", arg.TokenLiteral())
		}
		params = append(params, ident)
	}

	return &Lambda{
		Params: params,
		Body:   expr.Rest[1:],
		Env:    env,
	}
}

func evalDefine(expr *parser.Form, env *Env) Object {
	if len(expr.Rest) != 2 {
		return newError("define expects 2 arguments, got %d", len(expr.Rest))
	}
	name, ok := expr.Rest[0].(*parser.Identifier)
	if !ok {
		return newError("define expects first argument to be identifier, got %s", expr.Rest[0].TokenLiteral())
	}
	obj := Eval(expr.Rest[1], env)
	if isError(obj) {
		return obj
	}
	lambda, ok := obj.(*Lambda)
	if ok {
		proc := &Procedure{Name: name.Value, Value: lambda}
		env.Set(name.Value, proc)
		return proc
	} else {
		env.Set(name.Value, obj)
		return obj
	}
}

func evalComparison(ident *Builtin, expr *parser.Form, env *Env) Object {
	if len(expr.Rest) != 2 {
		return newError("%s expects 2 arguments, got %d", ident.Value, len(expr.Rest))
	}
	left := Eval(expr.Rest[0], env)
	if isError(left) {
		return left
	}
	right := Eval(expr.Rest[1], env)
	if isError(right) {
		return right
	}
	if left.Type() != right.Type() {
		return newError("type mismatch - %s and %s", left.Type(), right.Type())
	}
	switch left.(type) {
	case *Number:
		return evalNumberCmp(ident, left.(*Number), right.(*Number))
	case *String:
		return evalStringCmp(ident, left.(*String), right.(*String))
	}
	return newError("unsupported type %s", left.Type())
}

func evalStringCmp(op *Builtin, left *String, right *String) Object {
	switch op.Value {
	case lexer.LT:
		return &Boolean{Value: left.Value < right.Value}
	case lexer.GT:
		return &Boolean{Value: left.Value > right.Value}
	case lexer.LTE:
		return &Boolean{Value: left.Value <= right.Value}
	case lexer.GTE:
		return &Boolean{Value: left.Value >= right.Value}
	}
	return newError("unknown operator: %s", op)
}

func evalNumberCmp(op *Builtin, left *Number, right *Number) Object {
	switch op.Value {
	case lexer.LT:
		return &Boolean{Value: left.Value < right.Value}
	case lexer.GT:
		return &Boolean{Value: left.Value > right.Value}
	case lexer.LTE:
		return &Boolean{Value: left.Value <= right.Value}
	case lexer.GTE:
		return &Boolean{Value: left.Value >= right.Value}
	}
	return newError("unknown operator: %s", op)
}

func evalIf(expr *parser.Form, env *Env) Object {
	if len(expr.Rest) != 3 {
		return newError("if expects 3 arguments, got %d", len(expr.Rest))
	}
	cond := Eval(expr.Rest[0], env)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(expr.Rest[1], env)
	} else {
		return Eval(expr.Rest[2], env)
	}
}

func evalDisplay(expr *parser.Form, env *Env) Object {
	str := make([]string, 0)
	for _, expr := range expr.Rest {
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

func evalEq(expr *parser.Form, env *Env) Object {
	if len(expr.Rest) != 2 {
		return newError("eq expects 2 arguments, got %d", len(expr.Rest))
	}

	left := Eval(expr.Rest[0], env)
	if isError(left) {
		return left
	}
	right := Eval(expr.Rest[1], env)
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

func evalMath(op *Builtin, expr *parser.Form, env *Env) Object {
	objs := make([]*Number, 0)
	for _, expr := range expr.Rest {
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
		switch op.Value {
		case lexer.PLUS:
			val = val + obj.Value
		case lexer.MINUS:
			val = val - obj.Value
		case lexer.ASTERISK:
			val = val * obj.Value
		case lexer.SLASH:
			val = val / obj.Value
		default:
			return newError("unknown operator: %s", op.Value)
		}
	}
	return &Number{Value: val}
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
