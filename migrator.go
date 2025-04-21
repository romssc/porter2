package porter

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/xoticdsign/porter/internal/lib"
	"github.com/xoticdsign/porter/internal/utils"
)

var (
	ErrElasticsearchConnection    = fmt.Errorf("can't connect to elasticsearch")
	ErrIndexNameNotSpecified      = fmt.Errorf("index name must be specified")
	ErrMigrationsPathNotSpecified = fmt.Errorf("migrations path must be specified")
	ErrWrongLevel                 = fmt.Errorf("level can't be equal to 0 or greater than 2")
	ErrMigrationsNotFound         = fmt.Errorf("migrations not found")
	ErrUnexpected                 = fmt.Errorf("unexpected")
	ErrExistsCheck                = fmt.Errorf("can't check index existence")
	ErrCreateQuery                = fmt.Errorf("can't create an index")
	ErrBulkQuery                  = fmt.Errorf("can't bulk insert")
	ErrIndicesDelete              = fmt.Errorf("can't delete indice")
	ErrDocumentsDelete            = fmt.Errorf("can't delete documents")
	ErrIndexAlreadyExists         = fmt.Errorf("index already exists")
	ErrUnexpectedStatusCode       = fmt.Errorf("unxexpected status code received")
)

var (
	levelZero      = 0
	LevelIndex     = 1
	LevelDocuments = 2
)

var (
	AnalyzerLanguageSpecificArabic     = "arabic"
	AnalyzerLanguageSpecificCatalan    = "catalan"
	AnalyzerLanguageSpecificDanish     = "danish"
	AnalyzerLanguageSpecificDutch      = "dutch"
	AnalyzerLanguageSpecificEnglish    = "english"
	AnalyzerLanguageSpecificFinnish    = "finnish"
	AnalyzerLanguageSpecificFrench     = "french"
	AnalyzerLanguageSpecificGerman     = "german"
	AnalyzerLanguageSpecificGreek      = "greek"
	AnalyzerLanguageSpecificHungarian  = "hungarian"
	AnalyzerLanguageSpecificItalian    = "italian"
	AnalyzerLanguageSpecificLatvian    = "latvian"
	AnalyzerLanguageSpecificNorwegian  = "norwegian"
	AnalyzerLanguageSpecificPortuguese = "portuguese"
	AnalyzerLanguageSpecificRomanian   = "romanian"
	AnalyzerLanguageSpecificRussian    = "russian"
	AnalyzerLanguageSpecificSpanish    = "spanish"
	AnalyzerLanguageSpecificSwedish    = "swedish"
	AnalyzerLanguageSpecificTurkish    = "turkish"
)

/*






 */

type M struct {
	client searcher
	config Config
}

type searcher interface {
	IsIndexExist(ctx context.Context, name string) (*esapi.Response, error)
	CreateIndex(ctx context.Context, name string, body []byte) (*esapi.Response, error)
	CreateDocuments(ctx context.Context, name string, documents []byte) (*esapi.Response, error)
	DeleteIndex(ctx context.Context, name string) (*esapi.Response, error)
	DeleteDocuments(ctx context.Context, name string, query string) (*esapi.Response, error)
}

type client struct {
	*elasticsearch.Client
}

func (c client) IsIndexExist(ctx context.Context, name string) (*esapi.Response, error) {
	resp, err := c.Indices.Exists(
		[]string{name},
		c.Indices.Exists.WithContext(context.Background()),
		c.Indices.Exists.WithPretty(),
	)
	return resp, err
}

func (c client) CreateIndex(ctx context.Context, name string, body []byte) (*esapi.Response, error) {
	resp, err := c.Indices.Create(
		name,
		c.Indices.Create.WithContext(ctx),
		c.Indices.Create.WithBody(bytes.NewBuffer(body)),
		c.Indices.Create.WithPretty(),
	)
	return resp, err
}

func (c client) CreateDocuments(ctx context.Context, name string, documents []byte) (*esapi.Response, error) {
	resp, err := c.Bulk(
		bytes.NewBuffer(documents),
		c.Bulk.WithContext(context.Background()),
		c.Bulk.WithIndex(name),
		c.Bulk.WithPretty(),
	)
	return resp, err
}

func (c client) DeleteIndex(ctx context.Context, name string) (*esapi.Response, error) {
	resp, err := c.Indices.Delete(
		[]string{name},
		c.Indices.Delete.WithContext(context.Background()),
		c.Indices.Delete.WithPretty(),
	)
	return resp, err
}

func (c client) DeleteDocuments(ctx context.Context, name string, query string) (*esapi.Response, error) {
	resp, err := c.DeleteByQuery(
		[]string{name},
		strings.NewReader(query),
		c.DeleteByQuery.WithContext(context.Background()),
		c.DeleteByQuery.WithPretty(),
	)
	return resp, err
}

