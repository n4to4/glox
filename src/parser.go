package main

import (
	"errors"
	"log"
)

var (
	ErrParse = errors.New("parse error")
)

type Parser struct {
	tokens  []Token
	current int
}

func (p *Parser) Parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() Stmt {
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name, _ := p.consume(IDENTIFIER, "expect variable name.")

	var initializer *Expr
	if p.match(EQUAL) {
		init := p.Expression()
		initializer = &init
	}

	p.consume(SEMICOLON, "expect ';' aftter variable declaration.")
	return Var{name, initializer}
}

func (p *Parser) statement() Stmt {
	if p.match(IF) {
		return p.IfStatement()
	}

	if p.match(PRINT) {
		return p.PrintStatement()
	}
	if p.match(LEFT_BRACE) {
		return Block{p.block()}
	}

	return p.ExpressionStatement()
}

func (p *Parser) PrintStatement() Stmt {
	value := p.Expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return Print{value}
}

func (p *Parser) ExpressionStatement() Stmt {
	expr := p.Expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return Expression{expr}
}

func (p *Parser) IfStatement() Stmt {
	p.consume(LEFT_PAREN, "expect '(' after 'if'")
	condition := p.Expression()
	p.consume(RIGHT_PAREN, "expect ')' after if condition")

	var thenBranch Stmt = p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return If{condition, &thenBranch, &elseBranch}
}

func (p *Parser) block() []Stmt {
	statements := []Stmt{}

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(RIGHT_BRACE, "expect '}' after block")
	return statements
}

//
// Expr
//

func (p *Parser) Expression() Expr {
	return p.Assignment()
}

func (p *Parser) Assignment() Expr {
	expr := p.Equality()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.Assignment()

		if v, ok := expr.(Variable); ok {
			name := v.name
			return Assign{name, value}
		}

		log.Printf("Invalid assignment target %v", equals)
	}

	return expr
}

func (p *Parser) Equality() Expr {
	expr := p.Comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.Comparison()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) Comparison() Expr {
	expr := p.Term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.Term()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) Term() Expr {
	expr := p.Factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.Factor()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) Factor() Expr {
	expr := p.Unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.Unary()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) Unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.Unary()
		return Unary{operator, right}
	}

	return p.Primary()
}

func (p *Parser) Primary() Expr {
	switch {
	case p.match(FALSE):
		return Literal{false}
	case p.match(TRUE):
		return Literal{true}
	case p.match(NIL):
		return Literal{nil}
	case p.match(NUMBER, STRING):
		return Literal{p.previous().Literal}
	case p.match(IDENTIFIER):
		return Variable{p.previous()}
	case p.match(LEFT_PAREN):
		expr := p.Expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{expr}
	}

	panic("Expect expression.")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, ttype := range types {
		if p.check(ttype) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(ttype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Ttype == ttype
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Ttype == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(ttype TokenType, message string) (Token, error) {
	if p.check(ttype) {
		return p.advance(), nil
	}

	return Token{}, p.parseError(p.peek(), message)
}

func (p *Parser) parseError(token Token, message string) error {
	ReportError(token, message)
	return ErrParse
}

//func (p *Parser) synchronize() {
//	p.advance()
//
//	for !p.isAtEnd() {
//		if p.previous().Ttype == SEMICOLON {
//			return
//		}
//
//		switch p.peek().Ttype {
//		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
//			return
//		}
//
//		p.advance()
//	}
//}
