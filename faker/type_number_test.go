package faker

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntGenerator_DefaultRange(t *testing.T) {
	gen, err := newIntGenerator("field", map[string]string{})
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		val, err := gen.Generate()
		assert.NoError(t, err)
		n := val.(int64)
		assert.GreaterOrEqual(t, n, int64(0))
		assert.LessOrEqual(t, n, int64(100))
	}
}

func TestIntGenerator_WithRange(t *testing.T) {
	gen, err := newIntGenerator("field", map[string]string{"min": "10", "max": "20"})
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		val, err := gen.Generate()
		assert.NoError(t, err)
		n := val.(int64)
		assert.GreaterOrEqual(t, n, int64(10))
		assert.LessOrEqual(t, n, int64(20))
	}
}

func TestIntGenerator_PrefixSuffix(t *testing.T) {
	generator, _ := newIntGenerator("field", map[string]string{"min": "1", "max": "1", "suffix": "p", "prefix": "£"})
	val, err := generator.Generate()
	assert.NoError(t, err)

	strVal, ok := val.(string)
	assert.True(t, ok)
	assert.Equal(t, "£1p", strVal)
}

func TestIntGenerator_InvalidParams(t *testing.T) {
	_, err := newIntGenerator("field", map[string]string{"min": "a"})
	assert.Error(t, err)

	_, err = newIntGenerator("field", map[string]string{"min": "10", "max": "5"})
	assert.Error(t, err)

	_, err = newIntGenerator("field", map[string]string{"min": "10", "max": "a"})
	assert.Error(t, err)

	_, err = newIntGenerator("field", map[string]string{"min": "10", "max": "20", "left_zero_padding": "p"})
	assert.Error(t, err)
}

func TestFloatGenerator_StoreName(t *testing.T) {
	gen := intGenerator{
		min:             100,
		max:             9999,
		leftZeroPadding: 5,
		prefix:          "ca",
		suffix:          "bos01",
	}
	val, err := gen.Generate()

	assert.NoError(t, err)
	strVal, ok := val.(string)
	assert.True(t, ok)
	assert.Equal(t, 12, len(strVal), strVal+" was not 12 characters")
}

func TestFloatGenerator_DefaultRange(t *testing.T) {
	gen, err := newFloatGenerator("field", map[string]string{})
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		val, err := gen.Generate()
		assert.NoError(t, err)
		n := val.(float64)
		assert.GreaterOrEqual(t, n, 0.0)
		assert.LessOrEqual(t, n, 100.0)
	}
}

func TestFloatGenerator_WithRange(t *testing.T) {
	gen, err := newFloatGenerator("field", map[string]string{"min": "1.5", "max": "2.5"})
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		val, err := gen.Generate()
		assert.NoError(t, err)
		n := val.(float64)
		assert.GreaterOrEqual(t, n, 1.5)
		assert.LessOrEqual(t, n, 2.5)
	}
}

func TestFloatGenerator_WithDPS(t *testing.T) {
	gen, err := newFloatGenerator("field", map[string]string{"min": "1.0", "max": "2.0", "dps": "2"})
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		val, err := gen.Generate()
		assert.NoError(t, err)
		str := strconv.FormatFloat(val.(float64), 'f', -1, 64)
		parts := strings.Split(str, ".")
		if len(parts) == 2 {
			assert.LessOrEqual(t, len(parts[1]), 2)
		}
	}
}

func TestFloatGenerator_InvalidParams(t *testing.T) {
	_, err := newFloatGenerator("field", map[string]string{"min": "abc"})
	assert.Error(t, err)

	_, err = newFloatGenerator("field", map[string]string{"min": "5.0", "max": "4.0"})
	assert.Error(t, err)

	_, err = newFloatGenerator("field", map[string]string{"min": "5.0", "max": "a"})
	assert.Error(t, err)

	_, err = newFloatGenerator("field", map[string]string{"dps": "-1"})
	assert.Error(t, err)

	_, err = newFloatGenerator("field", map[string]string{"dps": "99"})
	assert.Error(t, err)

	_, err = newFloatGenerator("field", map[string]string{"dps": "not-a-number"})
	assert.Error(t, err)
}
