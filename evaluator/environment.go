package evaluator

type Env struct {
	store map[string]interface{}
	consts map[string]bool
	outer *Env
}

func NewEnv() *Env {
	return &Env{store: make(map[string]interface{}), consts: make(map[string]bool)}
}

func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}

func (e *Env) Get(name string) (interface{}, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return val, ok
}

func (e *Env) Set(name string, val interface{}, isConst bool) {
	e.store[name] = val
	if isConst { e.consts[name] = true }
}

func (e *Env) Update(name string, val interface{}) bool {
	if _, ok := e.store[name]; ok {
		if _, isConst := e.consts[name]; isConst {
			return false
		}
		e.store[name] = val
		return true
	}
	if e.outer != nil {
		return e.outer.Update(name, val)
	}
	return false
}
