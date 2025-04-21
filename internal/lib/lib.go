package lib

import "github.com/xoticdsign/porter/internal/utils"

type AnalyzerFunc = utils.MapperFunc

type StandardAnalyzerPropertiesFunc = utils.MapperFunc
type WhitespaceAnalyzerPropertiesFunc = utils.MapperFunc
type CustomAnalyzerPropertiesFunc = utils.MapperFunc
type StopAnalyzerPropertiesFunc = utils.MapperFunc
type KeywordAnalyzerPropertiesFunc = utils.MapperFunc
type PatternAnalyzerPropertiesFunc = utils.MapperFunc
type EdgeNGramAnalyzerPropertiesFunc = utils.MapperFunc
type LanguageSpecificAnalyzerPropertiesFunc = utils.MapperFunc

type NormalizerFunc = utils.MapperFunc

type CustomNormalizerPropertiesFunc = utils.MapperFunc

type PropertiesFunc = utils.MapperFunc

type KeywordPropertiesFunc = utils.MapperFunc
type TextPropertiesFunc = utils.MapperFunc
type IntegerPropertiesFunc = utils.MapperFunc
type LongPropertiesFunc = utils.MapperFunc
type FloatPropertiesFunc = utils.MapperFunc
type DoublePropertiesFunc = utils.MapperFunc
type ShortPropertiesFunc = utils.MapperFunc
type BytePropertiesFunc = utils.MapperFunc
type HalfFloatPropertiesFunc = utils.MapperFunc
type ScaledFloatPropertiesFunc = utils.MapperFunc
type DatePropertiesFunc = utils.MapperFunc
type DateNanosPropertiesFunc = utils.MapperFunc
type BooleanPropertiesFunc = utils.MapperFunc
type BinaryPropertiesFunc = utils.MapperFunc
type IPPropertiesFunc = utils.MapperFunc
type GeoPointPropertiesFunc = utils.MapperFunc
type GeoShapePropertiesFunc = utils.MapperFunc
type IntegerRangePropertiesFunc = utils.MapperFunc
type FloatRangePropertiesFunc = utils.MapperFunc
type DateRangePropertiesFunc = utils.MapperFunc
