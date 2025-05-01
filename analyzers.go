package porter

import (
	"strings"
)

/*

This file contains the full set of factory functions and types for constructing built-in and
custom analyzers. Each analyzer type is implemented as a composable function with optional
configuration helpers (e.g., WithStopwords(), WithPattern()).

The goal is to enable flexible, declarative definition of analyzers as reusable, type-safe Go
functions, which are ultimately compiled to Elasticsearch-compatible JSON objects.

*/

type analyzer struct {
	Standart    AnalyzerStandard
	Simple      AnalyzerSimple
	Whitespace  AnalyzerWhitespace
	Stop        AnalyzerStop
	Keyword     AnalyzerKeyword
	Pattern     AnalyzerPattern
	Language    AnalyzerLanguage
	Fingerprint AnalyzerFingerprint
	Custom      AnalyzerCustom
}

// NewAnalyzer() applies an AnalyzerFunc to return a map structure used in Elasticsearch settings.
func (a analysis) NewAnalyzer(analyzer AnalyzerFunc) map[string]interface{} {
	r := map[string]interface{}{}

	for k, v := range analyzer() {
		r[k] = v
	}

	return r
}

type AnalyzerFunc func() map[string]interface{}

// STANDART ANALYZER

type AnalyzerStandardProperties func() map[string]interface{}
type AnalyzerStandard func(name string, properties ...AnalyzerStandardProperties) AnalyzerFunc

func newAnalyzerStandard() AnalyzerStandard {
	return func(name string, properties ...AnalyzerStandardProperties) AnalyzerFunc {
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

			r["type"] = "standard"

			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithMaxTokenLength() sets the max_token_length for the standard analyzer.
func (s AnalyzerStandard) WithMaxTokenLength(value int) AnalyzerStandardProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"max_token_length": value,
		}
	}
}

// WithStopwords() sets the stopwords for the standard analyzer.
func (s AnalyzerStandard) WithStopwords(value []string) AnalyzerStandardProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords": value,
		}
	}
}

// WithStopwordsPath() sets an external stopwords file path.
func (s AnalyzerStandard) WithStopwordsPath(value string) AnalyzerStandardProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords_path": value,
		}
	}
}

// SIMPLE ANALYZER

type AnalyzerSimple func(name string) AnalyzerFunc

