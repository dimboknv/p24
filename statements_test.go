package p24

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestClient_GetStatements(t *testing.T) {
	merchant := Merchant{"id", "pass"}
	cases := []struct {
		opts     StatementsOpts
		reqBody  []byte
		respBody []byte
		expected Statements
		errMsg   string
	}{
		// valid
		{
			opts: StatementsOpts{
				StartDate:  time.Date(2021, 1, 1, 0, 0, 0, 0, kievLocation),
				EndDate:    time.Date(2021, 1, 2, 0, 0, 0, 0, kievLocation),
				CardNumber: "1234567890123456",
			},
			reqBody:  []byte(xml.Header + `<request version="1.0"><merchant><id>id</id><signature>6295880c80459b0b50d208de152dc1000bde1708</signature></merchant><data><payment id=""><prop name="sd" value="01.01.2021"></prop><prop name="ed" value="02.01.2021"></prop><prop name="card" value="1234567890123456"></prop></payment><oper>cmt</oper><wait>0</wait><test>0</test></data></request>`),
			respBody: []byte(`<?xml version="1.0" encoding="UTF-8"?><response version="1.0"><merchant><id>id</id><signature>68ca17bc2ca05d70ec51611dfd6a84cf1fcc388f</signature></merchant><data><oper>cmt</oper><info><statements status="excellent" credit="0.0" debet="5.5"><statement card="1234567890123456" appcode="12345" trandate="2021-01-01" trantime="05:05:05" amount="5.50 UAH" cardamount="-5.50 UAH" rest="10 UAH" terminal="PrivatBank, 123" description="test"/></statements></info></data></response>`),
			expected: Statements{
				Status: "excellent",
				Statements: []Statement{
					{
						Card:        "1234567890123456",
						Appcode:     "12345",
						Date:        time.Date(2021, 1, 1, 5, 5, 5, 0, kievLocation),
						Terminal:    "PrivatBank, 123",
						Description: "test",
						Amount: Funds{
							Currency: "UAH",
							Amount:   550,
						},
						CardAmount: Funds{
							Currency: "UAH",
							Amount:   -550,
						},
						Rest: Funds{
							Currency: "UAH",
							Amount:   1000,
						},
					},
				},
				Credit: 0,
				Debet:  550,
			},
		},

		// invalid opts
		{
			opts: StatementsOpts{
				StartDate:  time.Date(2021, 1, 1, 1, 1, 0, 0, kievLocation),
				EndDate:    time.Date(2021, 1, 1, 1, 1, 0, 0, kievLocation),
				CardNumber: "err",
			},
			errMsg: "invalid card number",
		},
		{
			opts: StatementsOpts{
				StartDate:  time.Date(2021, 1, 1, 0, 0, 0, 0, kievLocation),
				EndDate:    time.Date(2021, 1, 2, 0, 0, 0, 0, kievLocation),
				CardNumber: "1234567890123456",
			},
			errMsg:   "invalid signature",
			reqBody:  []byte(xml.Header + `<request version="1.0"><merchant><id>id</id><signature>6295880c80459b0b50d208de152dc1000bde1708</signature></merchant><data><payment id=""><prop name="sd" value="01.01.2021"></prop><prop name="ed" value="02.01.2021"></prop><prop name="card" value="1234567890123456"></prop></payment><oper>cmt</oper><wait>0</wait><test>0</test></data></request>`),
			respBody: []byte(`<?xml version="1.0" encoding="UTF-8"?><response version="1.0"><merchant><id>id</id><signature>61ca17bc2ca05d70ec51611dfd6a84cf1fcc388f</signature></merchant><data><oper>cmt</oper><info><statements status="excellent" credit="0.0" debet="5.5"><statement card="1234567890123456" appcode="12345" trandate="2021-01-01" trantime="05:05:05" amount="5.50 UAH" cardamount="-5.50 UAH" rest="10 UAH" terminal="PrivatBank, 123" description="test"/></statements></info></data></response>`),
		},
	}
	for i, c := range cases {
		c := c
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var do DoFunc = func(req *http.Request) (*http.Response, error) {
				reqBody, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				require.Equal(t, string(c.reqBody), string(reqBody))

				tr := httptest.NewRecorder()
				_, _ = tr.Write(c.respBody)
				return tr.Result(), nil
			}

			cli := Client{do, nil, merchant}
			actual, err := cli.GetStatements(context.Background(), c.opts)
			if err != nil {
				require.NotEmpty(t, c.errMsg)
				require.ErrorContains(t, err, c.errMsg)
				return
			}
			require.Empty(t, c.errMsg)
			require.Equal(t, c.expected, actual)
		})
	}
}
