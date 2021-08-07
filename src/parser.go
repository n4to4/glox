package main

type Parser struct {
	tokens  []Token
	current int
}

func (p *Parser) Expression() Expr {
	return p.Equality()
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
	case p.match(LEFT_PAREN):
		expr := p.Expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{expr}
	}

	return nil
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

func (p *Parser) consume(ttype TokenType, message string) Token {
	if p.check(ttype) {
		return p.advance()
	}

	panic("parse error")
}