func newAnalyzerSimple() AnalyzerSimple {
	return func(name string) AnalyzerFunc {
		return func() map[string]interface{} {
			r := map[string]interface{}{}

			r["type"] = "simple"

			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WHITESPACE ANALYZER

type AnalyzerWhitespace func(name string) AnalyzerFunc

func newAnalyzerWhitespace() AnalyzerWhitespace {
	return func(name string) AnalyzerFunc {
		return func() map[string]interface{} {
			r := map[string]interface{}{}

			r["type"] = "whitespace"

			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// STOP ANALYZER

type AnalyzerStopProperties func() map[string]interface{}
type AnalyzerStop func(name string, properties ...AnalyzerStopProperties) AnalyzerFunc

func newAnalyzerStop() AnalyzerStop {
	return func(name string, properties ...AnalyzerStopProperties) AnalyzerFunc {
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

			r["type"] = "stop"

			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithStopwords() sets the stopwords for the standard analyzer.
func (s AnalyzerStop) WithStopwords(value []string) AnalyzerStopProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords": value,
		}
	}
}

// WithStopwordsPath() sets an external stopwords file path.
func (s AnalyzerStop) WithStopwordsPath(value string) AnalyzerStopProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords_path": value,
		}
	}
}

// KEYWORD ANALYZER

type AnalyzerKeyword func(name string) AnalyzerFunc

func newAnalyzerKeyword() AnalyzerKeyword {
	return func(name string) AnalyzerFunc {
		return func() map[string]interface{} {
			r := map[string]interface{}{}

			r["type"] = "keyword"

			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// PATTERN ANALYZER

type AnalyzerPatternProperties func() map[string]interface{}
type AnalyzerPattern func(name string, properties ...AnalyzerPatternProperties) AnalyzerFunc

func newAnalyzerPattern() AnalyzerPattern {
	return func(name string, properties ...AnalyzerPatternProperties) AnalyzerFunc {
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

			r["type"] = "pattern"

			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// AnalyzerPatternPattern defines common regex patterns used in pattern-based tokenizers.
type AnalyzerPatternPattern string

var (
	AnalyzerPatternPatternNonWord     AnalyzerPatternPattern = `\W+`
	AnalyzerPatternPatternWhitespace  AnalyzerPatternPattern = `\s+`
	AnalyzerPatternPatternComma       AnalyzerPatternPattern = `,`
	AnalyzerPatternPatternPipe        AnalyzerPatternPattern = `\|`
	AnalyzerPatternPatternDot         AnalyzerPatternPattern = `\.`
	AnalyzerPatternPatternCustomWords AnalyzerPatternPattern = `[\s,;:\.\-]+`
)

// WithPattern() sets the pattern used to tokenize text.
func (p AnalyzerPattern) WithPattern(value AnalyzerPatternPattern) AnalyzerPatternProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"pattern": value,
		}
	}
}

// AnalyzerPatternFlags represents Java-compatible regex flags for pattern tokenizers.
type AnalyzerPatternFlags string

var (
	AnalyzerPatternFlagsCaseInsensitive AnalyzerPatternFlags = "CASE_INSENSITIVE"
	AnalyzerPatternFlagsComments        AnalyzerPatternFlags = "COMMENTS"
	AnalyzerPatternFlagsDotAll          AnalyzerPatternFlags = "DOTALL"
	AnalyzerPatternFlagsMultiline       AnalyzerPatternFlags = "MULTILINE"
	AnalyzerPatternFlagsUnicodeCase     AnalyzerPatternFlags = "UNICODE_CASE"
	AnalyzerPatternFlagsUnixLines       AnalyzerPatternFlags = "UNIX_LINES"
)

// WithFlags() configures regex flags (e.g., CASE_INSENSITIVE).
func (p AnalyzerPattern) WithFlags(value []AnalyzerPatternFlags) AnalyzerPatternProperties {
	return func() map[string]interface{} {
		flags := []string{}

		for _, v := range value {
			flags = append(flags, string(v))
		}

		f := strings.Join(flags, "|")

		return map[string]interface{}{
			"flags": f,
		}
	}
}

// WithLowercase() sets whether text should be lowercased.
func (p AnalyzerPattern) WithLowercase(enable bool) AnalyzerPatternProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"lowercase": enable,
		}
	}
}

// WithStopwords() sets the stopwords for the standard analyzer.
func (p AnalyzerPattern) WithStopwords(value []string) AnalyzerPatternProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords": value,
		}
	}
}

// WithStopwordsPath() sets an external stopwords file path.
func (p AnalyzerPattern) WithStopwordsPath(value string) AnalyzerPatternProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords_path": value,
		}
	}
}

// LANGUAGE ANALYZER

type AnalyzerLanguageProperties func() map[string]interface{}
type AnalyzerLanguage func(name string, language AnalyzerLanguageLanguage, properties ...AnalyzerLanguageProperties) AnalyzerFunc

// AnalyzerLanguageLanguage defines language-specific analyzers provided by Elasticsearch.
type AnalyzerLanguageLanguage string

