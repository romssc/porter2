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
		expectedErr     error
	}{
		{
			name: "happy case",
			in: porter.Config{
				Name: "porter",
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
			expectedErr:     nil,
		},
		{
			name: "no index case",
			in: porter.Config{
				Name: "porter",
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
			inWithIndex:     false,
			inWithDocuments: true,
			expectedErr:     porter.ErrMigrateDocuments,
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

			if err != nil {
				assert.Equal(t, cs.expectedErr, err)
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
		expectedErr     error
	}{
		{
			name: "happy case",
			in: porter.Config{
				Name: "porter",
				Definition: porter.DefinitionConfig{
					Mappings: &porter.MappingsConfig{
						Properties: s.Porter.Index.Mappings.NewFields(
							s.Porter.Index.Mappings.Properties.Keyword("keyword", porter.FakeCity),
						),
					},
				},
			},
			inWithIndex:     true,
			inWithDocuments: true,
			expectedErr:     nil,
		},
	}

	for _, cs := range cases {
		s.T.Run(cs.name, func(t *testing.T) {
			var err error

			err = s.Porter.MigrateUp(cs.in, s.Porter.Index.MigrateIndex(), s.Porter.Documents.MigrateDocuments(s.Porter.Documents.Origin.Generate(20)))
			if err != nil {
				assert.Equal(t, cs.expectedErr, err)
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

			if err != nil {
				assert.Equal(t, cs.expectedErr, err)
			}
		})
	}
}
