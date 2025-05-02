package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xoticdsign/porter2/internal/tests/suite"
)

func TestLocationFromFile_Functional(t *testing.T) {
	s, err := suite.New(t, true)
	if err != nil {
		panic(err)
	}

	cases := []struct {
		name        string
		in          string
		expectedErr bool
	}{
		{
			name:        "happy case",
			in:          "t.json",
			expectedErr: false,
		},
		{
			name:        "file does not exist case",
			in:          "wrong.json",
			expectedErr: true,
		},
	}

	for _, cs := range cases {
		s.T.Run(cs.name, func(t *testing.T) {
			_, err := os.Create("t.json")
			defer os.Remove("t.json")

			assert.NoError(t, err)

			fc := s.Porter.Documents.Origin.FromFile(cs.in)

			_, err = fc(s.Temp)

			switch {
			case cs.expectedErr:
				assert.Error(t, err)

			case !cs.expectedErr:
				assert.NoError(t, err)
			}
		})
	}
}
