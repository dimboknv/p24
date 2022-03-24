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
			withErr  bool
		}{
			{[]byte("0"), 0, false},

			{[]byte("0.01"), 1, false},
			{[]byte("0.20"), 20, false},
			{[]byte("21"), 2100, false},
			{[]byte("100.09"), 10009, false},
			{[]byte("100.90"), 10090, false},
			{[]byte("100.89"), 10089, false},

			{[]byte("-0.01"), -1, false},
			{[]byte("-0.20"), -20, false},
			{[]byte("-21"), -2100, false},
			{[]byte("-100.09"), -10009, false},
			{[]byte("-100.90"), -10090, false},
			{[]byte("-100.89"), -10089, false},
			{[]byte("1.15"), 115, false},

			{[]byte(""), 0, true},
			{[]byte("31.as"), 0, true},
			{[]byte("as.31"), 0, true},
			{[]byte("hello"), 0, true},
			{[]byte("-32.-89"), 0, true},
		}
		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				var actual Amount
				err := actual.UnmarshalText(c.text)
				require.True(t, c.withErr == (err != nil), err)
				require.Equal(t, c.expected, actual)
			})
		}
	})
}

func Test_Funds(t *testing.T) {
	t.Run("MarshalText", func(t *testing.T) {
		cases := []struct {
			expected []byte
			f        Funds
			withErr  bool
		}{
			{[]byte("0 UAH"), Funds{"UAH", 0}, false},
			{[]byte("0 "), Funds{"", 0}, false},

			{[]byte("21 UAH"), Funds{"UAH", 2100}, false},
			{[]byte("100.09 USD"), Funds{"USD", 10009}, false},
			{[]byte("100.90 USD"), Funds{"USD", 10090}, false},
			{[]byte("100.89 BTC"), Funds{"BTC", 10089}, false},

			{[]byte("-21 EUR"), Funds{"EUR", -2100}, false},
			{[]byte("-100.09 ETH"), Funds{"ETH", -10009}, false},
			{[]byte("-100.90 ETH"), Funds{"ETH", -10090}, false},
			{[]byte("-100.89 RUB"), Funds{"RUB", -10089}, false},
		}
		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				actual, err := c.f.MarshalText()
				require.True(t, c.withErr == (err != nil), err)
				require.Equal(t, string(c.expected), string(actual))
			})
		}
	})
	t.Run("UnmarshalText", func(t *testing.T) {
		cases := []struct {
			text     []byte
			expected Funds
			withErr  bool
		}{
			{[]byte("0 UAH"), Funds{"UAH", 0}, false},

			{[]byte("21 UAH"), Funds{"UAH", 2100}, false},
			{[]byte("100.09 USD"), Funds{"USD", 10009}, false},
			{[]byte("100.90 USD"), Funds{"USD", 10090}, false},
			{[]byte("100.89 BTC"), Funds{"BTC", 10089}, false},

			{[]byte("-21 EUR"), Funds{"EUR", -2100}, false},
			{[]byte("-100.09 ETH"), Funds{"ETH", -10009}, false},
			{[]byte("-100.90 ETH"), Funds{"ETH", -10090}, false},
			{[]byte("-100.89 RUB"), Funds{"RUB", -10089}, false},

			{[]byte(""), Funds{}, true},
			{[]byte("31.as USD"), Funds{}, true},
			{[]byte("as.31 USD"), Funds{}, true},
			{[]byte("hello USD"), Funds{}, true},
			{[]byte("-32.-89 USD"), Funds{}, true},
		}
		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				var actual Funds
				err := actual.UnmarshalText(c.text)
				require.True(t, c.withErr == (err != nil), err)
				require.Equal(t, c.expected, actual)
			})
		}
	})
}
