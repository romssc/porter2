package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xoticdsign/porter"
	"github.com/xoticdsign/porter/internal/tests/suite"
)

func TestMigrateUp_Integration(t *testing.T) {
	s, err := suite.New(t, false)
	if err != nil {
		panic(err)
	}

	cases := []struct {
		name            string
		in              porter.Config
		inWithIndex     bool
		inWithDocuments bool
		expectedErr     bool
	}{
		{
			name: "happy case",
			in: porter.Config{
				Name: "porter_happy",
				Definition: porter.DefinitionConfig{
					Settings: &porter.SettingsConfig{
						NumberOfShards:   1,
						NumberOfReplicas: 1,
						Analysis: &porter.AnalysisConfig{
							Analyzer: s.Porter.Index.Settings.Analysis.NewAnalyzer(s.Porter.Index.Settings.Analysis.Analyzer.Simple("analyzer")),
							Normalizer: s.Porter.Index.Settings.Analysis.NewNormalizer(s.Porter.Index.Settings.Analysis.Normalizer.Custom(
								"normalyzer",
								s.Porter.Index.Settings.Analysis.Normalizer.Custom.WithFilter([]porter.NormalizerCustomFilter{porter.NormalizerCustomFilterASCIIFolding}),
							)),
						},
					},
					Mappings: &porter.MappingsConfig{
						Properties: s.Porter.Index.Mappings.NewFields(
							s.Porter.Index.Mappings.Properties.Keyword(
								"keyword",
								porter.FakeCity,
								s.Porter.Index.Mappings.Properties.Keyword.WithStore(true),
							),
							s.Porter.Index.Mappings.Properties.Integer(
								"integer",
								porter.FakeIntegerInt,
								s.Porter.Index.Mappings.Properties.Integer.WithStore(true),
							),
						),
					},
				},
			},
			inWithIndex:     true,
			inWithDocuments: true,
			expectedErr:     false,
		},
		{
			name:            "empty config case",
			in:              porter.Config{},
			inWithIndex:     true,
			inWithDocuments: true,
			expectedErr:     true,
		},
		{
			name: "invalid field mapping",
			in: porter.Config{
				Name: "porter_invalid_field",
				Definition: porter.DefinitionConfig{
					Settings: &porter.SettingsConfig{
						NumberOfShards:   1,
						NumberOfReplicas: 1,
					},
					Mappings: &porter.MappingsConfig{
						Properties: map[string]interface{}{
							"broken": map[string]interface{}{
								"type": "invalid",
							},
						},
					},
				},
			},
			inWithIndex:     true,
			inWithDocuments: false,
			expectedErr:     true,
		},
		{
			name: "no index and no documents case",
			in: porter.Config{
				Name: "porter_nothing",
			},
			inWithIndex:     false,
			inWithDocuments: false,
			expectedErr:     false,
		},
	}

	for _, cs := range cases {
		s.T.Run(cs.name, func(t *testing.T) {
			var err error

			switch {
			case cs.inWithIndex && cs.inWithDocuments:
				err = s.Porter.MigrateUp(cs.in, s.Porter.Index.MigrateIndex(), s.Porter.Documents.MigrateDocuments(s.Porter.Documents.Origin.Generate(10)))

			case !cs.inWithIndex && cs.inWithDocuments:
				err = s.Porter.MigrateUp(cs.in, s.Porter.Index.NoIndex(), s.Porter.Documents.MigrateDocuments(s.Porter.Documents.Origin.Generate(10)))

			case cs.inWithIndex && !cs.inWithDocuments:
				err = s.Porter.MigrateUp(cs.in, s.Porter.Index.MigrateIndex(), s.Porter.Documents.NoDocuments())

			case !cs.inWithIndex && !cs.inWithDocuments:
				err = s.Porter.MigrateUp(cs.in, s.Porter.Index.NoIndex(), s.Porter.Documents.NoDocuments())

			default:
				assert.Fail(t, "wrong combination of migrate functions")
			}

			switch {
			case cs.expectedErr:
				assert.Error(t, err)

			case !cs.expectedErr:
				assert.NoError(t, err)
			}
		})
	}
}

