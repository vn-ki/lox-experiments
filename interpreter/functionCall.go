package interpreter

import "time"

type LoxCallable interface {
	Call(*Interpreter, []interface{}) interface{}
	Arity() int
}

/// Native Function: clock
type FnClock struct{}

func (f FnClock) Arity() int { return 0 }

func (f FnClock) Call(_ *Interpreter, _ []interface{}) interface{} {
	return float64(time.Now().UnixNano())
}

func (f FnClock) String() string { return "<clock native fn>" }
