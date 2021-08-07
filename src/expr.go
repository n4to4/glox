package main

type Expr interface {
	TokenLiteral() string
	Acceptor
}

type Visitor interface {
	VisitBinaryExpr(expr Binary) interface{}
	VisitGroupingExpr(expr Grouping) interface{}
	VisitLiteralExpr(expr Literal) interface{}
	VisitUnaryExpr(expr Unary) interface{}
}

type Acceptor interface {
	Accept(v Visitor) interface{}
}

type Binary struct {
	left Expr
	operator Token
	right Expr
}

func (x Binary) TokenLiteral() string { return "" }

func (x Binary) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(x)
}

type Grouping struct {
	expression Expr
}

func (x Grouping) TokenLiteral() string { return "" }

func (x Grouping) Accept(v Visitor) interface{} {
	return v.VisitGroupingExpr(x)
}

type Literal struct {
	value interface{}
}

func (x Literal) TokenLiteral() string { return "" }

func (x Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(x)
}

type Unary struct {
	operator Token
	right Expr
}

func (x Unary) TokenLiteral() string { return "" }

func (x Unary) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(x)
}
