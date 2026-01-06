package ast

import (
	"az-lang/token"
	"bytes"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier represents a variable name
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents a numeric value
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// StringLiteral represents a string value
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// BooleanLiteral represents a boolean value
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// ListLiteral represents a list
type ListLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (ll *ListLiteral) expressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }
func (ll *ListLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range ll.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("a list of ")
	out.WriteString(strings.Join(elements, " and "))
	return out.String()
}

// SetStatement represents: set x to 5
type SetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ss *SetStatement) statementNode()       {}
func (ss *SetStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *SetStatement) String() string {
	var out bytes.Buffer
	out.WriteString("set ")
	out.WriteString(ss.Name.String())
	out.WriteString(" to ")
	if ss.Value != nil {
		out.WriteString(ss.Value.String())
	}
	return out.String()
}

// ArithmeticExpression represents: x plus y, x minus y, x times y, x divided by y
type ArithmeticExpression struct {
	Token    token.Token
	Left     Expression
	Operator string // "plus", "minus", "times", "divided"
	Right    Expression
}

func (ae *ArithmeticExpression) expressionNode()      {}
func (ae *ArithmeticExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *ArithmeticExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Left.String())
	out.WriteString(" ")
	out.WriteString(ae.Operator)
	if ae.Operator == "divided" {
		out.WriteString(" by ")
	} else {
		out.WriteString(" ")
	}
	out.WriteString(ae.Right.String())
	return out.String()
}

// IncreaseStatement represents: increase x by 5
type IncreaseStatement struct {
	Token  token.Token
	Target *Identifier
	Amount Expression
}

func (is *IncreaseStatement) statementNode()       {}
func (is *IncreaseStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IncreaseStatement) String() string {
	var out bytes.Buffer
	out.WriteString("increase ")
	out.WriteString(is.Target.String())
	out.WriteString(" by ")
	out.WriteString(is.Amount.String())
	return out.String()
}

// DecreaseStatement represents: decrease x by 5
type DecreaseStatement struct {
	Token  token.Token
	Target *Identifier
	Amount Expression
}

func (ds *DecreaseStatement) statementNode()       {}
func (ds *DecreaseStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DecreaseStatement) String() string {
	var out bytes.Buffer
	out.WriteString("decrease ")
	out.WriteString(ds.Target.String())
	out.WriteString(" by ")
	out.WriteString(ds.Amount.String())
	return out.String()
}

// IfStatement represents: if x equals y then ... done otherwise ... done
type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(is.Condition.String())
	out.WriteString(" then ")
	out.WriteString(is.Consequence.String())
	if is.Alternative != nil {
		out.WriteString(" otherwise ")
		out.WriteString(is.Alternative.String())
	}
	return out.String()
}

// ComparisonExpression represents: x equals y, x is greater than y, etc.
type ComparisonExpression struct {
	Token    token.Token
	Left     Expression
	Operator string // "equals", "greater", "less"
	Right    Expression
}

func (ce *ComparisonExpression) expressionNode()      {}
func (ce *ComparisonExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *ComparisonExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Left.String())
	out.WriteString(" ")
	out.WriteString(ce.Operator)
	out.WriteString(" ")
	out.WriteString(ce.Right.String())
	return out.String()
}

// LogicalExpression represents: x and y, x or y, not x
type LogicalExpression struct {
	Token    token.Token
	Left     Expression
	Operator string // "and", "or", "not"
	Right    Expression
}

func (le *LogicalExpression) expressionNode()      {}
func (le *LogicalExpression) TokenLiteral() string { return le.Token.Literal }
func (le *LogicalExpression) String() string {
	var out bytes.Buffer
	if le.Left != nil {
		out.WriteString(le.Left.String())
		out.WriteString(" ")
	}
	out.WriteString(le.Operator)
	if le.Right != nil {
		out.WriteString(" ")
		out.WriteString(le.Right.String())
	}
	return out.String()
}

// WhileStatement represents: while x is less than 100 do ... done
type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" do ")
	out.WriteString(ws.Body.String())
	return out.String()
}

// ForStatement represents: for each item in items do ... done
type ForStatement struct {
	Token    token.Token
	Variable *Identifier
	Iterable Expression
	Body     *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for each ")
	out.WriteString(fs.Variable.String())
	out.WriteString(" in ")
	out.WriteString(fs.Iterable.String())
	out.WriteString(" do ")
	out.WriteString(fs.Body.String())
	return out.String()
}

