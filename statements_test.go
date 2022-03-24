package p24

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_StatementsOpts(t *testing.T) {
	cases := []struct {
		errMsg string
		opts   StatementsOpts
	}{
		{
			"date range should be no longer than 90 days",
			StatementsOpts{
				StartDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
				EndDate:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local),
				CardNumber: "1111111111111111",
			},
		},
		{
			"date range should be with start date <= end date",
			StatementsOpts{
				StartDate:  time.Date(2000, 3, 3, 0, 0, 0, 0, time.Local),
				EndDate:    time.Date(2000, 2, 1, 0, 0, 0, 0, time.Local),
				CardNumber: "1111111111111111",
			},
		},
		{
			"invalid card number: should be sixteen length",
			StatementsOpts{
				StartDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
				EndDate:    time.Date(2000, 2, 1, 0, 0, 0, 0, time.Local),
				CardNumber: "not a card",
			},
		},
		{
			"",
			StatementsOpts{
				StartDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
				EndDate:    time.Date(2000, 2, 1, 0, 0, 0, 0, time.Local),
				CardNumber: "0101011111111111",
			},
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("Validate/%d", i), func(t *testing.T) {
			if err := c.opts.Validate(); err != nil {
				require.Contains(t, err.Error(), c.errMsg, "errors not equal")
			} else {
				require.Empty(t, c.errMsg, "should be an error")
			}
		})
	}
}

func Test_Statement(t *testing.T) {
	t.Run("MarshalXML", func(t *testing.T) {
		cases := []struct {
			expected  []byte
			statement Statement
		}{
			{
				[]byte(`<statement trantime="21:34:00" trandate="2013-09-02" card="5168742060221193" appcode="801111" terminal="Пополнение мобильного" description="description" amount="0.10 UAH" cardamount="-0.10 UAH" rest="1.15 UAH"></statement>`),
				Statement{"5168742060221193", "801111", time.Date(2013, 9, 2, 21, 34, 0, 0, kievLocation), "Пополнение мобильного", "description", Funds{"UAH", Amount(10)}, Funds{"UAH", Amount(-10)}, Funds{"UAH", Amount(115)}},
			},
		}
		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				actual, err := xml.Marshal(c.statement)
				require.NoError(t, err)
				require.Equal(t, string(c.expected), string(actual))
			})
		}
	})

	t.Run("UnmarshalXML", func(t *testing.T) {
		cases := []struct {
			data     []byte
			expected Statement
			withErr  bool
		}{
			{
				[]byte(`<statement trandate="2013-09-02" trantime="21:34:00" card="5168742060221193" appcode="801111" terminal="Пополнение мобильного" description="description" amount="0.10 UAH" cardamount="-0.10 UAH" rest="1.15 UAH"></statement>`),
				Statement{"5168742060221193", "801111", time.Date(2013, 9, 2, 21, 34, 0, 0, kievLocation), "Пополнение мобильного", "description", Funds{"UAH", Amount(10)}, Funds{"UAH", Amount(-10)}, Funds{"UAH", Amount(115)}},
				false,
			},
			{
				[]byte(`<statement card="1" ></statement1>`),
				Statement{},
				true,
			},
			{
				[]byte(`<statement amount="0.10" card="51687420"></statement>`),
				Statement{},
				true,
			},
			{
				[]byte(``),
				Statement{},
				true,
			},
		}
		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				var actual Statement
				err := xml.Unmarshal(c.data, &actual)
				require.True(t, c.withErr == (err != nil), err)
				require.Equal(t, c.expected, actual)
			})
		}
	})
}
