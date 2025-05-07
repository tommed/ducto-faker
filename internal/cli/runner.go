package cli

import (
	"flag"
	"fmt"
	"github.com/tommed/ducto-faker/config"
	"github.com/tommed/ducto-faker/faker"
	"io"
	"os"
)

//goland:noinspection GoUnhandledErrorResult
func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	var configPath string
	fs := flag.NewFlagSet("cli", flag.ContinueOnError)
	fs.StringVar(&configPath, "config", "", "Path to config file")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load config: %v\n", err)
		return 1
	}

	// Resolve template paths and load content
	var loadedTemplates []faker.LoadedTemplate
	for _, t := range cfg.Templates {
		content, err := os.ReadFile(t.Path)
		if err != nil {
			fmt.Fprintf(stderr, "failed to read template '%s': %v\n", t.Path, err)
			return 1
		}

		placeholders, err := faker.ParseTemplate(t.Path)
		if err != nil {
			fmt.Fprintf(stderr, "failed to parse template '%s': %v\n", t.Path, err)
			return 1
		}

		loadedTemplates = append(loadedTemplates, faker.LoadedTemplate{
			Path:    t.Path,
			Weight:  t.Weight,
			Content: string(content),
			Fields:  placeholders,
		})
	}

	g, err := faker.NewGenerator(cfg.TotalRecords, cfg.CustomTypes, loadedTemplates)
	if err != nil {
		fmt.Fprintf(stderr, "failed to initialize generator: %v\n", err)
		return 1
	}

	records, err := g.GenerateAll()
	if err != nil {
		fmt.Fprintf(stderr, "failed to generate data: %v\n", err)
		return 1
	}

	for _, rec := range records {
		fmt.Fprintln(stdout, rec)
	}
	return 0
}
