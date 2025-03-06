package parser

import (
	"fmt"
	"strconv"
	"yap/ast"
	"yap/lexer"
	"yap/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedence = map[token.TokenType]int{
	token.EQ:        EQUALS,
	token.NEQ:       EQUALS,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.ASTERISK:  PRODUCT,
	token.DASH:      PRODUCT,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.TERNARY:   LESSGREATER,
	token.LPAREN:    CALL,
	token.POWER:     PRODUCT,
	token.MOD:       PRODUCT,
	token.DECREMENT: SUM,
	token.INCREMENT: SUM,
	token.LBRACKET:  INDEX,
}

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns  map[token.TokenType]prefixParseFn
	infixPerseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression) ast.Expression
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.DECREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.INCREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	p.infixPerseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.DASH, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.POWER, p.parseInfixExpression)
	p.registerInfix(token.TERNARY, p.parseTernaryExpression)
	p.registerInfix(token.COLON, p.parseTernaryExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)
	p.registerPostfix(token.INCREMENT, p.parsePostfixExpression)
	p.registerPostfix(token.DECREMENT, p.parsePostfixExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

// This create the error of mismatch token and log it into p.errors
func (p *Parser) peekError(token token.TokenType) {
	msg := fmt.Sprintf("expected next token to be '%s', got '%s' instead", token, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParserProgram() *ast.Program {
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
	switch {
	case p.curToken.Type == token.LET:
		return p.parseLetStatement()
	case p.curToken.Type == token.GLOBAL:
		return p.parseGlobalStatement()
	case p.curToken.Type == token.CONST:
		return p.parseConstStatement()
	case p.curToken.Type == token.RETURN:
		return p.parseReturnStatement()
	case p.curToken.Type == token.IDENT && p.peekTokenIs(token.ASSIGN):
		return p.parseIdentStatement()
	case p.curToken.Type == token.FOR:
		return p.parseForLiteral()
	default:
		return p.parseExpressionstatement()
	}
}

func (p *Parser) parseLetStatement() *ast.SayStatement {
	stmt := &ast.SayStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	// The current token.Type is actually IDENT here since the expectPeek moved the token up one

	//Make it so that you can create a null variable for now and use it later
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		stmt.Value = nil
		return stmt
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseGlobalStatement() *ast.GlobalStatement {
	stmt := &ast.GlobalStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if !p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConstStatement() *ast.ConstStaement {
	stmt := &ast.ConstStaement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if !p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) peekTokenIs(token token.TokenType) bool {
	return p.peekToken.Type == token
}

func (p *Parser) curTokenIs(token token.TokenType) bool {
	return p.curToken.Type == token
}

func (p *Parser) expectPeek(token token.TokenType) bool {
	if p.peekToken.Type == token {
		p.nextToken()
		return true
	} else {
		p.peekError(token)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixPerseFns[tokenType] = fn
}

func (p *Parser) registerPostfix(tokenType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}

func (p *Parser) parseExpressionstatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	//Skipping the semicolon token for each statement
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	if p.peekTokenIs(token.INCREMENT) || p.peekTokenIs(token.DECREMENT) {
		p.nextToken()
		postfix := p.postfixParseFns[p.curToken.Type]
		leftExp = postfix(leftExp)
	}

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixPerseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as a float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as an integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedence[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.PostfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}
	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	val := &ast.Boolean{Token: p.curToken}
	boolean, ok := strconv.ParseBool(string(p.curToken.Type))
	if ok != nil {
		msg := fmt.Sprintf("Could not parse %s as boolean value", p.curToken.Type)
		p.errors = append(p.errors, msg)
	}
	val.Value = boolean
	return val
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELIF) {
		p.nextToken()
		if !p.expectPeek(token.LPAREN) {
			return nil
		}
		expression.Elif = p.parseElifExpression()
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseElifExpression() []ast.ElifEpxression {
	var expression []ast.ElifEpxression
	position := 0
	for !p.peekTokenIs(token.ELSE) {
		p.nextToken()
		expression = append(expression, ast.ElifEpxression{Token: p.curToken})
		expression[position].Conditions = p.parseExpression(LOWEST)
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression[position].Consequences = p.parseBlockStatement()
		position++
		if p.peekTokenIs(token.ELIF) {
			p.nextToken()
			p.nextToken()
		} else {
			break
		}
	}
	return expression
}

func (p *Parser) parseTernaryExpression(condition ast.Expression) ast.Expression {
	expression := &ast.TernaryExpression{
		Token:     p.curToken,
		Condition: condition,
	}
	p.parseTernaryBlockStatement(expression)
	return expression
}

func (p *Parser) parseTernaryBlockStatement(ternary *ast.TernaryExpression) {
	consequence := &ast.BlockStatement{Token: p.curToken}
	consequence.Statements = []ast.Statement{}
	for !p.curTokenIs(token.COLON) {
		p.nextToken()
		stmt := p.parseStatement()
		if stmt != nil {
			consequence.Statements = append(consequence.Statements, stmt)
		}
		p.nextToken()
	}
	ternary.Consequence = consequence
	alternative := &ast.BlockStatement{Token: p.curToken}
	alternative.Statements = []ast.Statement{}
	for !p.curTokenIs(token.EOF) {
		p.nextToken()
		stmt := p.parseStatement()
		if stmt != nil {
			alternative.Statements = append(alternative.Statements, stmt)
		}
		p.nextToken()
	}
	ternary.Alternative = alternative
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionExpression{Token: p.curToken}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		literal.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	literal.Body = p.parseBlockStatement()
	return literal
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifier := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifier
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifier = append(identifier, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		identifier = append(identifier, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifier
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	exp.Arguments = p.parseCallArgument()

	return exp
}

func (p *Parser) parseCallArgument() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseIdentStatement() *ast.PotentialStatement {
	stmt := &ast.PotentialStatement{Token: token.Token{Type: token.LET, Literal: "propose"}}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Literal: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	if array.Elements == nil {
		return nil
	}
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		curExp := p.parseExpression(LOWEST)
		list = append(list, curExp)
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseForLiteral() *ast.ForExpression {
	forStat := &ast.ForExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		msg := fmt.Sprintf("Expecting \"(\", got %s", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	for !p.curTokenIs(token.RPAREN) {
		if p.peekTokenIs(token.LET) {
			p.nextToken()
			forStat.Identifier = p.parseLetStatement()
		} else {
			if p.curTokenIs(token.LPAREN) {
				p.nextToken()
			}
			forStat.Conditions = append(forStat.Conditions, p.parseExpression(LOWEST))
			p.nextToken()
		}
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		} else if !p.curTokenIs(token.RPAREN) {
			msg := fmt.Sprintf("Expected \")\" got %s", p.curToken.Literal)
			p.errors = append(p.errors, msg)
		}
	}
	p.nextToken()

	if !p.curTokenIs(token.LBRACE) {
		msg := fmt.Sprintf("Expected \"{\", got %s", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	forStat.Statements = p.parseBlockStatement()

	return forStat
}
