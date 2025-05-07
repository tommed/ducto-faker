<!--suppress HtmlDeprecatedAttribute -->
<p align="right">
    <a href="https://github.com/tommed" title="See Project Ducto">
        <img src="../assets/ducto-logo-small.png" alt="A part of Project Ducto"/>
    </a>
</p>

# Ducto Faker V1 Specification

## 1. Overview

`ducto-faker` is a Go-powered CLI and library for generating fake/test event data in newline-delimited formats. It is designed to be highly composable and configuration-driven, following idiomatic Go patterns and matching the architecture of other Ducto ecosystem components.

### Goals

* Provide a CLI under `./cmd/ducto-faker` that reads a YAML/JSON configuration and streams fake data to stdout.
* Support multiple output formats in the future (JSONL, CSV, TSV, PSV, Avro), with JSONL initially.
* Load templates from file paths (absolute or relative to the config file) using Go’s `text/template` engine.
* Support multiple templates with weighted generation and global record count.
* Allow inline modifiers in placeholders for type, min/max, probabilities, etc., and custom types via top-level definitions.
* Expose a Go library API for programmatic use.
* Follow TDD, 100% unit coverage, clean code, and Ducto conventions.

## 2. CLI Interface

```
Usage: ducto-faker [flags]

Flags:
  -c, --config string    Path to config file (yaml or json)
  -h, --help             Help for ducto-faker
```

Example:

```bash
# To push to a file
ducto-faker -c configs/sample.yaml > data.jsonl

# To pipe jsonl into another process' stdin
ducto-faker -c configs/sample.yaml | ducto-orchestrator -config config.yaml
```

## 3. Configuration Schema (v1)

```yaml
# Number of records to generate
total_records: 1000

# Optional: define custom types like enums or repeated constraints
custom_types:
  status:
    type: enum
    values: "OPEN,CLOSED,UNKNOWN"
  uk_latitude:
    type: float
    min: 49.9
    max: 59.0
  uk_longitude:
    type: float
    min: -8.6
    max: 1.8

# Templates: path can be absolute or relative to this config file
templates:
  - path: templates/login.json.tmpl
    weight: 50
  - path: templates/purchase.json.tmpl
    weight: 30
  - path: templates/error.json.tmpl
    weight: 20
```

### 3.1 Template Files

* Each template file is a `text/template`-compatible `.tmpl` file containing a JSON model with quoted placeholders.
* Placeholders use **Go's double-curly syntax**, with quoted inline-modifier strings:

  ```go
  {
    "lat":      {{"latitude:float:min=49.9:max=59.0"}},
    "lon":      {{"longitude:float:min=-8.6:max=1.8"}},
    "status":   {{"status:enum:values=OPEN,CLOSED,UNKNOWN"}},
    "user_id":  {{"user_id:uuid"}},
    "event_ts": {{"event_ts:datetime:min=2021-01-01T00:00:00Z:max=2021-12-31T23:59:59Z"}}
  }
  ```

* To reference a **custom type** defined in `custom_types`, omit inline constraints:

  ```go
  {
    "state": {{"status:enum"}}
  }
  ```

* During load, for each placeholder:
  1. Strip the `{{ }}` delimiters and extract the quoted string.
  2. Split on `:` to get `name`, `type`, and optional `key=value` pairs.
  3. If `type` matches a key in `custom_types`, load its config.
  4. Override any custom-type defaults with inline constraints.
  5. Build the corresponding `Field` generator.

## 4. Architecture & Packages

```
./
├── cmd/
│   └── ducto-faker/
│       └── main.go           # Entry point for CLI
├── config/
│   └── config.go             # Config structs & Unmarshal logic
├── faker/
│   ├── generator.go          # Core logic: template loading, placeholder parsing, record gen
│   ├── parser.go             # Template AST inspection & inline-modifier parser
│   ├── types.go              # Field type definitions & interfaces
│   └── plugins.go            # Plugin registry (faker engines)
├── output/
│   └── writer.go             # Write newline-delimited output (JSONL)
├── internal/
│   └── utils.go              # shared utilities
├── docs/
│   └── spec-v1.md            # This specification
└── go.mod
```

## 5. Placeholder & Type Resolution

* **Inline modifiers**: `{{"name:type:key=val,..."}}` parsed into a map of parameters.
* **Custom types**: top-level `custom_types` block; each entry populates defaults.
* Resolution order: defaults → custom-type defaults → inline overrides.

## 6. Weighted Sampling

* Build cumulative weights from templates’ `weight` values.
* For each record, pick a template via random sampling.

## 7. Output Writers

* **JSONL** (v1): render each filled template as one JSON object per line (no indent).
* Future: add CSV/TSV/PSV with configurable delimiter; Avro.
* Always write to `stdout`.

## 8. Testing Strategy

* TDD: tests for parser (inline modifiers), custom-type loading, template AST hooks, generator logic, and writer.
* Table-driven tests to cover various placeholder scenarios.
* Achieve 100% coverage with Go’s `testing` package.

## 9. CLI UX

* Simple, declarative flags in the config; no additional CLI flags.
* Exit codes: zero on success; non-zero on errors.
* Logging via `log`, with verbose flag in the future.

## 10. Roadmap & Future Enhancements

1. Support other formats: CSV, TSV, PSV, Avro.
2. Subject grouping: generate `x` records for `y` unique IDs.
3. External plugin system for custom generators.
4. Streaming sinks: HTTP, Kafka (via Ducto Orchestrator).

## 11. License

- Code is all licensed under [MIT](../LICENSE)
- The Ducto name, logos and robot illustrations (and likeness) are (C) Copyright 2025 Tom Medhurst, all rights reserved.
