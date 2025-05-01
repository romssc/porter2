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
	Keyword     FieldKeyword
	Text        FieldText
	Integer     FieldInteger
	Long        FieldLong
	Float       FieldFloat
	Double      FieldDouble
	Short       FieldShort
	Byte        FieldByte
	HalfFloat   FieldHalfFloat
	ScaledFloat FieldScaledFloat
	Date        FieldDate
	Boolean     FieldBoolean
	IP          FieldIP
}

// NewFields() generates a map of field names and their corresponding values.
func (m mappings) NewFields(fields ...FieldFunc) map[string]interface{} {
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

type FieldFunc func() map[string]interface{}

// KEYWORD

type FieldKeywordProperties func() map[string]interface{}
type FieldKeyword func(name string, fake Fake, properties ...FieldKeywordProperties) FieldFunc

func newFieldKeyword() FieldKeyword {
	return func(name string, fake Fake, properties ...FieldKeywordProperties) FieldFunc {
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

// WithDocValues() adds a "doc_values" property to a FieldKeyword.
func (k FieldKeyword) WithDocValues(enabled bool) FieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithEagerGlobalOrdinals() adds an "eager_global_ordinals" property to a FieldKeyword.
func (k FieldKeyword) WithEagerGlobalOrdinals(enabled bool) FieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"eager_global_ordinals": enabled,
		}
	}
}

/* TODO: WithField(...) */

// WithIgnoreAbove() adds an "ignore_above" property to a FieldKeyword.
func (k FieldKeyword) WithIgnoreAbove(value int) FieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_above": value,
		}
	}
}

// WithIndex() adds an "index" property to a FieldKeyword.
func (k FieldKeyword) WithIndex(enabled bool) FieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithIndexOptions(...) */

/* TODO: WithMeta(...) */

/* TODO: WithNorms(...) */

// WithNullValue() adds a "null_value" property to a FieldKeyword.
func (k FieldKeyword) WithNullValue(value string) FieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a FieldKeyword.
func (k FieldKeyword) WithStore(enabled bool) FieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithSimilarity(...) */

// WithNormalizer() adds a "normalizer" property to a FieldKeyword.
func (k FieldKeyword) WithNormalizer(value string) FieldKeywordProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"normalizer": value,
		}
	}
}

/* TODO: WithSplitQueriesOnWhitespace(...) */

/* TODO: WithTimeSeriesDimension(...) */

// TEXT

type FieldTextProperties func() map[string]interface{}
type FieldText func(name string, fake Fake, properties ...FieldTextProperties) FieldFunc

func newFieldText() FieldText {
	return func(name string, fake Fake, properties ...FieldTextProperties) FieldFunc {
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

// WithAnalyzer() adds an "analyzer" property to a FieldText.
func (t FieldText) WithAnalyzer(value string) FieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"analyzer": value,
		}
	}
}

// WithEagerGlobalOrdinals() adds an "eager_global_ordinals" property to a FieldText.
func (t FieldText) WithEagerGlobalOrdinals(enabled bool) FieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"eager_global_ordinals": enabled,
		}
	}
}

/* TODO: WithFieldData(...) */

/* TODO: WithFieldDataFrequencyFilter(...) */

/* TODO: WithField(...) */

// WithIndex() adds an "index" property to a FieldText.
func (t FieldText) WithIndex(enabled bool) FieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithIndexOptions(...) */

/* TODO: WithIndexPrefixes(...) */

/* TODO: WithIndexPhrases(...) */

// WithNorms() adds a "norms" property to a FieldText.
func (t FieldText) WithNorms(enabled bool) FieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"norms": enabled,
		}
	}
}

// WithPositionIncrementGap() adds a "position_increment_gap" property to a FieldText.
func (t FieldText) WithPositionIncrementGap(value int) FieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"position_increment_gap": value,
		}
	}
}

// WithStore() adds a "store" property to a FieldText.
func (t FieldText) WithStore(enabled bool) FieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

// WithSearchAnalyzer() adds a "search_analyzer" property to a FieldText.
func (t FieldText) WithSearchAnalyzer(value string) FieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"search_analyzer": value,
		}
	}
}

/* TODO: WithSearchQuoteAnalyzer(...) */

/* TODO: WithSimilarity(...) */

// WithTermVector() adds a "term_vector" property to a FieldText.
func (t FieldText) WithTermVector(value string) FieldTextProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"term_vector": value,
		}
	}
}

/* TODO: WithMeta(...) */

// INTEGER

type FieldIntegerProperties func() map[string]interface{}
type FieldInteger func(name string, fake FakeInteger, properties ...FieldIntegerProperties) FieldFunc