type Config struct {
	IndexDefinition IndexConfig
	MigrationsPath  string
}

type IndexConfig struct {
	Name   string
	Schema SchemaConfig
}

type SchemaConfig struct {
	Settings *SettingsConfig `json:"settings,omitempty"`
	Mappings *MappingsConfig `json:"mappings,omitempty"`
}

type SettingsConfig struct {
	NumberOfShards   int             `json:"number_of_shards,omitempty"`
	NumberOfReplicas int             `json:"number_of_replicas,omitempty"`
	Analysis         *AnalysisConfig `json:"analysis,omitempty"`
}

type AnalysisConfig struct {
	Analyzer   map[string]interface{} `json:"analyzer,omitempty"`
	Normalizer map[string]interface{} `json:"normalizer,omitempty"`
}

type MappingsConfig struct {
	Properties map[string]interface{} `json:"properties,omitempty"`
}

func New(cc *elasticsearch.Client, config Config) M {
	return M{
		client: client{cc},
		config: config,
	}
}

func (m M) MigrateUp(to int) error {
	if to == levelZero || to > LevelDocuments {
		return fmt.Errorf("%w\n\nwhat to do? -> use one of the provided level types", ErrWrongLevel)
	}

	if to >= LevelIndex {
		respExists, err := m.client.IsIndexExist(context.Background(), m.config.IndexDefinition.Name)
		if err != nil {
			return fmt.Errorf("%w\n\nerror -> %v", ErrElasticsearchConnection, err)
		}
		defer respExists.Body.Close()

		switch respExists.StatusCode {
		case http.StatusNotFound:
			respCreate, err := m.client.CreateIndex(context.Background(), m.config.IndexDefinition.Name, utils.MarshalJSON(m.config.IndexDefinition.Schema))
			if err != nil {
				return fmt.Errorf("%w\n\nerror -> %v", ErrElasticsearchConnection, err)
			}
			defer respCreate.Body.Close()

			r, ok := utils.ExtractError(respCreate.Body)
			if ok {
				return fmt.Errorf("%w\n\nerror -> %v", ErrCreateQuery, r)
			}

		case http.StatusOK:
			if to == LevelIndex {
				return fmt.Errorf("%w\n\nstatus -> %v", ErrIndexAlreadyExists, respExists.StatusCode)
			}

		default:
			return fmt.Errorf("%w\n\nstatus -> %v", ErrUnexpectedStatusCode, respExists.StatusCode)
		}
	}

	if to == LevelDocuments {
		migrations, err := utils.GetContents(m.config.MigrationsPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("%w\n\nerror -> %v", ErrMigrationsNotFound, err)
			}
			return fmt.Errorf("%w\n\nerror -> %v", ErrUnexpected, err)
		}
		respBulk, err := m.client.CreateDocuments(context.Background(), m.config.IndexDefinition.Name, migrations)
		if err != nil {
			return fmt.Errorf("%w\n\nerror -> %v", ErrElasticsearchConnection, err)
		}
		defer respBulk.Body.Close()

		r, ok := utils.ExtractBulkErrors(respBulk.Body)
		if ok {
			return fmt.Errorf("%w\n\nerror -> %v", ErrBulkQuery, r)
		}

		switch respBulk.StatusCode {
		case http.StatusOK:

		default:
			return fmt.Errorf("%w\n\nstatus -> %v", ErrUnexpectedStatusCode, respBulk.StatusCode)
		}
	}

	return nil
}

func (m M) MigrateDown(to int) error {
	if to == levelZero || to > LevelDocuments {
		return fmt.Errorf("%w\n\nwhat to do? -> use one of the provided level types", ErrWrongLevel)
	}

	if to <= LevelDocuments {
		respDocumentsDelete, err := m.client.DeleteDocuments(context.Background(), m.config.IndexDefinition.Name, `{"query": {"match_all": {}}}`)
		if err != nil {
			return fmt.Errorf("%w\n\nerror -> %v", ErrElasticsearchConnection, err)
		}
		defer respDocumentsDelete.Body.Close()

		r, ok := utils.ExtractError(respDocumentsDelete.Body)
		if ok {
			return fmt.Errorf("%w\n\nerror -> %v", ErrIndicesDelete, r)
		}
	}

	if to == LevelDocuments {
		return nil
	}

	if to == LevelIndex {
		respIndicesDelete, err := m.client.DeleteIndex(context.Background(), m.config.IndexDefinition.Name)
		if err != nil {
			return fmt.Errorf("%w\n\nerror -> %v", ErrElasticsearchConnection, err)
		}
		defer respIndicesDelete.Body.Close()

		r, ok := utils.ExtractError(respIndicesDelete.Body)
		if ok {
			return fmt.Errorf("%w\n\nerror -> %v", ErrDocumentsDelete, r)
		}
	}

	return nil
}

