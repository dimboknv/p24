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
		expectedErrMsg string
		resp           []byte
		withErr        bool
	}{
		{"error msg", []byte(`<response><data><error message="error msg"/></data></response>`), false},
		{"error msg", []byte(`<response><data><error message="error msg"></error></data></response>`), false},
		{"", []byte(`<response><data><noterror message ="data"/></data></response>`), true},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("UnmarshalXML/%d", i), func(t *testing.T) {
			actualErr := &respDataErr{}
			err := xml.Unmarshal(c.resp, actualErr)
			require.True(t, c.withErr == (err != nil), err)
			require.EqualError(t, actualErr, c.expectedErrMsg)
		})
	}
}

func Test_respDataInfoErr(t *testing.T) {
	cases := []struct {
		expectedErrMsg string
		resp           []byte
		withErr        bool
	}{
		{"error msg", []byte(`<response><data><oper>cmt</oper><info>error msg</info></data></response>`), false},
		{"", []byte(`<response><data><oper>cmt</oper><info></info></data></response>`), true},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("UnmarshalXML/%d", i), func(t *testing.T) {
			actualErr := &respDataInfoErr{}
			err := xml.Unmarshal(c.resp, actualErr)
			require.True(t, c.withErr == (err != nil), err)
			require.EqualError(t, actualErr, c.expectedErrMsg)
		})
	}
}

func Test_respErr(t *testing.T) {
	cases := []struct {
		expectedErrMsg string
		resp           []byte
		withErr        bool
	}{
		{`For input string: "error msg"`, []byte(`<error>For input string: "error msg"</error>`), false},
		{"", []byte(`<error>For input string: "error msg"</error1>`), true},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("UnmarshalXML/%d", i), func(t *testing.T) {
			actualErr := &respErr{}
			err := xml.Unmarshal(c.resp, actualErr)
			require.True(t, c.withErr == (err != nil), err)
			require.EqualError(t, actualErr, c.expectedErrMsg)
		})
	}
}
