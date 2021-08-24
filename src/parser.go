package main

import (
	"errors"
	"fmt"
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
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name, _ := p.consume(IDENTIFIER, "expect variable name.")

	var initializer *Expr
	if p.match(EQUAL) {
		init := p.expression()
		initializer = &init
	}

	p.consume(SEMICOLON, "expect ';' aftter variable declaration.")
	return Var{name, initializer}
}

func (p *Parser) function(kind string) Function {
	name, _ := p.consume(IDENTIFIER, fmt.Sprintf("expect %s name", kind))
	p.consume(LEFT_PAREN, fmt.Sprintf("expect '(' after %s name", kind))

	var parameters []Token
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				// todo: error
				panic("can't have more than 255 parameters")
			}

			param, _ := p.consume(IDENTIFIER, "expect parameter name")
			parameters = append(parameters, param)

			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "expect ')' after parameters")

	p.consume(LEFT_BRACE, fmt.Sprintf("expect '{' before %s body", kind))
	body := p.block()

	return Function{name, parameters, body}
}

func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	}

	if p.match(IF) {
		return p.ifStatement()
	}

	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(RETURN) {
		return p.returnStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	if p.match(LEFT_BRACE) {
		return Block{p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return Print{value}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value *Expr = nil
	if !p.check(SEMICOLON) {
		v := p.expression()
		value = &v
	}

	p.consume(SEMICOLON, "expect ';' after return value")
	return Return{keyword, value}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return Expression{expr}
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "expect '(' after 'if'")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "expect ')' after if condition")

	var thenBranch Stmt = p.statement()
	var elseBranch *Stmt = nil
	if p.match(ELSE) {
		els := p.statement()
		elseBranch = &els
	}

	return If{condition, &thenBranch, elseBranch}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "expect '(' after while")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "expect ')' after condition")
	body := p.statement()

	return While{condition, body}
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "expect '(' after 'for'")

	var initializer *Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		stmt := p.varDeclaration()
		initializer = &stmt
	} else {
		exp := p.expressionStatement()
		initializer = &exp
	}

	var condition *Expr = nil
	if !p.check(SEMICOLON) {
		cond := p.expression()
		condition = &cond
	}
	p.consume(SEMICOLON, "expect ';' after loop condition")

	var increment *Expr = nil
	if !p.check(RIGHT_PAREN) {
		exp := p.expression()
		increment = &exp
	}
	p.consume(RIGHT_PAREN, "expect ')' after ")

	body := p.statement()

	if increment != nil {
		body = Block{[]Stmt{body, Expression{*increment}}}
	}

	if condition == nil {
		var lit Expr = Literal{true}
		condition = &lit
	}
	body = While{*condition, body}

	if initializer != nil {
		body = Block{[]Stmt{*initializer, body}}
	}

	return body
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

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if v, ok := expr.(Variable); ok {
			name := v.name
			return Assign{name, value}
		}

		log.Printf("Invalid assignment target %v", equals)
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = Logical{expr, operator, right}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = Logical{expr, operator, right}
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return Unary{operator, right}
	}

	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := []Expr{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				// todo: no panic
				panic("Can't have more than 255 arguments")
			}

			arguments = append(arguments, p.expression())

			if !p.match(COMMA) {
				break
			}
		}
	}

	paren, _ := p.consume(RIGHT_PAREN, "expect ')' after arguments")

	return Call{callee, paren, arguments}
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(FALSE):
		return Literal{false}
	case p.match(TRUE):
		return Literal{true}
	case p.match(NIL):
		return Literal{nil}
	case p.match(NUMBER, STRING):
		return Literal{p.previous().literal}
	case p.match(IDENTIFIER):
		return Variable{p.previous()}
	case p.match(LEFT_PAREN):
		expr := p.expression()
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
	return p.peek().ttype == ttype
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().ttype == EOF
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
