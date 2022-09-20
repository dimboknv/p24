package p24

import (
	"encoding/xml"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_xmlResp(t *testing.T) {
	t.Run("CheckErr", func(t *testing.T) {
		cases := []struct {
			errMsg  string
			xmlResp xmlResp
		}{
			{`For input string: "err msg"`, []byte(`<error>For input string: "err msg"</error>`)},
			{"err msg", []byte(`<response><data><error message="err msg"></error></data></response>`)},
			{"err msg", []byte(`<response><data><oper>cmt</oper><info>err msg</info></data></response>`)},
			{"", []byte(`<response><noterror><oper>cmt</oper><info></info></noterror></response>`)},
		}

		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				if err := c.xmlResp.CheckErr(); err != nil {
					require.EqualError(t, err, c.errMsg)
				} else {
					require.Empty(t, c.errMsg)
				}
			})
		}
	})

	t.Run("CheckContent", func(t *testing.T) {
		cases := []struct {
			xmlResp xmlResp
			errMsg  string
		}{
			{[]byte(`<invalid_resp>invalid data</invalid_resp>`), "can`t unmarshal common"},
			{[]byte(`<response><data1><oper>cmt</oper><info>123</info></data1></response>`), "invalid '<data>' tag"},
			{[]byte(`<response><data><oper>cmt</oper><info>123</info></data></response>`), ""},
		}

		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				if err := c.xmlResp.CheckContent(); err != nil {
					require.NotEmpty(t, c.errMsg)
					require.ErrorContains(t, err, c.errMsg)
					return
				}
				require.Empty(t, c.errMsg)
			})
		}
	})

	t.Run("VerifySign", func(t *testing.T) {
		signer := Merchant{"id", "pass"}
		type info struct {
			A int `xml:"a"`
		}
		cases := []struct {
			signer Merchant
			data   []byte
			info   info
			errMsg string
		}{
			{signer, []byte("<info><a>1</a></info><oper></oper>"), info{1}, ""},
			{signer, []byte("<info><a>2</a></info><oper></oper>"), info{2}, ""},
			{signer, []byte("other expectedMsg"), info{1}, "invalid signature"},
			{Merchant{"id", "other pass"}, []byte("<info><a>1</a></info><oper></oper>"), info{1}, "invalid signature"},
			{Merchant{"other id", " ass"}, []byte("<info><a>1</a></info><oper></oper>"), info{1}, "invalid signature"},
		}

		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				data, err := xml.Marshal(Response{
					Data:         ResponseData{Info: c.info},
					MerchantSign: signer.Sign(c.data),
				})
				require.NoError(t, err)

				if err := xmlResp(data).VerifySign(c.signer); err != nil {
					require.NotEmpty(t, c.errMsg)
					require.ErrorContains(t, err, c.errMsg)
					return
				}
				require.Empty(t, c.errMsg)
			})
		}
	})
}
