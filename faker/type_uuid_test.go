package faker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUIDGenerator_Generate(t *testing.T) {
	gen, _ := newUUIDGenerator("", nil)
	value, err := gen.Generate()
	assert.NoError(t, err)
	assert.IsType(t, "", value)
	assert.Len(t, value.(string), 38) // UUID v4 length + double quotes
}

func TestRegisterAndGetGenerator(t *testing.T) {
	// Register a dummy type
	dummyCalled := false
	dummyFactory := func(fieldName string, params map[string]string) (FieldGenerator, error) {
		dummyCalled = true
		return &uuidGenerator{}, nil
	}

	RegisterGenerator("dummy_type_test", dummyFactory)

	gen, err := GetGenerator("dummy_type_test", "example", map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, gen)
	assert.True(t, dummyCalled)

	// Unknown type should return error
	_, err = GetGenerator("missing_type", "field", nil)
	assert.Error(t, err)
}

func TestRegisterGenerator_DuplicatePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on duplicate registration, got none")
		}
	}()
	// This should panic because "uuid" is already registered in init
	RegisterGenerator("uuid", func(string, map[string]string) (FieldGenerator, error) {
		return nil, nil
	})
}
