package p24

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Amount(t *testing.T) {
	t.Run("MarshalText", func(t *testing.T) {
		cases := []struct {
			expected []byte
			a        Amount
		}{
			{[]byte("0"), 0},

			{[]byte("0.01"), 1},
			{[]byte("0.20"), 20},
			{[]byte("21"), 2100},
			{[]byte("100.09"), 10009},
			{[]byte("100.90"), 10090},
			{[]byte("100.89"), 10089},

			{[]byte("-0.01"), -1},
			{[]byte("-0.20"), -20},
			{[]byte("-21"), -2100},
			{[]byte("-100.09"), -10009},
			{[]byte("-100.90"), -10090},
			{[]byte("-100.89"), -10089},
		}
		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				actual, err := c.a.MarshalText()
				require.NoError(t, err)
				require.Equal(t, string(c.expected), string(actual))
			})
		}
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		cases := []struct {
			text     []byte
			expected Amount
			errMsg   string
		}{
			{[]byte("0"), 0, ""},

			{[]byte("0.01"), 1, ""},
			{[]byte("0.20"), 20, ""},
			{[]byte("21"), 2100, ""},
			{[]byte("100.09"), 10009, ""},
			{[]byte("100.90"), 10090, ""},
			{[]byte("100.89"), 10089, ""},

			{[]byte("-0.01"), -1, ""},
			{[]byte("-0.20"), -20, ""},
			{[]byte("-21"), -2100, ""},
			{[]byte("-100.09"), -10009, ""},
			{[]byte("-100.90"), -10090, ""},
			{[]byte("-100.89"), -10089, ""},
			{[]byte("1.15"), 115, ""},

			{[]byte(""), 0, "parsing"},
			{[]byte("31.as"), 0, "parsing"},
			{[]byte("as.31"), 0, "parsing"},
			{[]byte("hello"), 0, "parsing"},
			{[]byte("-32.-89"), 0, "parsing"},
		}
		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				var actual Amount
				if err := actual.UnmarshalText(c.text); err != nil {
					require.NotEmpty(t, c.errMsg, err)
					require.ErrorContains(t, err, c.errMsg)
					return
				}
				require.Empty(t, c.errMsg)
				require.Equal(t, c.expected, actual)
			})
		}
	})
}

func Test_Funds(t *testing.T) {
	t.Run("MarshalText", func(t *testing.T) {
		cases := []struct {
			expected []byte
			funds    Funds
		}{
			{[]byte("0 UAH"), Funds{"UAH", 0}},
			{[]byte("0 "), Funds{"", 0}},

			{[]byte("21 UAH"), Funds{"UAH", 2100}},
			{[]byte("100.09 USD"), Funds{"USD", 10009}},
			{[]byte("100.90 USD"), Funds{"USD", 10090}},
			{[]byte("100.89 BTC"), Funds{"BTC", 10089}},

			{[]byte("-21 EUR"), Funds{"EUR", -2100}},
			{[]byte("-100.09 ETH"), Funds{"ETH", -10009}},
			{[]byte("-100.90 ETH"), Funds{"ETH", -10090}},
			{[]byte("-100.89 RUB"), Funds{"RUB", -10089}},
		}
		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				actual, _ := c.funds.MarshalText()
				require.Equal(t, string(c.expected), string(actual))
			})
		}
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		cases := []struct {
			text     []byte
			expected Funds
			errMsg   string
		}{
			{[]byte("0 UAH"), Funds{"UAH", 0}, ""},

			{[]byte("21 UAH"), Funds{"UAH", 2100}, ""},
			{[]byte("100.09 USD"), Funds{"USD", 10009}, ""},
			{[]byte("100.90 USD"), Funds{"USD", 10090}, ""},
			{[]byte("100.89 BTC"), Funds{"BTC", 10089}, ""},

			{[]byte("-21 EUR"), Funds{"EUR", -2100}, ""},
			{[]byte("-100.09 ETH"), Funds{"ETH", -10009}, ""},
			{[]byte("-100.90 ETH"), Funds{"ETH", -10090}, ""},
			{[]byte("-100.89 RUB"), Funds{"RUB", -10089}, ""},

			{[]byte(""), Funds{}, "parsing"},
			{[]byte("31.as USD"), Funds{}, "parsing"},
			{[]byte("as.31 USD"), Funds{}, "parsing"},
			{[]byte("hello USD"), Funds{}, "parsing"},
			{[]byte("-32.-89 USD"), Funds{}, "parsing"},
		}
		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				var actual Funds
				if err := actual.UnmarshalText(c.text); err != nil {
					require.NotEmpty(t, c.errMsg)
					require.ErrorContains(t, err, c.errMsg)
					return
				}
				require.Empty(t, c.errMsg)
				require.Equal(t, c.expected, actual)
			})
		}
	})
}
