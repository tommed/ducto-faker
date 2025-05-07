<!--suppress HtmlDeprecatedAttribute -->
<p align="right">
    <a href="https://github.com/tommed" title="See Project Ducto">
        <img src="./assets/ducto-logo-small.png" alt="A part of Project Ducto"/>
    </a>
</p>

# Ducto Faker

[![CI](https://github.com/tommed/ducto-faker/actions/workflows/ci.yml/badge.svg)](https://github.com/tommed/ducto-faker/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/tommed/ducto-faker/branch/main/graph/badge.svg)](https://codecov.io/gh/tommed/ducto-faker)

> Go-powered CLI and library for generating fake/test event data based on structured templates and inline modifiers.

---

## ✅ Features

- [x] CLI and library mode (`ducto-faker`) for streaming newline-delimited fake records
- [x] YAML/JSON-based config with template path resolution (relative to config)
- [x] Mustache-compatible templating with inline modifiers: `{{ "field:type:min=1:max=10" }}`
- [ ] Supports:
    - [ ] Primitive types: `int`, `float`, `uuid`, `datetime`
    - [x] Enums: `status`, `severity`, etc.
    - [x] Faker data: `first_name`, `last_name`, `phone`, `address_line1`, `postcode`, `country`, `zip_us`
- [x] Custom types and reusability via top-level `custom_types`
- [x] Weighted sampling between multiple templates
- [x] JSONL output streamed to `stdout` for easy piping into files or other services
- [ ] 100% unit test coverage, clean Go idioms, fully testable components
- [x] Fully MIT licensed

---

## 🚀 Getting Started

```bash
# Install
go install github.com/tommed/ducto-faker/cmd/ducto-faker@latest

# Run from config
ducto-faker -config ./configs/sample.yaml > output.jsonl
```

---

## 🧾 Example Template

```json
{
  "id": {{ "user_id:uuid" }},
  "first_name": {{ "first_name" }},
  "status": {{ "status:enum:values=ACTIVE,INACTIVE" }},
  "score": {{ "score:float:min=0:max=10:dps=2" }},
  "created": {{ "created:datetime:min=2021-01-01T00:00:00Z:max=2022-01-01T00:00:00Z" }}
}
```

---

## 📄 Configuration

```yaml
total_records: 1000

custom_types:
  status:
    type: enum
    values: 'OPEN,CLOSED,UNKNOWN'

templates:
  - path: templates/login.json.mustache
    weight: 60
  - path: templates/logout.json.mustache
    weight: 40
```

---

## 📚 Documentation

- See the [specification document](docs/spec-v1.md) for full syntax and implementation details.

---

## 🧑‍💻 Contributing

Pull requests are welcome! Please see our [Contributing Guide](./CONTRIBUTING.md).

---

## 🤖 Related Projects

- [ducto-dsl](https://github.com/tommed/ducto-dsl) — Declarative transformation engine
- [ducto-featureflags](https://github.com/tommed/ducto-featureflags) — Pluggable feature flag engine
- [ducto-orchestrator](https://github.com/tommed/ducto-orchestrator) — Event-driven transformation runtime

---

## 📜 License

- MIT licensed — see [LICENSE](./LICENSE)
- The Ducto name, logos and robot illustrations are © 2025 Tom Medhurst. All rights reserved.
