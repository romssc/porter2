# Porter

**Porter** is a flexible and composable migration toolkit for Elasticsearch, written in Go. It helps developers define, manage, and test index definitions and synthetic document generation.

```bash
go get github.com/xoticdsign/porter
```

## Features

- ✅ **Composable migration API** for index and document operations
- ⚙️ **Field-type factories** with fake data generators (`gofakeit`)
- 🔍 **Custom analyzers & normalizers** in a fluent functional style
- 📦 **Support for `bulk` document ingestion**
- 📂 **Load documents from files or dynamically generate them**
- 🧪 **Test-friendly design** (can run with real or mock Elasticsearch)
- 🔄 **Bidirectional migrations**: `.MigrateUp()` and `.MigrateDown()`

## Basic Usage

```go
import (
   "github.com/elastic/go-elasticsearch/v8"
   "github.com/xoticdsign/porter"
)

func main() {
   cc, _ := es.NewDefaultClient()

   p := porter.New(cc)

   c := porter.Config{
      Name: "porter",
      Definition: porter.DefinitionConfig{
         Mappings: &porter.MappingsConfig{
            Properties: migrator.Index.Mappings.NewFields(
               migrator.Index.Mappings.Properties.Keyword("city", porter.FakeCity),
               migrator.Index.Mappings.Properties.Integer("age", porter.FakeIntegerInt),
            ),
         },
      },
   }
}
```

### 2. Define Index Config

```go


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





