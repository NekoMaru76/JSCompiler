package js_to_go

type Variable struct {
	Value      Object
	IsConstant bool
}

type Context struct {
	Vars          map[string]*Variable
	VarsAliases   map[string]string
	ParentContext *Context
}

func (ctx *Context) HasVariable(name string) bool {
	_, ok := ctx.VarsAliases[name]

	return ok
}

func (ctx *Context) TestDefined(name string) {
	if ctx.HasVariable(name) {
		return
	}

	panic(name + " NOT DEFINED")
}

func (ctx *Context) SetVariable(name string, value Object) {
	ctx.TestDefined(name)

	alias := ctx.VarsAliases[name]

	if ctx.Vars[alias].IsConstant {
		panic(name + " IS CONSTANT")
	}

	ctx.Vars[alias].Value = value
}

func (ctx *Context) GetVariable(name string) *Object {
	ctx.TestDefined(name)
	return ctx.Vars[ctx.VarsAliases[name]].Value
}

func (ctx *Context) DefineVariable(name string, value Object, isConst bool) {
	varAlias, ok := ctx.VarsAliases[name]

	if ok {
		panic(name + " DEFINED ALREADY")
	}

	varAlias = name + "Var"
	ctx.VarsAliases[name] = varAlias
	ctx.Vars[varAlias] = &Variable{
		IsConstant: isConst,
		Value:      value,
	}
}

type Object interface {
	GetType() string
	SetProp(string, Object)
	GetProp(string) Object
	GetProps() map[string]Object
}

type BaseObject struct {
	Object

	Props  map[string]Object
	String string
}

func (v BaseObject) GetProp(name string) Object {
	val, ok := v.Props[name]

	if ok {
		return val
	}

	return &Undefined{}
}

func (BaseObject) GetType() string {
	return "object"
}

type Null struct {
	*BaseObject
}

type Undefined struct {
	*BaseObject
}

func (*Undefined) GetType() string {
	return "undefined"
}

type BaseFunction struct {
	*BaseObject
	Value  func(Context ...Object) Object
	String string
}

type Function struct {
	*BaseFunction
	Name string
}

func (*Function) GetType() string {
	return "function"
}

type ArrowFunction struct {
	*BaseFunction
	Context Context
}

func (*ArrowFunction) GetType() string {
	return "function"
}

type Number struct {
	*BaseObject
	Value float64
}

func (*Number) GetType() string {
	return "number"
}

type String struct {
	*BaseObject
	Value string
}

func (*String) GetType() string {
	return "string"
}
