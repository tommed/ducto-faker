package faker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writeTempTemplate(t *testing.T, content string) string {
	t.Helper()
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.tmpl")
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestParseTemplate_SinglePlaceholder(t *testing.T) {
	tpl := `{"id": {{"user_id:uuid"}}}`
	path := writeTempTemplate(t, tpl)
	result, err := ParseTemplate(path)
	assert.NoError(t, err)
	if assert.Len(t, result, 1) {
		assert.Equal(t, "user_id", result[0].FieldName)
		assert.Equal(t, "uuid", result[0].Type)
		assert.Empty(t, result[0].Params)
	}
}

func TestParseTemplate_WithParams(t *testing.T) {
	tpl := `{"lat": {{"latitude:float:min=49.9:max=59.0"}}}`
	path := writeTempTemplate(t, tpl)
	result, err := ParseTemplate(path)
	assert.NoError(t, err)
	if assert.Len(t, result, 1) {
		assert.Equal(t, "latitude", result[0].FieldName)
		assert.Equal(t, "float", result[0].Type)
		assert.Equal(t, map[string]string{"min": "49.9", "max": "59.0"}, result[0].Params)
	}
}

func TestParseTemplate_MultiplePlaceholders(t *testing.T) {
	tpl := `{"ts": {{"event_ts:datetime:min=2021-01-01:max=2022-01-01"}}, "status": {{"status:enum:values=OPEN,CLOSED"}}}`
	path := writeTempTemplate(t, tpl)
	result, err := ParseTemplate(path)
	assert.NoError(t, err)
	if assert.Len(t, result, 2) {
		assert.Equal(t, "event_ts", result[0].FieldName)
		assert.Equal(t, "status", result[1].FieldName)
	}
}

func TestParseTemplate_Empty(t *testing.T) {
	tpl := `{"static": "value"}`
	path := writeTempTemplate(t, tpl)
	result, err := ParseTemplate(path)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestParseTemplate_InvalidPlaceholder(t *testing.T) {
	tpl := `{"invalid": {{"justone"}}}` // too few parts
	path := writeTempTemplate(t, tpl)
	result, err := ParseTemplate(path)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestParseTemplate_IncorrectPath(t *testing.T) {
	_, err := ParseTemplate("./non-existent-file.template")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read template")
}

func TestParseTemplate_BadTemplateToParse(t *testing.T) {
	path, _ := filepath.Abs("../testdata/bad_template.txt")
	_, err := ParseTemplate(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse template")
}

func TestParseTemplate_IfRangeWithBlocks(t *testing.T) {
	tpl := `{{if true}}{{"one:type:foo=bar"}}{{else}}{{"two:type"}}{{end}}{{range .Items}}{{"three:type:min=1:max=2"}}{{else}}{{"four:type"}}{{end}}{{with .Inner}}{{"five:type"}}{{else}}{{"six:type"}}{{end}}`
	path := writeTempTemplate(t, tpl)
	result, err := ParseTemplate(path)
	assert.NoError(t, err)
	if assert.Len(t, result, 6) {
		assert.Equal(t, "one", result[0].FieldName)
		assert.Equal(t, "two", result[1].FieldName)
		assert.Equal(t, "three", result[2].FieldName)
		assert.Equal(t, "four", result[3].FieldName)
		assert.Equal(t, "five", result[4].FieldName)
		assert.Equal(t, "six", result[5].FieldName)
	}
}
