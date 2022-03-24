package p24

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CheckCardNumber(t *testing.T) {
	cases := []struct {
		errMsg string
		card   string
	}{
		{"should contains digits only", "11111not1a1card1"},
		{"should be sixteen length", "111111111111111"},
		{"should be sixteen length", "11111111111111111"},
		{"", "1111111111111111"},
	}

	for i, c := range cases {
		c := c
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := CheckCardNumber(c.card)
			if err != nil {
				require.EqualError(t, err, c.errMsg)
			} else {
				require.Empty(t, c.errMsg)
			}
		})
	}
}

func Test_dataTagContent(t *testing.T) {
	cases := []struct {
		errMsg   string
		payload  []byte
		expected []byte
	}{
		{"", []byte("<data>payload</data>"), []byte("payload")},
		{"", []byte("<data><data>payload</data></data>"), []byte("<data>payload</data>")},
		{"not found", []byte("payload</data></data>"), nil},
		{"not found", []byte("<data><data>payload"), nil},
	}

	for i, c := range cases {
		c := c
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual, err := dataTagContent(c.payload)
			if err != nil {
				require.EqualError(t, err, c.errMsg, "errors not equal")
			} else {
				require.Empty(t, c.errMsg, "should be an error")
				require.Equal(t, string(c.expected), string(actual), "payloads not equal")
			}
		})
	}
}
