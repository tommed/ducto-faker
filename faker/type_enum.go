package faker

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

func init() {
	RegisterGenerator("enum", newEnumGenerator)
}

type enumGenerator struct {
	values []string
}

func (g *enumGenerator) Generate() (any, error) {
	if len(g.values) == 0 {
		return nil, errors.New("enum has no values to choose from")
	}
	return fmt.Sprintf("%q", g.values[rand.Intn(len(g.values))]), nil
}

func newEnumGenerator(_ string, params map[string]string) (FieldGenerator, error) {
	raw, ok := params["values"]
	if !ok || strings.TrimSpace(raw) == "" {
		return nil, errors.New("enum requires a 'values' parameter")
	}
	values := strings.Split(raw, ",")
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}
	return &enumGenerator{values: values}, nil
}