/*






 */

type C struct {
	Settings settings
	Mappings mappings
}

type settings struct {
	Analysis analysis
}

type analysis struct {
	Analyzer   analyzer
	Normalizer normalizer
}

type analyzer struct {
	GeneralPurpose   generalPurposeAnalyzers
	LanguageSpecific languageSpecificAnalyzer
	Custom           customAnalyzer
}

type generalPurposeAnalyzers struct {
	Standart   standartAnalyzer
	Whitespace whitespaceAnalyzer
	Stop       stopAnalyzer
	Keyword    keywordAnalyzer
	Pattern    patternAnalyzer
	EngeNGram  edgeNGramAnalyzer
}

type normalizer struct {
	Predefined predefinedNormalizers
	Custom     customNormalizer
}

type predefinedNormalizers struct {
	Lowercase    lowercaseNormalizer
	ASCIIFolding asciiFoldingNormalizer
	Keyword      keywordNormalizer
}

type mappings struct {
	Properties properties
}

type properties struct {
	fieldTypes
}

type fieldTypes struct {
	Keyword      keyword
	Text         text
	Integer      integer
	Long         long
	Float        float
	Double       double
	Short        short
	Byte         bute
	HalfFloat    halfFloat
	ScaledFloat  scaledFloat
	Date         date
	DateNanos    dateNanos
	Boolean      boolean
	Binary       binary
	IP           ip
	GeoPoint     geoPoint
	GeoShape     geoShape
	IntegerRange integerRange
	FloatRange   floatRange
	DateRange    dateRange
}

func GetComponents() C {
	return C{
		Settings: settings{
			Analysis: analysis{
				Analyzer: analyzer{
					GeneralPurpose: generalPurposeAnalyzers{
						Standart:   newStandartAnalyzer(),
						Whitespace: newWhitespaceAnalyzer(),
						Stop:       newStopAnalyzer(),
						Keyword:    newKeywordAnalyzer(),
						Pattern:    newPatternAnalyzer(),
						EngeNGram:  newEdgeNGramAnalyzer(),
					},
					LanguageSpecific: newLanguageSpecificAnalyzer(),
					Custom:           newCustomAnalyzer(),
				},
				Normalizer: normalizer{
					Predefined: predefinedNormalizers{
						Lowercase:    newLowercaseNormalizer(),
						ASCIIFolding: newASCIIFoldingNormalizer(),
						Keyword:      newKeywordNormalizer(),
					},
					Custom: newCustomNormalizer(),
				},
			},
		},
		Mappings: mappings{
			Properties: properties{
				fieldTypes: fieldTypes{
					Keyword:      newKeyword(),
					Text:         newText(),
					Integer:      newInteger(),
					Long:         newLong(),
					Float:        newFloat(),
					Double:       newDouble(),
					Short:        newShort(),
					Byte:         newByte(),
					HalfFloat:    newHalfFloat(),
					ScaledFloat:  newScaledFloat(),
					Date:         newDate(),
					DateNanos:    newDateNanos(),
					Boolean:      newBoolean(),
					Binary:       newBinary(),
					IP:           newIP(),
					GeoPoint:     newGeoPoint(),
					GeoShape:     newGeoShape(),
					IntegerRange: newIntegerRange(),
					FloatRange:   newFloatRange(),
					DateRange:    newDateRange(),
				},
			},
		},
	}
}

func (c C) CreateAnalyzer(analyzer lib.AnalyzerFunc) map[string]interface{} {
	return utils.GetMap([]utils.MapperFunc{analyzer})
}

func (c C) CreateNormalizer(normalizer lib.NormalizerFunc) map[string]interface{} {
	return utils.GetMap([]utils.MapperFunc{normalizer})
}

func (c C) CreateProperties(fields ...lib.PropertiesFunc) map[string]interface{} {
	return utils.GetMap(fields)
}

/*






 */

/*

 ANALYZERS

*/

// STANDART ANALYZER TYPE

type standartAnalyzer func(name string, properties ...lib.StandardAnalyzerPropertiesFunc) lib.AnalyzerFunc

func newStandartAnalyzer() standartAnalyzer {
	return func(name string, properties ...lib.StandardAnalyzerPropertiesFunc) lib.AnalyzerFunc {
		m := utils.GetMap(properties)

		m["type"] = "standard"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (s standartAnalyzer) WithMaxTokenLength(value int) lib.StandardAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"max_token_length": value,
		}
	}
}

func (s standartAnalyzer) WithStopwords(value []string) lib.StandardAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords": value,
		}
	}
}

