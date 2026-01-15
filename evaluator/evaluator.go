package evaluator

import (
	"fmt"
	"luma/ast"
	"reflect"
	"sort"
	"strings"
)

type ReturnValue struct{ Value interface{} }

func FormatLumaValue(v interface{}) string {
	switch val := v.(type) {
	case int:
		return fmt.Sprintf("%d", val)
	case string:
		return val
	case bool:
		return fmt.Sprintf("%v", val)
	case []interface{}:
		var out strings.Builder
		out.WriteString("[")
		for i, el := range val {
			out.WriteString(FormatLumaValue(el))
			if i < len(val)-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString("]")
		return out.String()
	case map[string]interface{}:
		var out strings.Builder
		out.WriteString("{")
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			out.WriteString(k + ": ")
			v := val[k]
			if s, ok := v.(string); ok {
				out.WriteString("\"" + s + "\"")
			} else {
				out.WriteString(FormatLumaValue(v))
			}
			if i < len(keys)-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString("}")
		return out.String()
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func Eval(node ast.Node, env *Env) interface{} {
	if node == nil {
		return nil
	}
	val := reflect.ValueOf(node)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return nil
	}

	switch n := node.(type) {
	case *ast.Program:
		var result interface{}
		for _, stmt := range n.Statements {
			result = Eval(stmt, env)
			if rv, ok := result.(*ReturnValue); ok {
				return rv.Value
			}
		}
		return result
	case *ast.BlockStatement:
		var result interface{}
		for _, stmt := range n.Statements {
			result = Eval(stmt, env)
			if result != nil {
				if _, ok := result.(*ReturnValue); ok {
					return result
				}
			}
		}
		return result
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(n.Value, env)
		return &ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(n.Value, env)
		env.Set(n.Name.Value, val, n.IsConst)
		return val
	case *ast.LoopStatement:
		loopEnv := NewEnclosedEnv(env)
		Eval(n.Init, loopEnv)
		for {
			cond := Eval(n.Condition, loopEnv)
			condBool, ok := cond.(bool)
			if !ok || !condBool {
				break
			}
			Eval(n.Body, loopEnv)
			Eval(n.Post, loopEnv)
		}
		return nil

	case *ast.NumberLiteral:
		return n.Value
	case *ast.StringLiteral:
		return n.Value
	case *ast.BooleanLiteral:
		return n.Value
	case *ast.Identifier:
		val, ok := env.Get(n.Value)
		if !ok {
			if n.Value == "log" {
				return "BUILTIN_LOG"
			}
			if n.Value == "len" {
				return "BUILTIN_LEN"
			}
			fmt.Printf("ERROR: identifier not found: %s\n", n.Value)
			return nil
		}
		return val
	case *ast.ArrayLiteral:
		elements := []interface{}{}
		for _, el := range n.Elements {
			elements = append(elements, Eval(el, env))
		}
		return elements
	case *ast.ObjectLiteral:
		obj := make(map[string]interface{})
		for k, v := range n.Pairs {
			obj[k] = Eval(v, env)
		}
		return obj

	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if right == nil {
			return nil
		}
		if n.Operator == "-" {
			if r, ok := right.(int); ok {
				return -r
			}
		}
		return nil
	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		right := Eval(n.Right, env)
		if left == nil || right == nil {
			return nil
		}
		switch n.Operator {
		case "+":
			if l, ok := left.(string); ok {
				if r, ok := right.(string); ok {
					return l + r
				}
			}
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok {
					return l + r
				}
			}
			return nil
		case "-":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok {
					return l - r
				}
			}
			return nil
		case "*":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok {
					return l * r
				}
			}
			return nil
		case "/":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok {
					return l / r
				}
			}
			return nil
		case "=":
			if ident, ok := n.Left.(*ast.Identifier); ok {
				val := Eval(n.Right, env)
				env.Update(ident.Value, val)
				return val
			}
			return nil
		case "==":
			return left == right
		case "!=":
			return left != right
		case "<":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok {
					return l < r
				}
			}
			return false
		case ">":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok {
					return l > r
				}
			}
			return false
		case "<=":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok {
					return l <= r
				}
			}
			return false
		case ">=":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok {
					return l >= r
				}
			}
			return false
		}
	case *ast.IfExpression:
		cond := Eval(n.Condition, env)
		if condBool, ok := cond.(bool); ok && condBool {
			return Eval(n.Consequence, env)
		} else if n.Alternative != nil {
			return Eval(n.Alternative, env)
		}
		return nil
	case *ast.FunctionLiteral:
		return n
	case *ast.CallExpression:
		fn := Eval(n.Function, env)
		args := []interface{}{}
		for _, a := range n.Arguments {
			args = append(args, Eval(a, env))
		}

		if s, ok := fn.(string); ok {
			if s == "BUILTIN_LOG" {
				formatted := []interface{}{}
				for _, arg := range args {
					formatted = append(formatted, FormatLumaValue(arg))
				}
				fmt.Println(formatted...)
				return nil
			}
			if s == "BUILTIN_LEN" {
				if arr, ok := args[0].([]interface{}); ok {
					return len(arr)
				}
				if str, ok := args[0].(string); ok {
					return len(str)
				}
				return 0
			}
		}

		if fl, ok := fn.(*ast.FunctionLiteral); ok {
			childEnv := NewEnclosedEnv(env)
			for i, p := range fl.Parameters {
				childEnv.Set(p.Value, args[i], false)
			}
			res := Eval(fl.Body, childEnv)
			if rv, ok := res.(*ReturnValue); ok {
				return rv.Value
			}
			return res
		}
	case *ast.IndexExpression:
		left := Eval(n.Left, env)
		idx := Eval(n.Index, env)
		if arr, ok := left.([]interface{}); ok {
			if i, ok := idx.(int); ok && i >= 0 && i < len(arr) {
				return arr[i]
			}
		}
		return nil
	case *ast.MemberExpression:
		left := Eval(n.Left, env)
		if obj, ok := left.(map[string]interface{}); ok {
			return obj[n.Member.Value]
		}
		return nil
	}
	return nil
}
