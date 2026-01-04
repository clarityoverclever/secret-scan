package validators

type Validator interface {
	Validate(match string, context string) bool
}

type ValidatorFunc func(match string, context string) bool

func (f ValidatorFunc) Validate(match string, context string) bool {
	return f(match, context)
}

type Registry struct {
	validators map[string]Validator
}

func NewRegistry() *Registry {
	r := &Registry{validators: make(map[string]Validator)}

	//r.RegisterDefaults()
	return r
}

func (r *Registry) RegisterDefaults() {
	// todo r.Register("azure_context", AzureContextValidator)
}

func (r *Registry) Register(name string, validator Validator) {
	r.validators[name] = validator
}

func (r *Registry) Get(name string) Validator {
	return r.validators[name]
}
