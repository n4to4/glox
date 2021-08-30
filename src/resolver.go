package main

type scope map[string]bool

type Resolver struct {
	interpreter *Interpreter
	scopes      []scope
}

func NewResolver(interpreter *Interpreter) Resolver {
	var scopes []scope
	return Resolver{
		interpreter,
		scopes,
	}
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

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) beginScope() {
	s := make(scope, 1)
	r.scopes = append(r.scopes, s)
}

func (r *Resolver) endScope() {
	r.scopes[len(r.scopes)-1] = nil
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) VisitVarStmt(stmt Var) (interface{}, error) {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolveExpr(*stmt.initializer)
	}
	r.define(stmt.name)
	return nil, nil
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.lexeme] = true
}

func (r *Resolver) VisitVariableExpr(expr Variable) (interface{}, error) {
	if len(r.scopes) != 0 {
		scope := r.scopes[len(r.scopes)-1]
		if !scope[expr.name.lexeme] {
			panic("can't read local variable in its own initializer")
		}
	}
	r.resolveLocal(expr, expr.name)
	return nil, nil
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) VisitAssignExpr(expr Assign) (interface{}, error) {
	r.resolveExpr(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt Function) (interface{}, error) {
	r.declare(stmt.name)
	r.define(stmt.name)
	r.resolveFunction(stmt)
	return nil, nil
}

func (r *Resolver) resolveFunction(function Function) {
	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmts(function.body)
	r.endScope()
}

func (r *Resolver) VisitExpressionStmt(stmt Expression) (interface{}, error) {
	r.resolveExpr(stmt.expression)
	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt If) (interface{}, error) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(*stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStmt(*stmt.elseBranch)
	}
	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt Print) (interface{}, error) {
	r.resolveExpr(stmt.expression)
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt Return) (interface{}, error) {
	if stmt.value != nil {
		r.resolveExpr(*stmt.value)
	}
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt While) (interface{}, error) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr Binary) (interface{}, error) {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr Call) (interface{}, error) {
	r.resolveExpr(expr.callee)
	for _, argument := range expr.arguments {
		r.resolveExpr(argument)
	}
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr Grouping) (interface{}, error) {
	r.resolveExpr(expr.expression)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr Literal) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr Logical) (interface{}, error) {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr Unary) (interface{}, error) {
	r.resolveExpr(expr.right)
	return nil, nil
}
