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

func equalXML(t require.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) {
	actualXML, err := xml.Marshal(actual)
	require.NoError(t, err)
	expectedXML, err := xml.Marshal(expected)
	require.NoError(t, err)
	require.Equal(t, string(expectedXML), string(actualXML), msgAndArgs...)
}

func TestClient_DoContext(t *testing.T) {
	cases := []struct {
		expected Response
		m        Merchant
		body     []byte
		code     int
		withErr  bool
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
			body: []byte(xml.Header + `<response><data><info><test>test</test></info><oper>cmt</oper></data><merchant><id>id</id><signature>ad67cf1c11e0f87bedac2c9bb260e3abf54e9862</signature></merchant></response>`),
			m:    Merchant{"id", "pass"},
			code: 200,
		},
		{ // can`t unmarshal data
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
			body:    []byte(xml.Header + `<response><data><info><test1>test1</test1></info><oper>cmt</oper></data><merchant><id>id</id><signature>ad67cf1c11e0f87bedac2c9bb260e3abf54e9862</signature></merchant></response>`),
			m:       Merchant{"id", "pass"},
			withErr: true,
			code:    200,
		},
		{ // invalid signature
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
			body:    []byte(xml.Header + `<response><data><info><test>test</test></info><oper>cmt</oper></data><merchant><id>id</id><signature>ad67cf1c11e0f87bedac2c9bb260e3abf54e9862</signature></merchant></response>`),
			m:       Merchant{"id", "other pass"},
			withErr: true,
			code:    200,
		},
		{ // xml response with err
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
			body:    []byte(xml.Header + `<response><data><info>some error</info><oper>cmt</oper></data><merchant><id>id</id><signature>ad67cf1c11e0f87bedac2c9bb260e3abf54e9862</signature></merchant></response>`),
			m:       Merchant{"id", "pass"},
			withErr: true,
			code:    200,
		},
		{ // invalid http code
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
			body:    []byte(xml.Header + `<response><data><info><test>test</test></info><oper>cmt</oper></data><merchant><id>id</id><signature>ad67cf1c11e0f87bedac2c9bb260e3abf54e9862</signature></merchant></response>`),
			m:       Merchant{"id", "pass"},
			withErr: true,
			code:    400,
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

			cli := Client{do, nil, c.m}
			actual := c.expected
			err := cli.DoContext(context.Background(), url, method, req, &actual)
			require.True(t, c.withErr == (err != nil), err)
			require.Equal(t, c.expected, actual)
		})
	}
}
