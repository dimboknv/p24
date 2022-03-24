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
			{"", []byte(`<response><expectedErrMsg><oper>cmt</oper><info></info></expectedErrMsg></response>`)},
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
			withErr bool
		}{
			{[]byte(`<invalid_resp>invalid data</invalid_resp>`), true},
			{[]byte(`<response><data1><oper>cmt</oper><info>123</info></data1></response>`), true},
			{[]byte(`<response><data><oper>cmt</oper><info>123</info></data></response>`), false},
		}

		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				err := c.xmlResp.CheckContent()
				require.True(t, c.withErr == (err != nil), err)
			})
		}
	})

	t.Run("VerifySign", func(t *testing.T) {
		signer := Merchant{"id", "pass"}
		type info struct {
			A int `xml:"a"`
		}
		cases := []struct {
			signer  Merchant
			data    []byte
			info    info
			withErr bool
		}{
			{signer, []byte("<info><a>1</a></info><oper></oper>"), info{1}, false},
			{signer, []byte("<info><a>2</a></info><oper></oper>"), info{2}, false},
			{signer, []byte("other expectedErrMsg"), info{1}, true},
			{Merchant{"id", "other pass"}, []byte("<info><a>1</a></info><oper></oper>"), info{1}, true},
			{Merchant{"other id", " ass"}, []byte("<info><a>1</a></info><oper></oper>"), info{1}, true},
		}

		for i, c := range cases {
			c := c
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				data, err := xml.Marshal(Response{
					Data:         ResponseData{Info: c.info},
					MerchantSign: signer.Sign(c.data),
				})
				require.NoError(t, err)

				err = xmlResp(data).VerifySign(c.signer)
				require.True(t, c.withErr == (err != nil), err)
			})
		}
	})
}
