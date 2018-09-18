package eval

import (
	"io"
	"strconv"
)

const (
	STATUS = "?"
)

type Environment struct {
	store map[string]string
	outer *Environment

	In  io.Reader
	Out io.Writer
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

func (e *Environment) Stdin() io.Reader {
	if e.In != nil {
		return e.In
	}
	if e.outer != nil {
		return e.outer.Stdin()
	}
	return nil
}

func (e *Environment) Stdout() io.Writer {
	if e.Out != nil {
		return e.Out
	}
	if e.outer != nil {
		return e.outer.Stdout()
	}
	return nil
}

func (env *Environment) GetStatus() int {
	s, _ := env.Get(STATUS)
	if s == "" {
		return 0
	}
	status, _ := strconv.Atoi(s)
	return status
}

func (env *Environment) SetStatus(status int) {
	env.Set(STATUS, strconv.Itoa(status))
}
