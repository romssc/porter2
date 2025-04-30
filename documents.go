package porter

import (
	"fmt"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/xoticdsign/porter/internal/utils"
)

/*

This file contains the full set of types and factory functions for generating synthetic
documents using field-specific fake data generators.

It supports two sources of data: reading raw documents from files or dynamically generating
structured documents with realistic values using "gofakeit". Each field type maps to a
specific generator, ensuring values match Elasticsearch
index expectations.

The goal is to enable fast and accurate test data generation for Elasticsearch-based systems,
with type-safe mappings and reusable, composable generation logic.

*/

var (
	ErrOriginFromFile = fmt.Errorf("origin: failed to read documents from file")
)

type location struct {
	FromFile locationFromFile
	Generate locationGenerate
}

type originFunc func(t temp) ([]byte, error)

type locationFromFile func(path string) originFunc

func newLocationFromFile() locationFromFile {
	return func(path string) originFunc {
		return func(t temp) ([]byte, error) {
			contents, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("%w\n%v", ErrOriginFromFile, err)
			}
			return contents, nil
		}
	}
}

type locationGenerate func(amount int) originFunc

func newLocationGenerate() locationGenerate {
	return func(amount int) originFunc {
		return func(t temp) ([]byte, error) {
			var docs []byte

			for c := 1; c <= amount; c++ {
				m := map[string]interface{}{
					"index": map[string]interface{}{
						"_index": t.config.Name,
						"_id":    c,
					},
				}

				f := map[string]interface{}{}

				for k, v := range toGenerate {
					data := generateFakeData(v)

					f[k] = data
				}

				mB := utils.MarshalJSON(m)
				fB := utils.MarshalJSON(f)

				docs = append(docs, mB...)
				docs = append(docs, '\n')
				docs = append(docs, fB...)
				docs = append(docs, '\n')
			}

			return docs, nil
		}
	}
}

type Fake string

var (
	FakeEmail     Fake = "email"
	FakeFirstName Fake = "first_name"
	FakeLastName  Fake = "last_name"
	FakeFullName  Fake = "full_name"
	FakeUsername  Fake = "username"
	FakePhone     Fake = "phone"
	FakeCountry   Fake = "country"
	FakeCity      Fake = "city"
	FakeStreet    Fake = "street"
	FakeZip       Fake = "zip"
	FakeUUID      Fake = "uuid"
	FakeURL       Fake = "url"
	FakeCompany   Fake = "company"
	FakeJobTitle  Fake = "job_title"
	FakeColor     Fake = "color"
	FakeIPv4      Fake = "ipv4"
	FakeIPv6      Fake = "ipv6"
	FakeBool      Fake = "bool"
	FakeInt       Fake = "int"
	FakeFloat     Fake = "float"
	FakeDate      Fake = "date"
	FakeTimestamp Fake = "timestamp"
	FakeParagraph Fake = "paragraph"
)

type FakeBoolean string

var (
	FakeBooleanBool FakeBoolean = "bool"
)

type FakeInteger string

var (
	FakeIntegerInt FakeInteger = "int"
)

type FakeLong string

var (
	FakeLongInt FakeLong = "int"
)

type FakeShort string

var (
	FakeShortInt FakeShort = "int"
)

type FakeByte string

var (
	FakeByteInt FakeByte = "int"
)

type FakeFloats string

var (
	FakeFloatFloat FakeFloats = "float"
)

type FakeDouble string

var (
	FakeDoubleFloat FakeDouble = "float"
)

type FakeHalfFloat string

var (
	FakeHalfFloatFloat FakeHalfFloat = "float"
)

type FakeScaledFloat string

var (
	FakeScaledFloatFloat FakeScaledFloat = "float"
)

type FakeDates string

var (
	FakeDateDate      FakeDates = "date"
	FakeDateTimestamp FakeDates = "timestamp"
)

type FakeIP string

var (
	FakeIPIPv4 FakeIP = "ipv4"
	FakeIPIPv6 FakeIP = "ipv6"
)

var toGenerate = map[string]string{}

