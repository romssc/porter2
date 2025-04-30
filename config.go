package porter

/*

This file defines the configuration structure for managing Elasticsearch index settings and mappings.

The Config{} struct represents the complete configuration for an Elasticsearch index, including
the index name and its associated settings and mappings. The DefinitionConfig{} struct encapsulates
the settings and mappings configurations.

The SettingsConfig{} struct allows users to define index-specific settings, such as the number of
shards and replicas, as well as custom analysis configurations (analyzers and normalizers).

The AnalysisConfig{} struct defines the custom analyzers and normalizers used for text analysis in the
index.

The MappingsConfig{} struct contains the properties of the index, mapping each field name to its
definition and properties.

This structure is designed to provide an easy way to define and manage Elasticsearch index configurations
in a structured and flexible manner, facilitating the creation of indices with appropriate settings and mappings.

*/

// Config{} represents the overall configuration for an Elasticsearch index.
type Config struct {
	Name       string
	Definition DefinitionConfig
}

// DefinitionConfig{} contains the settings and mappings for an Elasticsearch index.
type DefinitionConfig struct {
	Settings *SettingsConfig `json:"settings,omitempty"`
	Mappings *MappingsConfig `json:"mappings,omitempty"`
}

// SettingsConfig{} defines the settings related to an Elasticsearch index, including the number of shards, replicas, and custom analysis configurations.
type SettingsConfig struct {
	NumberOfShards   int             `json:"number_of_shards,omitempty"`
	NumberOfReplicas int             `json:"number_of_replicas,omitempty"`
	Analysis         *AnalysisConfig `json:"analysis,omitempty"`
}

// AnalysisConfig{} holds custom analysis settings for the Elasticsearch index, including analyzers and normalizers to control text processing during indexing and searching.
type AnalysisConfig struct {
	Analyzer   map[string]interface{} `json:"analyzer,omitempty"`
	Normalizer map[string]interface{} `json:"normalizer,omitempty"`
}

// MappingsConfig{} defines the field mappings for an Elasticsearch index, including the types and properties for each field in the index.
type MappingsConfig struct {
	Properties map[string]interface{} `json:"properties,omitempty"`
}