// BlockStatement represents a block of statements ending with done
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
		out.WriteString(" ")
	}
	out.WriteString("done")
	return out.String()
}

// FunctionDefinition represents: to greet with name ... done
type FunctionDefinition struct {
	Token      token.Token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fd *FunctionDefinition) statementNode()       {}
func (fd *FunctionDefinition) TokenLiteral() string { return fd.Token.Literal }
func (fd *FunctionDefinition) String() string {
	var out bytes.Buffer
	out.WriteString("to ")
	out.WriteString(fd.Name.String())
	if len(fd.Parameters) > 0 {
		out.WriteString(" with ")
		params := []string{}
		for _, p := range fd.Parameters {
			params = append(params, p.String())
		}
		out.WriteString(strings.Join(params, " and "))
	}
	out.WriteString(" ")
	out.WriteString(fd.Body.String())
	return out.String()
}

// CallExpression represents: funcname with args
type CallExpression struct {
	Token     token.Token
	Function  *Identifier
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	if len(ce.Arguments) > 0 {
		out.WriteString(" with ")
		args := []string{}
		for _, a := range ce.Arguments {
			args = append(args, a.String())
		}
		out.WriteString(strings.Join(args, " and "))
	}
	return out.String()
}

// ReturnStatement represents: return x
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString("return ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	return out.String()
}

// SayStatement represents: say x
type SayStatement struct {
	Token token.Token
	Value Expression
}

func (ss *SayStatement) statementNode()       {}
func (ss *SayStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *SayStatement) String() string {
	var out bytes.Buffer
	out.WriteString("say ")
	out.WriteString(ss.Value.String())
	return out.String()
}

// AskStatement represents: ask into answer
type AskStatement struct {
	Token  token.Token
	Target *Identifier
}

func (as *AskStatement) statementNode()       {}
func (as *AskStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AskStatement) String() string {
	var out bytes.Buffer
	out.WriteString("ask into ")
	out.WriteString(as.Target.String())
	return out.String()
}

// LengthExpression represents: length of items (as expression)
type LengthExpression struct {
	Token token.Token
	List  Expression
}

func (le *LengthExpression) expressionNode()      {}
func (le *LengthExpression) TokenLiteral() string { return le.Token.Literal }
func (le *LengthExpression) String() string {
	return "length of " + le.List.String()
}

// AppendStatement represents: append value to items
type AppendStatement struct {
	Token token.Token
	Value Expression
	List  *Identifier
}

func (as *AppendStatement) statementNode()       {}
func (as *AppendStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AppendStatement) String() string {
	var out bytes.Buffer
	out.WriteString("append ")
	out.WriteString(as.Value.String())
	out.WriteString(" to ")
	out.WriteString(as.List.String())
	return out.String()
}

// IndexExpression represents: item N from list (as expression)
type IndexExpression struct {
	Token token.Token
	Index Expression
	List  Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return "item " + ie.Index.String() + " from " + ie.List.String()
}

// NegativeExpression represents: minus 5
type NegativeExpression struct {
	Token token.Token
	Value Expression
}

func (ne *NegativeExpression) expressionNode()      {}
func (ne *NegativeExpression) TokenLiteral() string { return ne.Token.Literal }
func (ne *NegativeExpression) String() string {
	return "minus " + ne.Value.String()
}

// === HTTP AST Nodes ===

// FetchStatement represents: fetch from "URL" into response
type FetchStatement struct {
	Token   token.Token
	URL     Expression
	Headers Expression   // optional: headers list
	Target  *Identifier
}

func (fs *FetchStatement) statementNode()       {}
func (fs *FetchStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FetchStatement) String() string {
	var out bytes.Buffer
	out.WriteString("fetch from ")
	out.WriteString(fs.URL.String())
	if fs.Headers != nil {
		out.WriteString(" with ")
		out.WriteString(fs.Headers.String())
	}
	out.WriteString(" into ")
	out.WriteString(fs.Target.String())
	return out.String()
}

// SendStatement represents: send "body" to "URL" into response (POST)
type SendStatement struct {
	Token   token.Token
	Body    Expression
	URL     Expression
	Headers Expression // optional
	Target  *Identifier
}

func (ss *SendStatement) statementNode()       {}
func (ss *SendStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *SendStatement) String() string {
	var out bytes.Buffer
	out.WriteString("send ")
	out.WriteString(ss.Body.String())
	out.WriteString(" to ")
	out.WriteString(ss.URL.String())
	if ss.Headers != nil {
		out.WriteString(" with ")
		out.WriteString(ss.Headers.String())
	}
	out.WriteString(" into ")
	out.WriteString(ss.Target.String())
	return out.String()
}

