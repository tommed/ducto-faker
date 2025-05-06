// faker/generator_test.go
package faker

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGenerator_Success(t *testing.T) {
	templates := []LoadedTemplate{
		{
			Path:    "fake.tmpl",
			Weight:  1,
			Content: `{"id": {{"user_id:uuid"}}}`,
			Fields: []Placeholder{
				{
					FieldName: "user_id",
					Type:      "uuid",
					Params:    map[string]string{},
				},
			},
		},
	}

	gen, err := NewGenerator(10, map[string]CustomType{}, templates)
	assert.NoError(t, err)
	assert.NotNil(t, gen)
	assert.Equal(t, 10, gen.TotalRecords)
	assert.Len(t, gen.Templates, 1)
}

func TestNewGenerator_InvalidTotal(t *testing.T) {
	_, err := NewGenerator(0, nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "total_records")
}

func TestNewGenerator_NoTemplates(t *testing.T) {
	_, err := NewGenerator(5, nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no templates")
}

func TestGenerator_GenerateAll_UUIDOnly(t *testing.T) {
	templates := []LoadedTemplate{
		{
			Path:    "uuid-only",
			Weight:  1,
			Content: `{"id": {{"user_id:uuid"}}}`,
			Fields: []Placeholder{
				{
					FieldName: "user_id",
					Type:      "uuid",
					Params:    map[string]string{},
				},
			},
		},
	}

	g, err := NewGenerator(5, map[string]CustomType{}, templates)
	assert.NoError(t, err)

	records, err := g.GenerateAll()
	assert.NoError(t, err)
	assert.Len(t, records, 5)
	for _, rec := range records {
		assert.True(t, strings.Contains(rec, "id"))
		assert.True(t, strings.Contains(rec, "-")) // UUID must have dashes
		assert.GreaterOrEqual(t, len(templates[0].Content)+20, len(rec))
	}
}
