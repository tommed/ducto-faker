package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/tommed/ducto-faker/faker"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	TotalRecords int                         `mapstructure:"total_records"`
	CustomTypes  map[string]faker.CustomType `mapstructure:"custom_types"`
	Templates    []templateDef               `mapstructure:"templates"`
}

type templateDef struct {
	Path   string `mapstructure:"path"`
	Weight int    `mapstructure:"weight"`
}

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, errors.New("config file path is required")
	}

	// read
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// parse
	var raw map[string]interface{}
	switch {
	case strings.HasSuffix(configPath, ".yaml"), strings.HasSuffix(configPath, ".yml"):
		err = yaml.Unmarshal(data, &raw)
	case strings.HasSuffix(configPath, ".json"):
		err = json.Unmarshal(data, &raw)
	default:
		err = errors.New("unsupported config format: must be .yaml or .json")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse raw config: %v", err)
	}

	// decode
	var cfg Config
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &cfg,
		TagName: "mapstructure",
	})
	err = decoder.Decode(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %v", err)
	}

	// resolve template paths relative to the config file
	cfgDir := filepath.Dir(configPath)
	for i, t := range cfg.Templates {
		relPath := filepath.Join(cfgDir, t.Path)
		if !filepath.IsAbs(relPath) {
			cfg.Templates[i].Path, _ = filepath.Abs(relPath)
		}
	}

	return &cfg, nil
}