// PutStatement represents: put "body" to "URL" into response
type PutStatement struct {
	Token   token.Token
	Body    Expression
	URL     Expression
	Headers Expression // optional
	Target  *Identifier
}

func (ps *PutStatement) statementNode()       {}
func (ps *PutStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PutStatement) String() string {
	var out bytes.Buffer
	out.WriteString("put ")
	out.WriteString(ps.Body.String())
	out.WriteString(" to ")
	out.WriteString(ps.URL.String())
	if ps.Headers != nil {
		out.WriteString(" with ")
		out.WriteString(ps.Headers.String())
	}
	out.WriteString(" into ")
	out.WriteString(ps.Target.String())
	return out.String()
}

// DeleteStatement represents: delete from "URL" into response
type DeleteStatement struct {
	Token   token.Token
	URL     Expression
	Headers Expression // optional
	Target  *Identifier
}

func (ds *DeleteStatement) statementNode()       {}
func (ds *DeleteStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DeleteStatement) String() string {
	var out bytes.Buffer
	out.WriteString("delete from ")
	out.WriteString(ds.URL.String())
	if ds.Headers != nil {
		out.WriteString(" with ")
		out.WriteString(ds.Headers.String())
	}
	out.WriteString(" into ")
	out.WriteString(ds.Target.String())
	return out.String()
}

// BodyOfExpression represents: body of response
type BodyOfExpression struct {
	Token    token.Token
	Response Expression
}

func (boe *BodyOfExpression) expressionNode()      {}
func (boe *BodyOfExpression) TokenLiteral() string { return boe.Token.Literal }
func (boe *BodyOfExpression) String() string {
	return "body of " + boe.Response.String()
}

// StatusOfExpression represents: status of response
type StatusOfExpression struct {
	Token    token.Token
	Response Expression
}

func (soe *StatusOfExpression) expressionNode()      {}
func (soe *StatusOfExpression) TokenLiteral() string { return soe.Token.Literal }
func (soe *StatusOfExpression) String() string {
	return "status of " + soe.Response.String()
}

// HeaderFromExpression represents: header "Name" from response
type HeaderFromExpression struct {
	Token      token.Token
	HeaderName Expression
	Response   Expression
}

func (hfe *HeaderFromExpression) expressionNode()      {}
func (hfe *HeaderFromExpression) TokenLiteral() string { return hfe.Token.Literal }
func (hfe *HeaderFromExpression) String() string {
	return "header " + hfe.HeaderName.String() + " from " + hfe.Response.String()
}

// === JSON AST Nodes ===

// ParseJsonStatement represents: parse X as json into Y
type ParseJsonStatement struct {
	Token  token.Token
	Source Expression
	Target *Identifier
}

func (pjs *ParseJsonStatement) statementNode()       {}
func (pjs *ParseJsonStatement) TokenLiteral() string { return pjs.Token.Literal }
func (pjs *ParseJsonStatement) String() string {
	var out bytes.Buffer
	out.WriteString("parse ")
	out.WriteString(pjs.Source.String())
	out.WriteString(" as json into ")
	out.WriteString(pjs.Target.String())
	return out.String()
}

// FieldFromExpression represents: field "name" from data
type FieldFromExpression struct {
	Token     token.Token
	FieldName Expression
	Source    Expression
}

func (ffe *FieldFromExpression) expressionNode()      {}
func (ffe *FieldFromExpression) TokenLiteral() string { return ffe.Token.Literal }
func (ffe *FieldFromExpression) String() string {
	return "field " + ffe.FieldName.String() + " from " + ffe.Source.String()
}

// EncodeJsonStatement represents: encode X as json into Y
type EncodeJsonStatement struct {
	Token  token.Token
	Source Expression
	Target *Identifier
}

func (ejs *EncodeJsonStatement) statementNode()       {}
func (ejs *EncodeJsonStatement) TokenLiteral() string { return ejs.Token.Literal }
func (ejs *EncodeJsonStatement) String() string {
	var out bytes.Buffer
	out.WriteString("encode ")
	out.WriteString(ejs.Source.String())
	out.WriteString(" as json into ")
	out.WriteString(ejs.Target.String())
	return out.String()
}

// === Web Server AST Nodes ===

// ServeStatement represents: serve on 8080 or serve on 8080 in background
type ServeStatement struct {
	Token      token.Token
	Port       Expression
	Background bool
}

