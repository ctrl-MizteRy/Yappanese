package object

func NewEncloseEnviroment(outer *Enviroment) *Enviroment {
	env := NewEnviroment()
	env.outer = outer
	return env
}
func NewEnviroment() *Enviroment {
	s := make(map[string]Object)
	return &Enviroment{store: s, outer: nil}
}

type Enviroment struct {
	store map[string]Object
	outer *Enviroment
}

func (e *Enviroment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Enviroment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Enviroment) Exist(name string) bool {
	_, ok := e.store[name]
	return ok
}

func (e *Enviroment) TypeComp(name string, valType ObjectType) bool {
	val := e.store[name]
	if val == nil {
		return true
	}
	ident := val.Type()
	return ident == valType
}

func (e *Enviroment) GetType(name string) Object {
	return e.store[name]
}