func TestMigrateDown_Integration(t *testing.T) {
	s, err := suite.New(t, false)
	if err != nil {
		panic(err)
	}

	cases := []struct {
		name            string
		in              porter.Config
		inWithIndex     bool
		inWithDocuments bool
		expectedErr     bool

		migratePrior bool
	}{
		{
			name: "happy case",
			in: porter.Config{
				Name: "porter_happy",
				Definition: porter.DefinitionConfig{
					Settings: &porter.SettingsConfig{
						NumberOfShards:   1,
						NumberOfReplicas: 1,
						Analysis: &porter.AnalysisConfig{
							Analyzer: s.Porter.Index.Settings.Analysis.NewAnalyzer(s.Porter.Index.Settings.Analysis.Analyzer.Simple("analyzer")),
							Normalizer: s.Porter.Index.Settings.Analysis.NewNormalizer(s.Porter.Index.Settings.Analysis.Normalizer.Custom(
								"normalyzer",
								s.Porter.Index.Settings.Analysis.Normalizer.Custom.WithFilter([]porter.NormalizerCustomFilter{porter.NormalizerCustomFilterASCIIFolding}),
							)),
						},
					},
					Mappings: &porter.MappingsConfig{
						Properties: s.Porter.Index.Mappings.NewFields(
							s.Porter.Index.Mappings.Properties.Keyword(
								"keyword",
								porter.FakeCity,
								s.Porter.Index.Mappings.Properties.Keyword.WithStore(true),
							),
							s.Porter.Index.Mappings.Properties.Integer(
								"integer",
								porter.FakeIntegerInt,
								s.Porter.Index.Mappings.Properties.Integer.WithStore(true),
							),
						),
					},
				},
			},
			inWithIndex:     true,
			inWithDocuments: true,
			expectedErr:     false,

			migratePrior: true,
		},
		{
			name: "no index and no documents case",
			in: porter.Config{
				Name: "porter_nothing",
			},
			inWithIndex:     false,
			inWithDocuments: false,
			expectedErr:     false,

			migratePrior: true,
		},
		{
			name: "non-existent index case",
			in: porter.Config{
				Name: "porter_non_existent",
			},
			inWithIndex:     true,
			inWithDocuments: true,
			expectedErr:     true,
		},
	}

	for _, cs := range cases {
		s.T.Run(cs.name, func(t *testing.T) {
			var err error

			if cs.migratePrior {
				err = s.Porter.MigrateUp(cs.in, s.Porter.Index.MigrateIndex(), s.Porter.Documents.MigrateDocuments(s.Porter.Documents.Origin.Generate(10)))
				if err != nil {
					assert.Equal(t, cs.expectedErr, err)
				}
			}

			switch {
			case cs.inWithIndex && cs.inWithDocuments:
				err = s.Porter.MigrateDown(cs.in, s.Porter.Documents.MigrateDocuments(nil), s.Porter.Index.MigrateIndex())

			case !cs.inWithIndex && cs.inWithDocuments:
				err = s.Porter.MigrateDown(cs.in, s.Porter.Documents.MigrateDocuments(nil), s.Porter.Index.NoIndex())

			case cs.inWithIndex && !cs.inWithDocuments:
				err = s.Porter.MigrateDown(cs.in, s.Porter.Documents.NoDocuments(), s.Porter.Index.MigrateIndex())

			case !cs.inWithIndex && !cs.inWithDocuments:
				err = s.Porter.MigrateDown(cs.in, s.Porter.Documents.NoDocuments(), s.Porter.Index.NoIndex())

			default:
				assert.Fail(t, "wrong combination of migrate functions")
			}

			switch {
			case cs.expectedErr:
				assert.Error(t, err, err)

			case !cs.expectedErr:
				assert.NoError(t, err, err)
			}
		})
	}
}