func (ss *ServeStatement) statementNode()       {}
func (ss *ServeStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *ServeStatement) String() string {
	var out bytes.Buffer
	out.WriteString("serve on ")
	out.WriteString(ss.Port.String())
	if ss.Background {
		out.WriteString(" in background")
	}
	return out.String()
}

// WhenRouteStatement represents: when get at "/path" using req do ... done
type WhenRouteStatement struct {
	Token      token.Token
	Method     string          // "" for any, "GET", "POST", etc.
	Path       Expression
	RequestVar *Identifier     // optional request variable
	Body       *BlockStatement
}

func (wr *WhenRouteStatement) statementNode()       {}
func (wr *WhenRouteStatement) TokenLiteral() string { return wr.Token.Literal }
func (wr *WhenRouteStatement) String() string {
	var out bytes.Buffer
	out.WriteString("when ")
	if wr.Method != "" {
		out.WriteString(strings.ToLower(wr.Method))
		out.WriteString(" ")
	} else {
		out.WriteString("request ")
	}
	out.WriteString("at ")
	out.WriteString(wr.Path.String())
	if wr.RequestVar != nil {
		out.WriteString(" using ")
		out.WriteString(wr.RequestVar.String())
	}
	out.WriteString(" do ")
	out.WriteString(wr.Body.String())
	return out.String()
}

// RouteToStatement represents: route "/path" to handlerFunc
type RouteToStatement struct {
	Token   token.Token
	Path    Expression
	Handler *Identifier
}

func (rt *RouteToStatement) statementNode()       {}
func (rt *RouteToStatement) TokenLiteral() string { return rt.Token.Literal }
func (rt *RouteToStatement) String() string {
	var out bytes.Buffer
	out.WriteString("route ")
	out.WriteString(rt.Path.String())
	out.WriteString(" to ")
	out.WriteString(rt.Handler.String())
	return out.String()
}

// HeaderPair represents a header name-value pair for responses
type HeaderPair struct {
	Name  Expression
	Value Expression
}

// ReplyStatement represents: reply with "data" with status 201 or reply with data as json
type ReplyStatement struct {
	Token      token.Token
	Body       Expression
	AsJson     bool           // if true, auto-encode body as JSON
	StatusCode Expression     // optional, defaults to 200
	Headers    []HeaderPair   // optional response headers
}

func (rs *ReplyStatement) statementNode()       {}
func (rs *ReplyStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReplyStatement) String() string {
	var out bytes.Buffer
	out.WriteString("reply with ")
	out.WriteString(rs.Body.String())
	if rs.AsJson {
		out.WriteString(" as json")
	}
	if rs.StatusCode != nil {
		out.WriteString(" with status ")
		out.WriteString(rs.StatusCode.String())
	}
	return out.String()
}

// StopServerStatement represents: stop server or stop server on 8080
type StopServerStatement struct {
	Token token.Token
	Port  Expression // optional, nil means stop all
}

func (ss *StopServerStatement) statementNode()       {}
func (ss *StopServerStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *StopServerStatement) String() string {
	var out bytes.Buffer
	out.WriteString("stop server")
	if ss.Port != nil {
		out.WriteString(" on ")
		out.WriteString(ss.Port.String())
	}
	return out.String()
}

// MethodOfExpression represents: method of req
type MethodOfExpression struct {
	Token   token.Token
	Request Expression
}

func (moe *MethodOfExpression) expressionNode()      {}
func (moe *MethodOfExpression) TokenLiteral() string { return moe.Token.Literal }
func (moe *MethodOfExpression) String() string {
	return "method of " + moe.Request.String()
}

// PathOfExpression represents: path of req
type PathOfExpression struct {
	Token   token.Token
	Request Expression
}

func (poe *PathOfExpression) expressionNode()      {}
func (poe *PathOfExpression) TokenLiteral() string { return poe.Token.Literal }
func (poe *PathOfExpression) String() string {
	return "path of " + poe.Request.String()
}

// QueryFromExpression represents: query "name" from req
type QueryFromExpression struct {
	Token     token.Token
	QueryName Expression
	Request   Expression
}

func (qfe *QueryFromExpression) expressionNode()      {}
func (qfe *QueryFromExpression) TokenLiteral() string { return qfe.Token.Literal }
func (qfe *QueryFromExpression) String() string {
	return "query " + qfe.QueryName.String() + " from " + qfe.Request.String()
}