func newFieldInteger() FieldInteger {
	return func(name string, fake FakeInteger, properties ...FieldIntegerProperties) FieldFunc {
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
func (i FieldInteger) WithCoerce(enabled bool) FieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to an integer field.
func (i FieldInteger) WithDocValues(enabled bool) FieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to an integer field.
func (i FieldInteger) WithIgnoreMalformed(enabled bool) FieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to an integer field.
func (i FieldInteger) WithIndex(enabled bool) FieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to an integer field.
func (i FieldInteger) WithNullValue(value int) FieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to an integer field.
func (i FieldInteger) WithStore(enabled bool) FieldIntegerProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// LONG

type FieldLongProperties func() map[string]interface{}
type FieldLong func(name string, fake FakeLong, properties ...FieldLongProperties) FieldFunc

func newFieldLong() FieldLong {
	return func(name string, fake FakeLong, properties ...FieldLongProperties) FieldFunc {
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
func (l FieldLong) WithCoerce(enabled bool) FieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a long field.
func (l FieldLong) WithDocValues(enabled bool) FieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a long field.
func (l FieldLong) WithIgnoreMalformed(enabled bool) FieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a long field.
func (l FieldLong) WithIndex(enabled bool) FieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a long field.
func (l FieldLong) WithNullValue(value int) FieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a long field.
func (l FieldLong) WithStore(enabled bool) FieldLongProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// FLOAT

type FieldFloatProperties func() map[string]interface{}
type FieldFloat func(name string, fake FakeFloats, properties ...FieldFloatProperties) FieldFunc

func newFieldFloat() FieldFloat {
	return func(name string, fake FakeFloats, properties ...FieldFloatProperties) FieldFunc {
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
func (f FieldFloat) WithCoerce(enabled bool) FieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a float field.
func (f FieldFloat) WithDocValues(enabled bool) FieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a float field.
func (f FieldFloat) WithIgnoreMalformed(enabled bool) FieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a float field.
func (f FieldFloat) WithIndex(enabled bool) FieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a float field.
func (f FieldFloat) WithNullValue(value int) FieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a float field.
func (f FieldFloat) WithStore(enabled bool) FieldFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// DOUBLE

type FieldDoubleProperties func() map[string]interface{}
type FieldDouble func(name string, fake FakeDouble, properties ...FieldDoubleProperties) FieldFunc

func newFieldDouble() FieldDouble {
	return func(name string, fake FakeDouble, properties ...FieldDoubleProperties) FieldFunc {
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
func (d FieldDouble) WithCoerce(enabled bool) FieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a double field.
func (d FieldDouble) WithDocValues(enabled bool) FieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a double field.
func (d FieldDouble) WithIgnoreMalformed(enabled bool) FieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a double field.
func (d FieldDouble) WithIndex(enabled bool) FieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a double field.
func (d FieldDouble) WithNullValue(value int) FieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a double field.
func (d FieldDouble) WithStore(enabled bool) FieldDoubleProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// SHORT

type FieldShortProperties func() map[string]interface{}
type FieldShort func(name string, fake FakeShort, properties ...FieldShortProperties) FieldFunc

func newFieldShort() FieldShort {
	return func(name string, fake FakeShort, properties ...FieldShortProperties) FieldFunc {
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
func (s FieldShort) WithCoerce(enabled bool) FieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a short field.
func (s FieldShort) WithDocValues(enabled bool) FieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a short field.
func (s FieldShort) WithIgnoreMalformed(enabled bool) FieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a short field.
func (s FieldShort) WithIndex(enabled bool) FieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a short field.
func (s FieldShort) WithNullValue(value int) FieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a short field.
func (s FieldShort) WithStore(enabled bool) FieldShortProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// BYTE

type FieldByteProperties func() map[string]interface{}
type FieldByte func(name string, fake FakeByte, properties ...FieldByteProperties) FieldFunc

func newFieldByte() FieldByte {
	return func(name string, fake FakeByte, properties ...FieldByteProperties) FieldFunc {
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
func (b FieldByte) WithCoerce(enabled bool) FieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a byte field.
func (b FieldByte) WithDocValues(enabled bool) FieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a byte field.
func (b FieldByte) WithIgnoreMalformed(enabled bool) FieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a byte field.
func (b FieldByte) WithIndex(enabled bool) FieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a byte field.
func (b FieldByte) WithNullValue(value int) FieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a byte field.
func (b FieldByte) WithStore(enabled bool) FieldByteProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// HALFFLOAT

type FieldHalfFloatProperties func() map[string]interface{}
type FieldHalfFloat func(name string, fake FakeHalfFloat, properties ...FieldHalfFloatProperties) FieldFunc

func newFieldHalfFloat() FieldHalfFloat {
	return func(name string, fake FakeHalfFloat, properties ...FieldHalfFloatProperties) FieldFunc {
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
func (h FieldHalfFloat) WithCoerce(enabled bool) FieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a half_float field.
func (h FieldHalfFloat) WithDocValues(enabled bool) FieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a half_float field.
func (h FieldHalfFloat) WithIgnoreMalformed(enabled bool) FieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a half_float field.
func (h FieldHalfFloat) WithIndex(enabled bool) FieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a half_float field.
func (h FieldHalfFloat) WithNullValue(value int) FieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a half_float field.
func (h FieldHalfFloat) WithStore(enabled bool) FieldHalfFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// SCALEDFLOAT

type FieldScaledFloatProperties func() map[string]interface{}
type FieldScaledFloat func(name string, fake FakeScaledFloat, properties ...FieldScaledFloatProperties) FieldFunc

func newFieldScaledFloat() FieldScaledFloat {
	return func(name string, fake FakeScaledFloat, properties ...FieldScaledFloatProperties) FieldFunc {
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
func (s FieldScaledFloat) WithCoerce(enabled bool) FieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"coerce": enabled,
		}
	}
}

// WithDocValues() adds a "doc_values" property to a scaled_float field.
func (s FieldScaledFloat) WithDocValues(enabled bool) FieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIgnoreMalformed() adds an "ignore_malformed" property to a scaled_float field.
func (s FieldScaledFloat) WithIgnoreMalformed(enabled bool) FieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a scaled_float field.
func (s FieldScaledFloat) WithIndex(enabled bool) FieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// WithNullValue() adds a "null_value" property to a scaled_float field.
func (s FieldScaledFloat) WithNullValue(value int) FieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a scaled_float field.
func (s FieldScaledFloat) WithStore(enabled bool) FieldScaledFloatProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */

// DATE

type FieldDateProperties func() map[string]interface{}
type FieldDate func(name string, fake FakeDates, properties ...FieldDateProperties) FieldFunc

func newFieldDate() FieldDate {
	return func(name string, fake FakeDates, properties ...FieldDateProperties) FieldFunc {
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
func (d FieldDate) WithDocValues(enabled bool) FieldDateProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithFormat() adds a "format" property to a date field.
func (d FieldDate) WithFormat(value string) FieldDateProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"format": value,
		}
	}
}

/* TODO: WithLocale(...) */

// WithIgnoreMalformed() adds an "ignore_malformed" property to a date field.
func (d FieldDate) WithIgnoreMalformed(enabled bool) FieldDateProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"ignore_malformed": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a date field.
func (d FieldDate) WithIndex(enabled bool) FieldDateProperties {
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
func (d FieldDate) WithStore(enabled bool) FieldDateProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

// BOOLEAN

type FieldBooleanProperties func() map[string]interface{}
type FieldBoolean func(name string, fake FakeBoolean, properties ...FieldBooleanProperties) FieldFunc

func newFieldBoolean() FieldBoolean {
	return func(name string, fake FakeBoolean, properties ...FieldBooleanProperties) FieldFunc {
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
func (b FieldBoolean) WithDocValues(enabled bool) FieldBooleanProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

// WithIndex() adds an "index" property to a boolean field.
func (b FieldBoolean) WithIndex(enabled bool) FieldBooleanProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

/* TODO: WithIgnoreMalformed(...) */

// WithNullValue() adds a "null_value" property to a boolean field.
func (b FieldBoolean) WithNullValue(value bool) FieldBooleanProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to a boolean field.
func (b FieldBoolean) WithStore(enabled bool) FieldBooleanProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithMeta(...) */

/* TODO: WithTimeSeriesDimension(...) */

// IP

type FieldIPProperties func() map[string]interface{}
type FieldIP func(name string, fake FakeIP, properties ...FieldIPProperties) FieldFunc

func newFieldIP() FieldIP {
	return func(name string, fake FakeIP, properties ...FieldIPProperties) FieldFunc {
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
func (i FieldIP) WithDocValues(enabled bool) FieldIPProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"doc_values": enabled,
		}
	}
}

/* TODO: WithIgnoreMalformed(...) */

// WithIndex() adds an "index" property to an IP field.
func (i FieldIP) WithIndex(enabled bool) FieldIPProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"index": enabled,
		}
	}
}

// WithNullValue() adds a "null_value" property to an IP field.
func (i FieldIP) WithNullValue(value string) FieldIPProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"null_value": value,
		}
	}
}

/* TODO: WithOnScriptError(...) */

/* TODO: WithScript(...) */

// WithStore() adds a "store" property to an IP field.
func (i FieldIP) WithStore(enabled bool) FieldIPProperties {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"store": enabled,
		}
	}
}

/* TODO: WithTimeSeriesDimension(...) */
