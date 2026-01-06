package object

import (
	"az-lang/ast"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	LIST_OBJ         = "LIST"
	RESPONSE_OBJ     = "RESPONSE"
	JSON_OBJ         = "JSON"
	REQUEST_OBJ      = "REQUEST"
	SERVER_OBJ       = "SERVER"
	REPLY_VALUE_OBJ  = "REPLY_VALUE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer represents an integer value
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// String represents a string value
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Boolean represents a boolean value
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Null represents a null value
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue wraps a value being returned from a function
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error represents a runtime error
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Function represents a user-defined function
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("to function")
	if len(params) > 0 {
		out.WriteString(" with ")
		out.WriteString(strings.Join(params, " and "))
	}
	out.WriteString(" z\n")
	out.WriteString(f.Body.String())

	return out.String()
}

// List represents a list/array of values
type List struct {
	Elements []Object
}

func (l *List) Type() ObjectType { return LIST_OBJ }
func (l *List) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range l.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

func (r *Response) Type() ObjectType { return RESPONSE_OBJ }
func (r *Response) Inspect() string {
	return fmt.Sprintf("Response{status: %d, body: %q}", r.StatusCode, r.Body)
}

// Json represents a parsed JSON value
type Json struct {
	Value interface{} // Holds parsed JSON (map[string]interface{} or []interface{})
}

func (j *Json) Type() ObjectType { return JSON_OBJ }
func (j *Json) Inspect() string {
	bytes, err := json.Marshal(j.Value)
	if err != nil {
		return "invalid json"
	}
	return string(bytes)
}

// Environment holds variable bindings
type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Request represents an incoming HTTP request
type Request struct {
	Method      string
	Path        string
	Body        string
	Headers     map[string]string
	QueryParams map[string]string
}

func (r *Request) Type() ObjectType { return REQUEST_OBJ }
func (r *Request) Inspect() string {
	return fmt.Sprintf("Request{method: %s, path: %s}", r.Method, r.Path)
}

// Server represents a running HTTP server (metadata only - actual server managed by interpreter)
type Server struct {
	Port    int
	Running bool
}

func (s *Server) Type() ObjectType { return SERVER_OBJ }
func (s *Server) Inspect() string {
	return fmt.Sprintf("Server{port: %d, running: %t}", s.Port, s.Running)
}

// ReplyValue represents a response to be sent
type ReplyValue struct {
	Body       string
	StatusCode int
	Headers    map[string]string
}

func (rv *ReplyValue) Type() ObjectType { return REPLY_VALUE_OBJ }
func (rv *ReplyValue) Inspect() string {
	return fmt.Sprintf("Reply{status: %d, body: %q}", rv.StatusCode, rv.Body)
}
