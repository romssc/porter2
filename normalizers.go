package porter

/*

This file includes the core definitions for building custom normalizers using a fluent and
composable function-based API. Each normalizer is defined with a name and optional configuration
functions (e.g., WithCharFilter(), WithFilter()), producing a map structure compatible with Elasticsearch.

*/

type normalizer struct {
	// Custom normalizer.
	Custom normalizerCustom
}

// NewNormalizer() takes a composed normalizer function and converts it into a map[string]interface{} structure compatible with Elasticsearch.
func (a analysis) NewNormalizer(normalizer normalizerFunc) map[string]interface{} {
	r := map[string]interface{}{}

	for k, v := range normalizer() {
		r[k] = v
	}

	return r
}

type normalizerFunc func() map[string]interface{}

// CUSTOM NORMALIZER

type normalizerCustomProperties func() map[string]interface{}
type normalizerCustom func(name string, properties ...normalizerCustomProperties) normalizerFunc

func newNormalizerCustom() normalizerCustom {
	return func(name string, properties ...normalizerCustomProperties) normalizerFunc {
		return func() map[string]interface{} {
			r := map[string]interface{}{}

			for _, fn := range properties {
				if fn == nil {
					continue
				}
				for k, v := range fn() {
					r[k] = v
				}
			}

			r["type"] = "custom"

			return map[string]interface{}{
				name: r,
			}
		}
	}
}

type NormalizerCustomCharFilter string

var (
	NormalizerCustomCharFilterHTMLStrip      NormalizerCustomCharFilter = "html_strip"
	NormalizerCustomCharFilterMapping        NormalizerCustomCharFilter = "mapping"
	NormalizerCustomCharFilterPatternReplace NormalizerCustomCharFilter = "pattern_replace"
)

// WithCharFilter() defines a list of character filters for the custom normalizer.
func (c normalizerCustom) WithCharFilter(value []NormalizerCustomCharFilter) normalizerCustomProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"char_filter": value,
		}
	}
}

type NormalizerCustomFilter string

var (
	NormalizerCustomFilterArabicNormalization  NormalizerCustomFilter = "arabic_normalization"
	NormalizerCustomFilterASCIIFolding         NormalizerCustomFilter = "asciifolding"
	NormalizerCustomFilterBengaliNormalization NormalizerCustomFilter = "bengali_normalization"
	NormalizerCustomFilterCJKWidth             NormalizerCustomFilter = "cjk_width"
	NormalizerCustomFilterDecimalDigit         NormalizerCustomFilter = "decimal_digit"
	NormalizerCustomFilterElision              NormalizerCustomFilter = "elision"
	NormalizerCustomFilterGermanNormalization  NormalizerCustomFilter = "german_normalization"
	NormalizerCustomFilterHindiNormalization   NormalizerCustomFilter = "hindi_normalization"
	NormalizerCustomFilterIndicNormalization   NormalizerCustomFilter = "indic_normalization"
	NormalizerCustomFilterLowercase            NormalizerCustomFilter = "lowercase"
	NormalizerCustomFilterPatternReplace       NormalizerCustomFilter = "pattern_replace"
	NormalizerCustomFilterPersianNormalization NormalizerCustomFilter = "persian_normalization"
	NormalizerCustomFilterScandinavianFolding  NormalizerCustomFilter = "scandinavian_folding"
	NormalizerCustomFilterSerbianNormalization NormalizerCustomFilter = "serbian_normalization"
	NormalizerCustomFilterSoraniNormalization  NormalizerCustomFilter = "sorani_normalization"
	NormalizerCustomFilterTrim                 NormalizerCustomFilter = "trim"
	NormalizerCustomFilterUppercase            NormalizerCustomFilter = "uppercase"
)

// WithFilter() defines a list of token filters to be used in the custom normalizer.
func (c normalizerCustom) WithFilter(value []NormalizerCustomFilter) normalizerCustomProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"filter": value,
		}
	}
}
