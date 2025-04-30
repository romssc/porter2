package porter

/*

This file defines a set of types, functions, and properties to facilitate the dynamic generation
of structured Elasticsearch-compatible field data. The core functionality includes the creation of
various field types (e.g., keyword, text, integer) along with customizable properties for each type.

The NewFields() function allows users to compose a map of fields using a flexible set of field-specific
functions. Each field type has a corresponding factory function that returns a function for generating
field definitions with customizable properties. These properties can be added dynamically using functional
composition.

The goal of this package is to provide an easy and extensible way to generate realistic, test-ready data
for Elasticsearch systems, ensuring compatibility with Elasticsearch's field mappings and type constraints.

*/

type fields struct {
	Keyword     fieldKeyword
	Text        fieldText
	Integer     fieldInteger
	Long        fieldLong
	Float       fieldFloat
	Double      fieldDouble
	Short       fieldShort
	Byte        fieldByte
	HalfFloat   fieldHalfFloat
	ScaledFloat fieldScaledFloat
	Date        fieldDate
	Boolean     fieldBoolean
	IP          fieldIP
}

// NewFields() generates a map of field names and their corresponding values.
func (m mappings) NewFields(fields ...fieldFunc) map[string]interface{} {
	r := map[string]interface{}{}

	for _, fn := range fields {
		if fn == nil {
			continue
		}
		for k, v := range fn() {
			r[k] = v
		}
	}

	return r
}

type fieldFunc func() map[string]interface{}

// KEYWORD

type fieldKeywordProperties func() map[string]interface{}
type fieldKeyword func(name string, fake Fake, properties ...fieldKeywordProperties) fieldFunc

func newFieldKeyword() fieldKeyword {
	return func(name string, fake Fake, properties ...fieldKeywordProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "keyword"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithDocValues() adds a "doc_values" property to a fieldKeyword.
func (k fieldKeyword) WithDocValues(enabled bool) fieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithEagerGlobalOrdinals() adds an "eager_global_ordinals" property to a fieldKeyword.
func (k fieldKeyword) WithEagerGlobalOrdinals(enabled bool) fieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"eager_global_ordinals": enabled,
		}
	}
}

/* TODO: WithField(...) */

// WithIgnoreAbove() adds an "ignore_above" property to a fieldKeyword.
func (k fieldKeyword) WithIgnoreAbove(value int) fieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_above": value,
		}
	}
}

// WithIndex() adds an "index" property to a fieldKeyword.
func (k fieldKeyword) WithIndex(enabled bool) fieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithIndexOptions(...) */

/* TODO: WithMeta(...) */

/* TODO: WithNorms(...) */

// WithNullValue() adds a "null_value" property to a fieldKeyword.
func (k fieldKeyword) WithNullValue(value string) fieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a fieldKeyword.
func (k fieldKeyword) WithStore(enabled bool) fieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithSimilarity(...) */

// WithNormalizer() adds a "normalizer" property to a fieldKeyword.
func (k fieldKeyword) WithNormalizer(value string) fieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"normalizer": value,
		}
	}
}

/* TODO: WithSplitQueriesOnWhitespace(...) */

/* TODO: WithTimeSeriesDimension(...) */

// TEXT

type fieldTextProperties func() map[string]interface{}
type fieldText func(name string, fake Fake, properties ...fieldTextProperties) fieldFunc

func newFieldText() fieldText {
	return func(name string, fake Fake, properties ...fieldTextProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "text"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithAnalyzer() adds an "analyzer" property to a fieldText.
func (t fieldText) WithAnalyzer(value string) fieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"analyzer": value,
		}
	}
}

// WithEagerGlobalOrdinals() adds an "eager_global_ordinals" property to a fieldText.
func (t fieldText) WithEagerGlobalOrdinals(enabled bool) fieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"eager_global_ordinals": enabled,
		}
	}
}

/* TODO: WithFieldData(...) */

/* TODO: WithFieldDataFrequencyFilter(...) */

/* TODO: WithField(...) */

// WithIndex() adds an "index" property to a fieldText.
func (t fieldText) WithIndex(enabled bool) fieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithIndexOptions(...) */

/* TODO: WithIndexPrefixes(...) */

/* TODO: WithIndexPhrases(...) */

// WithNorms() adds a "norms" property to a fieldText.
func (t fieldText) WithNorms(enabled bool) fieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"norms": enabled,
		}
	}
}

// WithPositionIncrementGap() adds a "position_increment_gap" property to a fieldText.
func (t fieldText) WithPositionIncrementGap(value int) fieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"position_increment_gap": value,
		}
	}
}

// WithStore() adds a "store" property to a fieldText.
func (t fieldText) WithStore(enabled bool) fieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

// WithSearchAnalyzer() adds a "search_analyzer" property to a fieldText.
func (t fieldText) WithSearchAnalyzer(value string) fieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"search_analyzer": value,
		}
	}
}

/* TODO: WithSearchQuoteAnalyzer(...) */

/* TODO: WithSimilarity(...) */

