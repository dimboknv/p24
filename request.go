package p24

import (
	"encoding/xml"
)

const (
	defaultOper = "cmt"
)

// CommonOpts store common request options
// used across all p24 requests
type CommonOpts struct {
	Oper string `xml:"oper"`
	Wait int    `xml:"wait"`
	Test int    `xml:"test"`
}

// DefaultCommonOpts returns default CommonOpts
func DefaultCommonOpts() CommonOpts {
	return CommonOpts{Oper: defaultOper}
}

// RequestData is struct for mapping p24 request data
type RequestData struct {
	XMLName xml.Name `xml:"data"`
	Payment struct {
		ID   string `xml:"id,attr"`
		Prop []struct {
			Name  string `xml:"name,attr"`
			Value string `xml:"value,attr"`
		} `xml:"prop"`
	} `xml:"payment"`
	CommonOpts
}

// Request is struct for mapping p24 api request
type Request struct {
	XMLName      xml.Name     `xml:"request"`
	Version      string       `xml:"version,attr"`
	MerchantSign MerchantSign `xml:"merchant"`
	Data         RequestData  `xml:"data"`
}

// NewRequest returns Request with MerchantSign of reqData
func NewRequest(m Merchant, reqData RequestData) Request {
	if zero := (CommonOpts{}); reqData.CommonOpts == zero {
		reqData.CommonOpts = DefaultCommonOpts()
	}

	xmlData, _ := xml.Marshal(reqData)
	dataTag, _ := dataTagContent(xmlData)
	return Request{
		Version:      "1.0",
		MerchantSign: m.Sign(dataTag),
		Data:         reqData,
	}
}
