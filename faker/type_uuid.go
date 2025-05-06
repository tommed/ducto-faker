package faker

import (
	"github.com/google/uuid"
)

func init() {
	RegisterGenerator("uuid", newUUIDGenerator)
}

type uuidGenerator struct{}

func (g *uuidGenerator) Generate() (any, error) {
	return QuoteString(uuid.NewString()), nil
}

func newUUIDGenerator(_ string, _ map[string]string) (FieldGenerator, error) {
	return &uuidGenerator{}, nil
}
