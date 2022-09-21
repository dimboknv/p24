package p24

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Merchant(t *testing.T) {
	t.Run("Sign", func(t *testing.T) {
		cases := []struct {
			pass     string
			data     string
			expected string
		}{
			{"some pass", "some data", "b015cc47af24cee37c979fbe1744a9c8eda825d7"},
			{"some pass", "some data1", "9ebbe65a8ce739bcd47d7907f45475681c830fb9"},
			{"some pass1", "some data", "27dcc7a4625675bf4cb0de0d3231ee8c5598a49d"},
		}

		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				require.Equal(t, c.expected, Merchant{Pass: c.pass}.Sign([]byte(c.data)).Sign)
			})
		}
	})
	t.Run("VerifySign", func(t *testing.T) {
		cases := []struct {
			merchant Merchant
			data     string
			dataSign MerchantSign
			errMsg   string
		}{
			{Merchant{"id", "some pass"}, "some data", MerchantSign{"id", "b015cc47af24cee37c979fbe1744a9c8eda825d7"}, ""},
			{Merchant{"id", "some pass"}, "some data", MerchantSign{"notvalidid", "b015cc47af24cee37c979fbe1744a9c8eda825d7"}, "invalid signature"},
			{Merchant{"id", "some pass"}, "some data", MerchantSign{"id", "notvalidsignature"}, "invalid signature"},
		}

		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				err := c.merchant.VerifySign([]byte(c.data), c.dataSign)
				if err != nil {
					require.EqualError(t, err, c.errMsg)
				} else {
					require.Empty(t, c.errMsg)
				}
			})
		}
	})
}