func (s standartAnalyzer) WithStopwordsPath(value string) lib.StandardAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords_path": value,
		}
	}
}

// WHITESPACE ANALYZER TYPE

type whitespaceAnalyzer func(name string, properties ...lib.WhitespaceAnalyzerPropertiesFunc) lib.AnalyzerFunc

func newWhitespaceAnalyzer() whitespaceAnalyzer {
	return func(name string, properties ...lib.WhitespaceAnalyzerPropertiesFunc) lib.AnalyzerFunc {
		m := utils.GetMap(properties)
		m["type"] = "whitespace"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (w whitespaceAnalyzer) WithMaxTokenLength(value int) lib.WhitespaceAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"max_token_length": value,
		}
	}
}

// STOP ANALYZER TYPE

type stopAnalyzer func(name string, properties ...lib.StopAnalyzerPropertiesFunc) lib.AnalyzerFunc

func newStopAnalyzer() stopAnalyzer {
	return func(name string, properties ...lib.StopAnalyzerPropertiesFunc) lib.AnalyzerFunc {
		m := utils.GetMap(properties)
		m["type"] = "stop"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (s stopAnalyzer) WithStopwords(value []string) lib.StopAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords": value,
		}
	}
}

func (s stopAnalyzer) WithStopwordsPath(value string) lib.StopAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords_path": value,
		}
	}
}

// KEYWORD ANALYZER TYPE

type keywordAnalyzer func(name string, properties ...lib.KeywordAnalyzerPropertiesFunc) lib.AnalyzerFunc

func newKeywordAnalyzer() keywordAnalyzer {
	return func(name string, properties ...lib.KeywordAnalyzerPropertiesFunc) lib.AnalyzerFunc {
		m := utils.GetMap(properties)
		m["type"] = "keyword"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (k keywordAnalyzer) WithIgnoreAbove(value int) lib.KeywordAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_above": value,
		}
	}
}

// PATTERN ANALYZER TYPE

type patternAnalyzer func(name string, properties ...lib.PatternAnalyzerPropertiesFunc) lib.AnalyzerFunc

func newPatternAnalyzer() patternAnalyzer {
	return func(name string, properties ...lib.PatternAnalyzerPropertiesFunc) lib.AnalyzerFunc {
		m := utils.GetMap(properties)
		m["type"] = "pattern"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (p patternAnalyzer) WithPattern(value string) lib.PatternAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"pattern": value,
		}
	}
}

func (p patternAnalyzer) WithGroupName(value string) lib.PatternAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"group_name": value,
		}
	}
}

// EDGENGRAM ANALYZER TYPE

type edgeNGramAnalyzer func(name string, properties ...lib.EdgeNGramAnalyzerPropertiesFunc) lib.AnalyzerFunc

func newEdgeNGramAnalyzer() edgeNGramAnalyzer {
	return func(name string, properties ...lib.EdgeNGramAnalyzerPropertiesFunc) lib.AnalyzerFunc {
		m := utils.GetMap(properties)
		m["type"] = "edge_ngram"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (e edgeNGramAnalyzer) WithMinGram(value int) lib.EdgeNGramAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"min_gram": value,
		}
	}
}

func (e edgeNGramAnalyzer) WithMaxGram(value int) lib.EdgeNGramAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"max_gram": value,
		}
	}
}

// CUSTOM ANALYZER TYPE

type customAnalyzer func(name string, properties ...lib.CustomAnalyzerPropertiesFunc) lib.AnalyzerFunc

func newCustomAnalyzer() customAnalyzer {
	return func(name string, properties ...lib.CustomAnalyzerPropertiesFunc) lib.AnalyzerFunc {
		m := utils.GetMap(properties)
		m["type"] = "custom"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (c customAnalyzer) WithTokenizer(value string) lib.CustomAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"tokenizer": value,
		}
	}
}

func (c customAnalyzer) WithFilters(value []string) lib.CustomAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"filter": value,
		}
	}
}

func (c customAnalyzer) WithCharFilter(value []string) lib.CustomAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"char_filter": value,
		}
	}
}

// LANGUAGESPECIFIC ANALYZER TYPE

type languageSpecificAnalyzer func(name string, language string, properties ...lib.LanguageSpecificAnalyzerPropertiesFunc) lib.AnalyzerFunc

func newLanguageSpecificAnalyzer() languageSpecificAnalyzer {
	return func(name string, language string, properties ...lib.LanguageSpecificAnalyzerPropertiesFunc) lib.AnalyzerFunc {
		m := utils.GetMap(properties)
		m["type"] = language

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (l languageSpecificAnalyzer) WithStopwords(value []string) lib.LanguageSpecificAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stopwords": value,
		}
	}
}

