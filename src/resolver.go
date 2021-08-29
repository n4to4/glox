package main

type scope map[string]bool

type Resolver struct {
	interpreter Interpreter
	scopes      []scope
}

func NewResolver(i Interpreter) Resolver {
	var scopes []scope
	resolver := Resolver{
		i, scopes,
	}
	return resolver
}

func (r *Resolver) VisitBlockStmt(stmt Block) (interface{}, error) {
	r.beginScope()
	r.resolveStmts(stmt.statements)
	r.endScope()
	return nil, nil
}

func (r *Resolver) resolveStmts(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) beginScope() {
	s := make(scope, 1)
	r.scopes = append(r.scopes, s)
}

func (r *Resolver) endScope() {
	r.scopes[len(r.scopes)-1] = nil
	r.scopes = r.scopes[:len(r.scopes)-1]
}

/*

type ExprVisitor interface {
	VisitAssignExpr(expr Assign) (interface{}, error)
	VisitBinaryExpr(expr Binary) (interface{}, error)
	VisitCallExpr(expr Call) (interface{}, error)
	VisitGroupingExpr(expr Grouping) (interface{}, error)
	VisitLiteralExpr(expr Literal) (interface{}, error)
	VisitLogicalExpr(expr Logical) (interface{}, error)
	VisitUnaryExpr(expr Unary) (interface{}, error)
	VisitVariableExpr(expr Variable) (interface{}, error)
}

type StmtVisitor interface {
	VisitBlockStmt(stmt Block) (interface{}, error)
	VisitExpressionStmt(stmt Expression) (interface{}, error)
	VisitFunctionStmt(stmt Function) (interface{}, error)
	VisitIfStmt(stmt If) (interface{}, error)
	VisitPrintStmt(stmt Print) (interface{}, error)
	VisitReturnStmt(stmt Return) (interface{}, error)
	VisitVarStmt(stmt Var) (interface{}, error)
	VisitWhileStmt(stmt While) (interface{}, error)
}

*/