func storeToGenerate(name string, t string) {
	toGenerate[name] = t
}

func generateFakeData(t string) string {
	fakeFuncs := map[string]fakeFunc{
		string(FakeEmail):     fakeEmail(),
		string(FakeFirstName): fakeFirstName(),
		string(FakeLastName):  fakeLastName(),
		string(FakeFullName):  fakeFullName(),
		string(FakeUsername):  fakeUsername(),
		string(FakePhone):     fakePhone(),
		string(FakeCountry):   fakeCountry(),
		string(FakeCity):      fakeCity(),
		string(FakeStreet):    fakeStreet(),
		string(FakeZip):       fakeZIP(),
		string(FakeUUID):      fakeUUID(),
		string(FakeURL):       fakeURL(),
		string(FakeCompany):   fakeCompany(),
		string(FakeJobTitle):  fakeJobTitle(),
		string(FakeColor):     fakeColor(),
		string(FakeIPv4):      fakeIPv4(),
		string(FakeIPv6):      fakeIPv6(),
		string(FakeBool):      fakeBoolean(),
		string(FakeInt):       fakeInteger(),
		string(FakeFloat):     fakeFloat64(),
		string(FakeDate):      fakeDate(),
		string(FakeTimestamp): fakeTimestamp(),
		string(FakeParagraph): fakeParagraph(),
	}

	fc := fakeFuncs[t]

	return fc()
}

type fakeFunc func() string

func fakeEmail() fakeFunc {
	return func() string {
		return gofakeit.Email()
	}
}

func fakeFirstName() fakeFunc {
	return func() string {
		return gofakeit.FirstName()
	}
}

func fakeLastName() fakeFunc {
	return func() string {
		return gofakeit.LastName()
	}
}
func fakeFullName() fakeFunc {
	return func() string {
		return gofakeit.Name()
	}
}

func fakeUsername() fakeFunc {
	return func() string {
		return gofakeit.Username()
	}
}

func fakePhone() fakeFunc {
	return func() string {
		return gofakeit.Phone()
	}
}

func fakeCountry() fakeFunc {
	return func() string {
		return gofakeit.Country()
	}
}

func fakeCity() fakeFunc {
	return func() string {
		return gofakeit.City()
	}
}

func fakeStreet() fakeFunc {
	return func() string {
		return gofakeit.Street()
	}
}

func fakeZIP() fakeFunc {
	return func() string {
		return gofakeit.Zip()
	}
}

func fakeUUID() fakeFunc {
	return func() string {
		return gofakeit.UUID()
	}
}

func fakeURL() fakeFunc {
	return func() string {
		return gofakeit.URL()
	}
}

func fakeCompany() fakeFunc {
	return func() string {
		return gofakeit.Company()
	}
}

func fakeJobTitle() fakeFunc {
	return func() string {
		return gofakeit.JobTitle()
	}
}

func fakeColor() fakeFunc {
	return func() string {
		return gofakeit.Color()
	}
}

func fakeIPv4() fakeFunc {
	return func() string {
		return gofakeit.IPv4Address()
	}
}

func fakeIPv6() fakeFunc {
	return func() string {
		return gofakeit.IPv6Address()
	}
}

func fakeBoolean() fakeFunc {
	return func() string {
		return fmt.Sprintf("%t", gofakeit.Bool())
	}
}

func fakeInteger() fakeFunc {
	return func() string {
		return fmt.Sprintf("%d", gofakeit.IntRange(0, 1000))
	}
}
func fakeFloat64() fakeFunc {
	return func() string {
		return fmt.Sprintf("%.2f", gofakeit.Float64Range(0.0, 1000.0))
	}
}

func fakeDate() fakeFunc {
	return func() string {
		return gofakeit.Date().Format("2006-01-02")
	}
}

func fakeTimestamp() fakeFunc {
	return func() string {
		return gofakeit.Date().Format(time.RFC3339)
	}
}

func fakeParagraph() fakeFunc {
	return func() string {
		return gofakeit.Paragraph(1, 3, 10, " ")
	}
}
