// faker/type_enum_test.go
package faker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnumGenerator_Generate(t *testing.T) {
	gen, err := newEnumGenerator("status", map[string]string{
		"values": "OPEN,CLOSED,UNKNOWN",
	})
	assert.NoError(t, err)
	assert.NotNil(t, gen)

	for i := 0; i < 10; i++ {
		val, err := gen.Generate()
		assert.NoError(t, err)
		assert.Contains(t, []string{"\"OPEN\"", "\"CLOSED\"", "\"UNKNOWN\""}, val)
	}
}

func TestEnumGenerator_MissingValues(t *testing.T) {
	_, err := newEnumGenerator("status", map[string]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires a 'values'")
}
