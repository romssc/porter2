package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xoticdsign/porter"
	"github.com/xoticdsign/porter/internal/tests/suite"
)

func TestNewNormalizer_Functional(t *testing.T) {
	s, err := suite.New(t, true)
	if err != nil {
		panic(err)
	}

	cases := []struct {
		name     string
		in       porter.NormalizerFunc
		expected map[string]interface{}
	}{
		{
			name: "only filter case",
			in: s.Porter.Index.Settings.Analysis.Normalizer.Custom(
				"normalizer",
				s.Porter.Index.Settings.Analysis.Normalizer.Custom.WithFilter([]porter.NormalizerCustomFilter{
					porter.NormalizerCustomFilterASCIIFolding,
					porter.NormalizerCustomFilterLowercase,
				}),
			),
			expected: map[string]interface{}{
				"normalizer": map[string]interface{}{
					"type": "custom",
					"filter": []porter.NormalizerCustomFilter{
						porter.NormalizerCustomFilterASCIIFolding,
						porter.NormalizerCustomFilterLowercase,
					},
				},
			},
		},
		{
			name: "only char filter case",
			in: s.Porter.Index.Settings.Analysis.Normalizer.Custom(
				"normalizer",
				s.Porter.Index.Settings.Analysis.Normalizer.Custom.WithCharFilter([]porter.NormalizerCustomCharFilter{
					porter.NormalizerCustomCharFilterHTMLStrip,
				}),
			),
			expected: map[string]interface{}{
				"normalizer": map[string]interface{}{
					"type": "custom",
					"char_filter": []porter.NormalizerCustomCharFilter{
						porter.NormalizerCustomCharFilterHTMLStrip,
					},
				},
			},
		},
		{
			name: "filter and char filter case",
			in: s.Porter.Index.Settings.Analysis.Normalizer.Custom(
				"normalizer",
				s.Porter.Index.Settings.Analysis.Normalizer.Custom.WithFilter([]porter.NormalizerCustomFilter{
					porter.NormalizerCustomFilterUppercase,
				}),
				s.Porter.Index.Settings.Analysis.Normalizer.Custom.WithCharFilter([]porter.NormalizerCustomCharFilter{
					porter.NormalizerCustomCharFilterMapping,
				}),
			),
			expected: map[string]interface{}{
				"normalizer": map[string]interface{}{
					"type": "custom",
					"filter": []porter.NormalizerCustomFilter{
						porter.NormalizerCustomFilterUppercase,
					},
					"char_filter": []porter.NormalizerCustomCharFilter{
						porter.NormalizerCustomCharFilterMapping,
					},
				},
			},
		},
		{
			name: "empty normalizer case",
			in:   s.Porter.Index.Settings.Analysis.Normalizer.Custom("normalizer"),
			expected: map[string]interface{}{
				"normalizer": map[string]interface{}{
					"type": "custom",
				},
			},
		},
	}

	for _, cs := range cases {
		s.T.Run(cs.name, func(t *testing.T) {
			n := s.Porter.Index.Settings.Analysis.NewNormalizer(cs.in)

			assert.Equal(t, cs.expected, n)
		})
	}
}