func (l languageSpecificAnalyzer) WithStemmer(value string) lib.LanguageSpecificAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"stemmer": value,
		}
	}
}

func (l languageSpecificAnalyzer) WithCharFilter(value []string) lib.LanguageSpecificAnalyzerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"char_filter": value,
		}
	}
}

/*

NORMALIZERS

*/

// LOWERCASE NORMALIZER TYPE

type lowercaseNormalizer func(name string) lib.NormalizerFunc

func newLowercaseNormalizer() lowercaseNormalizer {
	return func(name string) lib.NormalizerFunc {
		m := map[string]interface{}{
			"type": "lowercase",
		}

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

// ASCIIFOLDING NORMALIZER TYPE

type asciiFoldingNormalizer func(name string) lib.NormalizerFunc

func newASCIIFoldingNormalizer() asciiFoldingNormalizer {
	return func(name string) lib.NormalizerFunc {
		m := map[string]interface{}{
			"type": "asciifolding",
		}

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

// KEYWORD NORMALIZER TYPE

type keywordNormalizer func(name string) lib.NormalizerFunc

func newKeywordNormalizer() keywordNormalizer {
	return func(name string) lib.NormalizerFunc {
		m := map[string]interface{}{
			"type": "keyword",
		}

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

// CUSTOM NORMALIZER TYPE

type customNormalizer func(name string, properties ...lib.CustomNormalizerPropertiesFunc) lib.NormalizerFunc

func newCustomNormalizer() customNormalizer {
	return func(name string, properties ...lib.CustomNormalizerPropertiesFunc) lib.NormalizerFunc {
		m := utils.GetMap(properties)
		m["type"] = "custom"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (c customNormalizer) WithCharFilter(value []string) lib.CustomNormalizerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"char_filter": value,
		}
	}
}

func (c customNormalizer) WithFilter(value []string) lib.CustomNormalizerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"filter": value,
		}
	}
}

/*

 FIELD TYPES

*/

// KEYWORD TYPE

type keyword func(name string, properties ...lib.KeywordPropertiesFunc) lib.PropertiesFunc

func newKeyword() keyword {
	return func(name string, properties ...lib.KeywordPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "keyword"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (k keyword) WithNormalizer(value string) lib.KeywordPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"normalizer": value,
		}
	}
}

func (k keyword) WithIgnoreAbove(value int) lib.KeywordPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_above": value,
		}
	}
}

func (k keyword) WithBoost(value float64) lib.KeywordPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (k keyword) WithIndex(enabled bool) lib.KeywordPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (k keyword) WithDocValues(enabled bool) lib.KeywordPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (k keyword) WithEagerGlobalOrdinals(enabled bool) lib.KeywordPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"eager_global_ordinals": enabled,
		}
	}
}

func (k keyword) WithNullValue(value string) lib.KeywordPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

func (k keyword) WithStore(enabled bool) lib.KeywordPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

// TEXT TYPE

type text func(name string, properties ...lib.TextPropertiesFunc) lib.PropertiesFunc

func newText() text {
	return func(name string, properties ...lib.TextPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "text"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (t text) WithAnalyzer(value string) lib.TextPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"analyzer": value,
		}
	}
}

func (t text) WithSearchAnalyzer(value string) lib.TextPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"search_analyzer": value,
		}
	}
}

func (t text) WithBoost(value float64) lib.TextPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (t text) WithEagerGlobalOrdinals(enabled bool) lib.TextPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"eager_global_ordinals": enabled,
		}
	}
}

func (t text) WithIndex(enabled bool) lib.TextPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (t text) WithStore(enabled bool) lib.TextPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (t text) WithTermVector(value string) lib.TextPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"term_vector": value,
		}
	}
}

func (t text) WithPositionIncrementGap(value int) lib.TextPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"position_increment_gap": value,
		}
	}
}

func (t text) WithNorms(enabled bool) lib.TextPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"norms": enabled,
		}
	}
}

// INTEGER TYPE

type integer func(name string, properties ...lib.IntegerPropertiesFunc) lib.PropertiesFunc

func newInteger() integer {
	return func(name string, properties ...lib.IntegerPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "integer"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (i integer) WithCoerce(enabled bool) lib.IntegerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

func (i integer) WithBoost(value float64) lib.IntegerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (i integer) WithDocValues(enabled bool) lib.IntegerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (i integer) WithIgnoreMalformed(enabled bool) lib.IntegerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

func (i integer) WithIndex(enabled bool) lib.IntegerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (i integer) WithStore(enabled bool) lib.IntegerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (i integer) WithNullValue(value int) lib.IntegerPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

// LONG TYPE

type long func(name string, properties ...lib.LongPropertiesFunc) lib.PropertiesFunc

func newLong() long {
	return func(name string, properties ...lib.LongPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "long"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (l long) WithCoerce(value bool) lib.LongPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": value,
		}
	}
}

func (l long) WithBoost(value float64) lib.LongPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (l long) WithDocValues(enabled bool) lib.LongPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (l long) WithIndex(enabled bool) lib.LongPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (l long) WithNullValue(value int64) lib.LongPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

func (l long) WithStore(enabled bool) lib.LongPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (l long) WithOnScriptError(value string) lib.LongPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"on_script_error": value,
		}
	}
}

