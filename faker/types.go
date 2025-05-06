package faker

type CustomType struct {
	Type   string            `mapstructure:"type"`    // enum, float, datetime, etc.
	Params map[string]string `mapstructure:",remain"` // all additional fields like min, max, values
}

// Placeholder represents a single interpolation tag in a template.
type Placeholder struct {
	FieldName string
	Type      string
	Params    map[string]string
}

// FieldGenerator produces a value for a given placeholder.
type FieldGenerator interface {
	Generate() (any, error)
}
