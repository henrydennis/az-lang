package parser

import (
	"az-lang/ast"
	"az-lang/lexer"
	"az-lang/token"
	"fmt"
	"strconv"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("line %d: expected next token to be %s, got %s instead",
		p.peekToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.SET:
		return p.parseSetStatement()
	case token.INCREASE:
		return p.parseIncreaseStatement()
	case token.DECREASE:
		return p.parseDecreaseStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.TO:
		return p.parseFunctionDefinition()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.SAY:
		return p.parseSayStatement()
	case token.ASK:
		return p.parseAskStatement()
	case token.APPEND:
		return p.parseAppendStatement()
	case token.FETCH:
		return p.parseFetchStatement()
	case token.SEND:
		return p.parseSendStatement()
	case token.PUT:
		return p.parsePutStatement()
	case token.DELETE:
		return p.parseDeleteStatement()
	case token.PARSE:
		return p.parseParseJsonStatement()
	case token.ENCODE:
		return p.parseEncodeJsonStatement()
	case token.SERVE:
		return p.parseServeStatement()
	case token.WHEN:
		return p.parseWhenRouteStatement()
	case token.ROUTE:
		return p.parseRouteToStatement()
	case token.REPLY:
		return p.parseReplyStatement()
	case token.STOP:
		return p.parseStopServerStatement()
	default:
		return nil
	}
}

// parseSetStatement parses: set x to 5
func (p *Parser) parseSetStatement() *ast.SetStatement {
	stmt := &ast.SetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.TO) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression()

	return stmt
}

// parseIncreaseStatement parses: increase x by 5
func (p *Parser) parseIncreaseStatement() *ast.IncreaseStatement {
	stmt := &ast.IncreaseStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Target = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.BY) {
		return nil
	}

	p.nextToken()
	stmt.Amount = p.parseExpression()

	return stmt
}

// parseDecreaseStatement parses: decrease x by 5
func (p *Parser) parseDecreaseStatement() *ast.DecreaseStatement {
	stmt := &ast.DecreaseStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Target = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.BY) {
		return nil
	}

	p.nextToken()
	stmt.Amount = p.parseExpression()

	return stmt
}