func (l long) WithScript(value map[string]interface{}) lib.LongPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"script": value,
		}
	}
}

func (l long) WithTimeSeriesMetric(value string) lib.LongPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"time_series_metric": value,
		}
	}
}

// FLOAT TYPE

type float func(name string, properties ...lib.FloatPropertiesFunc) lib.PropertiesFunc

func newFloat() float {
	return func(name string, properties ...lib.FloatPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "float"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (f float) WithCoerce(value bool) lib.FloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": value,
		}
	}
}

func (f float) WithBoost(value float64) lib.FloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (f float) WithDocValues(enabled bool) lib.FloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (f float) WithIndex(enabled bool) lib.FloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (f float) WithNullValue(value float64) lib.FloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

func (f float) WithStore(enabled bool) lib.FloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (f float) WithOnScriptError(value string) lib.FloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"on_script_error": value,
		}
	}
}

func (f float) WithScript(value map[string]interface{}) lib.FloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"script": value,
		}
	}
}

func (f float) WithTimeSeriesMetric(value string) lib.FloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"time_series_metric": value,
		}
	}
}

// DOUBLE TYPE

type double func(name string, properties ...lib.DoublePropertiesFunc) lib.PropertiesFunc

func newDouble() double {
	return func(name string, properties ...lib.DoublePropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "double"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (d double) WithCoerce(value bool) lib.DoublePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": value,
		}
	}
}

func (d double) WithBoost(value float64) lib.DoublePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (d double) WithDocValues(enabled bool) lib.DoublePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (d double) WithIndex(enabled bool) lib.DoublePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (d double) WithNullValue(value float64) lib.DoublePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

func (d double) WithStore(enabled bool) lib.DoublePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (d double) WithOnScriptError(value string) lib.DoublePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"on_script_error": value,
		}
	}
}

func (d double) WithScript(value map[string]interface{}) lib.DoublePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"script": value,
		}
	}
}

func (d double) WithTimeSeriesMetric(value string) lib.DoublePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"time_series_metric": value,
		}
	}
}

// SHORT TYPE

type short func(name string, properties ...lib.ShortPropertiesFunc) lib.PropertiesFunc

func newShort() short {
	return func(name string, properties ...lib.ShortPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "short"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (s short) WithCoerce(value bool) lib.ShortPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": value,
		}
	}
}

func (s short) WithBoost(value float64) lib.ShortPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (s short) WithDocValues(enabled bool) lib.ShortPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (s short) WithIndex(enabled bool) lib.ShortPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (s short) WithNullValue(value int16) lib.ShortPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

func (s short) WithStore(enabled bool) lib.ShortPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (s short) WithOnScriptError(value string) lib.ShortPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"on_script_error": value,
		}
	}
}

func (s short) WithScript(value map[string]interface{}) lib.ShortPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"script": value,
		}
	}
}

func (s short) WithTimeSeriesMetric(value string) lib.ShortPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"time_series_metric": value,
		}
	}
}

// BYTE TYPE

type bute func(name string, properties ...lib.BytePropertiesFunc) lib.PropertiesFunc

func newByte() bute {
	return func(name string, properties ...lib.BytePropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "byte"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (b bute) WithCoerce(value bool) lib.BytePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": value,
		}
	}
}

func (b bute) WithBoost(value float64) lib.BytePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (b bute) WithDocValues(enabled bool) lib.BytePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (b bute) WithIndex(enabled bool) lib.BytePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (b bute) WithNullValue(value int8) lib.BytePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

func (b bute) WithStore(enabled bool) lib.BytePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (b bute) WithOnScriptError(value string) lib.BytePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"on_script_error": value,
		}
	}
}

func (b bute) WithScript(value map[string]interface{}) lib.BytePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"script": value,
		}
	}
}

func (b bute) WithTimeSeriesMetric(value string) lib.BytePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"time_series_metric": value,
		}
	}
}

// HALFFLOAT TYPE

type halfFloat func(name string, properties ...lib.HalfFloatPropertiesFunc) lib.PropertiesFunc

