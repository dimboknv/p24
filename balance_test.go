package p24

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_CardBalance(t *testing.T) {
	t.Run("MarshalXML", func(t *testing.T) {
		cases := []struct {
			expected []byte
			cb       CardBalance
		}{
			{
				[]byte(`<cardbalance><bal_date>02.09.13 21:34</bal_date><bal_dyn>dyn</bal_dyn><card><account>acc</account><card_number>num</card_number><acc_name>name</acc_name><acc_type>acctype</acc_type><currency>UAH</currency><card_type>type</card_type><main_card_number>main</main_card_number><card_stat>Status</card_stat><src>src</src></card><av_balance>1.23</av_balance><balance>3.21</balance><fin_limit>0.10</fin_limit><trade_limit>0.01</trade_limit></cardbalance>`),
				CardBalance{Date: time.Date(2013, 9, 2, 21, 34, 0, 0, kievLocation), Dyn: "dyn", Card: Card{Account: "acc", Number: "num", AccName: "name", AccType: "acctype", Currency: "UAH", Type: "type", MainCard: "main", Status: "Status", Src: "src"}, Available: 123, Balance: 321, FinLimit: 10, TradeLimit: 1},
			},
		}
		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				actual, err := xml.Marshal(c.cb)
				require.NoError(t, err)
				require.Equal(t, string(c.expected), string(actual))
			})
		}
	})

	t.Run("UnmarshalXML", func(t *testing.T) {
		cases := []struct {
			data     []byte
			expected CardBalance
			errMsg   string
		}{
			{
				[]byte(`<cardbalance><bal_date>02.09.13 21:34</bal_date><bal_dyn>dyn</bal_dyn><card><account>acc</account><card_number>num</card_number><acc_name>name</acc_name><acc_type>acctype</acc_type><currency>UAH</currency><card_type>type</card_type><main_card_number>main</main_card_number><card_stat>Status</card_stat><src>src</src></card><av_balance>1.23</av_balance><balance>3.21</balance><fin_limit>0.10</fin_limit><trade_limit>0.01</trade_limit></cardbalance>`),
				CardBalance{Date: time.Date(2013, 9, 2, 21, 34, 0, 0, kievLocation), Dyn: "dyn", Card: Card{Account: "acc", Number: "num", AccName: "name", AccType: "acctype", Currency: "UAH", Type: "type", MainCard: "main", Status: "Status", Src: "src"}, Available: 123, Balance: 321, FinLimit: 10, TradeLimit: 1},
				"",
			},
			{
				[]byte(`<cardbalance><bal_date>02.09.13 21:34</bal_date><bal_dyn>dyn</bal_dyn><card><acc`),
				CardBalance{},
				"unexpected EOF",
			},
			{
				[]byte(`<cardbalance><bal_date>02-09-13 21:34</bal_date></cardbalance>`),
				CardBalance{},
				"parsing time",
			},
		}
		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				var actual CardBalance
				if err := xml.Unmarshal(c.data, &actual); err != nil {
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

func TestClient_GetCardBalance(t *testing.T) {
	m := Merchant{"id", "pass"}
	cases := []struct {
		opts     BalanceOpts
		reqBody  []byte
		respBody []byte
		expected CardBalance
		errMsg   string
	}{
		{
			opts: BalanceOpts{
				CardNumber: "1234567890123456",
				Country:    "USA",
				CommonOpts: DefaultCommonOpts(),
			},
			reqBody:  []byte(xml.Header + `<request version="1.0"><merchant><id>id</id><signature>7a0d071af8d0ebf513bb30ab74e1fd5f172abe82</signature></merchant><data><payment id=""><prop name="cardnum" value="1234567890123456"></prop><prop name="country" value="USA"></prop></payment><oper>cmt</oper><wait>0</wait><test>0</test></data></request>`),
			respBody: []byte(`<?xml version="1.0" encoding="UTF-8"?><response version="1.0"><merchant><id>id</id><signature>13dccaec0c5303ae43217901d9a61cb94a132c19</signature></merchant><data><oper>cmt</oper><info><cardbalance><bal_date>01.01.21 01:01</bal_date><bal_dyn></bal_dyn><card><account></account><card_number>1234567890123456</card_number><acc_name></acc_name><acc_type></acc_type><currency></currency><card_type></card_type><main_card_number>1234567890123456</main_card_number><card_stat></card_stat><src></src></card><av_balance>0.01</av_balance><balance>0</balance><fin_limit>0</fin_limit><trade_limit>0.02</trade_limit></cardbalance></info></data></response>`),
			expected: CardBalance{
				Card: Card{
					Number:   "1234567890123456",
					MainCard: "1234567890123456",
				},
				Date:       time.Date(2021, 1, 1, 1, 1, 0, 0, kievLocation),
				Available:  1,
				TradeLimit: 2,
			},
		},
		{
			opts: BalanceOpts{
				CardNumber: "sdalkfj",
				Country:    "USA",
			},
			errMsg: "invalid card number",
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

			cli := Client{do, nil, m}
			actual, err := cli.GetCardBalance(context.Background(), c.opts)
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
