# Porter2

**Porter** is a flexible and composable migration toolkit for Elasticsearch, written in Go. It helps developers define, manage, and test index definitions and synthetic document generation.

- 📦 Define Elasticsearch index settings and mappings with a fluent DSL
- 📁 Support for reading documents from files or generating in-memory

## ⚙️ Usage

This section outlines the typical flow for using Porter in your application: install the package, set up the Elasticsearch client, create a migration instance, define a configuration, and run the migration.

### 1. Get Package

Install **Porter** using `go get`.

```bash
go get github.com/xoticdsign/porter2
```

### 2. Configure Elasticsearch

Establish a connection to your Elasticsearch instance using the **[official Go client](https://github.com/elastic/go-elasticsearch)**.

```go
cc, err := elasticsearch.NewDefaultClient()
if err != nil {
   panic(err)
}
```

### 3. Create new Porter instance

Initialize a new **Porter** migrator with the Elasticsearch client.

```go
p := porter.New(cc)
```

### 4. Create and define Porter config

Define your index name, settings, and field mappings using Porter’s fluent config API.

```go
c := porter.Config{
   Name: ...
   Definition: ...
}
```

### 5. Migrate

Run `.MigrateUp()` to create an index and insert fake documents or load them from your own file, and `.MigrateDown()` to remove partially or completely.

```go
// Migrate Up
err := migrator.MigrateUp(
   config,
   migrator.Index.MigrateIndex(),
   migrator.Documents.MigrateDocuments(migrator.Documents.Origin.Generate(100)),
)
if err != nil {
   panic(err)
}

// Migrate Down
err := migrator.MigrateDown(
      config,
      migrator.Documents.MigrateDocuments(nil),
      migrator.Index.MigrateIndex(),
   )
if err != nil {
   panic(err)
}
```

## 📘 Examples

These examples demonstrate real-world usage of Porter for creating and deleting Elasticsearch indices and documents.

### Migrating Up

This example demonstrates how to create a new index and insert **100 fake documents**. It defines a basic Config object with field mappings. Then it applies `.MigrateUp()`, which:
- Creates the index using the provided settings/mappings
- Inserts 100 generated documents using the configured field types

```go
package main

import (
   "github.com/elastic/go-elasticsearch/v8"
   porter "github.com/xoticdsign/porter2"
)

func main() {
   cc, _ := es.NewDefaultClient()

   p := porter.New(cc)

   c := porter.Config{
      Name: "index",
      Definition: porter.DefinitionConfig{
         Mappings: &porter.MappingsConfig{
            Properties: migrator.Index.Mappings.NewFields(
               migrator.Index.Mappings.Properties.Keyword("keyword", porter.FakeCity),
               migrator.Index.Mappings.Properties.Integer("integer", porter.FakeIntegerInt),
            ),
         },
      },
   }

   err := migrator.MigrateUp(
      config,
      migrator.Index.MigrateIndex(),
      migrator.Documents.MigrateDocuments(migrator.Documents.Origin.Generate(100)),
   )
   if err != nil {
      panic(err)
   }
}
```

### Migrating Down

This example demonstrates how to delete documents and the index. `.MigrateDown()` will:
- Delete all documents using a match_all query
- Drop the index itself

```go
package main

import (
   "github.com/elastic/go-elasticsearch/v8"
   porter "github.com/xoticdsign/porter2"
)

func main() {
   cc, _ := es.NewDefaultClient()

   p := porter.New(cc)

   c := porter.Config{
      Name: "index",
      Definition: porter.DefinitionConfig{
         Mappings: &porter.MappingsConfig{
            Properties: migrator.Index.Mappings.NewFields(
               migrator.Index.Mappings.Properties.Keyword("keyword", porter.FakeCity),
               migrator.Index.Mappings.Properties.Integer("integer", porter.FakeIntegerInt),
            ),
         },
      },
   }

   err := migrator.MigrateDown(
      config,
      migrator.Documents.MigrateDocuments(nil),
      migrator.Index.MigrateIndex(),
   )
   if err != nil {
      panic(err)
   }
}
```

## 🧠 Porter API Reference

This section describes the core building blocks of the **Porter** toolkit, including its primary functions, index/document operations, and configuration options.

### Porter main functions

These are the top-level functions for initializing and running migrations.

| Function                                                                        | Description                             |
|---------------------------------------------------------------------------------|-----------------------------------------|
| `porter.New(< Elasticsearch client >)`                                          | Initializes a new `Porter` migrator     |
| `.MigrateUp(< Porter config >, < Index operation >, < Documents operation >)`   | Creates an index and inserts documents  |
| `.MigrateDown(< Porter config >, < Documents operation >, < Index operation >)` | Deletes documents and the index         |

### Index operations

Functions related to creating or skipping index operations during migration.

| Function          | Description                                     |
|-------------------|-------------------------------------------------|
| `.MigrateIndex()` | Creates or deletes the index based on direction |
| `.NoIndex()`      | Skip index operations                           |

### Documents operations

Functions related to inserting or skipping document operations during migration.

| Function                                   | Description                 |
|--------------------------------------------|-----------------------------|
| `.MigrateDocuments(< Origin operation >)`  | Adds or deletes documents   |
| `.NoDocuments()`                           | Skip document operations    |

### Origin operations

**Origin operations** define where the documents should come from.

| Function                                          | Description                                                   |
|---------------------------------------------------|---------------------------------------------------------------|
| `.Generate(< Amount of documents to generate >)`  | Dynamically generates documents using configured field fakes. |
| `.FromFile(< Path to File to with migrations >)`  | Loads raw JSON-formatted documents from a file.               |

## 🛠 Configuring Porter

**Porter** configuration is done using the porter.Config struct:
- `Name`: Name of the Elasticsearch index
- `Definition.Settings`: Defines shards, replicas, analyzers, and normalizers
- `Definition.Mappings`: Defines field properties like type, storage, analyzers, etc.

### Defining Field Types

Field types are created using fluent builder functions under p.Index.Mappings.Properties. Each type has optional configuration methods to customize it's behavior.

```go
Properties: p.Index.Mappings.NewFields(
   p.Index.Mappings.Properties.Keyword("keyword", porter.FakeCity,
      p.Index.Mappings.Properties.Keyword.WithStore(ture),
      p.Index.Mappings.Properties.Keyword.WithNormalizer("normalizer"),
   ),
   p.Index.Mappings.Properties.Integer("integer", porter.FakeIntegerInt,
      p.Index.Mappings.Properties.Integer.WithStore(true),
      p.Index.Mappings.Properties.Integer.WithNullValue(0),
   ),
),
```

The value generators (like `porter.FakeCity`, `porter.FakeIntegerInt`) are used when generating documents dynamically with `.Origin.Generate(...)`.

### Supported Field Types

You can use the following field types with corresponding builder functions:

- `Keyword`
- `Text`
- `Integer`
- `Long`
- `Short`
- `Byte`
- `Float`
- `Double`
- `HalfFloat`
- `ScaledFloat`
- `Date`
- `Boolean`
- `IP`

Each type has dedicated `.With*()` helpers (e.g. `.WithIndex(...)`, `.WithStore(...)`, `.WithCoerce(...)`, `.WithNullValue(...)`, etc.).

### Defining Analyzers

**Analyzers** are configured inside `Settings.Analysis.Analyzer` using built-in or custom types. Here's how to define a simple custom analyzer:

```go
Analysis: &porter.AnalysisConfig{
   Analyzer: p.Index.Settings.Analysis.NewAnalyzer(
      p.Index.Settings.Analysis.Analyzer.Custom("analyzer",
         p.Index.Settings.Analysis.Analyzer.Custom.WithTokenizer("tokenizer"),
         p.Index.Settings.Analysis.Analyzer.Custom.WithFilter([]porter.AnalyzerCustomFilter{
            porter.AnalyzerCustomFilterLowercase,
            porter.AnalyzerCustomFilterStop,
         }),
      ),
   ),
},
```

You can also use built-in **analyzers** like:

```go
p.Index.Settings.Analysis.NewAnalyzer(
   p.Index.Settings.Analysis.Analyzer.Simple("analyzer"),
)
```

### Defining Normalizers

**Normalizers** work similarly to analyzers but are applied to keyword fields. You define them using `Settings.Analysis.Normalizer`:

```go
Normalizer: p.Index.Settings.Analysis.NewNormalizer(
   p.Index.Settings.Analysis.Normalizer.Custom("normalizer",
      p.Index.Settings.Analysis.Normalizer.Custom.WithFilter([]porter.NormalizerCustomFilter{
         porter.NormalizerCustomFilterASCIIFolding,
         porter.NormalizerCustomFilterLowercase,
      }),
   ),
),
```

This **normalizer** can now be referenced by any keyword field via `.WithNormalizer(< Normalizer name >)`.

## 🤝 Contribution

Contributions are welcome! If you’d like to improve the toolkit, fix bugs, or add features:
- Fork this repository
- Create your feature branch: `git checkout -b feature/my-feature`
- Commit your changes: `git commit -am "Add my feature"`
- Push to the branch: `git push origin feature/my-feature`
- Open a **pull request**
- Please ensure your code is clean, covered by tests, and adheres to idiomatic Go practices.

If you have ideas, feedback, or questions—feel free to open an issue or start a discussion.

## 📄 License

[MIT](https://mit-license.org/)