// WithTermVector() adds a "term_vector" property to a fieldText.
func (t fieldText) WithTermVector(value string) fieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"term_vector": value,
		}
	}
}

/* TODO: WithMeta(...) */

// INTEGER

type fieldIntegerProperties func() map[string]interface{}
type fieldInteger func(name string, fake FakeInteger, properties ...fieldIntegerProperties) fieldFunc

func newFieldInteger() fieldInteger {
	return func(name string, fake FakeInteger, properties ...fieldIntegerProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "integer"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithCoerce() adds a "coerce" property to an integer field.
func (i fieldInteger) WithCoerce(enabled bool) fieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to an integer field.
func (i fieldInteger) WithDocValues(enabled bool) fieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to an integer field.
func (i fieldInteger) WithIgnoreMalformed(enabled bool) fieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to an integer field.
func (i fieldInteger) WithIndex(enabled bool) fieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to an integer field.
func (i fieldInteger) WithNullValue(value int) fieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to an integer field.
func (i fieldInteger) WithStore(enabled bool) fieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// LONG

type fieldLongProperties func() map[string]interface{}
type fieldLong func(name string, fake FakeLong, properties ...fieldLongProperties) fieldFunc

func newFieldLong() fieldLong {
	return func(name string, fake FakeLong, properties ...fieldLongProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "long"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithCoerce() adds a "coerce" property to a long field.
func (l fieldLong) WithCoerce(enabled bool) fieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a long field.
func (l fieldLong) WithDocValues(enabled bool) fieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a long field.
func (l fieldLong) WithIgnoreMalformed(enabled bool) fieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a long field.
func (l fieldLong) WithIndex(enabled bool) fieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a long field.
func (l fieldLong) WithNullValue(value int) fieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a long field.
func (l fieldLong) WithStore(enabled bool) fieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// FLOAT

type fieldFloatProperties func() map[string]interface{}
type fieldFloat func(name string, fake FakeFloats, properties ...fieldFloatProperties) fieldFunc

func newFieldFloat() fieldFloat {
	return func(name string, fake FakeFloats, properties ...fieldFloatProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "float"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithCoerce() adds a "coerce" property to a float field.
func (f fieldFloat) WithCoerce(enabled bool) fieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a float field.
func (f fieldFloat) WithDocValues(enabled bool) fieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a float field.
func (f fieldFloat) WithIgnoreMalformed(enabled bool) fieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a float field.
func (f fieldFloat) WithIndex(enabled bool) fieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a float field.
func (f fieldFloat) WithNullValue(value int) fieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a float field.
func (f fieldFloat) WithStore(enabled bool) fieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// DOUBLE

type fieldDoubleProperties func() map[string]interface{}
type fieldDouble func(name string, fake FakeDouble, properties ...fieldDoubleProperties) fieldFunc

func newFieldDouble() fieldDouble {
	return func(name string, fake FakeDouble, properties ...fieldDoubleProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "double"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithCoerce() adds a "coerce" property to a double field.
func (d fieldDouble) WithCoerce(enabled bool) fieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a double field.
func (d fieldDouble) WithDocValues(enabled bool) fieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a double field.
func (d fieldDouble) WithIgnoreMalformed(enabled bool) fieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a double field.
func (d fieldDouble) WithIndex(enabled bool) fieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a double field.
func (d fieldDouble) WithNullValue(value int) fieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a double field.
func (d fieldDouble) WithStore(enabled bool) fieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// SHORT

type fieldShortProperties func() map[string]interface{}
type fieldShort func(name string, fake FakeShort, properties ...fieldShortProperties) fieldFunc

func newFieldShort() fieldShort {
	return func(name string, fake FakeShort, properties ...fieldShortProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "short"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithCoerce() adds a "coerce" property to a short field.
func (s fieldShort) WithCoerce(enabled bool) fieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a short field.
func (s fieldShort) WithDocValues(enabled bool) fieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a short field.
func (s fieldShort) WithIgnoreMalformed(enabled bool) fieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a short field.
func (s fieldShort) WithIndex(enabled bool) fieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a short field.
func (s fieldShort) WithNullValue(value int) fieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a short field.
func (s fieldShort) WithStore(enabled bool) fieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// BYTE

type fieldByteProperties func() map[string]interface{}
type fieldByte func(name string, fake FakeByte, properties ...fieldByteProperties) fieldFunc

func newFieldByte() fieldByte {
	return func(name string, fake FakeByte, properties ...fieldByteProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "byte"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithCoerce() adds a "coerce" property to a byte field.
func (b fieldByte) WithCoerce(enabled bool) fieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a byte field.
func (b fieldByte) WithDocValues(enabled bool) fieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a byte field.
func (b fieldByte) WithIgnoreMalformed(enabled bool) fieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a byte field.
func (b fieldByte) WithIndex(enabled bool) fieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a byte field.
func (b fieldByte) WithNullValue(value int) fieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a byte field.
func (b fieldByte) WithStore(enabled bool) fieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// HALFFLOAT

type fieldHalfFloatProperties func() map[string]interface{}
type fieldHalfFloat func(name string, fake FakeHalfFloat, properties ...fieldHalfFloatProperties) fieldFunc

func newFieldHalfFloat() fieldHalfFloat {
	return func(name string, fake FakeHalfFloat, properties ...fieldHalfFloatProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "half_float"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithCoerce() adds a "coerce" property to a half_float field.
func (h fieldHalfFloat) WithCoerce(enabled bool) fieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a half_float field.
func (h fieldHalfFloat) WithDocValues(enabled bool) fieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a half_float field.
func (h fieldHalfFloat) WithIgnoreMalformed(enabled bool) fieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a half_float field.
func (h fieldHalfFloat) WithIndex(enabled bool) fieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a half_float field.
func (h fieldHalfFloat) WithNullValue(value int) fieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a half_float field.
func (h fieldHalfFloat) WithStore(enabled bool) fieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// SCALEDFLOAT

type fieldScaledFloatProperties func() map[string]interface{}
type fieldScaledFloat func(name string, fake FakeScaledFloat, properties ...fieldScaledFloatProperties) fieldFunc

func newFieldScaledFloat() fieldScaledFloat {
	return func(name string, fake FakeScaledFloat, properties ...fieldScaledFloatProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "scaled_float"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithCoerce() adds a "coerce" property to a scaled_float field.
func (s fieldScaledFloat) WithCoerce(enabled bool) fieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a scaled_float field.
func (s fieldScaledFloat) WithDocValues(enabled bool) fieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a scaled_float field.
func (s fieldScaledFloat) WithIgnoreMalformed(enabled bool) fieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a scaled_float field.
func (s fieldScaledFloat) WithIndex(enabled bool) fieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a scaled_float field.
func (s fieldScaledFloat) WithNullValue(value int) fieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a scaled_float field.
func (s fieldScaledFloat) WithStore(enabled bool) fieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// DATE

type fieldDateProperties func() map[string]interface{}
type fieldDate func(name string, fake FakeDates, properties ...fieldDateProperties) fieldFunc

func newFieldDate() fieldDate {
	return func(name string, fake FakeDates, properties ...fieldDateProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "date"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithDocValues() adds a "doc_values" property to a date field.
func (d fieldDate) WithDocValues(enabled bool) fieldDateProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithFormat() adds a "format" property to a date field.
func (d fieldDate) WithFormat(value string) fieldDateProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"format": value,
		}
	}
}

/* TODO: WithLocale(...) */

// WithIgnoreMalformed() adds an "ignore_malformed" property to a date field.
func (d fieldDate) WithIgnoreMalformed(enabled bool) fieldDateProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a date field.
func (d fieldDate) WithIndex(enabled bool) fieldDateProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithNullValue(...) */

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a date field.
func (d fieldDate) WithStore(enabled bool) fieldDateProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// BOOLEAN

type fieldBooleanProperties func() map[string]interface{}
type fieldBoolean func(name string, fake FakeBoolean, properties ...fieldBooleanProperties) fieldFunc

func newFieldBoolean() fieldBoolean {
	return func(name string, fake FakeBoolean, properties ...fieldBooleanProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "boolean"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithDocValues() adds a "doc_values" property to a boolean field.
func (b fieldBoolean) WithDocValues(enabled bool) fieldBooleanProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a boolean field.
func (b fieldBoolean) WithIndex(enabled bool) fieldBooleanProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithIgnoreMalformed(...) */

// WithNullValue() adds a "null_value" property to a boolean field.
func (b fieldBoolean) WithNullValue(value bool) fieldBooleanProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a boolean field.
func (b fieldBoolean) WithStore(enabled bool) fieldBooleanProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

/* TODO: WithTimeSeriesDimension(...) */

// IP

type fieldIPProperties func() map[string]interface{}
type fieldIP func(name string, fake FakeIP, properties ...fieldIPProperties) fieldFunc

func newFieldIP() fieldIP {
	return func(name string, fake FakeIP, properties ...fieldIPProperties) fieldFunc {
		storeToGenerate(name, string(fake))

		r := map[string]interface{}{}

		for _, fn := range properties {
			if fn == nil {
				continue
			}
			for k, v := range fn() {
				r[k] = v
			}
		}

		r["type"] = "ip"

		return func() map[string]interface{} {
			return map[string]interface{}{
				name: r,
			}
		}
	}
}

// WithDocValues() adds a "doc_values" property to an IP field.
func (i fieldIP) WithDocValues(enabled bool) fieldIPProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

/* TODO: WithIgnoreMalformed(...) */

// WithIndex() adds an "index" property to an IP field.
func (i fieldIP) WithIndex(enabled bool) fieldIPProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

// WithNullValue() adds a "null_value" property to an IP field.
func (i fieldIP) WithNullValue(value string) fieldIPProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to an IP field.
func (i fieldIP) WithStore(enabled bool) fieldIPProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */
