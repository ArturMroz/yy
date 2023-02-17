package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: map[string]Object{}}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{
		store: map[string]Object{},
		outer: outer,
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Set(name string, val Object) {
	e.store[name] = val
}

// Update tries to update a value for a given name. It first checks if the given name exsist. If it
// does, it updates the value and returns true. If it doesn't, it returns false.
func (e *Environment) Update(name string, val Object) bool {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return true
	}

	if e.outer != nil {
		return e.outer.Update(name, val)
	}

	return false
}

func (e *Environment) YoloMode() bool {
	_, ok := e.Get(YoloKey)
	return ok
}
