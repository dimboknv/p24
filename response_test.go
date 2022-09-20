package p24

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Response(t *testing.T) {
	responses := []Response{
		{
			Version: "1.0",
			XMLName: xml.Name{Local: "response"},
			MerchantSign: MerchantSign{
				ID:   "id",
				Sign: "sig",
			},
			Data: ResponseData{
				Info: struct {
					Statements Statements `xml:"statements"`
				}{
					Statements: Statements{
						Status: "ok",
						Credit: 1,
					},
				},
			},
		},
		{
			Version: "1.0",
			XMLName: xml.Name{Local: "response"},
			MerchantSign: MerchantSign{
				ID:   "id",
				Sign: "sig",
			},
			Data: ResponseData{
				Info: struct {
					Statements Statements `xml:"statements"`
				}{
					Statements: Statements{
						Status: "ok",
						Credit: 1,
					},
				},
			},
		},
		{
			Version: "1.0",
			XMLName: xml.Name{Local: "response"},
			MerchantSign: MerchantSign{
				ID:   "id",
				Sign: "sig",
			},
			Data: ResponseData{
				Info: "string",
			},
		},
		{
			Version: "1.0",
			XMLName: xml.Name{Local: "response"},
			MerchantSign: MerchantSign{
				ID:   "id",
				Sign: "sig",
			},
			Data: ResponseData{
				Info: struct {
					Test string `xml:"test"`
				}{
					Test: "test",
				},
			},
		},
	}

	for i, expected := range responses {
		expected := expected
		t.Run(fmt.Sprintf("UnmarshalXML/%d", i), func(t *testing.T) {
			data, err := xml.Marshal(expected)
			require.NoError(t, err, "response Marshal error")

			actual := Response{}
			actual.Data.Info = expected.Data.Info
			require.NoError(t, xml.Unmarshal(data, &actual), "response Unmarshal error")
			require.EqualValues(t, expected, actual, "responses not equal")
		})
	}
}

func Test_respDataErr(t *testing.T) {
	cases := []struct {
		expectedMsg string
		resp        []byte
	}{
		{"error msg", []byte(`<response><data><error message="error msg"/></data></response>`)},
		{"error msg", []byte(`<response><data><error message="error msg"></error></data></response>`)},
		{"", []byte(`<response><data><noterror message ="data"/></data></response>`)},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("UnmarshalXML/%d", i), func(t *testing.T) {
			actual := &respDataErr{}
			if err := xml.Unmarshal(c.resp, actual); err != nil {
				require.Empty(t, c.expectedMsg)
			}
			require.EqualError(t, actual, c.expectedMsg)
		})
	}
}

func Test_respDataInfoErr(t *testing.T) {
	cases := []struct {
		expectedMsg string
		resp        []byte
	}{
		{"error msg", []byte(`<response><data><oper>cmt</oper><info>error msg</info></data></response>`)},
		{"", []byte(`<response><data><oper>cmt</oper><info></info></data></response>`)},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("UnmarshalXML/%d", i), func(t *testing.T) {
			actual := &respDataInfoErr{}
			if err := xml.Unmarshal(c.resp, actual); err != nil {
				require.Empty(t, c.expectedMsg)
			}
			require.EqualError(t, actual, c.expectedMsg)
		})
	}
}

func Test_respErr(t *testing.T) {
	cases := []struct {
		expectedMsg string
		resp        []byte
	}{
		{`For input string: "error msg"`, []byte(`<error>For input string: "error msg"</error>`)},
		{"", []byte(`<error>For input string: "error msg"</error1>`)},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("UnmarshalXML/%d", i), func(t *testing.T) {
			actual := &respErr{}
			if err := xml.Unmarshal(c.resp, actual); err != nil {
				require.Empty(t, c.expectedMsg)
			}
			require.EqualError(t, actual, c.expectedMsg)
		})
	}
}
