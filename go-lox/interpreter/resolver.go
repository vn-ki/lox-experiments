package interpreter

//
// import (
// 	"fmt"
//
// 	"github.com/vn-ki/go-lox/ast"
// )
//
// type Scope map[string]bool
//
// func newScope() Scope { return make(Scope) }
//
// type Stack struct {
// 	stack []Scope
// }
//
// func (s *Stack) Push(scope Scope) {
// 	s.stack = append(s.stack, scope)
// }
//
// func (s *Stack) Pop() (Scope, bool) {
// 	ret := s.Head()
// 	if ret == nil {
// 		return nil, false
// 	}
// 	last := len(s.stack) - 1
// 	s.stack = s.stack[last:]
// 	return ret, true
// }
//
// func (s *Stack) Head() Scope {
// 	last := len(s.stack) - 1
// 	if last < 0 {
// 		return nil
// 	}
// 	return s.stack[last]
// }
//
// // Resolver
//
// type Resolver struct {
// 	scopes Stack
// 	i      *Interpreter
// }
//
// func (r *Resolver) VisitBlock(s ast.Sblock) interface{} {
// 	r.beginScope()
//
// 	for _, stmt := range s.Stmts {
// 		r.resolveStmt(stmt)
// 	}
//
// 	r.endScope()
// 	return nil
// }
//
// func (r *Resolver) resolveStmt(s ast.Stmt) {
// 	s.Accept(r)
// }
//
// func (r *Resolver) resolveExpr(e ast.Expr) {
// 	e.Accept(r)
// }
//
// func (r *Resolver) beginScope() {
// 	r.scopes.Push(newScope())
// }
//
// func (r *Resolver) endScope() {
// 	r.scopes.Pop()
// }
//
// func (r *Resolver) VisitVar(s ast.Svar) interface{} {
// 	r.declare(s.Name)
//
// 	if s.Expression != nil {
// 		r.resolveExpr(s.Expression)
// 	}
//
// 	r.define(s.Name)
// 	return nil
// }