// parseIfStatement parses: if x equals y then ... done otherwise ... done
func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseCondition()

	if !p.expectPeek(token.THEN) {
		return nil
	}

	p.nextToken() // move past THEN
	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.OTHERWISE) {
		p.nextToken() // consume OTHERWISE
		p.nextToken() // move to first statement of alternative
		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

// parseCondition parses comparison and logical expressions
func (p *Parser) parseCondition() ast.Expression {
	return p.parseLogicalOr()
}

// parseLogicalOr handles: x or y
func (p *Parser) parseLogicalOr() ast.Expression {
	left := p.parseLogicalAnd()

	for p.peekTokenIs(token.OR) {
		p.nextToken() // consume OR
		opToken := p.curToken
		p.nextToken()
		right := p.parseLogicalAnd()
		left = &ast.LogicalExpression{
			Token:    opToken,
			Left:     left,
			Operator: "or",
			Right:    right,
		}
	}

	return left
}

// parseLogicalAnd handles: x and y
func (p *Parser) parseLogicalAnd() ast.Expression {
	left := p.parseLogicalNot()

	for p.peekTokenIs(token.AND) {
		p.nextToken() // consume AND
		opToken := p.curToken
		p.nextToken()
		right := p.parseLogicalNot()
		left = &ast.LogicalExpression{
			Token:    opToken,
			Left:     left,
			Operator: "and",
			Right:    right,
		}
	}

	return left
}

// parseLogicalNot handles: not x
func (p *Parser) parseLogicalNot() ast.Expression {
	if p.curTokenIs(token.NOT) {
		opToken := p.curToken
		p.nextToken()
		right := p.parseLogicalNot()
		return &ast.LogicalExpression{
			Token:    opToken,
			Left:     nil,
			Operator: "not",
			Right:    right,
		}
	}

	return p.parseComparison()
}

// parseComparison handles: x equals y, x is greater than y, x is less than y
func (p *Parser) parseComparison() ast.Expression {
	left := p.parseArithmeticExpression()

	// Check for comparison operators
	if p.peekTokenIs(token.EQUALS) {
		p.nextToken() // consume EQUALS
		opToken := p.curToken
		p.nextToken()
		right := p.parseArithmeticExpression()
		return &ast.ComparisonExpression{
			Token:    opToken,
			Left:     left,
			Operator: "equals",
			Right:    right,
		}
	}

	if p.peekTokenIs(token.IS) {
		p.nextToken() // consume IS

		if p.peekTokenIs(token.GREATER) {
			p.nextToken() // consume GREATER
			opToken := p.curToken
			if !p.expectPeek(token.THAN) {
				return nil
			}
			p.nextToken()
			right := p.parseArithmeticExpression()
			return &ast.ComparisonExpression{
				Token:    opToken,
				Left:     left,
				Operator: "greater",
				Right:    right,
			}
		}

		if p.peekTokenIs(token.LESS) {
			p.nextToken() // consume LESS
			opToken := p.curToken
			if !p.expectPeek(token.THAN) {
				return nil
			}
			p.nextToken()
			right := p.parseArithmeticExpression()
			return &ast.ComparisonExpression{
				Token:    opToken,
				Left:     left,
				Operator: "less",
				Right:    right,
			}
		}
	}

	return left
}

// parseArithmeticExpression handles: x plus y, x minus y, x times y, x divided by y
func (p *Parser) parseArithmeticExpression() ast.Expression {
	left := p.parseTerm()

	for p.peekTokenIs(token.PLUS) || p.peekTokenIs(token.MINUS) {
		p.nextToken() // consume operator
		opToken := p.curToken
		op := opToken.Literal
		p.nextToken()
		right := p.parseTerm()
		left = &ast.ArithmeticExpression{
			Token:    opToken,
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left
}

// parseTerm handles: x times y, x divided by y
func (p *Parser) parseTerm() ast.Expression {
	left := p.parsePrimary()

	for p.peekTokenIs(token.TIMES) || p.peekTokenIs(token.DIVIDED) {
		p.nextToken() // consume operator
		opToken := p.curToken
		op := opToken.Literal

		if op == "divided" {
			// Expect "by" after "divided"
			if !p.expectPeek(token.BY) {
				return nil
			}
		}

		p.nextToken()
		right := p.parsePrimary()
		left = &ast.ArithmeticExpression{
			Token:    opToken,
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left
}

// parsePrimary handles primary expressions
func (p *Parser) parsePrimary() ast.Expression {
	// Handle negative numbers: minus 5
	if p.curTokenIs(token.MINUS) {
		negToken := p.curToken
		p.nextToken()
		value := p.parsePrimary()
		return &ast.NegativeExpression{Token: negToken, Value: value}
	}

	// Handle quoted strings
	if p.curTokenIs(token.STRING) {
		return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	}

	// Handle "a list of" list literals
	if p.curTokenIs(token.A) && p.peekTokenIs(token.LIST) {
		return p.parseListLiteral()
	}

	// Handle "length of" expression
	if p.curTokenIs(token.LENGTH) && p.peekTokenIs(token.OF) {
		return p.parseLengthExpression()
	}

	// Handle "item N from list" expression
	if p.curTokenIs(token.ITEM) {
		return p.parseIndexExpression()
	}

	// Handle "body of" expression
	if p.curTokenIs(token.BODY) && p.peekTokenIs(token.OF) {
		return p.parseBodyOfExpression()
	}

	// Handle "status of" expression
	if p.curTokenIs(token.STATUS) && p.peekTokenIs(token.OF) {
		return p.parseStatusOfExpression()
	}

	// Handle "header X from" expression
	if p.curTokenIs(token.HEADER) {
		return p.parseHeaderFromExpression()
	}

	// Handle "field X from" expression
	if p.curTokenIs(token.FIELD) {
		return p.parseFieldFromExpression()
	}

	// Handle "method of" expression
	if p.curTokenIs(token.METHOD) && p.peekTokenIs(token.OF) {
		return p.parseMethodOfExpression()
	}

	// Handle "path of" expression
	if p.curTokenIs(token.PATH) && p.peekTokenIs(token.OF) {
		return p.parsePathOfExpression()
	}

	// Handle "query X from" expression
	if p.curTokenIs(token.QUERY) {
		return p.parseQueryFromExpression()
	}

	// Handle number words
	if token.IsNumberWord(p.curToken.Type) {
		return p.parseNumberWord()
	}

	// Handle numeric literals
	if p.curTokenIs(token.NUMBER) {
		value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("could not parse %q as integer", p.curToken.Literal))
			return nil
		}
		return &ast.IntegerLiteral{Token: p.curToken, Value: value}
	}

	// Handle identifiers (including function calls)
	if p.curTokenIs(token.IDENT) {
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		// Check if this is a function call: funcname with args
		if p.peekTokenIs(token.WITH) {
			return p.parseCallExpression(ident)
		}

		return ident
	}

	return nil
}

// parseCallExpression parses: funcname with arg1 and arg2
func (p *Parser) parseCallExpression(fn *ast.Identifier) *ast.CallExpression {
	call := &ast.CallExpression{Token: fn.Token, Function: fn}
	call.Arguments = []ast.Expression{}

	p.nextToken() // consume WITH
	p.nextToken() // move to first argument

	arg := p.parsePrimary()
	call.Arguments = append(call.Arguments, arg)

	// Handle multiple arguments with "and"
	for p.peekTokenIs(token.AND) {
		p.nextToken() // consume AND
		p.nextToken() // move to argument
		arg := p.parsePrimary()
		call.Arguments = append(call.Arguments, arg)
	}

	return call
}

// parseWhileStatement parses: while x is less than 100 do ... done
func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseCondition()

	if !p.expectPeek(token.DO) {
		return nil
	}

	p.nextToken() // move past DO
	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseForStatement parses: for each item in items do ... done
func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	if !p.expectPeek(token.EACH) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Variable = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.IN) {
		return nil
	}

	p.nextToken()
	stmt.Iterable = p.parsePrimary()

	if !p.expectPeek(token.DO) {
		return nil
	}

	p.nextToken() // move past DO
	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseBlockStatement parses statements until "done"
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	for !p.curTokenIs(token.DONE) && !p.curTokenIs(token.OTHERWISE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseFunctionDefinition parses: to funcname with param1 and param2 ... done
func (p *Parser) parseFunctionDefinition() *ast.FunctionDefinition {
	stmt := &ast.FunctionDefinition{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	stmt.Parameters = []*ast.Identifier{}

	// Check for parameters
	if p.peekTokenIs(token.WITH) {
		p.nextToken() // consume WITH
		p.nextToken() // move to first parameter

		param := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		stmt.Parameters = append(stmt.Parameters, param)

		// Handle multiple parameters with "and"
		for p.peekTokenIs(token.AND) {
			p.nextToken() // consume AND
			p.nextToken() // move to parameter
			param := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			stmt.Parameters = append(stmt.Parameters, param)
		}
	}

	p.nextToken() // move to body
	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseReturnStatement parses: return x
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// Check if there's a return value (not at end of block or EOF)
	if !p.curTokenIs(token.DONE) && !p.curTokenIs(token.EOF) {
		stmt.ReturnValue = p.parseExpression()
	}

	return stmt
}

// parseSayStatement parses: say x or say "hello"
func (p *Parser) parseSayStatement() *ast.SayStatement {
	stmt := &ast.SayStatement{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression()

	return stmt
}

// parseAskStatement parses: ask into answer
func (p *Parser) parseAskStatement() *ast.AskStatement {
	stmt := &ast.AskStatement{Token: p.curToken}

	if !p.expectPeek(token.INTO) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Target = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

// parseLengthExpression parses: length of items
func (p *Parser) parseLengthExpression() *ast.LengthExpression {
	expr := &ast.LengthExpression{Token: p.curToken}

	p.nextToken() // consume LENGTH, now at OF
	p.nextToken() // consume OF, now at list expression

	expr.List = p.parsePrimary()

	return expr
}

// parseIndexExpression parses: item N from list
func (p *Parser) parseIndexExpression() *ast.IndexExpression {
	expr := &ast.IndexExpression{Token: p.curToken}

	p.nextToken() // move past ITEM
	expr.Index = p.parsePrimary()

	if !p.expectPeek(token.FROM) {
		return nil
	}

	p.nextToken()
	expr.List = p.parsePrimary()

	return expr
}

// parseAppendStatement parses: append value to items
func (p *Parser) parseAppendStatement() *ast.AppendStatement {
	stmt := &ast.AppendStatement{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression()

	if !p.expectPeek(token.TO) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.List = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

// parseExpression is the main entry point for expression parsing
func (p *Parser) parseExpression() ast.Expression {
	return p.parseArithmeticExpression()
}

// parseListLiteral parses: a list of 1 and 2 and 3
func (p *Parser) parseListLiteral() *ast.ListLiteral {
	list := &ast.ListLiteral{Token: p.curToken}
	list.Elements = []ast.Expression{}

	// curToken is A, advance to LIST
	if !p.expectPeek(token.LIST) {
		return nil
	}

	// curToken is LIST, advance to OF
	if !p.expectPeek(token.OF) {
		return nil
	}

	// curToken is OF, advance to first element
	p.nextToken()
	elem := p.parsePrimary()
	list.Elements = append(list.Elements, elem)

	for p.peekTokenIs(token.AND) {
		p.nextToken() // consume AND
		p.nextToken() // move to next element
		elem := p.parsePrimary()
		list.Elements = append(list.Elements, elem)
	}

	return list
}

// parseNumberWord parses English number words like "forty two"
func (p *Parser) parseNumberWord() *ast.IntegerLiteral {
	startToken := p.curToken
	value := p.parseCompoundNumber()
	return &ast.IntegerLiteral{Token: startToken, Value: value}
}

// parseCompoundNumber handles compound numbers like "forty two", "one hundred twenty three"
func (p *Parser) parseCompoundNumber() int64 {
	var total int64 = 0
	var current int64 = 0

	for token.IsNumberWord(p.curToken.Type) {
		wordValue := token.NumberWordValue(p.curToken.Type)

		if token.IsMultiplier(p.curToken.Type) {
			if current == 0 {
				current = 1
			}
			if p.curToken.Type == token.MILLION {
				total += current * wordValue
				current = 0
			} else if p.curToken.Type == token.THOUSAND {
				total += current * wordValue
				current = 0
			} else if p.curToken.Type == token.HUNDRED {
				current *= wordValue
			}
		} else {
			current += wordValue
		}

		if !token.IsNumberWord(p.peekToken.Type) {
			break
		}
		p.nextToken()
	}

	return total + current
}

// === HTTP Parser Functions ===

// parseFetchStatement parses: fetch from "URL" into response
// or: fetch from "URL" with headers into response
func (p *Parser) parseFetchStatement() *ast.FetchStatement {
	stmt := &ast.FetchStatement{Token: p.curToken}

	if !p.expectPeek(token.FROM) {
		return nil
	}

	p.nextToken()
	stmt.URL = p.parseExpression()

	// Check for optional "with headers"
	if p.peekTokenIs(token.WITH) {
		p.nextToken() // consume WITH
		p.nextToken() // move to headers expression
		stmt.Headers = p.parseExpression()
	}

	if !p.expectPeek(token.INTO) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Target = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

// parseSendStatement parses: send "body" to "URL" into response
func (p *Parser) parseSendStatement() *ast.SendStatement {
	stmt := &ast.SendStatement{Token: p.curToken}

	p.nextToken()
	stmt.Body = p.parseExpression()

	if !p.expectPeek(token.TO) {
		return nil
	}

	p.nextToken()
	stmt.URL = p.parseExpression()

	// Check for optional "with headers"
	if p.peekTokenIs(token.WITH) {
		p.nextToken() // consume WITH
		p.nextToken() // move to headers expression
		stmt.Headers = p.parseExpression()
	}

	if !p.expectPeek(token.INTO) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Target = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

// parsePutStatement parses: put "body" to "URL" into response
func (p *Parser) parsePutStatement() *ast.PutStatement {
	stmt := &ast.PutStatement{Token: p.curToken}

	p.nextToken()
	stmt.Body = p.parseExpression()

	if !p.expectPeek(token.TO) {
		return nil
	}

	p.nextToken()
	stmt.URL = p.parseExpression()

	// Check for optional "with headers"
	if p.peekTokenIs(token.WITH) {
		p.nextToken() // consume WITH
		p.nextToken() // move to headers expression
		stmt.Headers = p.parseExpression()
	}

	if !p.expectPeek(token.INTO) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Target = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

// parseDeleteStatement parses: delete from "URL" into response
func (p *Parser) parseDeleteStatement() *ast.DeleteStatement {
	stmt := &ast.DeleteStatement{Token: p.curToken}

	if !p.expectPeek(token.FROM) {
		return nil
	}

	p.nextToken()
	stmt.URL = p.parseExpression()

	// Check for optional "with headers"
	if p.peekTokenIs(token.WITH) {
		p.nextToken() // consume WITH
		p.nextToken() // move to headers expression
		stmt.Headers = p.parseExpression()
	}

	if !p.expectPeek(token.INTO) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Target = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

// parseBodyOfExpression parses: body of response
func (p *Parser) parseBodyOfExpression() *ast.BodyOfExpression {
	expr := &ast.BodyOfExpression{Token: p.curToken}

	p.nextToken() // consume BODY, now at OF
	p.nextToken() // consume OF, now at response expression

	expr.Response = p.parsePrimary()

	return expr
}

// parseStatusOfExpression parses: status of response
func (p *Parser) parseStatusOfExpression() *ast.StatusOfExpression {
	expr := &ast.StatusOfExpression{Token: p.curToken}

	p.nextToken() // consume STATUS, now at OF
	p.nextToken() // consume OF, now at response expression

	expr.Response = p.parsePrimary()

	return expr
}

// parseHeaderFromExpression parses: header "Name" from response
func (p *Parser) parseHeaderFromExpression() *ast.HeaderFromExpression {
	expr := &ast.HeaderFromExpression{Token: p.curToken}

	p.nextToken() // move past HEADER to header name
	expr.HeaderName = p.parsePrimary()

	if !p.expectPeek(token.FROM) {
		return nil
	}

	p.nextToken()
	expr.Response = p.parsePrimary()

	return expr
}

// === JSON Parser Functions ===

// parseParseJsonStatement parses: parse X as json into Y
func (p *Parser) parseParseJsonStatement() *ast.ParseJsonStatement {
	stmt := &ast.ParseJsonStatement{Token: p.curToken}

	p.nextToken()
	stmt.Source = p.parseExpression()

	if !p.expectPeek(token.AS) {
		return nil
	}

	if !p.expectPeek(token.JSON) {
		return nil
	}

	if !p.expectPeek(token.INTO) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Target = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

// parseFieldFromExpression parses: field "name" from data
func (p *Parser) parseFieldFromExpression() *ast.FieldFromExpression {
	expr := &ast.FieldFromExpression{Token: p.curToken}

	p.nextToken() // move past FIELD to field name
	expr.FieldName = p.parsePrimary()

	if !p.expectPeek(token.FROM) {
		return nil
	}

	p.nextToken()
	expr.Source = p.parsePrimary()

	return expr
}

// parseEncodeJsonStatement parses: encode X as json into Y
func (p *Parser) parseEncodeJsonStatement() *ast.EncodeJsonStatement {
	stmt := &ast.EncodeJsonStatement{Token: p.curToken}

	p.nextToken()
	stmt.Source = p.parseExpression()

	if !p.expectPeek(token.AS) {
		return nil
	}

	if !p.expectPeek(token.JSON) {
		return nil
	}

	if !p.expectPeek(token.INTO) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Target = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

// === Web Server Parser Functions ===

// parseServeStatement parses: serve on 8080 or serve on 8080 in background
func (p *Parser) parseServeStatement() *ast.ServeStatement {
	stmt := &ast.ServeStatement{Token: p.curToken}

	if !p.expectPeek(token.ON) {
		return nil
	}

	p.nextToken()
	stmt.Port = p.parseExpression()

	// Check for optional "in background"
	if p.peekTokenIs(token.IN) {
		p.nextToken() // consume IN
		if p.peekTokenIs(token.BACKGROUND) {
			p.nextToken() // consume BACKGROUND
			stmt.Background = true
		}
	}

	return stmt
}

// parseWhenRouteStatement parses:
// - when request at "/path" do ... done
// - when request at "/path" using req do ... done
// - when get at "/path" using req do ... done
func (p *Parser) parseWhenRouteStatement() *ast.WhenRouteStatement {
	stmt := &ast.WhenRouteStatement{Token: p.curToken}

	p.nextToken() // move past WHEN

	// Check for method-specific route or generic "request"
	if p.curTokenIs(token.REQUEST) {
		stmt.Method = ""
		p.nextToken() // consume REQUEST
	} else if p.curTokenIs(token.GET) {
		stmt.Method = "GET"
		p.nextToken()
	} else if p.curTokenIs(token.SEND) {
		// "send" is used for POST in existing HTTP client syntax
		// For routes, we'll treat it as POST
		stmt.Method = "POST"
		p.nextToken()
	} else if p.curTokenIs(token.PUT) {
		stmt.Method = "PUT"
		p.nextToken()
	} else if p.curTokenIs(token.DELETE) {
		stmt.Method = "DELETE"
		p.nextToken()
	} else if p.curTokenIs(token.FETCH) {
		// "fetch" can be used as alias for GET in routes
		stmt.Method = "GET"
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("line %d: expected request or HTTP method after 'when', got %s", p.curToken.Line, p.curToken.Type))
		return nil
	}

	// Expect AT
	if !p.curTokenIs(token.AT) {
		p.errors = append(p.errors, fmt.Sprintf("line %d: expected 'at', got %s", p.curToken.Line, p.curToken.Type))
		return nil
	}

	p.nextToken() // move past AT
	stmt.Path = p.parseExpression()

	// Check for optional "using reqVar"
	if p.peekTokenIs(token.USING) {
		p.nextToken() // consume USING
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		stmt.RequestVar = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.DO) {
		return nil
	}

	p.nextToken() // move past DO
	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseRouteToStatement parses: route "/path" to handlerFunc
func (p *Parser) parseRouteToStatement() *ast.RouteToStatement {
	stmt := &ast.RouteToStatement{Token: p.curToken}

	p.nextToken()
	stmt.Path = p.parseExpression()

	if !p.expectPeek(token.TO) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Handler = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

// parseReplyStatement parses: reply with "data" or reply with "data" as json with status 201
func (p *Parser) parseReplyStatement() *ast.ReplyStatement {
	stmt := &ast.ReplyStatement{Token: p.curToken}

	if !p.expectPeek(token.WITH) {
		return nil
	}

	p.nextToken()
	stmt.Body = p.parseExpression()

	// Check for "as json" modifier
	if p.peekTokenIs(token.AS) {
		p.nextToken() // consume AS
		if p.peekTokenIs(token.JSON) {
			p.nextToken() // consume JSON
			stmt.AsJson = true
		}
	}

	// Parse optional modifiers: with status N, with header "X" as "Y"
	for p.peekTokenIs(token.WITH) {
		p.nextToken() // consume WITH
		p.nextToken() // move to modifier type

		if p.curTokenIs(token.STATUS) {
			p.nextToken()
			stmt.StatusCode = p.parseExpression()
		} else if p.curTokenIs(token.HEADER) {
			p.nextToken()
			headerName := p.parseExpression()
			if !p.expectPeek(token.AS) {
				return nil
			}
			p.nextToken()
			headerValue := p.parseExpression()
			stmt.Headers = append(stmt.Headers, ast.HeaderPair{Name: headerName, Value: headerValue})
		}
	}

	return stmt
}

// parseStopServerStatement parses: stop server or stop server on 8080
func (p *Parser) parseStopServerStatement() *ast.StopServerStatement {
	stmt := &ast.StopServerStatement{Token: p.curToken}

	if !p.expectPeek(token.SERVER) {
		return nil
	}

	// Check for optional "on port"
	if p.peekTokenIs(token.ON) {
		p.nextToken() // consume ON
		p.nextToken()
		stmt.Port = p.parseExpression()
	}

	return stmt
}

// parseMethodOfExpression parses: method of req
func (p *Parser) parseMethodOfExpression() *ast.MethodOfExpression {
	expr := &ast.MethodOfExpression{Token: p.curToken}

	p.nextToken() // consume METHOD, now at OF
	p.nextToken() // consume OF, now at request expression

	expr.Request = p.parsePrimary()

	return expr
}

// parsePathOfExpression parses: path of req
func (p *Parser) parsePathOfExpression() *ast.PathOfExpression {
	expr := &ast.PathOfExpression{Token: p.curToken}

	p.nextToken() // consume PATH, now at OF
	p.nextToken() // consume OF, now at request expression

	expr.Request = p.parsePrimary()

	return expr
}

// parseQueryFromExpression parses: query "name" from req
func (p *Parser) parseQueryFromExpression() *ast.QueryFromExpression {
	expr := &ast.QueryFromExpression{Token: p.curToken}

	p.nextToken() // move past QUERY to query name
	expr.QueryName = p.parsePrimary()

	if !p.expectPeek(token.FROM) {
		return nil
	}

	p.nextToken()
	expr.Request = p.parsePrimary()

	return expr
}
