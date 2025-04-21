package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xoticdsign/porter"
	"github.com/xoticdsign/porter/internal/lib"
	"github.com/xoticdsign/porter/internal/tests/suite"
)

// FUNCTIONAL TESTS

func TestCreateAnalyzer_Functional(t *testing.T) {
	s := suite.New(t, false)

	cases := []struct {
		name     string
		in       lib.AnalyzerFunc
		expected map[string]interface{}
	}{
		{
			name: "custom full case",
			in: s.Migrator.Components.Settings.Analysis.Analyzer.Custom(
				"custom",
				s.Migrator.Components.Settings.Analysis.Analyzer.Custom.WithTokenizer("tokenizer"),
				s.Migrator.Components.Settings.Analysis.Analyzer.Custom.WithFilters([]string{"word", "sentence"}),
				s.Migrator.Components.Settings.Analysis.Analyzer.Custom.WithCharFilter([]string{"q", "w", "e"}),
			),
			expected: map[string]interface{}{
				"custom": map[string]interface{}{
					"type":        "custom",
					"tokenizer":   "tokenizer",
					"filter":      []string{"word", "sentence"},
					"char_filter": []string{"q", "w", "e"},
				},
			},
		},
		{
			name: "custom empty case",
			in:   s.Migrator.Components.Settings.Analysis.Analyzer.Custom("custom_empty"),
			expected: map[string]interface{}{
				"custom_empty": map[string]interface{}{
					"type": "custom",
				},
			},
		},
		{
			name: "russian case",
			in: s.Migrator.Components.Settings.Analysis.Analyzer.LanguageSpecific(
				"russian_analyzer",
				porter.AnalyzerLanguageSpecificRussian,
			),
			expected: map[string]interface{}{
				"russian_analyzer": map[string]interface{}{
					"type": "russian",
				},
			},
		},
		{
			name: "english case",
			in: s.Migrator.Components.Settings.Analysis.Analyzer.LanguageSpecific(
				"english_analyzer",
				porter.AnalyzerLanguageSpecificEnglish,
			),
			expected: map[string]interface{}{
				"english_analyzer": map[string]interface{}{
					"type": "english",
				},
			},
		},
		{
			name: "standard with max_token_length case",
			in: s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.Standart(
				"standard_custom",
				s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.Standart.WithMaxTokenLength(255),
			),
			expected: map[string]interface{}{
				"standard_custom": map[string]interface{}{
					"type":             "standard",
					"max_token_length": 255,
				},
			},
		},
		{
			name: "whitespace case",
			in:   s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.Whitespace("whitespace_custom"),
			expected: map[string]interface{}{
				"whitespace_custom": map[string]interface{}{
					"type": "whitespace",
				},
			},
		},
		{
			name: "stop with stopwords case",
			in: s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.Stop(
				"stop_custom",
				s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.Stop.WithStopwords([]string{"_english_"}),
			),
			expected: map[string]interface{}{
				"stop_custom": map[string]interface{}{
					"type":      "stop",
					"stopwords": []string{"_english_"},
				},
			},
		},
		{
			name: "keyword case",
			in:   s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.Keyword("keyword_custom"),
			expected: map[string]interface{}{
				"keyword_custom": map[string]interface{}{
					"type": "keyword",
				},
			},
		},
		{
			name: "pattern with pattern case",
			in: s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.Pattern(
				"pattern_custom",
				s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.Pattern.WithPattern("\\W+"),
			),
			expected: map[string]interface{}{
				"pattern_custom": map[string]interface{}{
					"type":    "pattern",
					"pattern": "\\W+",
				},
			},
		},
		{
			name: "edge ngram with min and max case",
			in: s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.EngeNGram(
				"edge_ngram_custom",
				s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.EngeNGram.WithMinGram(2),
				s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.EngeNGram.WithMaxGram(5),
			),
			expected: map[string]interface{}{
				"edge_ngram_custom": map[string]interface{}{
					"type":     "edge_ngram",
					"min_gram": 2,
					"max_gram": 5,
				},
			},
		},
		{
			name: "custom missing name case",
			in: s.Migrator.Components.Settings.Analysis.Analyzer.Custom(
				"",
				s.Migrator.Components.Settings.Analysis.Analyzer.Custom.WithTokenizer("std"),
			),
			expected: map[string]interface{}{
				"": map[string]interface{}{
					"type":      "custom",
					"tokenizer": "std",
				},
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {

			m := s.Migrator.Components.CreateAnalyzer(cs.in)

			assert.Equal(t, cs.expected, m)
		})
	}
}

func TestCreateNormalizer_Functional(t *testing.T) {
	s := suite.New(t, false)

	cases := []struct {
		name     string
		in       lib.NormalizerFunc
		expected map[string]interface{}
	}{
		{
			name: "lowercase case",
			in:   s.Migrator.Components.Settings.Analysis.Normalizer.Predefined.Lowercase("lowercase_norm"),
			expected: map[string]interface{}{
				"lowercase_norm": map[string]interface{}{
					"type": "lowercase",
				},
			},
		},
		{
			name: "asciifolding case",
			in:   s.Migrator.Components.Settings.Analysis.Normalizer.Predefined.ASCIIFolding("ascii_norm"),
			expected: map[string]interface{}{
				"ascii_norm": map[string]interface{}{
					"type": "asciifolding",
				},
			},
		},
		{
			name: "keyword case",
			in:   s.Migrator.Components.Settings.Analysis.Normalizer.Predefined.Keyword("keyword_norm"),
			expected: map[string]interface{}{
				"keyword_norm": map[string]interface{}{
					"type": "keyword",
				},
			},
		},
		{
			name: "custom with filters case",
			in: s.Migrator.Components.Settings.Analysis.Normalizer.Custom(
				"custom_filters",
				s.Migrator.Components.Settings.Analysis.Normalizer.Custom.WithFilter([]string{"lowercase", "asciifolding"}),
			),
			expected: map[string]interface{}{
				"custom_filters": map[string]interface{}{
					"type":   "custom",
					"filter": []string{"lowercase", "asciifolding"},
				},
			},
		},
		{
			name: "custom with char filters case",
			in: s.Migrator.Components.Settings.Analysis.Normalizer.Custom(
				"custom_char_filters",
				s.Migrator.Components.Settings.Analysis.Normalizer.Custom.WithCharFilter([]string{"html_strip"}),
			),
			expected: map[string]interface{}{
				"custom_char_filters": map[string]interface{}{
					"type":        "custom",
					"char_filter": []string{"html_strip"},
				},
			},
		},
		{
			name: "custom full case",
			in: s.Migrator.Components.Settings.Analysis.Normalizer.Custom(
				"custom_full",
				s.Migrator.Components.Settings.Analysis.Normalizer.Custom.WithFilter([]string{"lowercase"}),
				s.Migrator.Components.Settings.Analysis.Normalizer.Custom.WithCharFilter([]string{"html_strip"}),
			),
			expected: map[string]interface{}{
				"custom_full": map[string]interface{}{
					"type":        "custom",
					"filter":      []string{"lowercase"},
					"char_filter": []string{"html_strip"},
				},
			},
		},
		{
			name: "custom empty case",
			in:   s.Migrator.Components.Settings.Analysis.Normalizer.Custom("custom_empty"),
			expected: map[string]interface{}{
				"custom_empty": map[string]interface{}{
					"type": "custom",
				},
			},
		},
		{
			name: "custom missing name case",
			in: s.Migrator.Components.Settings.Analysis.Normalizer.Custom(
				"",
				s.Migrator.Components.Settings.Analysis.Normalizer.Custom.WithFilter([]string{"lowercase"}),
			),
			expected: map[string]interface{}{
				"": map[string]interface{}{
					"type":   "custom",
					"filter": []string{"lowercase"},
				},
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {

			m := s.Migrator.Components.CreateNormalizer(cs.in)

			assert.Equal(t, cs.expected, m)
		})
	}
}

func TestCreateProperties_Functional(t *testing.T) {
	s := suite.New(t, false)

	cases := []struct {
		name     string
		in       lib.PropertiesFunc
		expected map[string]interface{}
	}{
		{
			name: "keyword case",
			in: s.Migrator.Components.Mappings.Properties.Keyword(
				"status",
				s.Migrator.Components.Mappings.Properties.Keyword.WithIgnoreAbove(256),
			),
			expected: map[string]interface{}{
				"status": map[string]interface{}{
					"type":         "keyword",
					"ignore_above": 256,
				},
			},
		},
		{
			name: "text case",
			in: s.Migrator.Components.Mappings.Properties.Text(
				"description",
				s.Migrator.Components.Mappings.Properties.Text.WithAnalyzer("standard"),
			),
			expected: map[string]interface{}{
				"description": map[string]interface{}{
					"type":     "text",
					"analyzer": "standard",
				},
			},
		},
		{
			name: "integer case",
			in:   s.Migrator.Components.Mappings.Properties.Integer("age"),
			expected: map[string]interface{}{
				"age": map[string]interface{}{
					"type": "integer",
				},
			},
		},
		{
			name: "boolean case",
			in:   s.Migrator.Components.Mappings.Properties.Boolean("is_active"),
			expected: map[string]interface{}{
				"is_active": map[string]interface{}{
					"type": "boolean",
				},
			},
		},
		{
			name: "date case",
			in: s.Migrator.Components.Mappings.Properties.Date(
				"created_at",
				s.Migrator.Components.Mappings.Properties.Date.WithFormat("strict_date_optional_time||epoch_millis"),
			),
			expected: map[string]interface{}{
				"created_at": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_optional_time||epoch_millis",
				},
			},
		},
		{
			name: "float case",
			in:   s.Migrator.Components.Mappings.Properties.Float("rating"),
			expected: map[string]interface{}{
				"rating": map[string]interface{}{
					"type": "float",
				},
			},
		},
		{
			name: "ip case",
			in:   s.Migrator.Components.Mappings.Properties.IP("client_ip"),
			expected: map[string]interface{}{
				"client_ip": map[string]interface{}{
					"type": "ip",
				},
			},
		},
		{
			name: "geopoint case",
			in:   s.Migrator.Components.Mappings.Properties.GeoPoint("location"),
			expected: map[string]interface{}{
				"location": map[string]interface{}{
					"type": "geo_point",
				},
			},
		},
		{
			name: "binary case",
			in:   s.Migrator.Components.Mappings.Properties.Binary("file_data"),
			expected: map[string]interface{}{
				"file_data": map[string]interface{}{
					"type": "binary",
				},
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {

			m := s.Migrator.Components.CreateProperties(cs.in)

			assert.Equal(t, cs.expected, m)
		})
	}
}

// INTEGRATION TESTS

func TestMigrateUp_Integration(t *testing.T) {
	s := suite.New(t, true)
	defer s.Elasticsearch.Container.Terminate(context.Background())

	cases := []struct {
		name     string
		in       porter.Config
		to       int
		expected error
	}{
		{
			name: "standard analyzer with text field case",
			in: porter.Config{
				IndexDefinition: porter.IndexConfig{
					Name: "standard-analyzer-index",
					Schema: porter.SchemaConfig{
						Settings: &porter.SettingsConfig{
							Analysis: &porter.AnalysisConfig{
								Analyzer: s.Migrator.Components.CreateAnalyzer(
									s.Migrator.Components.Settings.Analysis.Analyzer.GeneralPurpose.Standart("standart"),
								),
							},
						},
						Mappings: &porter.MappingsConfig{
							Properties: s.Migrator.Components.CreateProperties(
								s.Migrator.Components.Mappings.Properties.Text("description"),
							),
						},
					},
				},
				MigrationsPath: "",
			},
			to:       porter.LevelIndex,
			expected: nil,
		},
		{
			name: "lowercase normalizer with keyword field case",
			in: porter.Config{
				IndexDefinition: porter.IndexConfig{
					Name: "lowercase-normalizer-index",
					Schema: porter.SchemaConfig{
						Settings: &porter.SettingsConfig{
							Analysis: &porter.AnalysisConfig{
								Normalizer: s.Migrator.Components.CreateNormalizer(
									s.Migrator.Components.Settings.Analysis.Normalizer.Predefined.Lowercase("lowercase_custom"),
								),
							},
						},
						Mappings: &porter.MappingsConfig{
							Properties: s.Migrator.Components.CreateProperties(
								s.Migrator.Components.Mappings.Properties.Keyword(
									"email",
									s.Migrator.Components.Mappings.Properties.Keyword.WithNormalizer("lowercase_custom"),
								),
							),
						},
					},
				},
				MigrationsPath: "",
			},
			to:       porter.LevelIndex,
			expected: nil,
		},
		{
			name: "custom analyzer with char filter and filters case",
			in: porter.Config{
				IndexDefinition: porter.IndexConfig{
					Name: "custom-analyzer-index",
					Schema: porter.SchemaConfig{
						Settings: &porter.SettingsConfig{
							Analysis: &porter.AnalysisConfig{
								Analyzer: s.Migrator.Components.CreateAnalyzer(
									s.Migrator.Components.Settings.Analysis.Analyzer.Custom(
										"my_custom",
										s.Migrator.Components.Settings.Analysis.Analyzer.Custom.WithTokenizer("standard"),
										s.Migrator.Components.Settings.Analysis.Analyzer.Custom.WithCharFilter([]string{"html_strip"}),
										s.Migrator.Components.Settings.Analysis.Analyzer.Custom.WithFilters([]string{"lowercase", "asciifolding"}),
									),
								),
							},
						},
						Mappings: &porter.MappingsConfig{
							Properties: s.Migrator.Components.CreateProperties(
								s.Migrator.Components.Mappings.Properties.Text(
									"content",
									s.Migrator.Components.Mappings.Properties.Text.WithAnalyzer("my_custom"),
								),
							),
						},
					},
				},
				MigrationsPath: "",
			},
			to:       porter.LevelIndex,
			expected: nil,
		},
	}

	for _, cs := range cases {
		s.T.Run(cs.name, func(t *testing.T) {
			m := porter.New(s.Elasticsearch.Client, cs.in)

			err := m.MigrateUp(cs.to)
			if err != nil {
				assert.Equal(t, cs.expected, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
