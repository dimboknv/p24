package p24

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_DoContext(t *testing.T) {
	cases := []struct {
		expected Response
		merchant Merchant
		body     []byte
		code     int
		errMsg   string
	}{
		{ // ok
			expected: Response{
				XMLName:      xml.Name{Local: "response"},
				MerchantSign: MerchantSign{"id", "ad67cf1c11e0f87bedac2c9bb260e3abf54e9862"},
				Data: ResponseData{
					Oper: defaultOper,
					Info: struct {
						Test string `xml:"test"`
					}{Test: "test"},
				},
			},
			body:     []byte(xml.Header + `<response><data><info><test>test</test></info><oper>cmt</oper></data><merchant><id>id</id><signature>ad67cf1c11e0f87bedac2c9bb260e3abf54e9862</signature></merchant></response>`),
			merchant: Merchant{"id", "pass"},
			code:     200,
		},
		{
			expected: Response{
				XMLName:      xml.Name{Local: "response"},
				MerchantSign: MerchantSign{"id", "ad67cf1c11e0f87bedac2c9bb260e3abf54e9862"},
				Data: ResponseData{
					Oper: defaultOper,
					Info: struct {
						Test string `xml:"test1"`
					}{Test: "test"},
				},
			},
			body:     []byte(xml.Header + `<response><data><info><test>test</test></info><oper>cmt</oper></data><merchant><id>id</id><signature>ad67cf1c11e0f87bedac2c9bb260e3abf54e9862</signature></merchant></response>`),
			merchant: Merchant{"id", "pass"},
			code:     200,
			errMsg:   "can`t unmarshal xml response",
		},
		{
			expected: Response{
				XMLName:      xml.Name{Local: "response"},
				MerchantSign: MerchantSign{"id", "ad67cf1c11e0f87bedac2c9bb260e3abf54e9862"},
				Data: ResponseData{
					Oper: defaultOper,
					Info: struct {
						Test string `xml:"test"`
					}{Test: "test"},
				},
			},
			body:     []byte(xml.Header + `<response><data><info><test>test</test></info><oper>cmt</oper><merchant><id>id</id><signature>ad67cf1c11e0f87bedac2c9bb260e3abf54e9862</signature></merchant></response>`),
			merchant: Merchant{"id", "pass"},
			code:     200,
			errMsg:   "unexpected xml response content",
		},
		{
			expected: Response{
				XMLName:      xml.Name{Local: "response"},
				MerchantSign: MerchantSign{"id", "ad67cf1c11e0f87bedac2c9bb260e3abf54e9862"},
				Data: ResponseData{
					Oper: defaultOper,
					Info: struct {
						Test string `xml:"test"`
					}{Test: "test"},
				},
			},
			body:     []byte(xml.Header + `<response><data><info><test>test</test></info><oper>cmt</oper></data><merchant><id>id</id><signature>ad67cf1c11e0f87bedac2c9bb260e3abf54e9862</signature></merchant></response>`),
			merchant: Merchant{"id", "other pass"},
			errMsg:   "invalid signature",
			code:     200,
		},
		{
			expected: Response{
				XMLName:      xml.Name{Local: "response"},
				MerchantSign: MerchantSign{"id", "ad67cf1c11e0f87bedac2c9bb260e3abf54e9862"},
				Data: ResponseData{
					Oper: defaultOper,
					Info: struct {
						Test string `xml:"test"`
					}{Test: "test"},
				},
			},
			body:     []byte(xml.Header + `<response><data><info>some error</info><oper>cmt</oper></data><merchant><id>id</id><signature>ad67cf1c11e0f87bedac2c9bb260e3abf54e9862</signature></merchant></response>`),
			merchant: Merchant{"id", "pass"},
			errMsg:   "xml response with err",
			code:     200,
		},
		{
			expected: Response{
				XMLName:      xml.Name{Local: "response"},
				MerchantSign: MerchantSign{"id", "ad67cf1c11e0f87bedac2c9bb260e3abf54e9862"},
				Data: ResponseData{
					Oper: defaultOper,
					Info: struct {
						Test string `xml:"test"`
					}{Test: "test"},
				},
			},
			body:     []byte(xml.Header + `<response><data><info><test>test</test></info><oper>cmt</oper></data><merchant><id>id</id><signature>ad67cf1c11e0f87bedac2c9bb260e3abf54e9862</signature></merchant></response>`),
			merchant: Merchant{"id", "pass"},
			errMsg:   "unexpected http status code",
			code:     400,
		},
	}
	url, method, req := "http://localhost", "POST", Request{}
	exceptedBody, _ := xml.Marshal(req)
	exceptedBody = []byte(xml.Header + string(exceptedBody))

	for i, c := range cases {
		c := c
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var do DoFunc = func(req *http.Request) (*http.Response, error) {
				require.Equal(t, url, req.URL.String())
				require.Equal(t, method, req.Method)

				actualBody, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				require.Equal(t, string(exceptedBody), string(actualBody))

				tr := httptest.NewRecorder()
				_, _ = tr.Write(c.body)
				tr.Code = c.code
				return tr.Result(), nil
			}

			cli := Client{do, nil, c.merchant}
			actual := c.expected
			if err := cli.DoContext(context.Background(), url, method, req, &actual); err != nil {
				require.NotEmpty(t, c.errMsg)
				require.ErrorContains(t, err, c.errMsg)
				return
			}
			require.Empty(t, c.errMsg)
			require.Equal(t, c.expected, actual)
		})
	}
}
