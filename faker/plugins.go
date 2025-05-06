package faker

import (
	"fmt"
)

// GeneratorFactory builds a FieldGenerator from a field name and params.
type GeneratorFactory func(fieldName string, params map[string]string) (FieldGenerator, error)

var generatorRegistry = make(map[string]GeneratorFactory)

// RegisterGenerator installs a new type generator.
func RegisterGenerator(typeName string, factory GeneratorFactory) {
	if _, exists := generatorRegistry[typeName]; exists {
		panic("duplicate generator type registered: " + typeName)
	}
	generatorRegistry[typeName] = factory
}

// GetGenerator returns a generator by type name.
func GetGenerator(typeName, fieldName string, params map[string]string) (FieldGenerator, error) {
	factory, ok := generatorRegistry[typeName]
	if !ok {
		return nil, fmt.Errorf("unknown generator type: %s", typeName)
	}
	return factory(fieldName, params)
}
