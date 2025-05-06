package faker

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"text/template"
)

// LoadedTemplate represents a template ready for generation.
type LoadedTemplate struct {
	Path    string
	Weight  int
	Content string
	Fields  []Placeholder
}

// Generator is the main structure used to emit synthetic data.
type Generator struct {
	Templates    []LoadedTemplate
	CustomTypes  map[string]CustomType
	TotalRecords int
}

// NewGenerator builds and validates a generator instance from raw config data.
func NewGenerator(total int, customTypes map[string]CustomType, templates []LoadedTemplate) (*Generator, error) {
	if total <= 0 {
		return nil, fmt.Errorf("total_records must be > 0")
	}
	if len(templates) == 0 {
		return nil, fmt.Errorf("no templates provided")
	}
	return &Generator{
		Templates:    templates,
		CustomTypes:  customTypes,
		TotalRecords: total,
	}, nil
}

// GenerateAll returns a slice of generated records as raw JSON strings (one per record).
func (g *Generator) GenerateAll() ([]string, error) {
	var out []string

	// Precompute cumulative weights
	var cumulative []int
	totalWeight := 0
	for _, t := range g.Templates {
		totalWeight += t.Weight
		cumulative = append(cumulative, totalWeight)
	}

	for i := 0; i < g.TotalRecords; i++ {
		tmpl := g.pickTemplate(cumulative, totalWeight)

		// Resolve values
		values := make(map[string]any)
		for _, field := range tmpl.Fields {
			typeName, params := mergeParams(field.Type, g.CustomTypes, field)
			gen, err := GetGenerator(typeName, field.FieldName, params)
			if err != nil {
				return nil, fmt.Errorf("field '%s': %w", field.FieldName, err)
			}
			val, err := gen.Generate()
			if err != nil {
				return nil, fmt.Errorf("generate '%s': %w", field.FieldName, err)
			}
			values[field.FieldName] = val
		}

		// correct template so it speaks `text/template` instead of our custom
		// mark-up with the inline modifiers
		runtime := tmpl.Content
		re := regexp.MustCompile(`{{\s*"(.*?)"\s*}}`)
		runtime = re.ReplaceAllStringFunc(runtime, func(match string) string {
			inner := re.FindStringSubmatch(match)
			if len(inner) != 2 {
				return match // leave untouched if malformed
			}

			parts := strings.Split(inner[1], ":")
			if len(parts) < 2 {
				return match
			}

			fieldName := parts[0]
			return fmt.Sprintf("{{ .%s }}", fieldName)
		})

		tpl, err := template.New("record").Parse(runtime)
		if err != nil {
			return nil, fmt.Errorf("parse template: %w", err)
		}

		var sb strings.Builder
		err = tpl.Execute(&sb, values)
		if err != nil {
			return nil, fmt.Errorf("render template: %w", err)
		}

		// Check it's a valid JSON object produced
		var jsonObj map[string]interface{}
		if err := json.Unmarshal([]byte(sb.String()), &jsonObj); err != nil {
			return nil, fmt.Errorf("invalid json produced: %w", err)
		}
		jsonLine, _ := json.Marshal(jsonObj)

		out = append(out, string(jsonLine))
	}

	return out, nil
}

func (g *Generator) pickTemplate(cumulative []int, total int) LoadedTemplate {
	pick := rand.Intn(total)
	for i, threshold := range cumulative {
		if pick < threshold {
			return g.Templates[i]
		}
	}
	return g.Templates[len(g.Templates)-1] // fallback
}

func mergeParams(typeName string, custom map[string]CustomType, field Placeholder) (string, map[string]string) {
	out := map[string]string{}
	if def, ok := custom[field.Type]; ok {
		typeName = def.Type
		for k, v := range def.Params {
			out[k] = v
		}
	}
	for k, v := range field.Params {
		out[k] = v // override
	}
	return typeName, out
}
