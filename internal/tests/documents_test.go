package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xoticdsign/porter/internal/tests/suite"
)

func TestLocationFromFile_Functional(t *testing.T) {
	s, err := suite.New(t, true)
	if err != nil {
		panic(err)
	}

	cases := []struct {
		name        string
		in          string
		expectedErr error
	}{
		{},
	}

	for _, cs := range cases {
		s.T.Run(cs.name, func(t *testing.T) {
			fc := s.Porter.Documents.Origin.FromFile(cs.in)

			fc()

			assert.Equal(t, cs.expected, n)
		})
	}
}
