# porter

**Porter** is a flexible and composable migration toolkit for Elasticsearch, written in Go. It helps developers define, manage, and test index definitions and synthetic document generation.

---

## ✨ Features

- ✅ **Composable migration API** for index and document operations
- ⚙️ **Field-type factories** with fake data generators (`gofakeit`)
- 🔍 **Custom analyzers & normalizers** in a fluent functional style
- 📦 **Support for `bulk` document ingestion**
- 📂 **Load documents from files or dynamically generate them**
- 🧪 **Test-friendly design** (can run with real or mock Elasticsearch)
- 🔄 **Bidirectional migrations**: `.MigrateUp()` and `.MigrateDown()`

---

## 📦 Installation

```bash
go get github.com/xoticdsign/porter
```

## 🧰 Usage

### 1. Initialize Porter

```go
import (
    "github.com/elastic/go-elasticsearch/v8"
	"github.com/xoticdsign/porter"
)

func main() {
    client, _ := es.NewDefaultClient()
    migrator := porter.New(client)
}
```

### 2. Define Index Config

```go

config := porter.Config{
	Name: "my-index",
	Definition: porter.DefinitionConfig{
		Settings: &porter.SettingsConfig{
			NumberOfShards:   1,
			NumberOfReplicas: 0,
			Analysis: &porter.AnalysisConfig{
				Analyzer: migrator.Index.Settings.Analysis.NewAnalyzer(
					migrator.Index.Settings.Analysis.Analyzer.Simple("my-analyzer"),
				),
				Normalizer: migrator.Index.Settings.Analysis.NewNormalizer(
					migrator.Index.Settings.Analysis.Normalizer.Custom(
						"my-normalizer",
						migrator.Index.Settings.Analysis.Normalizer.Custom.WithFilter([]porter.NormalizerCustomFilter{
							porter.NormalizerCustomFilterASCIIFolding,
							porter.NormalizerCustomFilterLowercase,
						}),
					),
				),
			},
		},
		Mappings: &porter.MappingsConfig{
			Properties: migrator.Index.Mappings.NewFields(
				migrator.Index.Mappings.Properties.Keyword("city", porter.FakeCity),
				migrator.Index.Mappings.Properties.Integer("age", porter.FakeIntegerInt),
			),
		},
	},
}
```

### 3. Migrate Up

```go
err := migrator.MigrateUp(
	config,
	migrator.Index.MigrateIndex(),
	migrator.Documents.MigrateDocuments(
		migrator.Documents.Origin.Generate(100),
	),
)
```

### 4. Migrate Down

```go
err := migrator.MigrateDown(
	config,
	migrator.Documents.MigrateDocuments(nil),
	migrator.Index.MigrateIndex(),
)
```

## 📁 Document Origins

Porter supports two sources for documents:
    - Generate using fake data: `.Generate(n)`
    - Load from file: `.FromFile("data.json")`

## 🧱 Field & Analyzer DSL

✅ Keyword, Text, Integer, Float, Date, IP, etc.
🔤 Analyzers: standard, simple, whitespace, custom, language, etc.
🧹 Filters: asciifolding, lowercase, stopwords, etc.

## 🔧 Contributing

Feel free to open issues or PRs for:
    - Adding more field types and filter support
    - Enhancing analyzer/normalizer composition
    - Improving test coverage and fixtures

## 📄 License

(MIT)[]





