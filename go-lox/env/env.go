package env

import (
	"log"
	"strings"
)

type Environemnt struct {
	values    map[string]interface{}
	Enclosing *Environemnt
}

func NewEnvironment(enclosing *Environemnt) *Environemnt {
	return &Environemnt{make(map[string]interface{}), enclosing}
}

func (e *Environemnt) Define(key string, val interface{}) {
	e.values[key] = val
}

func (e *Environemnt) Get(key string) (interface{}, bool) {
	val, ok := e.values[key]
	if !ok && e.Enclosing != nil {
		return e.Enclosing.Get(key)
	}
	return val, ok
}

func (e *Environemnt) Assign(key string, value interface{}) bool {
	_, ok := e.values[key]
	if ok {
		e.Define(key, value)
	} else if e.Enclosing != nil {
		ok = e.Enclosing.Assign(key, value)
	}

	return ok
}

func (e *Environemnt) DumpEnv() {
	e.dumpEnv(0)
}

func (e *Environemnt) dumpEnv(depth int) {
	log.Printf(strings.Repeat(">", depth)+"env: %v\n", e.values)
	if e.Enclosing != nil {
		e.Enclosing.dumpEnv(depth + 1)
	}
}