func newHalfFloat() halfFloat {
	return func(name string, properties ...lib.HalfFloatPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "half_float"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (hf halfFloat) WithCoerce(value bool) lib.HalfFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": value,
		}
	}
}

func (hf halfFloat) WithBoost(value float64) lib.HalfFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (hf halfFloat) WithDocValues(enabled bool) lib.HalfFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (hf halfFloat) WithIndex(enabled bool) lib.HalfFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (hf halfFloat) WithNullValue(value float64) lib.HalfFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

func (hf halfFloat) WithStore(enabled bool) lib.HalfFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (hf halfFloat) WithOnScriptError(value string) lib.HalfFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"on_script_error": value,
		}
	}
}

func (hf halfFloat) WithScript(value map[string]interface{}) lib.HalfFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"script": value,
		}
	}
}

func (hf halfFloat) WithTimeSeriesMetric(value string) lib.HalfFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"time_series_metric": value,
		}
	}
}

// SCALEDFLOAT TYPE

type scaledFloat func(name string, properties ...lib.ScaledFloatPropertiesFunc) lib.PropertiesFunc

func newScaledFloat() scaledFloat {
	return func(name string, properties ...lib.ScaledFloatPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "scaled_float"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (sf scaledFloat) WithScalingFactor(value float64) lib.ScaledFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"scaling_factor": value,
		}
	}
}

func (sf scaledFloat) WithCoerce(value bool) lib.ScaledFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": value,
		}
	}
}

func (sf scaledFloat) WithBoost(value float64) lib.ScaledFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (sf scaledFloat) WithDocValues(enabled bool) lib.ScaledFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (sf scaledFloat) WithIndex(enabled bool) lib.ScaledFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (sf scaledFloat) WithNullValue(value float64) lib.ScaledFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

func (sf scaledFloat) WithStore(enabled bool) lib.ScaledFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (sf scaledFloat) WithOnScriptError(value string) lib.ScaledFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"on_script_error": value,
		}
	}
}

func (sf scaledFloat) WithScript(value map[string]interface{}) lib.ScaledFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"script": value,
		}
	}
}

func (sf scaledFloat) WithTimeSeriesMetric(value string) lib.ScaledFloatPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"time_series_metric": value,
		}
	}
}

// DATE TYPE

type date func(name string, properties ...lib.DatePropertiesFunc) lib.PropertiesFunc

func newDate() date {
	return func(name string, properties ...lib.DatePropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "date"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (d date) WithFormat(value string) lib.DatePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"format": value,
		}
	}
}

func (d date) WithBoost(value float64) lib.DatePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (d date) WithDocValues(enabled bool) lib.DatePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (d date) WithIndex(enabled bool) lib.DatePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (d date) WithStore(enabled bool) lib.DatePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (d date) WithIgnoreMalformed(enabled bool) lib.DatePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

func (d date) WithOnScriptError(value string) lib.DatePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"on_script_error": value,
		}
	}
}

func (d date) WithScript(value map[string]interface{}) lib.DatePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"script": value,
		}
	}
}

// DATENANOS TYPE

type dateNanos func(name string, properties ...lib.DateNanosPropertiesFunc) lib.PropertiesFunc

func newDateNanos() dateNanos {
	return func(name string, properties ...lib.DateNanosPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "date_nanos"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (dn dateNanos) WithFormat(value string) lib.DateNanosPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"format": value,
		}
	}
}

func (dn dateNanos) WithBoost(value float64) lib.DateNanosPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (dn dateNanos) WithDocValues(enabled bool) lib.DateNanosPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (dn dateNanos) WithIndex(enabled bool) lib.DateNanosPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (dn dateNanos) WithStore(enabled bool) lib.DateNanosPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (dn dateNanos) WithIgnoreMalformed(enabled bool) lib.DateNanosPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

func (dn dateNanos) WithOnScriptError(value string) lib.DateNanosPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"on_script_error": value,
		}
	}
}

func (dn dateNanos) WithScript(value map[string]interface{}) lib.DateNanosPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"script": value,
		}
	}
}

// BOOLEAN TYPE

type boolean func(name string, properties ...lib.BooleanPropertiesFunc) lib.PropertiesFunc

func newBoolean() boolean {
	return func(name string, properties ...lib.BooleanPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "boolean"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (b boolean) WithBoost(value float64) lib.BooleanPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (b boolean) WithDocValues(enabled bool) lib.BooleanPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (b boolean) WithIndex(enabled bool) lib.BooleanPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (b boolean) WithNullValue(value bool) lib.BooleanPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

func (b boolean) WithStore(enabled bool) lib.BooleanPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (b boolean) WithOnScriptError(value string) lib.BooleanPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"on_script_error": value,
		}
	}
}