var (
	AnalyzerLanguageArabic     AnalyzerLanguageLanguage = "arabic"
	AnalyzerLanguageArmenian   AnalyzerLanguageLanguage = "armenian"
	AnalyzerLanguageBasque     AnalyzerLanguageLanguage = "basque"
	AnalyzerLanguageBengali    AnalyzerLanguageLanguage = "bengali"
	AnalyzerLanguageBrazilian  AnalyzerLanguageLanguage = "brazilian"
	AnalyzerLanguageBulgarian  AnalyzerLanguageLanguage = "bulgarian"
	AnalyzerLanguageCatalan    AnalyzerLanguageLanguage = "catalan"
	AnalyzerLanguageCJK        AnalyzerLanguageLanguage = "cjk"
	AnalyzerLanguageCzech      AnalyzerLanguageLanguage = "czech"
	AnalyzerLanguageDanish     AnalyzerLanguageLanguage = "danish"
	AnalyzerLanguageDutch      AnalyzerLanguageLanguage = "dutch"
	AnalyzerLanguageEnglish    AnalyzerLanguageLanguage = "english"
	AnalyzerLanguageEstonian   AnalyzerLanguageLanguage = "estonian"
	AnalyzerLanguageFinnish    AnalyzerLanguageLanguage = "finnish"
	AnalyzerLanguageFrench     AnalyzerLanguageLanguage = "french"
	AnalyzerLanguageGalician   AnalyzerLanguageLanguage = "galician"
	AnalyzerLanguageGerman     AnalyzerLanguageLanguage = "german"
	AnalyzerLanguageGreek      AnalyzerLanguageLanguage = "greek"
	AnalyzerLanguageHindi      AnalyzerLanguageLanguage = "hindi"
	AnalyzerLanguageHungarian  AnalyzerLanguageLanguage = "hungarian"
	AnalyzerLanguageIndonesian AnalyzerLanguageLanguage = "indonesian"
	AnalyzerLanguageIrish      AnalyzerLanguageLanguage = "irish"
	AnalyzerLanguageItalian    AnalyzerLanguageLanguage = "italian"
	AnalyzerLanguageLatvian    AnalyzerLanguageLanguage = "latvian"
	AnalyzerLanguageLithuanian AnalyzerLanguageLanguage = "lithuanian"
	AnalyzerLanguageNorwegian  AnalyzerLanguageLanguage = "norwegian"
	AnalyzerLanguagePersian    AnalyzerLanguageLanguage = "persian"
	AnalyzerLanguagePortuguese AnalyzerLanguageLanguage = "portuguese"
	AnalyzerLanguageRomanian   AnalyzerLanguageLanguage = "romanian"
	AnalyzerLanguageRussian    AnalyzerLanguageLanguage = "russian"
	AnalyzerLanguageSerbian    AnalyzerLanguageLanguage = "serbian"
	AnalyzerLanguageSorani     AnalyzerLanguageLanguage = "sorani"
	AnalyzerLanguageSpanish    AnalyzerLanguageLanguage = "spanish"
	AnalyzerLanguageSwedish    AnalyzerLanguageLanguage = "swedish"
	AnalyzerLanguageTurkish    AnalyzerLanguageLanguage = "turkish"
	AnalyzerLanguageThai       AnalyzerLanguageLanguage = "thai"
)

func newAnalyzerLanguage() AnalyzerLanguage {
	return func(name string, language AnalyzerLanguageLanguage, properties ...AnalyzerLanguageProperties) AnalyzerFunc {
		return func() map[string]interface{} {
			stemExclusionAllow := map[string]string{
				"arabic":     "",
				"armenian":   "",
				"basque":     "",
				"bengali":    "",
				"bulgarian":  "",
				"catalan":    "",
				"czech":      "",
				"dutch":      "",
				"english":    "",
				"finnish":    "",
				"french":     "",
				"galician":   "",
				"german":     "",
				"hindi":      "",
				"hungarian":  "",
				"indonesian": "",
				"irish":      "",
				"italian":    "",
				"latvian":    "",
				"lithuanian": "",
				"norwegian":  "",
				"portuguese": "",
				"romanian":   "",
				"russian":    "",
				"serbian":    "",
				"sorani":     "",
				"spanish":    "",
				"swedish":    "",
				"turkish":    "",
			}

			r := map[string]interface{}{}

			for _, fn := range properties {
				if fn == nil {
					continue
				}
				for k, v := range fn() {
					if k == "stem_exclusion" {
						_, ok := stemExclusionAllow[string(language)]
						if !ok {
							continue
						}
					}

					r[k] = v
				}
			}

			r["type"] = string(language)

			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithStopwords() sets the stopwords for the standard analyzer.
func (l AnalyzerLanguage) WithStopwords(value []string) AnalyzerLanguageProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords": value,
		}
	}
}

// WithStopwordsPath() sets an external stopwords file path.
func (l AnalyzerLanguage) WithStopwordsPath(value string) AnalyzerLanguageProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords_path": value,
		}
	}
}

// WithStemExclusion() defines a list of terms that should not be stemmed during analysis.
func (l AnalyzerLanguage) WithStemExclusion(value []string) AnalyzerLanguageProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stem_exclusion": value,
		}
	}
}

// FINGERPRINT ANALYZER

type AnalyzerFingerprintProperties func() map[string]interface{}
type AnalyzerFingerprint func(name string, properties ...AnalyzerFingerprintProperties) AnalyzerFunc

func newAnalyzerFingerprint() AnalyzerFingerprint {
	return func(name string, properties ...AnalyzerFingerprintProperties) AnalyzerFunc {
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

			r["type"] = "fingerprint"

			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithSeparator() sets the string used to join tokens into a fingerprint.
func (f AnalyzerFingerprint) WithSeparator(value string) AnalyzerFingerprintProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"separator": value,
		}
	}
}

// WithMaxOutputSize() sets max length of the resulting fingerprint string.
func (f AnalyzerFingerprint) WithMaxOutputSize(value int) AnalyzerFingerprintProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"max_output_size": value,
		}
	}
}

