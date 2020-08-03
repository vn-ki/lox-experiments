package env

type Environemnt struct {
	values map[string]interface{}
}

func NewEnvironment() *Environemnt {
	return &Environemnt{make(map[string]interface{})}
}

func (e *Environemnt) Define(key string, val interface{}) {
	e.values[key] = val
}

func (e *Environemnt) Get(key string) (interface{}, bool) {
	val, ok := e.values[key]
	return val, ok
}
