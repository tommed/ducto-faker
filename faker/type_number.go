package faker

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

type intGenerator struct {
	min, max, leftZeroPadding int64
	prefix, suffix            string
}

func (g *intGenerator) Generate() (any, error) {
	var val any = rand.Int63n(g.max-g.min+1) + g.min
	if g.prefix != "" || g.suffix != "" || g.leftZeroPadding > 0 {
		paddingFormat := "%0" + strconv.Itoa(int(g.leftZeroPadding)) + "d"
		return g.prefix + fmt.Sprintf(paddingFormat, val) + g.suffix, nil
	}
	return val, nil
}

func newIntGenerator(_ string, params map[string]string) (FieldGenerator, error) {
	minVal := int64(0)
	maxVal := int64(100)
	leftZeroPadding := int64(0)
	prefix := params["prefix"]
	suffix := params["suffix"]
	var err error
	if val, ok := params["min"]; ok {
		minVal, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min param: %w", err)
		}
	}
	if val, ok := params["max"]; ok {
		maxVal, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid max param: %w", err)
		}
	}
	if val, ok := params["left_zero_padding"]; ok {
		leftZeroPadding, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid left_zero_padding param: %w", err)
		}
	}
	if minVal > maxVal {
		return nil, errors.New("min cannot be greater than max")
	}
	return &intGenerator{minVal, maxVal, leftZeroPadding, prefix, suffix}, nil
}

type floatGenerator struct {
	min, max float64
	dps      int // decimal places
}

func (g *floatGenerator) Generate() (any, error) {
	val := g.min + rand.Float64()*(g.max-g.min)
	if g.dps > 0 {
		factor := math.Pow(10, float64(g.dps))
		return math.Round(val*factor) / factor, nil
	}
	return val, nil
}

func newFloatGenerator(_ string, params map[string]string) (FieldGenerator, error) {
	min := float64(0.0)
	max := float64(100.0)
	dps := 0
	var err error
	if val, ok := params["min"]; ok {
		min, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min param: %w", err)
		}
	}
	if val, ok := params["max"]; ok {
		max, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid max param: %w", err)
		}
	}
	if val, ok := params["dps"]; ok {
		dps, err = strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid dps param: %w", err)
		}
		if dps < 0 || dps > 10 {
			return nil, errors.New("dps must be between 0 and 10")
		}
	}
	if min > max {
		return nil, errors.New("min cannot be greater than max")
	}
	return &floatGenerator{min, max, dps}, nil
}

func init() {
	RegisterGenerator("int", newIntGenerator)
	RegisterGenerator("float", newFloatGenerator)
}