// WithStopwords() sets the stopwords for the standard analyzer.
func (f AnalyzerFingerprint) WithStopwords(value []string) AnalyzerFingerprintProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords": value,
		}
	}
}

// WithStopwordsPath() sets an external stopwords file path.
func (f AnalyzerFingerprint) WithStopwordsPath(value string) AnalyzerFingerprintProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords_path": value,
		}
	}
}

// STOP ANALYZER

type AnalyzerCustomProperties func() map[string]interface{}
type AnalyzerCustom func(name string, properties ...AnalyzerCustomProperties) AnalyzerFunc

func newAnalyzerCustom() AnalyzerCustom {
	return func(name string, properties ...AnalyzerCustomProperties) AnalyzerFunc {
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

// AnalyzerCustomTokenizer defines available tokenizer types for custom analyzers.
type AnalyzerCustomTokenizer string

var (
	AnalyzerCustomTokenizerStandard           AnalyzerCustomTokenizer = "standard"
	AnalyzerCustomTokenizerLetter             AnalyzerCustomTokenizer = "letter"
	AnalyzerCustomTokenizerLowercase          AnalyzerCustomTokenizer = "lowercase"
	AnalyzerCustomTokenizerWhitespace         AnalyzerCustomTokenizer = "whitespace"
	AnalyzerCustomTokenizerUAXURLEmail        AnalyzerCustomTokenizer = "uax_url_email"
	AnalyzerCustomTokenizerClassic            AnalyzerCustomTokenizer = "classic"
	AnalyzerCustomTokenizerThai               AnalyzerCustomTokenizer = "thai"
	AnalyzerCustomTokenizerNGram              AnalyzerCustomTokenizer = "ngram"
	AnalyzerCustomTokenizerEdgeNGram          AnalyzerCustomTokenizer = "edge_ngram"
	AnalyzerCustomTokenizerKeyword            AnalyzerCustomTokenizer = "keyword"
	AnalyzerCustomTokenizerPattern            AnalyzerCustomTokenizer = "pattern"
	AnalyzerCustomTokenizerSimplePattern      AnalyzerCustomTokenizer = "simple_pattern"
	AnalyzerCustomTokenizerCharGroup          AnalyzerCustomTokenizer = "char_group"
	AnalyzerCustomTokenizerSimplePatternSplit AnalyzerCustomTokenizer = "simple_pattern_split"
	AnalyzerCustomTokenizerPathHierarchy      AnalyzerCustomTokenizer = "path_hierarchy"
)

// WithTokenizer() sets the tokenizer for a custom analyzer.
func (c AnalyzerCustom) WithTokenizer(value AnalyzerCustomTokenizer) AnalyzerCustomProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"tokenizer": value,
		}
	}
}

// AnalyzerCustomCharFilter defines character filters that modify text before tokenization.
type AnalyzerCustomCharFilter string

var (
	AnalyzerCustomCharFilterHTMLStrip      AnalyzerCustomCharFilter = "html_strip"
	AnalyzerCustomCharFilterMapping        AnalyzerCustomCharFilter = "mapping"
	AnalyzerCustomCharFilterPatternReplace AnalyzerCustomCharFilter = "pattern_replace"
)

// WithCharFilter() adds one or more character filters to a custom analyzer.
func (c AnalyzerCustom) WithCharFilter(value []AnalyzerCustomCharFilter) AnalyzerCustomProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"char_filter": value,
		}
	}
}

// AnalyzerCustomFilter defines token filters that modify or enhance token streams after tokenization.
type AnalyzerCustomFilter string

