package shell

type Environment struct {
	store map[string]string
	outer *Environment
}

func (e *Environment) Get(name string) (string, bool) {
	if e.store != nil {
		s, ok := e.store[name]
		if ok {
			return s, ok
		}
	}
	if e.outer != nil {
		return e.outer.Get(name)
	}
	return "", false
}

func (e *Environment) Set(name, val string) {
	if e.store == nil {
		e.store = map[string]string{}
	}
	e.store[name] = val
}