func (b boolean) WithScript(value map[string]interface{}) lib.BooleanPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"script": value,
		}
	}
}

// BINARY TYPE

type binary func(name string, properties ...lib.BinaryPropertiesFunc) lib.PropertiesFunc

func newBinary() binary {
	return func(name string, properties ...lib.BinaryPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "binary"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (b binary) WithDocValues(enabled bool) lib.BinaryPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (b binary) WithStore(enabled bool) lib.BinaryPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

// IP TYPE

type ip func(name string, properties ...lib.IPPropertiesFunc) lib.PropertiesFunc

func newIP() ip {
	return func(name string, properties ...lib.IPPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "ip"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (i ip) WithBoost(value float64) lib.IPPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"boost": value,
		}
	}
}

func (i ip) WithDocValues(enabled bool) lib.IPPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (i ip) WithIndex(enabled bool) lib.IPPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (i ip) WithNullValue(value string) lib.IPPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

func (i ip) WithStore(enabled bool) lib.IPPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (i ip) WithOnScriptError(value string) lib.IPPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"on_script_error": value,
		}
	}
}

func (i ip) WithScript(value map[string]interface{}) lib.IPPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"script": value,
		}
	}
}

// GEOPOINT TYPE

type geoPoint func(name string, properties ...lib.GeoPointPropertiesFunc) lib.PropertiesFunc

func newGeoPoint() geoPoint {
	return func(name string, properties ...lib.GeoPointPropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "geo_point"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (gp geoPoint) WithIgnoreMalformed(value bool) lib.GeoPointPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": value,
		}
	}
}

func (gp geoPoint) WithIgnoreZValue(value bool) lib.GeoPointPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_z_value": value,
		}
	}
}

func (gp geoPoint) WithStore(enabled bool) lib.GeoPointPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (gp geoPoint) WithDocValues(enabled bool) lib.GeoPointPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

func (gp geoPoint) WithIndex(enabled bool) lib.GeoPointPropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

// GEOSHAPE TYPE

type geoShape func(name string, properties ...lib.GeoShapePropertiesFunc) lib.PropertiesFunc

func newGeoShape() geoShape {
	return func(name string, properties ...lib.GeoShapePropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "geo_shape"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (gs geoShape) WithIgnoreMalformed(value bool) lib.GeoShapePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": value,
		}
	}
}

func (gs geoShape) WithIgnoreZValue(value bool) lib.GeoShapePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_z_value": value,
		}
	}
}

func (gs geoShape) WithStore(enabled bool) lib.GeoShapePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (gs geoShape) WithIndex(enabled bool) lib.GeoShapePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

func (gs geoShape) WithCoerce(enabled bool) lib.GeoShapePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

func (gs geoShape) WithOrientation(value string) lib.GeoShapePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"orientation": value,
		}
	}
}

func (gs geoShape) WithStrategy(value string) lib.GeoShapePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"strategy": value,
		}
	}
}

func (gs geoShape) WithTree(value string) lib.GeoShapePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"tree": value,
		}
	}
}

func (gs geoShape) WithPrecision(value string) lib.GeoShapePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"precision": value,
		}
	}
}

func (gs geoShape) WithTreeLevels(value string) lib.GeoShapePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"tree_levels": value,
		}
	}
}

// INTERGERRANGE TYPE

type integerRange func(name string, properties ...lib.IntegerRangePropertiesFunc) lib.PropertiesFunc

func newIntegerRange() integerRange {
	return func(name string, properties ...lib.IntegerRangePropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "integer_range"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (ir integerRange) WithStore(enabled bool) lib.IntegerRangePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (ir integerRange) WithIndex(enabled bool) lib.IntegerRangePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

// FLOATRANGE TYPE

type floatRange func(name string, properties ...lib.FloatRangePropertiesFunc) lib.PropertiesFunc

func newFloatRange() floatRange {
	return func(name string, properties ...lib.FloatRangePropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "float_range"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (fr floatRange) WithStore(enabled bool) lib.FloatRangePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (fr floatRange) WithIndex(enabled bool) lib.FloatRangePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

// DATERANGE TYPE

type dateRange func(name string, properties ...lib.DateRangePropertiesFunc) lib.PropertiesFunc

func newDateRange() dateRange {
	return func(name string, properties ...lib.DateRangePropertiesFunc) lib.PropertiesFunc {
		m := utils.GetMap(properties)

		m["type"] = "date_range"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: m,
			}
		}
	}
}

func (dr dateRange) WithFormat(value string) lib.DateRangePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"format": value,
		}
	}
}

func (dr dateRange) WithStore(enabled bool) lib.DateRangePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

func (dr dateRange) WithIndex(enabled bool) lib.DateRangePropertiesFunc {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}