var (
	AnalyzerCustomFilterApostrophe            AnalyzerCustomFilter = "apostrophe"
	AnalyzerCustomFilterASCIIFolding          AnalyzerCustomFilter = "asciifolding"
	AnalyzerCustomFilterCJKBigram             AnalyzerCustomFilter = "cjk_bigram"
	AnalyzerCustomFilterCJKWidth              AnalyzerCustomFilter = "cjk_width"
	AnalyzerCustomFilterClassic               AnalyzerCustomFilter = "classic"
	AnalyzerCustomFilterCommonGrams           AnalyzerCustomFilter = "common_grams"
	AnalyzerCustomFilterConditional           AnalyzerCustomFilter = "condition"
	AnalyzerCustomFilterDecimalDigit          AnalyzerCustomFilter = "decimal_digit"
	AnalyzerCustomFilterDelimitedPayload      AnalyzerCustomFilter = "delimited_payload"
	AnalyzerCustomFilterDictionaryDecompound  AnalyzerCustomFilter = "dictionary_decompounder"
	AnalyzerCustomFilterEdgeNGram             AnalyzerCustomFilter = "edge_ngram"
	AnalyzerCustomFilterElision               AnalyzerCustomFilter = "elision"
	AnalyzerCustomFilterFingerprint           AnalyzerCustomFilter = "fingerprint"
	AnalyzerCustomFilterFlattenGraph          AnalyzerCustomFilter = "flatten_graph"
	AnalyzerCustomFilterHunspell              AnalyzerCustomFilter = "hunspell"
	AnalyzerCustomFilterHyphenationDecompound AnalyzerCustomFilter = "hyphenation_decompounder"
	AnalyzerCustomFilterKeepTypes             AnalyzerCustomFilter = "keep_types"
	AnalyzerCustomFilterKeepWords             AnalyzerCustomFilter = "keep_words"
	AnalyzerCustomFilterKeywordMarker         AnalyzerCustomFilter = "keyword_marker"
	AnalyzerCustomFilterKeywordRepeat         AnalyzerCustomFilter = "keyword_repeat"
	AnalyzerCustomFilterKStem                 AnalyzerCustomFilter = "kstem"
	AnalyzerCustomFilterLength                AnalyzerCustomFilter = "length"
	AnalyzerCustomFilterLimitTokenCount       AnalyzerCustomFilter = "limit"
	AnalyzerCustomFilterLowercase             AnalyzerCustomFilter = "lowercase"
	AnalyzerCustomFilterMinHash               AnalyzerCustomFilter = "min_hash"
	AnalyzerCustomFilterMultiplexer           AnalyzerCustomFilter = "multiplexer"
	AnalyzerCustomFilterNGram                 AnalyzerCustomFilter = "ngram"
	AnalyzerCustomFilterNormalization         AnalyzerCustomFilter = "normalization"
	AnalyzerCustomFilterPatternCapture        AnalyzerCustomFilter = "pattern_capture"
	AnalyzerCustomFilterPatternReplace        AnalyzerCustomFilter = "pattern_replace"
	AnalyzerCustomFilterPhonetic              AnalyzerCustomFilter = "phonetic"
	AnalyzerCustomFilterPorterStem            AnalyzerCustomFilter = "porter_stem"
	AnalyzerCustomFilterPredicateScript       AnalyzerCustomFilter = "predicate_script"
	AnalyzerCustomFilterRemoveDuplicates      AnalyzerCustomFilter = "remove_duplicates"
	AnalyzerCustomFilterReverse               AnalyzerCustomFilter = "reverse"
	AnalyzerCustomFilterShingle               AnalyzerCustomFilter = "shingle"
	AnalyzerCustomFilterSnowball              AnalyzerCustomFilter = "snowball"
	AnalyzerCustomFilterStemmer               AnalyzerCustomFilter = "stemmer"
	AnalyzerCustomFilterStemmerOverride       AnalyzerCustomFilter = "stemmer_override"
	AnalyzerCustomFilterStop                  AnalyzerCustomFilter = "stop"
	AnalyzerCustomFilterSynonym               AnalyzerCustomFilter = "synonym"
	AnalyzerCustomFilterSynonymGraph          AnalyzerCustomFilter = "synonym_graph"
	AnalyzerCustomFilterTrim                  AnalyzerCustomFilter = "trim"
	AnalyzerCustomFilterTruncate              AnalyzerCustomFilter = "truncate"
	AnalyzerCustomFilterUnique                AnalyzerCustomFilter = "unique"
	AnalyzerCustomFilterUppercase             AnalyzerCustomFilter = "uppercase"
	AnalyzerCustomFilterWordDelimiter         AnalyzerCustomFilter = "word_delimiter"
	AnalyzerCustomFilterWordDelimiterGraph    AnalyzerCustomFilter = "word_delimiter_graph"
)

// WithFilter() sets the list of token filters to apply to the output of the tokenizer.
func (c AnalyzerCustom) WithFilter(value []AnalyzerCustomFilter) AnalyzerCustomProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"filter": value,
		}
	}
}

// WithPositionIncrementGap() sets the position increment gap.
func (c AnalyzerCustom) WithPositionIncrementGap(value int) AnalyzerCustomProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"position_increment_gap": value,
		}
	}
}
