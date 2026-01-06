package interpreter

import (
	"az-lang/ast"
	"az-lang/object"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.SetStatement:
		return evalSetStatement(node, env)
	case *ast.IncreaseStatement:
		return evalIncreaseStatement(node, env)
	case *ast.DecreaseStatement:
		return evalDecreaseStatement(node, env)
	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.WhileStatement:
		return evalWhileStatement(node, env)
	case *ast.ForStatement:
		return evalForStatement(node, env)
	case *ast.FunctionDefinition:
		return evalFunctionDefinition(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.SayStatement:
		return evalSayStatement(node, env)
	case *ast.AskStatement:
		return evalAskStatement(node, env)
	case *ast.AppendStatement:
		return evalAppendStatement(node, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.NegativeExpression:
		return evalNegativeExpression(node, env)
	case *ast.ListLiteral:
		return evalListLiteral(node, env)
	case *ast.ComparisonExpression:
		return evalComparisonExpression(node, env)
	case *ast.LogicalExpression:
		return evalLogicalExpression(node, env)
	case *ast.ArithmeticExpression:
		return evalArithmeticExpression(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.LengthExpression:
		return evalLengthExpression(node, env)
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)

	// HTTP Statements
	case *ast.FetchStatement:
		return evalFetchStatement(node, env)
	case *ast.SendStatement:
		return evalSendStatement(node, env)
	case *ast.PutStatement:
		return evalPutStatement(node, env)
	case *ast.DeleteStatement:
		return evalDeleteStatement(node, env)

	// HTTP Response Expressions
	case *ast.BodyOfExpression:
		return evalBodyOfExpression(node, env)
	case *ast.StatusOfExpression:
		return evalStatusOfExpression(node, env)
	case *ast.HeaderFromExpression:
		return evalHeaderFromExpression(node, env)

	// JSON Statements
	case *ast.ParseJsonStatement:
		return evalParseJsonStatement(node, env)
	case *ast.EncodeJsonStatement:
		return evalEncodeJsonStatement(node, env)

	// JSON Expressions
	case *ast.FieldFromExpression:
		return evalFieldFromExpression(node, env)

	// Web Server Statements
	case *ast.ServeStatement:
		return evalServeStatement(node, env)
	case *ast.WhenRouteStatement:
		return evalWhenRouteStatement(node, env)
	case *ast.RouteToStatement:
		return evalRouteToStatement(node, env)
	case *ast.ReplyStatement:
		return evalReplyStatement(node, env)
	case *ast.StopServerStatement:
		return evalStopServerStatement(node, env)

	// Web Server Request Expressions
	case *ast.MethodOfExpression:
		return evalMethodOfExpression(node, env)
	case *ast.PathOfExpression:
		return evalPathOfExpression(node, env)
	case *ast.QueryFromExpression:
		return evalQueryFromExpression(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalSetStatement(ss *ast.SetStatement, env *object.Environment) object.Object {
	val := Eval(ss.Value, env)
	if isError(val) {
		return val
	}
	env.Set(ss.Name.Value, val)
	return val
}

func evalIncreaseStatement(is *ast.IncreaseStatement, env *object.Environment) object.Object {
	currentVal, ok := env.Get(is.Target.Value)
	if !ok {
		return newError("undefined variable: %s", is.Target.Value)
	}

	currentInt, ok := currentVal.(*object.Integer)
	if !ok {
		return newError("increase requires an integer variable, got %s", currentVal.Type())
	}

	amount := Eval(is.Amount, env)
	if isError(amount) {
		return amount
	}

	amountInt, ok := amount.(*object.Integer)
	if !ok {
		return newError("increase amount must be an integer, got %s", amount.Type())
	}

	result := &object.Integer{Value: currentInt.Value + amountInt.Value}
	env.Set(is.Target.Value, result)
	return result
}

func evalDecreaseStatement(ds *ast.DecreaseStatement, env *object.Environment) object.Object {
	currentVal, ok := env.Get(ds.Target.Value)
	if !ok {
		return newError("undefined variable: %s", ds.Target.Value)
	}

	currentInt, ok := currentVal.(*object.Integer)
	if !ok {
		return newError("decrease requires an integer variable, got %s", currentVal.Type())
	}

	amount := Eval(ds.Amount, env)
	if isError(amount) {
		return amount
	}

	amountInt, ok := amount.(*object.Integer)
	if !ok {
		return newError("decrease amount must be an integer, got %s", amount.Type())
	}

	result := &object.Integer{Value: currentInt.Value - amountInt.Value}
	env.Set(ds.Target.Value, result)
	return result
}

func evalArithmeticExpression(ae *ast.ArithmeticExpression, env *object.Environment) object.Object {
	left := Eval(ae.Left, env)
	if isError(left) {
		return left
	}

	right := Eval(ae.Right, env)
	if isError(right) {
		return right
	}

	// Handle string concatenation with "plus"
	if ae.Operator == "plus" {
		leftStr, leftIsString := left.(*object.String)
		rightStr, rightIsString := right.(*object.String)
		if leftIsString || rightIsString {
			// Convert both to strings and concatenate
			var leftVal, rightVal string
			if leftIsString {
				leftVal = leftStr.Value
			} else {
				leftVal = left.Inspect()
			}
			if rightIsString {
				rightVal = rightStr.Value
			} else {
				rightVal = right.Inspect()
			}
			return &object.String{Value: leftVal + rightVal}
		}
	}

	leftVal, ok := left.(*object.Integer)
	if !ok {
		return newError("arithmetic operations require integers, got %s", left.Type())
	}

	rightVal, ok := right.(*object.Integer)
	if !ok {
		return newError("arithmetic operations require integers, got %s", right.Type())
	}

	var result int64
	switch ae.Operator {
	case "plus":
		result = leftVal.Value + rightVal.Value
	case "minus":
		result = leftVal.Value - rightVal.Value
	case "times":
		result = leftVal.Value * rightVal.Value
	case "divided":
		if rightVal.Value == 0 {
			return newError("division by zero")
		}
		result = leftVal.Value / rightVal.Value
	}

	return &object.Integer{Value: result}
}

func evalLogicalExpression(le *ast.LogicalExpression, env *object.Environment) object.Object {
	switch le.Operator {
	case "not":
		right := Eval(le.Right, env)
		if isError(right) {
			return right
		}
		return nativeBoolToBooleanObject(!isTruthy(right))

	case "and":
		left := Eval(le.Left, env)
		if isError(left) {
			return left
		}
		if !isTruthy(left) {
			return FALSE
		}
		right := Eval(le.Right, env)
		if isError(right) {
			return right
		}
		return nativeBoolToBooleanObject(isTruthy(right))

	case "or":
		left := Eval(le.Left, env)
		if isError(left) {
			return left
		}
		if isTruthy(left) {
			return TRUE
		}
		right := Eval(le.Right, env)
		if isError(right) {
			return right
		}
		return nativeBoolToBooleanObject(isTruthy(right))
	}

	return FALSE
}

func evalIfStatement(is *ast.IfStatement, env *object.Environment) object.Object {
	condition := Eval(is.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(is.Consequence, env)
	} else if is.Alternative != nil {
		return Eval(is.Alternative, env)
	}
	return NULL
}

func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	for {
		condition := Eval(ws.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(ws.Body, env)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	iterable := Eval(fs.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	list, ok := iterable.(*object.List)
	if !ok {
		return newError("for each requires a list, got %s", iterable.Type())
	}

	var result object.Object = NULL

	for _, element := range list.Elements {
		env.Set(fs.Variable.Value, element)
		result = Eval(fs.Body, env)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalFunctionDefinition(fd *ast.FunctionDefinition, env *object.Environment) object.Object {
	fn := &object.Function{
		Parameters: fd.Parameters,
		Body:       fd.Body,
		Env:        env,
	}
	env.Set(fd.Name.Value, fn)
	return fn
}

func evalCallExpression(ce *ast.CallExpression, env *object.Environment) object.Object {
	fnObj, ok := env.Get(ce.Function.Value)
	if !ok {
		return newError("function not defined: %s", ce.Function.Value)
	}

	fn, ok := fnObj.(*object.Function)
	if !ok {
		return newError("%s is not a function", ce.Function.Value)
	}

	// Evaluate arguments
	args := []object.Object{}
	for _, arg := range ce.Arguments {
		evaluated := Eval(arg, env)
		if isError(evaluated) {
			return evaluated
		}
		args = append(args, evaluated)
	}

	// Create new environment for function
	extendedEnv := object.NewEnclosedEnvironment(fn.Env)

	// Bind parameters
	for i, param := range fn.Parameters {
		if i < len(args) {
			extendedEnv.Set(param.Value, args[i])
		}
	}

	// Execute function body
	result := Eval(fn.Body, extendedEnv)

	// Unwrap return value
	if returnValue, ok := result.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return result
}

func evalReturnStatement(rs *ast.ReturnStatement, env *object.Environment) object.Object {
	if rs.ReturnValue == nil {
		return &object.ReturnValue{Value: NULL}
	}

	val := Eval(rs.ReturnValue, env)
	if isError(val) {
		return val
	}
	return &object.ReturnValue{Value: val}
}

func evalSayStatement(ss *ast.SayStatement, env *object.Environment) object.Object {
	val := Eval(ss.Value, env)
	if isError(val) {
		return val
	}
	fmt.Println(val.Inspect())
	return NULL
}

func evalAskStatement(as *ast.AskStatement, env *object.Environment) object.Object {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return newError("error reading input: %s", err)
	}

	// Remove newline
	if len(input) > 0 && input[len(input)-1] == '\n' {
		input = input[:len(input)-1]
	}
	if len(input) > 0 && input[len(input)-1] == '\r' {
		input = input[:len(input)-1]
	}

	result := &object.String{Value: input}
	env.Set(as.Target.Value, result)
	return result
}

func evalLengthExpression(le *ast.LengthExpression, env *object.Environment) object.Object {
	val := Eval(le.List, env)
	if isError(val) {
		return val
	}

	switch v := val.(type) {
	case *object.List:
		return &object.Integer{Value: int64(len(v.Elements))}
	case *object.String:
		return &object.Integer{Value: int64(len(v.Value))}
	default:
		return newError("length requires a list or string, got %s", val.Type())
	}
}

func evalIndexExpression(ie *ast.IndexExpression, env *object.Environment) object.Object {
	index := Eval(ie.Index, env)
	if isError(index) {
		return index
	}

	list := Eval(ie.List, env)
	if isError(list) {
		return list
	}

	idx, ok := index.(*object.Integer)
	if !ok {
		return newError("index must be an integer, got %s", index.Type())
	}

	switch l := list.(type) {
	case *object.List:
		if idx.Value < 1 || idx.Value > int64(len(l.Elements)) {
			return newError("index out of bounds: %d (list has %d elements)", idx.Value, len(l.Elements))
		}
		return l.Elements[idx.Value-1] // 1-indexed
	case *object.String:
		if idx.Value < 1 || idx.Value > int64(len(l.Value)) {
			return newError("index out of bounds: %d (string has %d characters)", idx.Value, len(l.Value))
		}
		return &object.String{Value: string(l.Value[idx.Value-1])} // 1-indexed
	default:
		return newError("indexing requires a list or string, got %s", list.Type())
	}
}

func evalAppendStatement(as *ast.AppendStatement, env *object.Environment) object.Object {
	value := Eval(as.Value, env)
	if isError(value) {
		return value
	}

	listObj, ok := env.Get(as.List.Value)
	if !ok {
		return newError("undefined variable: %s", as.List.Value)
	}

	list, ok := listObj.(*object.List)
	if !ok {
		return newError("append requires a list, got %s", listObj.Type())
	}

	list.Elements = append(list.Elements, value)
	return NULL
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	// Handle special keywords
	if node.Value == "null" {
		return NULL
	}
	if node.Value == "true" {
		return TRUE
	}
	if node.Value == "false" {
		return FALSE
	}

	val, ok := env.Get(node.Value)
	if !ok {
		return newError("undefined variable: %s", node.Value)
	}
	return val
}

func evalNegativeExpression(ne *ast.NegativeExpression, env *object.Environment) object.Object {
	val := Eval(ne.Value, env)
	if isError(val) {
		return val
	}

	intVal, ok := val.(*object.Integer)
	if !ok {
		return newError("minus requires an integer, got %s", val.Type())
	}

	return &object.Integer{Value: -intVal.Value}
}

func evalListLiteral(ll *ast.ListLiteral, env *object.Environment) object.Object {
	elements := []object.Object{}
	for _, elem := range ll.Elements {
		evaluated := Eval(elem, env)
		if isError(evaluated) {
			return evaluated
		}
		elements = append(elements, evaluated)
	}
	return &object.List{Elements: elements}
}

func evalComparisonExpression(ce *ast.ComparisonExpression, env *object.Environment) object.Object {
	left := Eval(ce.Left, env)
	if isError(left) {
		return left
	}

	right := Eval(ce.Right, env)
	if isError(right) {
		return right
	}

	switch ce.Operator {
	case "equals":
		return evalEquals(left, right)
	case "greater":
		return evalGreater(left, right)
	case "less":
		return evalLess(left, right)
	}

	return FALSE
}

func evalEquals(left, right object.Object) object.Object {
	// Handle null comparison
	if left.Type() == object.NULL_OBJ && right.Type() == object.NULL_OBJ {
		return TRUE
	}
	if left.Type() == object.NULL_OBJ || right.Type() == object.NULL_OBJ {
		return FALSE
	}

	switch l := left.(type) {
	case *object.Integer:
		if r, ok := right.(*object.Integer); ok {
			return nativeBoolToBooleanObject(l.Value == r.Value)
		}
	case *object.String:
		if r, ok := right.(*object.String); ok {
			return nativeBoolToBooleanObject(l.Value == r.Value)
		}
	case *object.Boolean:
		if r, ok := right.(*object.Boolean); ok {
			return nativeBoolToBooleanObject(l.Value == r.Value)
		}
	}
	return FALSE
}

func evalGreater(left, right object.Object) object.Object {
	leftInt, ok := left.(*object.Integer)
	if !ok {
		return newError("comparison requires integers, got %s", left.Type())
	}

	rightInt, ok := right.(*object.Integer)
	if !ok {
		return newError("comparison requires integers, got %s", right.Type())
	}

	return nativeBoolToBooleanObject(leftInt.Value > rightInt.Value)
}

func evalLess(left, right object.Object) object.Object {
	leftInt, ok := left.(*object.Integer)
	if !ok {
		return newError("comparison requires integers, got %s", left.Type())
	}

	rightInt, ok := right.(*object.Integer)
	if !ok {
		return newError("comparison requires integers, got %s", right.Type())
	}

	return nativeBoolToBooleanObject(leftInt.Value < rightInt.Value)
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Null:
		return false
	case *object.Boolean:
		return obj.Value
	case *object.Integer:
		return obj.Value != 0
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// HTTP Interpreter Functions

func evalFetchStatement(node *ast.FetchStatement, env *object.Environment) object.Object {
	url := Eval(node.URL, env)
	if isError(url) {
		return url
	}

	urlStr, ok := url.(*object.String)
	if !ok {
		return newError("fetch URL must be a string, got %s", url.Type())
	}

	var headers *object.List
	if node.Headers != nil {
		headersObj := Eval(node.Headers, env)
		if isError(headersObj) {
			return headersObj
		}
		headers, ok = headersObj.(*object.List)
		if !ok {
			return newError("headers must be a list, got %s", headersObj.Type())
		}
	}

	response, err := executeRequest("GET", urlStr.Value, "", headers)
	if err != nil {
		return newError("fetch failed: %s", err.Error())
	}

	env.Set(node.Target.Value, response)
	return response
}

func evalSendStatement(node *ast.SendStatement, env *object.Environment) object.Object {
	body := Eval(node.Body, env)
	if isError(body) {
		return body
	}

	bodyStr, ok := body.(*object.String)
	if !ok {
		return newError("send body must be a string, got %s", body.Type())
	}

	url := Eval(node.URL, env)
	if isError(url) {
		return url
	}

	urlStr, ok := url.(*object.String)
	if !ok {
		return newError("send URL must be a string, got %s", url.Type())
	}

	var headers *object.List
	if node.Headers != nil {
		headersObj := Eval(node.Headers, env)
		if isError(headersObj) {
			return headersObj
		}
		headers, ok = headersObj.(*object.List)
		if !ok {
			return newError("headers must be a list, got %s", headersObj.Type())
		}
	}

	response, err := executeRequest("POST", urlStr.Value, bodyStr.Value, headers)
	if err != nil {
		return newError("send failed: %s", err.Error())
	}

	env.Set(node.Target.Value, response)
	return response
}

func evalPutStatement(node *ast.PutStatement, env *object.Environment) object.Object {
	body := Eval(node.Body, env)
	if isError(body) {
		return body
	}

	bodyStr, ok := body.(*object.String)
	if !ok {
		return newError("put body must be a string, got %s", body.Type())
	}

	url := Eval(node.URL, env)
	if isError(url) {
		return url
	}

	urlStr, ok := url.(*object.String)
	if !ok {
		return newError("put URL must be a string, got %s", url.Type())
	}

	var headers *object.List
	if node.Headers != nil {
		headersObj := Eval(node.Headers, env)
		if isError(headersObj) {
			return headersObj
		}
		headers, ok = headersObj.(*object.List)
		if !ok {
			return newError("headers must be a list, got %s", headersObj.Type())
		}
	}

	response, err := executeRequest("PUT", urlStr.Value, bodyStr.Value, headers)
	if err != nil {
		return newError("put failed: %s", err.Error())
	}

	env.Set(node.Target.Value, response)
	return response
}

func evalDeleteStatement(node *ast.DeleteStatement, env *object.Environment) object.Object {
	url := Eval(node.URL, env)
	if isError(url) {
		return url
	}

	urlStr, ok := url.(*object.String)
	if !ok {
		return newError("delete URL must be a string, got %s", url.Type())
	}

	var headers *object.List
	if node.Headers != nil {
		headersObj := Eval(node.Headers, env)
		if isError(headersObj) {
			return headersObj
		}
		headers, ok = headersObj.(*object.List)
		if !ok {
			return newError("headers must be a list, got %s", headersObj.Type())
		}
	}

	response, err := executeRequest("DELETE", urlStr.Value, "", headers)
	if err != nil {
		return newError("delete failed: %s", err.Error())
	}

	env.Set(node.Target.Value, response)
	return response
}

func evalBodyOfExpression(node *ast.BodyOfExpression, env *object.Environment) object.Object {
	respObj := Eval(node.Response, env)
	if isError(respObj) {
		return respObj
	}

	// Handle both Response and Request objects
	switch obj := respObj.(type) {
	case *object.Response:
		return &object.String{Value: obj.Body}
	case *object.Request:
		return &object.String{Value: obj.Body}
	default:
		return newError("body of requires a response or request, got %s", respObj.Type())
	}
}

func evalStatusOfExpression(node *ast.StatusOfExpression, env *object.Environment) object.Object {
	respObj := Eval(node.Response, env)
	if isError(respObj) {
		return respObj
	}

	response, ok := respObj.(*object.Response)
	if !ok {
		return newError("status of requires a response, got %s", respObj.Type())
	}

	return &object.Integer{Value: int64(response.StatusCode)}
}

func evalHeaderFromExpression(node *ast.HeaderFromExpression, env *object.Environment) object.Object {
	headerName := Eval(node.HeaderName, env)
	if isError(headerName) {
		return headerName
	}

	headerStr, ok := headerName.(*object.String)
	if !ok {
		return newError("header name must be a string, got %s", headerName.Type())
	}

	respObj := Eval(node.Response, env)
	if isError(respObj) {
		return respObj
	}

	// Handle both Response and Request objects
	var headers map[string]string
	switch obj := respObj.(type) {
	case *object.Response:
		headers = obj.Headers
	case *object.Request:
		headers = obj.Headers
	default:
		return newError("header from requires a response or request, got %s", respObj.Type())
	}

	value, exists := headers[headerStr.Value]
	if !exists {
		return NULL
	}

	return &object.String{Value: value}
}

// HTTP Helper Functions

func executeRequest(method, url, body string, headers *object.List) (*object.Response, error) {
	var req *http.Request
	var err error

	if body != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	if headers != nil {
		applyHeaders(req, headers)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respHeaders := make(map[string]string)
	for key := range resp.Header {
		respHeaders[key] = resp.Header.Get(key)
	}

	return &object.Response{
		StatusCode: resp.StatusCode,
		Body:       string(respBody),
		Headers:    respHeaders,
	}, nil
}

func applyHeaders(req *http.Request, headers *object.List) {
	for _, elem := range headers.Elements {
		if str, ok := elem.(*object.String); ok {
			parts := strings.SplitN(str.Value, ":", 2)
			if len(parts) == 2 {
				req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}
}

// JSON Interpreter Functions

func evalParseJsonStatement(node *ast.ParseJsonStatement, env *object.Environment) object.Object {
	source := Eval(node.Source, env)
	if isError(source) {
		return source
	}

	sourceStr, ok := source.(*object.String)
	if !ok {
		return newError("parse json requires a string, got %s", source.Type())
	}

	var result interface{}
	if err := json.Unmarshal([]byte(sourceStr.Value), &result); err != nil {
		return newError("invalid JSON: %s", err.Error())
	}

	jsonObj := &object.Json{Value: result}
	env.Set(node.Target.Value, jsonObj)
	return jsonObj
}

func evalEncodeJsonStatement(node *ast.EncodeJsonStatement, env *object.Environment) object.Object {
	source := Eval(node.Source, env)
	if isError(source) {
		return source
	}

	var value interface{}

	switch src := source.(type) {
	case *object.Json:
		value = src.Value
	case *object.String:
		value = src.Value
	case *object.Integer:
		value = src.Value
	case *object.Boolean:
		value = src.Value
	case *object.List:
		arr := make([]interface{}, len(src.Elements))
		for i, elem := range src.Elements {
			arr[i] = objectToInterface(elem)
		}
		value = arr
	default:
		return newError("cannot encode %s as json", source.Type())
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return newError("json encoding failed: %s", err.Error())
	}

	result := &object.String{Value: string(bytes)}
	env.Set(node.Target.Value, result)
	return result
}

func evalFieldFromExpression(node *ast.FieldFromExpression, env *object.Environment) object.Object {
	fieldName := Eval(node.FieldName, env)
	if isError(fieldName) {
		return fieldName
	}

	fieldStr, ok := fieldName.(*object.String)
	if !ok {
		return newError("field name must be a string, got %s", fieldName.Type())
	}

	source := Eval(node.Source, env)
	if isError(source) {
		return source
	}

	jsonObj, ok := source.(*object.Json)
	if !ok {
		return newError("field from requires a json object, got %s", source.Type())
	}

	result := getJsonField(jsonObj.Value, fieldStr.Value)
	return result
}

// JSON Helper Functions

func getJsonField(data interface{}, path string) object.Object {
	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return NULL
			}
		default:
			return NULL
		}
	}

	return interfaceToObject(current)
}

func interfaceToObject(val interface{}) object.Object {
	switch v := val.(type) {
	case nil:
		return NULL
	case bool:
		return nativeBoolToBooleanObject(v)
	case float64:
		return &object.Integer{Value: int64(v)}
	case string:
		return &object.String{Value: v}
	case []interface{}:
		elements := make([]object.Object, len(v))
		for i, elem := range v {
			elements[i] = interfaceToObject(elem)
		}
		return &object.List{Elements: elements}
	case map[string]interface{}:
		return &object.Json{Value: v}
	default:
		return NULL
	}
}

func objectToInterface(obj object.Object) interface{} {
	switch o := obj.(type) {
	case *object.Integer:
		return o.Value
	case *object.String:
		return o.Value
	case *object.Boolean:
		return o.Value
	case *object.Null:
		return nil
	case *object.List:
		arr := make([]interface{}, len(o.Elements))
		for i, elem := range o.Elements {
			arr[i] = objectToInterface(elem)
		}
		return arr
	case *object.Json:
		return o.Value
	default:
		return nil
	}
}

// === Web Server Implementation ===

// ServerInfo holds information about a running server
type ServerInfo struct {
	Port       int
	Server     *http.Server
	Running    bool
}

// RouteHandler holds route handler information
type RouteHandler struct {
	Method     string // "" for any method
	Path       string
	Body       *ast.BlockStatement
	RequestVar string
	HandlerEnv *object.Environment
	HandlerFn  *object.Function // for function reference handlers
}

// Server and route registries
var (
	serverRegistry = make(map[int]*ServerInfo)
	routeRegistry  = make(map[int][]RouteHandler)
	registryMu     sync.RWMutex
	defaultPort    = 8080
)

// evalServeStatement starts an HTTP server
func evalServeStatement(node *ast.ServeStatement, env *object.Environment) object.Object {
	portObj := Eval(node.Port, env)
	if isError(portObj) {
		return portObj
	}

	portInt, ok := portObj.(*object.Integer)
	if !ok {
		return newError("serve port must be an integer, got %s", portObj.Type())
	}

	port := int(portInt.Value)

	registryMu.Lock()
	if _, exists := serverRegistry[port]; exists {
		registryMu.Unlock()
		return newError("server already running on port %d", port)
	}
	registryMu.Unlock()

	mux := http.NewServeMux()

	// Set up a catch-all handler that dispatches to registered routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleIncomingRequest(w, r, port, env)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	serverInfo := &ServerInfo{
		Port:    port,
		Server:  server,
		Running: true,
	}

	registryMu.Lock()
	serverRegistry[port] = serverInfo
	// Copy any routes registered to defaultPort to this port if different
	if port != defaultPort {
		if existingRoutes, ok := routeRegistry[defaultPort]; ok && len(existingRoutes) > 0 {
			routeRegistry[port] = append(routeRegistry[port], existingRoutes...)
		}
	}
	defaultPort = port
	registryMu.Unlock()

	serverObj := &object.Server{
		Port:    port,
		Running: true,
	}

	if node.Background {
		go func() {
			fmt.Printf("Server started in background on port %d\n", port)
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				fmt.Printf("Server error on port %d: %s\n", port, err)
			}
			registryMu.Lock()
			if info, exists := serverRegistry[port]; exists {
				info.Running = false
			}
			registryMu.Unlock()
		}()
		return serverObj
	} else {
		fmt.Printf("Server starting on port %d (foreground)...\n", port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return newError("server error: %s", err)
		}
		return NULL
	}
}

// evalWhenRouteStatement registers an inline route handler
func evalWhenRouteStatement(node *ast.WhenRouteStatement, env *object.Environment) object.Object {
	pathObj := Eval(node.Path, env)
	if isError(pathObj) {
		return pathObj
	}

	pathStr, ok := pathObj.(*object.String)
	if !ok {
		return newError("route path must be a string, got %s", pathObj.Type())
	}

	requestVarName := ""
	if node.RequestVar != nil {
		requestVarName = node.RequestVar.Value
	}

	handler := RouteHandler{
		Method:     node.Method,
		Path:       pathStr.Value,
		Body:       node.Body,
		RequestVar: requestVarName,
		HandlerEnv: env,
	}

	registryMu.Lock()
	routeRegistry[defaultPort] = append(routeRegistry[defaultPort], handler)
	registryMu.Unlock()

	return NULL
}

// evalRouteToStatement registers a function reference as route handler
func evalRouteToStatement(node *ast.RouteToStatement, env *object.Environment) object.Object {
	pathObj := Eval(node.Path, env)
	if isError(pathObj) {
		return pathObj
	}

	pathStr, ok := pathObj.(*object.String)
	if !ok {
		return newError("route path must be a string, got %s", pathObj.Type())
	}

	fnObj, ok := env.Get(node.Handler.Value)
	if !ok {
		return newError("handler function not defined: %s", node.Handler.Value)
	}

	fn, ok := fnObj.(*object.Function)
	if !ok {
		return newError("%s is not a function", node.Handler.Value)
	}

	handler := RouteHandler{
		Method:     "", // any method
		Path:       pathStr.Value,
		HandlerFn:  fn,
		HandlerEnv: env,
	}

	registryMu.Lock()
	routeRegistry[defaultPort] = append(routeRegistry[defaultPort], handler)
	registryMu.Unlock()

	return NULL
}

// evalReplyStatement creates a response object
func evalReplyStatement(node *ast.ReplyStatement, env *object.Environment) object.Object {
	bodyObj := Eval(node.Body, env)
	if isError(bodyObj) {
		return bodyObj
	}

	var bodyStr string
	headers := make(map[string]string)

	if node.AsJson {
		// Auto-encode body as JSON
		jsonBytes, err := json.Marshal(objectToInterface(bodyObj))
		if err != nil {
			return newError("failed to encode as JSON: %s", err)
		}
		bodyStr = string(jsonBytes)
		headers["Content-Type"] = "application/json"
	} else {
		switch b := bodyObj.(type) {
		case *object.String:
			bodyStr = b.Value
		case *object.Json:
			jsonBytes, _ := json.Marshal(b.Value)
			bodyStr = string(jsonBytes)
		default:
			bodyStr = bodyObj.Inspect()
		}
	}

	statusCode := 200
	if node.StatusCode != nil {
		statusObj := Eval(node.StatusCode, env)
		if isError(statusObj) {
			return statusObj
		}
		if sc, ok := statusObj.(*object.Integer); ok {
			statusCode = int(sc.Value)
		}
	}

	// Process additional headers
	for _, hp := range node.Headers {
		nameObj := Eval(hp.Name, env)
		valueObj := Eval(hp.Value, env)
		if nameStr, ok := nameObj.(*object.String); ok {
			if valueStr, ok := valueObj.(*object.String); ok {
				headers[nameStr.Value] = valueStr.Value
			}
		}
	}

	return &object.ReplyValue{
		Body:       bodyStr,
		StatusCode: statusCode,
		Headers:    headers,
	}
}

// evalStopServerStatement stops a running server
func evalStopServerStatement(node *ast.StopServerStatement, env *object.Environment) object.Object {
	registryMu.Lock()
	defer registryMu.Unlock()

	if node.Port != nil {
		portObj := Eval(node.Port, env)
		if isError(portObj) {
			return portObj
		}

		portInt, ok := portObj.(*object.Integer)
		if !ok {
			return newError("stop server port must be an integer, got %s", portObj.Type())
		}

		port := int(portInt.Value)
		serverInfo, exists := serverRegistry[port]
		if !exists {
			return newError("no server running on port %d", port)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := serverInfo.Server.Shutdown(ctx); err != nil {
			return newError("error stopping server: %s", err)
		}

		delete(serverRegistry, port)
		delete(routeRegistry, port)
		fmt.Printf("Server on port %d stopped\n", port)
		return NULL
	}

	// Stop all servers
	for port, serverInfo := range serverRegistry {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		serverInfo.Server.Shutdown(ctx)
		cancel()
		delete(serverRegistry, port)
		delete(routeRegistry, port)
		fmt.Printf("Server on port %d stopped\n", port)
	}

	return NULL
}

// evalMethodOfExpression extracts method from request
func evalMethodOfExpression(node *ast.MethodOfExpression, env *object.Environment) object.Object {
	reqObj := Eval(node.Request, env)
	if isError(reqObj) {
		return reqObj
	}

	req, ok := reqObj.(*object.Request)
	if !ok {
		return newError("method of requires a request, got %s", reqObj.Type())
	}

	return &object.String{Value: req.Method}
}

// evalPathOfExpression extracts path from request
func evalPathOfExpression(node *ast.PathOfExpression, env *object.Environment) object.Object {
	reqObj := Eval(node.Request, env)
	if isError(reqObj) {
		return reqObj
	}

	req, ok := reqObj.(*object.Request)
	if !ok {
		return newError("path of requires a request, got %s", reqObj.Type())
	}

	return &object.String{Value: req.Path}
}

// evalQueryFromExpression extracts query parameter from request
func evalQueryFromExpression(node *ast.QueryFromExpression, env *object.Environment) object.Object {
	queryName := Eval(node.QueryName, env)
	if isError(queryName) {
		return queryName
	}

	queryStr, ok := queryName.(*object.String)
	if !ok {
		return newError("query name must be a string, got %s", queryName.Type())
	}

	reqObj := Eval(node.Request, env)
	if isError(reqObj) {
		return reqObj
	}

	req, ok := reqObj.(*object.Request)
	if !ok {
		return newError("query from requires a request, got %s", reqObj.Type())
	}

	value, exists := req.QueryParams[queryStr.Value]
	if !exists {
		return NULL
	}

	return &object.String{Value: value}
}

// handleIncomingRequest dispatches incoming HTTP requests to registered handlers
func handleIncomingRequest(w http.ResponseWriter, r *http.Request, port int, env *object.Environment) {
	registryMu.RLock()
	routes := routeRegistry[port]
	registryMu.RUnlock()

	// Find matching route
	for _, route := range routes {
		if matchRoute(route, r) {
			// Build Request object
			body, _ := io.ReadAll(r.Body)
			headers := make(map[string]string)
			for key := range r.Header {
				headers[key] = r.Header.Get(key)
			}
			queryParams := make(map[string]string)
			for key, values := range r.URL.Query() {
				if len(values) > 0 {
					queryParams[key] = values[0]
				}
			}

			reqObj := &object.Request{
				Method:      r.Method,
				Path:        r.URL.Path,
				Body:        string(body),
				Headers:     headers,
				QueryParams: queryParams,
			}

			var result object.Object

			if route.HandlerFn != nil {
				// Function reference handler
				extendedEnv := object.NewEnclosedEnvironment(route.HandlerFn.Env)
				if len(route.HandlerFn.Parameters) > 0 {
					extendedEnv.Set(route.HandlerFn.Parameters[0].Value, reqObj)
				}
				result = Eval(route.HandlerFn.Body, extendedEnv)
				if returnValue, ok := result.(*object.ReturnValue); ok {
					result = returnValue.Value
				}
			} else {
				// Inline block handler
				handlerScope := object.NewEnclosedEnvironment(route.HandlerEnv)
				if route.RequestVar != "" {
					handlerScope.Set(route.RequestVar, reqObj)
				}
				result = Eval(route.Body, handlerScope)
				if returnValue, ok := result.(*object.ReturnValue); ok {
					result = returnValue.Value
				}
			}

			// Send response
			if rv, ok := result.(*object.ReplyValue); ok {
				for name, value := range rv.Headers {
					w.Header().Set(name, value)
				}
				w.WriteHeader(rv.StatusCode)
				w.Write([]byte(rv.Body))
				return
			}

			// Default response for non-reply returns
			w.WriteHeader(200)
			if result != nil {
				w.Write([]byte(result.Inspect()))
			}
			return
		}
	}

	// No matching route
	w.WriteHeader(404)
	w.Write([]byte("Not Found"))
}

// matchRoute checks if a route matches the incoming request
func matchRoute(route RouteHandler, r *http.Request) bool {
	// Method matching
	if route.Method != "" && route.Method != r.Method {
		return false
	}

	// Simple path matching (exact match)
	return route.Path == r.URL.Path
}
